package monitoring

import (
	"fmt"

	"github.com/openzipkin/zipkin-go"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

const endpointURL = "http://localhost:9411/api/v2/spans"
const serviceHostPort = "localhost:60061"
const serviceName = "auth"

//NewZipkinTracer creates a new tracer for zipkin
func NewZipkinTracer() (*zipkin.Tracer, error) {
	// The reporter sends traces to zipkin server
	reporter := reporterhttp.NewReporter(endpointURL)

	// Local endpoint represent the local service information
	localEndpoint, err := zipkin.NewEndpoint(serviceName, serviceHostPort)

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	t, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("tracer: ", t)

	return t, err
}
