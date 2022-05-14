package main

import(
	"github.com/Aleqxan/E-commerce-project/controllers"
	"github.com/Aleqxan/E-commerce-project/database"
	"github.com/Aleqxan/E-commerce-project/middleware"
	"github.com/Aleqxan/E-commerce-project/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	port := os.GetEnv("PORT")
	if port == ""(
		port = "8000"
	)

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.Userdata(database.Client), "Users")

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.GET("/addtocart", app.AddToCart())
	routes.GET("/removeitem", app.RemoveItem())
	routes.GET("/cartcheckout", app.BuyFromCart())
	routes.GET("/instantbuy", app.InstantNuy())

	log.fatal(router.Run(":" + port))
}