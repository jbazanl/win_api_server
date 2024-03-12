package controllers

import (
	"database/sql"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"api-server/initializers"
	"api-server/models"

	"github.com/go-playground/form"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"net/url"

	"github.com/dgrijalva/jwt-go"
)

type PhoneData struct {
	Phone string `form:"phone"`
}

func GenerateJWT(userEmail string) (string, error) {
	var mySigningKey = []byte("W1n#Ov500@786") // Usa una clave secreta segura

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = userEmail
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expira en 24 horas

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ProcessClientHandler(c *fiber.Ctx) error {
	var payload models.CreateClienteSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	//Conexión a la BD MySQL
	connStr := "ovswitch:ovswitch123@tcp(127.0.0.1:3306)/switch?charset=utf8mb4"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("error al conectarse a la bd MySQL")
	}
	defer db.Close()

	// Procesar la información del archivo
	error := createRecord(db, payload, 0)
	if error.Tipo == "Error" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Registro procesado con error",
			"error":   error,
		})
	} else {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"línea creada exitosamente": payload.MSISDN}})
	}
}

func DeleteClientHandler(c *fiber.Ctx) error {
	bodyBytes := c.Body()
	values, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Error al parsear request"})
	}

	var phoneData PhoneData
	decoder := form.NewDecoder()
	_ = decoder.Decode(&phoneData, values)

	phone := phoneData.Phone

	if phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Se requiere parametro phone"})
	}

	//Conexión a la BD MySQL
	connStr := "ovswitch:ovswitch123@tcp(127.0.0.1:3306)/switch?charset=utf8mb4"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("error al conectarse a la bd MySQL")
	}
	defer db.Close()

	// Procesar DELETE
	error := deleteRecord(db, phone, 0)
	if error.Tipo == "Error" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Solicitud procesada con error",
			"error":   error,
		})
	} else {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"línea eliminada exitosamente": phone}})
	}
}

func FindClientes(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var clientes []models.Cliente
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&clientes)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(clientes), "clientes": clientes})
}

func FindCustomers(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var customers []models.Customer
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&customers)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(customers), "customers": customers})
}

func FindAccounts(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var accounts []models.Account
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&accounts)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(accounts), "accounts": accounts})
}

func FindUsers(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var users []models.User
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&users)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(users), "customers": users})
}

func UpdateCliente(c *fiber.Ctx) error {
	clienteId := c.Params("clienteId")

	var payload *models.UpdateClienteSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var cliente models.Cliente
	result := initializers.DB.First(&cliente, "id = ?", clienteId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No note with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Plan != "" {
		updates["plan"] = payload.Plan
	}
	if payload.Username != "" {
		updates["username"] = payload.Username
	}
	if payload.Password != "" {
		updates["password"] = payload.Password
	}
	if payload.Dominio != "" {
		updates["dominio"] = payload.Dominio
	}
	if payload.Departamento != "" {
		updates["departamento"] = payload.Departamento
	}
	if payload.Provincia != "" {
		updates["provincia"] = payload.Provincia
	}
	if payload.Distrito != "" {
		updates["distrito"] = payload.Distrito
	}

	if payload.IsRegistered != nil {
		updates["isRegistered"] = payload.IsRegistered
	}

	if payload.IsBlocked != nil {
		updates["isBlocked"] = payload.IsBlocked
	}

	updates["updated_at"] = time.Now()

	initializers.DB.Model(&cliente).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"cliente": cliente}})
}

func FindClienteById(c *fiber.Ctx) error {
	clienteId := c.Params("clienteId")

	var cliente models.Cliente
	result := initializers.DB.First(&cliente, "id = ?", clienteId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No cliente with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"cliente": cliente}})
}

func DeleteCliente(c *fiber.Ctx) error {
	clienteId := c.Params("clienteId")

	result := initializers.DB.Delete(&models.Cliente{}, "id = ?", clienteId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No cliente with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
