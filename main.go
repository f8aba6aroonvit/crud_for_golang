package main

import (
	"train_golang/configs"
	"train_golang/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	configs.ConnectDB() // เชื่อมต่อ DB
	app := fiber.New()
	routes.UserRoute(app) // กำหนด Route สำหรับ User
	app.Listen(":6000")
}
