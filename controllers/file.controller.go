package controllers

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"api-server/initializers"
	"api-server/models"

	"github.com/gofiber/fiber/v2"
)

func utf8_decode(str string) string {
	var result string
	for i := range str {
		result += string(str[i])
	}
	return result
}

type Resultado struct {
	Tipo    string `json:"tipo"`
	Numero  string `json:"numero"`
	Mensaje string `json:"mensaje"`
	Linea   int    `json:"linea"`
}

func generateSecureKey(length int) (string, error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b), nil
}

func generateUserKey(db *sql.DB, table, uniqueFieldName string, maxKey int) (string, error) {
	// Calcula el nuevo ID basado en maxKey.
	newKeyInt := maxKey + 1
	newKey := fmt.Sprintf("UC%08d", newKeyInt) // Asegura que el ID tenga una longitud fija con el formato deseado.

	// Prepara la consulta para verificar la unicidad del nuevo ID.
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = ? AND user_type='CUSTOMERADMIN')", table, uniqueFieldName)
	stmt, err := db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	// Itera hasta encontrar un ID único.
	for {
		var exists bool
		err := stmt.QueryRow(newKey).Scan(&exists)
		if err != nil {
			return "", err
		}

		if !exists {
			return newKey, nil // Retorna el nuevo ID si es único.
		}

		// Incrementa y genera un nuevo ID si el actual ya existe.
		newKeyInt++
		newKey = fmt.Sprintf("UC%08d", newKeyInt)
	}
}

func generateVoipKey(db *sql.DB, table, uniqueFieldName string, maxKey int) (string, error) {
	// Calcula el nuevo ID basado en maxKey.
	newKeyInt := maxKey + 1
	newKey := fmt.Sprintf("CVM%09d", newKeyInt) // Asegura que el ID tenga una longitud fija con el formato deseado.

	// Prepara la consulta para verificar la unicidad del nuevo ID.
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = ? AND account_type='CUSTOMER')", table, uniqueFieldName)
	stmt, err := db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	return newKey, nil // Retorna el nuevo ID si es único

	/*
		// Itera hasta encontrar un ID único.
		for {
			var existe bool
			err := stmt.QueryRow(newKey).Scan(&existe)
			if err != nil {
				return "", err
			}

			if !existe {
				return newKey, nil // Retorna el nuevo ID si es único.
			}

			// Incrementa y genera un nuevo ID si el actual ya existe.
			newKeyInt++
			newKey = fmt.Sprintf("CVM%09d", newKeyInt)
		}
	*/
}

