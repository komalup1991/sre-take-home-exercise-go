# Fetch SRE Take-Home Exercise in GO
The code is of a CLI tool that periodically checks the health of HTTP endpoints defined in a YAML config file and logs the availability of each domain.
## Features
- Reads a YAML config describing endpoints to check
- Every 15 seconds:
  - Sends requests to all endpoints
  - Records availability by domain
- A domain is considered **available** if:
  - The endpoint returns a status code in the 200–299 range
  - The response is received in ≤ 500ms
- Ignores port numbers when grouping domains
- Includes request timing & status logs for better observability

## How to install
**Clone the repo**

```bash
git clone https://github.com/komalup1991/sre-take-home-exercise-go.git
cd sre-take-home-exercise-go
```
## How to run

**Prepare the config file**

Example (`config.yaml`):

```yaml
- name: sample up
  url: https://example.com/
- name: post example
  url: https://example.com/post
  method: POST
  headers:
    content-type: application/json
  body: '{"foo": "bar"}'
```
**Initialize Go module for the project**

```go mod init sre-take-home```

**To clean up your go.mod and go.sum files**

```go mod tidy ```

**Run the tool**

```go run main.go <config_file>```

Example: go run main.go config.yaml

## YAML Format

- `name` (required): Label for the endpoint
- `url` (required): Full HTTP/HTTPS URL
- `method` (optional): HTTP method, defaults to `GET`
- `headers` (optional): Key-value pairs
- `body` (optional): JSON string

## Fixes

### Bug fixes as per missing from requirements

| Requirement | Issues Identified | Fix |
|-------|-----|-------------------------|
| Endpoint considered available only if Status code is between 200 and 299 and response <= 500ms | Only status code was being checked. |Updated condition to also check latency.|
|Must ignore port numbers when determining domain | Port number was not stripped.|Rewrote domain extraction using net/url to ignore ports.|
|If method field is omitted, the default is GET. |No logic for empty method.| Added default method = "GET" for empty method.


### Best practices as in good to have

| Issue | Reason for Change| Fix |
|-------|-----|-------------------------|
| Used string splitting for domain parsing | Fragile and error-prone with edge cases. |Replaced with url.Parse() for robust domain extraction.|
| Hardcoded Timeouts|Hard to maintain or tweak. |Introduced constants like timeout = 500 * time.Millisecond.|
| Deprecated "io/ioutil"|Deprecated: As of Go 1.16| The same functionality is now provided by package [io] or package [os] using os as it is already used.|