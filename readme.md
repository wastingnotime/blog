

| Command                      | Purpose                                         |
| ---------------------------- | ----------------------------------------------- |
| `go run ./cmd/build`         | Generate static site into `public/`             |
| `go run ./cmd/serve`         | Serve that folder locally                       |
| `WATCH=1 go run ./cmd/serve` | Serve **and** auto-rebuild when content changes |


```bash
export SITE_BASE_URL=http://localhost:8080/
go run ./cmd/build
```