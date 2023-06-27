build:
	GOOS=linux GOARCH=amd64 go build -o binlambda-logs-loki main.go

test: build
	go test