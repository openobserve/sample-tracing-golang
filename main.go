package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zinclabs/otel-example/pkg/tel"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
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

	router.Use(otelgin.Middleware("otel1-gin-gonic"))

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

	span.AddEvent("GetUserDetails called")

	return "Hello User Details from Go microservice"
}