func checkUserKeyExists(db *sql.DB, key string, accountType string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = ? AND account_type = ?`
	row := db.QueryRow(query, key, accountType)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func generateCustomerKey(db *sql.DB, table string, uniqueFieldName string, maxKey int) (string, error) {
	// Calcula el nuevo ID basado en maxKey.
	newKeyInt := maxKey + 1
	newKey := fmt.Sprintf("STC%06d", newKeyInt) // Asegura que el ID tenga una longitud fija con el formato deseado.

	// Prepara la consulta para verificar la unicidad del nuevo ID.
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = ?)", table, uniqueFieldName)
	stmt, err := db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	return newKey, nil // Retorna el nuevo ID si es único
	/*
		// Itera hasta encontrar un ID único.
		for {
			var exists bool
			err := stmt.QueryRow(newKey).Scan(&exists)
			if err != nil {
				return "", err
			}

			if !exists {
				return newKey, nil // Retorna el nuevo ID si es único.
			}

			// Incrementa y genera un nuevo ID si el actual ya existe.
			newKeyInt++
			newKey = fmt.Sprintf("STC%06d", newKeyInt)
		}
	*/
}

func AddProcessFileHandler(c *fiber.Ctx) error {

	var payload models.CreateClienteSchema
	var api_results []Resultado

	// Obtener el archivo del cuerpo de la solicitud
	file, err := c.FormFile("archivo")
	fmt.Println(file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("error al obtener el archivo CSV")
	}

	// Abrir el archivo
	fileOpened, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("error al abrir el archivo CSV")
	}
	defer fileOpened.Close()

	// Crea un lector CSV
	reader := csv.NewReader(fileOpened)
	reader.Comma = ';'

	// Lee la primera línea (encabezados) y descártala
	firstLine, err := reader.Read()
	log.Println(firstLine)
	if err != nil {
		fmt.Println("Error al leer la primera línea:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error al leer la primera línea")
	}

	//Conexión a la BD MySQL
	connStr := "ovswitch:ovswitch123@tcp(127.0.0.1:3306)/switch?charset=utf8mb4"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("error al conectarse a la bd MySQL")
	}
	defer db.Close()

	// Procesar todas las líneas del archivo
	processedLine := 2
	for {
		row, err := reader.Read()
		log.Println(row)
		if err == io.EOF {
			break // fin del archivo
		} else if err != nil {
			fmt.Println("Error al leer el archivo CSV:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("error al leer el archivo CSV")
		}

		// Asignar valores a la estructura
		payload.Nombre = row[0]
		payload.MSISDN = row[1]
		payload.Username = row[2]
		payload.Password = row[3]
		payload.Dominio = row[4]
		payload.Plan = row[5]
		payload.Departamento = row[6]
		payload.Provincia = row[7]
		payload.Distrito = row[8]

		// Validar la información del archivo
		payload_errors := models.ValidateStruct(payload)
		if payload_errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(payload_errors)
		}

		// Procesar la información del archivo
		resultado := createRecord(db, payload, processedLine)
		if resultado.Tipo != "" {
			api_results = append(api_results, resultado)
		}
		processedLine += 1
	}

	if len(api_results) == 0 {
		//return c.SendString("archivo recibido y procesado correctamente")
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"no se actualizo ningun registro": file}})
	} else {
		// Retornar la información del archivo temporal
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Archivo procesado por API",
			"errores": api_results,
		})
	}
}

func createRecord(db *sql.DB, payload models.CreateClienteSchema, processedLine int) Resultado {
	now := time.Now()

	update_by := "API"
	msisdn := payload.MSISDN
	plan := payload.Plan
	address := payload.Departamento + " - " + payload.Provincia + " - " + payload.Distrito
	//accountType := "CUSTOMERADMIN" // Puedes cambiar esto según sea necesario
	name := payload.Nombre

	// Iniciar una transacción
	tx, err := db.Begin()
	if err != nil {
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al iniciar transacción en BD", Linea: processedLine}
	}

	// Obtener el máximo valor actual del campo único
	var maxCustomerKey int
	query := fmt.Sprintf("SELECT MAX(%s) FROM %s", "customer_id", "customers")
	err = db.QueryRow(query).Scan(&maxCustomerKey)
	if err != nil && err != sql.ErrNoRows {
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al obtener MAX customer key", Linea: processedLine}
	}

	customerKey, err := generateCustomerKey(db, "customers", "customer_id", maxCustomerKey)
	if err != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al generar new customer key", Linea: processedLine}
	}
	fmt.Println("customerKey: ", customerKey)

	claveWeb, _ := generateSecureKey(8) // Longitud de la clave establecida en 8

	var maxUserKey int
	query = fmt.Sprintf("SELECT MAX(CAST(SUBSTRING(%s, 3) AS INT)) FROM %s WHERE user_type='CUSTOMERADMIN'", "user_id", "users")
	//err = db.QueryRow(query).Scan(&maxUserKey)
	err = db.QueryRow(query).Scan(&maxUserKey)
	if err != nil && err != sql.ErrNoRows {
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al obtener MAX user key", Linea: processedLine}
	}
	userKey, _ := generateUserKey(db, "users", "user_id", maxUserKey)
	fmt.Println("userKey: ", userKey)

	var maxVoipKey int
	query = fmt.Sprintf("SELECT MAX(CAST(SUBSTRING(%s, 4) AS INT)) FROM %s WHERE account_type='CUSTOMER'", "customer_voipminute_id", "customer_voipminuts")
	//err = db.QueryRow(query).Scan(&maxUserKey)
	err = db.QueryRow(query).Scan(&maxVoipKey)
	if err != nil && err != sql.ErrNoRows {
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al obtener MAX Voip key", Linea: processedLine}
	}
	voipKey, _ := generateVoipKey(db, "customer_voipminuts", "customer_voipminute_id", maxVoipKey)
	fmt.Println("voipKey: ", voipKey)

	newCliente := models.Cliente{
		AccountId:    customerKey,
		Nombre:       payload.Nombre,
		Plan:         payload.Plan,
		MSISDN:       payload.MSISDN,
		Username:     payload.Username,
		Password:     payload.Password,
		Dominio:      payload.Dominio,
		Departamento: payload.Departamento,
		Provincia:    payload.Provincia,
		Distrito:     payload.Distrito,
		UpdateBy:     update_by,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result := initializers.DB.Create(&newCliente)
	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "msisdn already exists", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al insertar new msisdn", Linea: processedLine}
	}

	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	//if err != nil {
	//	tx.Rollback()
	//	return Error{Tipo: "Ejecución", Numero: msisdn, Mensaje: "error al generar hash password", Linea: processedLine}
	//}

	newCustomerSipAccount := models.CustomerSipAccount{
		Username: payload.MSISDN,
		//Secret:      string(hashedPassword),
		Secret:      "W1nn3t*2022#",
		Status:      "0",
		AccountId:   customerKey,
		Ipauthfrom:  "NO",
		CreatedBy:   update_by,
		CreatedDt:   now,
		UpdatedBy:   update_by,
		UserType:    "SWITCH",
		ExtensionId: "",
		Name:        name,
	}

	fmt.Printf("S.A.M. Intentando insertar new customer_sip_account: %+v\n", newCustomerSipAccount)
	result = initializers.DB.Create(&newCustomerSipAccount)
	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "customer SIP Account already exists", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: result.Error.Error() + "error al insertar new customer SIP account", Linea: processedLine}
	}

	newAccount := models.Account{
		AccountId:       customerKey,
		ParentAccountId: "",
		AccountType:     "CUSTOMER",
		StatusId:        "1",
		AccountLevel:    1,
		CreateBy:        update_by,
		CreateDt:        now,
		UpdateBy:        update_by,
	}

	result = initializers.DB.Create(&newAccount)
	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "account already exists", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al insertar new account", Linea: processedLine}
	}

	newUser := models.User{
		UserId:       userKey,
		AccountId:    customerKey,
		Username:     payload.Username,
		Secret:       claveWeb,
		Name:         payload.Nombre,
		EmailAddress: "jesus.bazan@gmail.com",
		Address:      address,
		UpdateBy:     update_by,
		CountryId:    169,
		StatusId:     1,
	}

	result = initializers.DB.Create(&newUser)
	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: result.Error.Error() + " user already exists", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al insertar new user", Linea: processedLine}
	}

	newCustomer := models.Customer{
		AccountId:    customerKey,
		CompanyName:  payload.Nombre,
		ContactName:  "administrador",
		Address:      address,
		CountryId:    169,
		StateCodeId:  1,
		Phone:        msisdn,
		EmailAddress: "jesus.bazan@gmail.com",
		UpdatedBy:    update_by,
		UpdatedDt:    now,
	}

	result = initializers.DB.Create(&newCustomer)
	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "customer already exists", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al insertar new Customer", Linea: processedLine}
	}

	tariff_id := "CUSTOMER54"
	if plan == "WINPHONE 110" {
		tariff_id = "TFWINPHO43"
	} else if plan == "WINPHONE 100" {
		tariff_id = "TFWINPHO54"
	} else if plan == "WINPHONE 10" {
		tariff_id = "TFWINPHO38"
	}

	newCustomerVoipminute := models.CustomerVoipminute{
		CustomerVoipminuteId: voipKey,
		AccountId:            customerKey,
		AccountType:          "CUSTOMER",
		TariffId:             tariff_id,
		Status:               "1",
		CreatedBy:            update_by,
		CreatedDt:            now,
		UpdatedBy:            update_by,
	}

	result = initializers.DB.Create(&newCustomerVoipminute)
	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "customer Voip Minute already exists", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al insertar new customer Voip Minute", Linea: processedLine}
	}

	saldo := 100.00
	if plan == "WINPHONE 110" {
		saldo = 110.00
	} else if plan == "WINPHONE 100" {
		saldo = 100.00
	} else if plan == "WINPHONE 10" {
		saldo = 10.00
	}

	newCustomerBalance := models.CustomerBalance{
		AccountId:      customerKey,
		ServiceType:    "SWITCH",
		CreditLimit:    0.00,
		MaxcreditLimit: 0.00,
		Balance:        float64(-saldo),
		UpdateDt:       now,
	}

	result = initializers.DB.Create(&newCustomerBalance)
	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "customer Balance already exists", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al insertar new customer balance", Linea: processedLine}
	}

	//INBOUND customer dialpattern
	newCustomerDialpattern := models.CustomerDialpattern{
		AccountId:     customerKey,
		MachingString: "%",
		AddString:     "%",
		DisplayString: "%=>%",
		Route:         "INBOUND",
		CreatedBy:     update_by,
		CreatedDt:     now,
		UpdatedBy:     update_by,
	}

	result = initializers.DB.Create(&newCustomerDialpattern)
	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "Inbound customer dialpattern already exists", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al insertar new Inbound customer dialpattern", Linea: processedLine}
	}

	//OUTBOUND customer dialpattern
	newCustomerDialpattern = models.CustomerDialpattern{
		AccountId:     customerKey,
		MachingString: "%",
		AddString:     "%",
		DisplayString: "%=>%",
		Route:         "OUTBOUND",
		CreatedBy:     update_by,
		CreatedDt:     now,
		UpdatedBy:     update_by,
	}

	result = initializers.DB.Create(&newCustomerDialpattern)
	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "Outbound customer dialpattern already exists", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al insertar new Outbound customer dialpattern", Linea: processedLine}
	}

	//customer dialplan
	newCustomerDialplan := models.CustomerDialplan{
		AccountId:     customerKey,
		DialplanId:    "RTMETASW35",
		MachingString: "%",
		DisplayString: "%=>RTMETASW35%",
		CreatedBy:     update_by,
		CreatedDt:     now,
		UpdatedBy:     update_by,
	}

	result = initializers.DB.Create(&newCustomerDialplan)
	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "customer dialplan already exists", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error al insertar new customer dialplan", Linea: processedLine}
	}

	// Confirmar la transacción
	err = tx.Commit()
	if err != nil {
		return Resultado{Tipo: "Error", Numero: msisdn, Mensaje: "error de commit en BD MySQL", Linea: processedLine}
	}

	return Resultado{Tipo: "Exitoso", Numero: msisdn, Mensaje: "se inserto registro en BD MySQL", Linea: processedLine}
}

func DeleteProcessFileHandler(c *fiber.Ctx) error {

	var api_results []Resultado

	// Obtener el archivo csv del cuerpo de la solicitud
	file, err := c.FormFile("archivo")
	fmt.Println(file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("error al obtener el archivo CSV")
	}

	// Abrir el archivo
	fileOpened, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("error al abrir el archivo CSV")
	}
	defer fileOpened.Close()

	// Crea un lector CSV
	reader := csv.NewReader(fileOpened)
	reader.Comma = ';'

	// Lee la primera línea (encabezados) y descártala
	firstLine, err := reader.Read()
	log.Println(firstLine)
	if err != nil {
		fmt.Println("Error al leer la primera línea:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("error al leer la primera línea")
	}

	//Conexión a la BD MySQL
	connStr := "ovswitch:ovswitch123@tcp(127.0.0.1:3306)/switch?charset=utf8mb4"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("error al conectarse a la bd MySQL")
	}
	defer db.Close()

	// Procesar todas las líneas del archivo
	processedLine := 2
	for {
		row, err := reader.Read()
		log.Println(row)
		if err == io.EOF {
			break // fin del archivo
		} else if err != nil {
			fmt.Println("Error al leer el archivo CSV:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("error al leer el archivo CSV")
		}

		// Asignar valores a la estructura
		phone := row[0]

		// Validar phone
		// ....

		// Procesar la información del archivo
		resultado := deleteRecord(db, phone, processedLine)
		if resultado.Tipo != "" {
			api_results = append(api_results, resultado)
		}
		processedLine += 1
	}

	if len(api_results) == 0 {
		//return c.SendString("archivo recibido y procesado correctamente")
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"no se actualizo ningun registro": file}})
	} else {
		// Retornar la información del archivo temporal
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Archivo procesado por API",
			"errores": api_results,
		})
	}
}

func deleteRecord(db *sql.DB, phone string, processedLine int) Resultado {
	var cliente models.Cliente

	// Iniciar una transacción
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "error al iniciar transacción en BD", Linea: processedLine}
	}

	//Encontrar campo account_id
	query := "SELECT * FROM clientes WHERE msisdn = ? LIMIT 1"
	result := initializers.DB.Raw(query, phone).First(&cliente)

	if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "error al buscar account_id en tabla clientes: " + result.Error.Error(), Linea: processedLine}
	}
	accountID := cliente.AccountId

	//Delete cliente
	result = initializers.DB.Delete(&models.Cliente{}, "msisdn = ?", phone)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "No existe un cliente con el phone indicado", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: result.Error.Error(), Linea: processedLine}
	}

	//Delete customer_sip_account
	result = initializers.DB.Delete(&models.CustomerSipAccount{}, "account_id = ?", accountID)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "No existe un customerSipAccount para el username(phone) indicado", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: result.Error.Error(), Linea: processedLine}
	}

	//Delete Account
	result = initializers.DB.Delete(&models.Account{}, "account_id = ?", accountID)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "No existe un Account para el username(phone) indicado", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: result.Error.Error(), Linea: processedLine}
	}

	//Delete User
	result = initializers.DB.Delete(&models.User{}, "account_id = ?", accountID)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "No existe un User para el username(phone) indicado", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: result.Error.Error(), Linea: processedLine}
	}

	//Delete Customer
	result = initializers.DB.Delete(&models.Customer{}, "account_id = ?", accountID)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "No existe un Customer para el username(phone) indicado", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: result.Error.Error(), Linea: processedLine}
	}

	//Delete CustomerVoipminute
	result = initializers.DB.Delete(&models.CustomerVoipminute{}, "account_id = ?", accountID)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "No existe un CustomerVoipMinute para el username(phone) indicado", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: result.Error.Error(), Linea: processedLine}
	}

	//Delete CustomerBalance
	result = initializers.DB.Delete(&models.CustomerBalance{}, "account_id = ?", accountID)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "No existe un CustomerBalance para el username(phone) indicado", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: result.Error.Error(), Linea: processedLine}
	}

	//Delete CustomerDialpattern INBOUND & OUTBOUND
	result = initializers.DB.Delete(&models.CustomerDialpattern{}, "account_id = ?", accountID)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "No existe un CustomerDialpattern para el username(phone) indicado", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: result.Error.Error(), Linea: processedLine}
	}

	//Delete CustomerDialplan
	result = initializers.DB.Delete(&models.CustomerDialplan{}, "account_id = ?", accountID)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "No existe un CustomerDialplan para el username(phone) indicado", Linea: processedLine}
	} else if result.Error != nil {
		tx.Rollback()
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: result.Error.Error(), Linea: processedLine}
	}

	// Confirmar la transacción
	err = tx.Commit()
	if err != nil {
		return Resultado{Tipo: "Error", Numero: phone, Mensaje: "error de commit en BD MySQL", Linea: processedLine}
	}

	return Resultado{Tipo: "Exitoso", Numero: phone, Mensaje: "se elimino registro en BD MySQL", Linea: processedLine}
}
