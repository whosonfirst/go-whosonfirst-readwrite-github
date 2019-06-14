fmt:
	go fmt cmd/*.go
	go fmt reader/*.go
	go fmt writer/*.go

tools:
	go build -mod vendor -o bin/api-reader cmd/api-reader/main.go
