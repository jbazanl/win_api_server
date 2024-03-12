package main

import (
	"api-server/models"
	"fmt"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MiError struct {
	Codigo  int
	Mensaje string
}

func (e *MiError) Error() string {
	return fmt.Sprintf("Código: %d, Mensaje: %s", e.Codigo, e.Mensaje)
}

type Resultado struct {
	Estado  string `json:"estado"`
	Numero  string `json:"numero"`
	Plan    string `json:"plan"`
	Mensaje string `json:"mensaje"`
}

func ResetBalanceClientes(db *gorm.DB) error {
	var api_results []Resultado
	var clientes []models.Cliente
	if err := db.Find(&clientes).Error; err != nil {
		return &MiError{Codigo: 500, Mensaje: "No se puede consultar tabla de clientes"}
	}

	for _, cliente := range clientes {
		plan := cliente.Plan
		account_id := cliente.AccountId
		phone := cliente.MSISDN
		saldo := 100.00 // Valor por defecto, ajusta según tu lógica
		if plan == "WINPHONE 110" {
			saldo = 110.00
		} else if plan == "WINPHONE 100" {
			saldo = 100.00
		} else if plan == "WINPHONE 10" {
			saldo = 10.00
		}

		resultado := resetSaldo(db, account_id, phone, saldo)
		if resultado.Estado != "" {
			api_results = append(api_results, resultado)
		}
	}

	if len(api_results) == 0 {
		return &MiError{Codigo: 500, Mensaje: "No se actualizo saldo a ningun cliente"}
	}

	// Puedes verificar cuántos registros se actualizaron, si es necesario
	fmt.Printf("Saldo actualizado para %v registros.\n", len(api_results))

	// Convertir api_results a JSON para incluir en el mensaje de error
	mensaje := "Se actualizo saldo de " + strconv.Itoa(len(api_results)) + " clientes"
	return &MiError{Codigo: 200, Mensaje: mensaje}
}

func resetSaldo(db *gorm.DB, account_id string, phone string, saldo float64) Resultado {
	now := time.Now()
	fmt.Println(now, phone, account_id)

	// Asumiendo que aquí se realiza alguna operación de actualización...
	// Encuentra y actualiza el saldo del cliente en la tabla customer_balance
	result := db.Model(&models.CustomerBalance{}).
		Where("account_id = ?", account_id).
		Update("balance", float64(-saldo))

	if result.Error != nil {
		fmt.Println("Error al actualizar el saldo:", result.Error)
		return Resultado{Estado: "ERROR", Numero: phone, Mensaje: "Error al actualizar saldo en BD MySQL"}
	}

	return Resultado{Estado: "OK", Numero: phone, Mensaje: "Saldo actualizado correctamente en BD MySQL"}
}

func main() {
	dsn := "ovswitch:ovswitch123@tcp(127.0.0.1:3306)/switch?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(&MiError{Codigo: 500, Mensaje: "Error al conectarse a la base de datos"})
		return
	}
	err = ResetBalanceClientes(db)
	if err != nil {
		fmt.Println(err)
	}
}
