package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type LambdaTelemetryApi struct {
}

func NewLambdaTelemetryApi() *LambdaTelemetryApi {
	return &LambdaTelemetryApi{}
}

// Subscribe to the telemetry API, Incoming messages are returned in the channel
func (api *LambdaTelemetryApi) Subscribe(ctx context.Context) {
	apiEndpoint := fmt.Sprintf("http://%s/2022-07-01/telemetry/", os.Getenv("AWS_LAMBDA_RUNTIME_API"))
	stdout.Debug("registering telemetry api callback", "telemetryendpoint", apiEndpoint)

	httpClient := http.Client{}
	dio := []byte{byte("s")}
	req, err := http.NewRequest(http.MethodPut, apiEndpoint)

}

func startHttpEchoServer() {

}
