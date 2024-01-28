package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello World!")
}
func uploadImage(c *fiber.Ctx) error {
	// Read file from request
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// save the file to the server
	err = c.SaveFile(file, "./uploads/"+file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("File uploaded successfully : ")
}

func renderTemplate(c *fiber.Ctx) error {
	// Render the template with variable data
	return c.Render("template", fiber.Map{
		"Name": "World",
	})
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getConfig(c *fiber.Ctx) error {
	// example : return a configuration value from environment variabled
	secretKey := getEnv("SECRET_KEY", "defaultSecret")
	return c.JSON(fiber.Map{
		"secret_key": secretKey,
	})
}

func isAdmin(c *fiber.Ctx) error {
	user := c.Locals(userContextKey).(*UserData)
	if user.Role != "admin" {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}

func main() {
	app := fiber.New()
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	// apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // adjust this to be more restrictive if needed
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATH",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))

	// middleware
	// use the logging middleware
	app.Use(loggingMiddleware)
	// setup routes
	// app.Get("/middleware", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello , world")
	// })
	// jwt secret key
	// secretKey := "secret"
	// // login route
	// app.Post("/login", login(secretKey))

	// Setup router
	app.Get("/books", GetBooks)
	app.Get("/book/:id", GetBook)
	app.Post("/", CreateBook)
	app.Put("/book/:id", UpdateBook)
	app.Delete("/book/:id", DeleteBook)
	app.Post("/upload", uploadImage)

	// JWT Middleware
	//  app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: []byte(secretKey),
	// }))

	// middleware to extract user data from jwt
	// app.Use(extractUserFromJWT)

	// group routes under /book
	// bookGroup := app.Group("/book")

	// apply the isAdmin middleware only to the /book
	// bookGroup.Use(isAdmin)
	// bookGroup.Get("/", GetBooks)

	// initialize standard Go html template engine
	// engine := html.New("example/views", ".html")

	// Pass the engine to Fiber
	// app = fiber.New(fiber.Config{
	// 	Views: engine,
	// })

	// use the environment varible for the port , default to 8080 if not set
	port := getEnv("Port : ", "8080")

	// set up route
	// app.Get("/", renderTemplate)
	// app.Get("/api/config", getConfig)

	// load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.Listen(": " + port)
}
