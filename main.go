package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/zinclabs/otel-example/pkg/tel"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("github.com/zinclabs/otel-example")

func main() {
	// tp := tel.InitTracerGRPC()
	tp := tel.InitTracerHTTP()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Println("Error shutting down tracer provider: ", err)
		}
	}()

	router := gin.Default()

	router.Use(otelgin.Middleware(""))

	router.GET("/", GetUser)

	router.Run(":8080")

}

func GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	// sleep for 1 second to simulate a slow request
	time.Sleep(1 * time.Second)

	childCtx, span := tracer.Start(ctx, "GetUser")
	defer span.End()

	details := GetUserDetails(childCtx)
	c.String(http.StatusOK, details)
}

func GetUserDetails(ctx context.Context) string {
	_, span := tracer.Start(ctx, "GetUserDetails")
	defer span.End()
	// sleep for 500 ms to simulate a slow request
	time.Sleep(500 * time.Millisecond)

	// log a message to stdout with the traceID and spanID
	log.Info().Str("traceID", span.SpanContext().TraceID().String()).Str("spanID", span.SpanContext().SpanID().String()).Msg("Log message for user details")

	span.AddEvent("GetUserDetails called")

	childCtx, spanHobby := tracer.Start(ctx, "GetHobbiesCall")
	defer spanHobby.End()

	hobbies, err := GetHobbies(childCtx)
	if err != nil {
		span.RecordError(err)
	} else {
		span.SetAttributes(attribute.String("hobbies", hobbies))
	}

	return "Hello User Details from Go microservice"
}

func GetHobbies(ctx context.Context) (string, error) {
	_, span := tracer.Start(ctx, "GetHobbies")

	defer func() {
		// We recover from any panics here and add the the stack trace to the span to allow for debugging
		r := recover()

		if r != nil {
			fmt.Println("recovered: ", r)

			// get error from panic
			e1 := r.(error)
			stackTrace := debug.Stack()
			runtime.Stack(stackTrace, true)

			// add attributes to the span event
			attributes := []attribute.KeyValue{
				attribute.String("user", "useremail@userdomain.com"),
				attribute.String("exception.escaped", "false"),
				attribute.String("exception.stacktrace", string(stackTrace)),
			}
			options := trace.WithAttributes(attributes...)

			// add error to span event
			span.RecordError(e1, options)

			// set event status to error
			span.SetStatus(codes.Error, e1.Error())
		}

		span.End()
	}()

	// sleep for 500 ms to simulate a slow request
	time.Sleep(500 * time.Millisecond)

	span.AddEvent("GetHobbies called")

	// We will cause a divide by zero error here
	a := 0
	b := 3
	c := b / a

	return strconv.Itoa(c), nil

}
