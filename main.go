package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ittechman101/go-pos/models"
	"github.com/ittechman101/go-pos/routes"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	models.ConnectDB()
	//	defer database.DB.Close()

	routes.Register(app, models.DB)

	log.Fatal(app.Listen(":3030"))
}

//docker build . --target prod -t jetdev7/posapp-be-test
//docker run -p 3030:3030 -e MYSQL_HOST=10.10.13.183 -e MYSQL_USER=root -e MYSQL_PASSWORD=rkdtjdeornr107 -e MYSQL_DBNAME=go-pos-db jetdev7/posapp-be-test
