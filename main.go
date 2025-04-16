package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

type DomainStats struct {
	Success int
	Total   int
}

var stats = make(map[string]*DomainStats)

// Replacing hardcoded values with constants
const (
	RequestTimeout    = 500 * time.Millisecond
	CheckInterval     = 15 * time.Second
	DefaultHTTPMethod = "GET"
)

func checkHealth(endpoint Endpoint) {
	var client = &http.Client{
		Timeout: RequestTimeout,
	}

	bodyBytes, err := json.Marshal(endpoint)
	if err != nil {
		return
	}

	//validation check if no method is given to clarify default is GET
	if endpoint.Method == "" {
		endpoint.Method = DefaultHTTPMethod
	}
	reqBody := bytes.NewReader(bodyBytes)

	req, err := http.NewRequest(endpoint.Method, endpoint.URL, reqBody)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
	}

	start := time.Now()
	resp, err := client.Do(req)
	// to calculate total time passed since req
	duration := time.Since(start)
	domain := extractDomain(endpoint.URL)

	stats[domain].Total++
	if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 && duration <= RequestTimeout {
		stats[domain].Success++
	} else if err != nil {
		//Extra logging incase of error or other response
		log.Printf("Request to %s failed: %v with duration: %v)\n", endpoint.URL, err, duration)
	} else {
		log.Printf("Request to %s returned %d in %v (unavailable)\n", endpoint.URL, resp.StatusCode, duration)
	}
}

// Using url.Parse for better parsing
func extractDomain(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	domain := parsed.Host
	//ignoring port no
	if host, _, err := net.SplitHostPort(domain); err == nil {
		return host
	}
	return domain
}

func monitorEndpoints(endpoints []Endpoint) {
	for _, endpoint := range endpoints {
		domain := extractDomain(endpoint.URL)
		if stats[domain] == nil {
			stats[domain] = &DomainStats{}
		}
	}

	for {
		for _, endpoint := range endpoints {
			checkHealth(endpoint)
		}
		logResults()
		time.Sleep(CheckInterval)
	}
}

func logResults() {
	for domain, stat := range stats {
		percentage := int(math.Round(100 * float64(stat.Success) / float64(stat.Total)))
		fmt.Printf("%s has %d%% availability\n", domain, percentage)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <config_file>")
	}

	filePath := os.Args[1]
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var endpoints []Endpoint
	if err := yaml.Unmarshal(data, &endpoints); err != nil {
		log.Fatal("Error parsing YAML:", err)
	}

	monitorEndpoints(endpoints)
}
