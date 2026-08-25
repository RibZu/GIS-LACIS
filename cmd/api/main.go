package main

import (
	"PaginaSEG/api"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	api.InitRoutes(r)

	if err := r.Run(":8080"); err != nil {
		panic(fmt.Errorf("Error al intentar iniciar el servidor: %v", err))
	}

}
