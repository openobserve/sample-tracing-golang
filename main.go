package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var tracer = otel.Tracer("github.com/zinclabs/otel-example")

// provider := sdktrace.NewTracerProvider(
// 	sdktrace.WithSampler(sdktrace.AlwaysSample()),
// )

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Println("Error shutting down tracer provider: ", err)
		}
	}()

	router := gin.Default()

	// This handler will match /user/john but will not match /user/ or /user
	router.GET("/", GetUser)

	router.Run(":8080")

}

func GetUser(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "GetUser")
	defer span.End()
	c.String(http.StatusOK, "Hello User")
}

func initTracer() *sdktrace.TracerProvider {
	OTEL_OTLP_HTTP_ENDPOINT := os.Getenv("OTEL_OTLP_HTTP_ENDPOINT")

	if OTEL_OTLP_HTTP_ENDPOINT == "" {
		OTEL_OTLP_HTTP_ENDPOINT = "otel.dev2.zinclabs.dev"
	}

	otlptracehttp.NewClient()

	otlpHTTPExporter, err := otlptracehttp.New(context.TODO(),
		// otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(OTEL_OTLP_HTTP_ENDPOINT),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": "Basic YWRtaW46Q29tcGxleHBhc3MjMTIz",
		}),
	)

	if err != nil {
		fmt.Println("Error creating HTTP OTLP exporter: ", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		// the service name used to display traces in backends
		semconv.ServiceNameKey.String("otel1-gin-gonic"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(otlpHTTPExporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
