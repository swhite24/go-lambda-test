package main

import (
	"github.com/apex/gateway"
	"github.com/gin-gonic/gin"
	"github.com/swhite24/go-lambda-test/pkg/controllers"
)

func main() {
	g := gin.New()

	controllers.ServeUsers(g, "golang-test")

	gateway.ListenAndServe(":3000", g)
}
