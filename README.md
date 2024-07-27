# Thread Pool based TCP server

- TCP server that has a thread pool based thread per client implementation that comes with rate-limiting by default.

- We can specify the number of workers.

- `go run main.go --workers 4`