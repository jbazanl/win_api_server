package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Abrir conexión con la base de datos switch
	switchDB, err := sql.Open("mysql", "ovswitch:ovswitch123@tcp(localhost:3306)/switch")
	if err != nil {
		log.Fatal("Error al conectar a la base de datos switch:", err)
	}
	defer switchDB.Close()

	// Abrir conexión con la base de datos kamailio
	kamailioDB, err := sql.Open("mysql", "ovswitch:ovswitch123@tcp(localhost:3306)/kamailio")
	if err != nil {
		log.Fatal("Error al conectar a la base de datos kamailio:", err)
	}
	defer kamailioDB.Close()

	// Consultar registros de la tabla customer_sip_account en la BD switch
	rows, err := switchDB.Query("SELECT username FROM customer_sip_account WHERE status = '1'")
	if err != nil {
		log.Fatal("Error al ejecutar la consulta en la tabla customer_sip_account:", err)
	}
	defer rows.Close()

	// Iterar sobre los registros
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			log.Fatal("Error al escanear fila:", err)
		}

		// Consultar en la tabla location de la BD kamailio usando el campo username
		var locationExists bool
		err := kamailioDB.QueryRow("SELECT COUNT(*) FROM location WHERE username = ?", username).Scan(&locationExists)
		if err != nil {
			log.Printf("Error al consultar en la tabla location para el usuario %s: %v", username, err)
			continue
		}

		// Actualizar el campo status en la tabla customer_sip_account si no se encuentra un registro en la tabla location
		if !locationExists {
			_, err := switchDB.Exec("UPDATE customer_sip_account SET status = '0' WHERE username = ?", username)
			if err != nil {
				log.Printf("Error al actualizar el campo status para el usuario %s: %v", username, err)
				continue
			}
			log.Printf("El usuario %s no tiene un registro en la tabla location. Se actualizó el campo status a 0.", username)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error al iterar sobre filas:", err)
	}
}

