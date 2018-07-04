package main

import (
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/gin-gonic/gin"
	"github.com/swhite24/go-lambda-test/pkg/controllers"
)

func main() {
	addr := ":3000"
	mode := os.Getenv("GIN_MODE")
	g := gin.New()

	controllers.ServeUsers(g, "golang-test")

	if mode == "release" {
		log.Fatal(gateway.ListenAndServe(addr, g))
	} else {
		log.Fatal(http.ListenAndServe(addr, g))
	}
}
