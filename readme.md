

| Command                                | Purpose                                         |
|----------------------------------------|-------------------------------------------------|
| `go run ./cmd/migrate/main.go check`   | validate markdown frontmatter                   |
| `go run ./cmd/build`                   | Generate static site into `public/`             |
| `go run ./cmd/serve`                   | Serve that folder locally                       |
| `WATCH=1 go run ./cmd/serve`           | Serve **and** auto-rebuild when content changes |
| `go run ./cmd/analytics-consumer`      | load events and send to plausible               |



## Example usage of build

For local test inform SITE_BASE_URL as localhost:8080  
that is the same address from ./cmd/serve 
```bash
export SITE_BASE_URL=http://localhost:8080/

go run ./cmd/build
```



## Example usage of consumer
load events and send to plausible  
aws cli should be already configured 

```bash
export AWS_REGION=<queueRegion>
export SQS_QUEUE_URL=<queueUrl>
export PLAUSIBLE_URL="http://localhost:8000/api/event"

go run ./cmd/analytics-consumer
```


## trying on docker

```bash
docker build . -t wastingnotime/blog:local
 
docker run -p 8080:80 wastingnotime/blog:local

```