package main

import (
	"PaginaSEG/api"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	api.InitRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	if err := r.Run(port); err != nil {
		panic(fmt.Errorf("Error al intentar iniciar el servidor: %v", err))
	}

}

