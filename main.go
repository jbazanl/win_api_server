package main

import (
	"log"
	"strings"
	"time"

	"api-server/controllers"
	"api-server/initializers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	config, err := initializers.LoadConfig("/app")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(&config)
}

type JwtToken struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Token     string
	ClientID  string // Identificador del cliente
	ExpiresAt time.Time
}

func GetToken(db *gorm.DB, tokenString string) (*JwtToken, error) {
	var token JwtToken
	result := db.Unscoped().Where("token = ?", tokenString).First(&token)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// El token no fue encontrado
			return nil, result.Error
		}
		// Otro error ocurrió al realizar la consulta
		log.Printf("Error al buscar el token: %v\n", result.Error)
		return nil, result.Error
	}

	return &token, nil
}

func JWTMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token faltante o inválido"})
		}

		// Utiliza GetToken para validar el token
		token, err := GetToken(db, tokenString)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al verificar el token"})
		}

		// Puedes verificar adicionalmente si el token ha expirado
		if token.ExpiresAt.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expirado"})
		}

		// Guardar la información del token en el contexto si es necesario
		clientID := token.ClientID
		c.Locals("clientID", clientID)

		return c.Next()
	}
}

func main() {
	app := fiber.New()
	micro := fiber.New()

	app.Mount("/api", micro)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	dsn := "ovswitch:ovswitch123@tcp(127.0.0.1:3306)/switch?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	micro.Route("/clientes", func(router fiber.Router) {
		router.Use(JWTMiddleware(db))
		router.Post("/", controllers.ProcessClientHandler)
		router.Delete("/", controllers.DeleteClientHandler)
		//router.Get("", controllers.FindClientes)
	})

	//micro.Route("/customers", func(router fiber.Router) {
	//	router.Get("", controllers.FindCustomers)
	//})
	//micro.Route("/accounts", func(router fiber.Router) {
	//router.Get("", controllers.FindAccounts)
	//})
	//micro.Route("/users", func(router fiber.Router) {
	//	router.Get("", controllers.FindUsers)
	//})
	micro.Route("/archivo", func(router fiber.Router) {
		router.Use(JWTMiddleware(db))
		router.Post("/", controllers.AddProcessFileHandler)
		router.Delete("/", controllers.DeleteProcessFileHandler)
	})

	//	micro.Route("/clientes/:clienteId", func(router fiber.Router) {
	//		router.Delete("", controllers.DeleteCliente)
	//		router.Get("", controllers.FindClienteById)
	//		router.Patch("", controllers.UpdateCliente)
	//	})

	micro.Get("/api/provision", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to Golang, Fiber, MySQL, and GORM",
		})
	})

	log.Fatal(app.Listen(":5000"))
}
