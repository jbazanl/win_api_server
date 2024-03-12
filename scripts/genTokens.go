package main

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type JwtToken struct {
	ID        uint `gorm:"primaryKey"`
	Token     string
	ClientID  string // Identificador del cliente
	ExpiresAt time.Time
}

func GenerateAndStoreTokens(db *gorm.DB, numberOfTokens int) error {
	var mySigningKey = []byte("WinPhone")

	for i := 0; i < numberOfTokens; i++ {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)

		claims["client_id"] = "Cliente_" + strconv.Itoa(i)
		claims["exp"] = time.Now().Add(time.Hour * 24 * 90).Unix() // Expira en 90 dias

		tokenString, err := token.SignedString(mySigningKey)
		if err != nil {
			return err
		}

		jwtToken := JwtToken{
			Token:     tokenString,
			ClientID:  "Cliente_" + strconv.Itoa(i),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 90),
		}

		db.Create(&jwtToken)
	}

	return nil
}

func main() {
	dsn := "ovswitch:ovswitch123@tcp(127.0.0.1:3306)/switch?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error al conectarse a la base de datos: " + err.Error())
	}
	db.AutoMigrate(&JwtToken{})
	GenerateAndStoreTokens(db, 5)
}
