package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// loggingMiddleware logs the processing time for each request
func loggingMiddleware(c *fiber.Ctx) error {
	// start timer
	start := time.Now()

	// process request
	err := c.Next()

	// calculate processing time
	duration := time.Since(start)

	// log the information
	fmt.Printf("Request URL: %s - Method: %s - Duration: %s\n", c.OriginalURL(), c.Method(), duration)

	return err
}

type User struct {
	Email    string
	Password string
}

// Dummy user for example
var memberUser = User{
	Email:    "user@example.com",
	Password: "password123",
}

func login(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		// check credentials - in real world , you should check against a database
		if req.Email != memberUser.Email || req.Password != memberUser.Password {
			return fiber.ErrUnauthorized
		}
		// create token
		token := jwt.New(jwt.SigningMethodHS256)

		// set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["email"] = memberUser.Email
		claims["role"] = "admin"
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// generate encoded token
		t, err := token.SignedString([]byte(secretKey))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(fiber.Map{"token": t})
	}
}

// UserData represent the user data extracted from the jwt token
type UserData struct {
	Email string
	Role  string
}

// userContextKey is the key used to store user data in the
const userContextKey = "user"

// extracUserFromJwt is a middleware that extracts user data from the JWT token
func extractUserFromJWT(c *fiber.Ctx) error {
	user := &UserData{}
	// Extract the token from the Fiber context (inserted by the JWT middleware)
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	fmt.Println(claims)
	user.Email = claims["email"].(string)
	user.Role = claims["role"].(string)

	// store the user data in the Fiber context
	c.Locals(userContextKey, user)
	return c.Next()
}
