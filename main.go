package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zinclabs/otel-example/pkg/tel"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/zinclabs/otel-example")

func main() {
	tp := tel.InitTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Println("Error shutting down tracer provider: ", err)
		}
	}()

	router := gin.Default()

	router.GET("/", GetUser)

	router.Run(":8080")

}

func GetUser(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "GetUser")
	defer span.End()

	details := GetUserDetails(c.Request.Context())
	c.String(http.StatusOK, details)
}

func GetUserDetails(ctx context.Context) string {
	_, span := tracer.Start(ctx, "GetUserDetails")
	defer span.End()
	return "Hello User Details"
}
