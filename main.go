package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/initializers"
	"github.com/pt-sinan-akbar/controllers"
	_ "github.com/pt-sinan-akbar/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title Split Bill API
// @version 1.0
// @description Split Bill Swagger Documentation
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/


func init(){
	config, err := initializers.LoadConfig(".")	
	if err != nil {
        log.Fatal("? Could not load environment variables")
    }

    initializers.ConnectDB(&config)
}


func main() {
    r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		bills := v1.Group("/bills")
		{
			bills.GET("", controllers.GetAll)
			bills.GET(":id", controllers.GetByID)
		}
	}

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    r.Run(":8080")
}