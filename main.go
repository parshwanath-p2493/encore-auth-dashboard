package main

import (
	"log"
	"net/http"
	"os"

	"encore.app/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	//database.OpenCollection("Users")
	r := fiber.New()
	err := godotenv.Load(".env") //mandatory to load the .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.Connection()
	PORT := os.Getenv("PORT")
	log.Println(PORT)
	if PORT == "" {
		PORT = "3000"
	}
	r.Get("/hi", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(&fiber.Map{"message": "Everything is ok SERVER is RUNNING on port " + PORT})
	})
	log.Println(PORT)
	//c := http.Handler
	//log.Fatal(http.ListenAndServe(":" + PORT,c))
	log.Fatal(r.Listen(":" + PORT))
}
