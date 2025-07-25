package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pt-sinan-akbar/controllers"
	_ "github.com/pt-sinan-akbar/docs"
	"github.com/pt-sinan-akbar/helpers"
	"github.com/pt-sinan-akbar/initializers"
	"github.com/pt-sinan-akbar/manager"
	"github.com/pt-sinan-akbar/routes"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

//	@title			Split Bill API
//	@version		1.0
//	@description	Split Bill Swagger Documentation
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

var (
	server *gin.Engine

	//Controllers
	BillController       controllers.BillController
	BillOwnerController  controllers.BillOwnerController
	BillMemberController controllers.BillMemberController
	BillDataController   controllers.BillDataController
	BillItemController   controllers.BillItemController

	//Routes
	BillRouterController       routes.BillRouterController
	BillOwnerRouterController  routes.BillOwnerRouterController
	BillMemberRouterController routes.BillMemberRouterController
	BillDataRouterController   routes.BillDataRouterController
	BillItemRouterController   routes.BillItemRouterController

	// Manager
	BillManager           manager.BillManager
	BillItemManager       manager.BillItemManager
	BillDataManager       manager.BillDataManager
	BillMemberManager     manager.BillMemberManager
	BillOwnerManager      manager.BillOwnerManager
	BillMemberItemManager manager.BillMemberItemManager
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables")
	}

	initializers.ConnectDB(&config)
	// Managers
	BillItemManager = manager.NewBillItemManager(initializers.DB)
	BillDataManager = manager.NewBillDataManager(initializers.DB)
	BillMemberItemManager = manager.NewBillMemberItemManager(initializers.DB)
	BillMemberManager = manager.NewBillMemberManager(initializers.DB)
	BillOwnerManager = manager.NewBillOwnerManager(initializers.DB)
	BillManager = manager.NewBillManager(initializers.DB, &BillItemManager, &BillDataManager, &BillMemberItemManager, &BillMemberManager, &BillOwnerManager)

	//Controllers
	BillController = controllers.NewBillController(&BillManager)
	BillOwnerController = controllers.NewBillOwnerController(&BillOwnerManager)
	BillMemberController = controllers.NewBillMemberController(&BillMemberManager)
	BillDataController = controllers.NewBillDataController(&BillDataManager)
	BillItemController = controllers.NewBillItemController(&BillItemManager)

	//Routes
	BillRouterController = routes.NewBillRouterController(BillController)
	BillOwnerRouterController = routes.NewBillOwnerRouterController(BillOwnerController)
	BillMemberRouterController = routes.NewBillMemberRouterController(BillMemberController)
	BillDataRouterController = routes.NewBillDataRouterController(BillDataController)
	BillItemRouterController = routes.NewBillItemRouterController(BillItemController)

	gin.SetMode(config.GinMode)
	server = gin.Default()
}

func main() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}

	server.Use(cors.New(corsConfig))
	router := server.Group("/api/v1")

	BillRouterController.BillRouter(router)
	BillOwnerRouterController.BillOwnerRouter(router)
	BillMemberRouterController.BillMemberRouter(router)
	BillDataRouterController.BillDataRouter(router)
	BillItemRouterController.BillItemRouter(router)

	//Ini buat apaan
	server.GET("/generate-image", func(c *gin.Context) {
		bytes, err := helpers.GenerateInitialsImage("Sinan")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate image",
				"message": err.Error()})
			return
		}

		c.Header("Content-Type", "image/png")
		_, err = c.Writer.Write(bytes)
		if err != nil {
			return
		}
	})
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	host := ""
	if initializers.ConfigSetting.GinMode == gin.ReleaseMode {
		host = "0.0.0.0"
	} else {
		host = "127.0.0.1"
	}
	log.Fatal(server.Run(host + ":" + initializers.ConfigSetting.ServerPort))
}
