package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/zinclabs/otel-example")

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
