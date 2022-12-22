package tel

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracerHTTP() *sdktrace.TracerProvider {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	OTEL_OTLP_HTTP_ENDPOINT := os.Getenv("OTEL_OTLP_HTTP_ENDPOINT")

	if OTEL_OTLP_HTTP_ENDPOINT == "" {
		OTEL_OTLP_HTTP_ENDPOINT = "<host>:<port>" //without trailing slash
	}

	// OTEL_OTLP_HTTP_ENDPOINT = "ziox2.dev.zincsearch.com"
	OTEL_OTLP_HTTP_ENDPOINT = "otel.dev2.zinclabs.dev"
	// OTEL_OTLP_HTTP_ENDPOINT = "httptester.dev2.zinclabs.dev"

	otlptracehttp.NewClient()

	otlpHTTPExporter, err := otlptracehttp.New(context.TODO(),
		// otlptracehttp.WithInsecure(), // use http & not https
		otlptracehttp.WithEndpoint(OTEL_OTLP_HTTP_ENDPOINT),
		// otlptracehttp.WithURLPath("/api/prabhat-org2/traces"),
		otlptracehttp.WithURLPath("/v1/traces"),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": "Basic YWRtaW46Q29tcGxleHBhc3MjMTIz",
		}),
	)

	// stdExporter, _ := stdouttrace.New(
	// 	stdouttrace.WithWriter(io.Writer(os.Stdout)),
	// 	stdouttrace.WithPrettyPrint(),
	// )

	if err != nil {
		fmt.Println("Error creating HTTP OTLP exporter: ", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		// the service name used to display traces in backends
		semconv.ServiceNameKey.String("otel1-gin-gonic"),
		semconv.ServiceVersionKey.String("0.0.1"),
		attribute.String("environment", "test"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(otlpHTTPExporter),
		// sdktrace.WithBatcher(stdExporter),
	)
	otel.SetTracerProvider(tp)

	return tp
}
