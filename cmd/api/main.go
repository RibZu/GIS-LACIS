package main

import (
	"PaginaSEG/api"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()
	api.InitRoutes(r)

	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080"
	}

	if err := r.Run(":" + puerto); err != nil {
		panic(fmt.Errorf("Error al intentar iniciar el servidor: %v", err))
	}

}
