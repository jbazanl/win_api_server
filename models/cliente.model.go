package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cliente struct {
	Id            string    `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	AccountId     string    `gorm:"varchar(25);not null" json:"accountId,omitempty"`
	Nombre        string    `gorm:"varchar(255)" json:"numero,omitempty"`
	Plan          string    `gorm:"varchar(25);not null" json:"plan,omitempty"`
	TipoDocumento string    `gorm:"varchar(25)" json:"tipoDocumento,omitempty"`
	Documento     string    `gorm:"varchar(25)" json:"documento,omitempty"`
	MSISDN        string    `gorm:"type:varchar(16);uniqueIndex:idx_cliente_msisdn,LENGTH(16);not null" json:"msisdn,omitempty"`
	Username      string    `gorm:"varchar(25);not null" json:"username,omitempty"`
	Password      string    `gorm:"varchar(25);not null" json:"password,omitempty"`
	Dominio       string    `gorm:"varchar(36);not null" json:"dominio,omitempty"`
	IsRegistered  bool      `gorm:"default:false;not null" json:"isRegistered"`
	IsBlocked     bool      `gorm:"default:false;not null" json:"isBlocked"`
	Departamento  string    `gorm:"varchar(36);not null" json:"departamento,omitempty"`
	Provincia     string    `gorm:"varchar(36);not null" json:"provincia,omitempty"`
	Distrito      string    `gorm:"varchar(36);not null" json:"distrito,omitempty"`
	CreatedAt     time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt     time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
	UpdateBy      string    `gorm:"varchar(30);not null" json:"updatedBy,omitempty"`
}

func (cliente *Cliente) BeforeCreate(tx *gorm.DB) (err error) {
	cliente.Id = uuid.New().String()
	return nil
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

type CreateClienteSchema struct {
	Nombre    string `json:"nombre" validate:"required"`
	Plan      string `json:"plan" validate:"required"`
	AccountId string `json:"accountId,omitempty"`
	MSISDN    string `json:"msisdn" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Dominio   string `json:"dominio" validate:"required"`
	//TipoDocumento  string `json:"tipoDocumento,omitempty"`
	//Documento  string `json:"documento,omitempty"`
	//IsRegistered bool   `json:"isRegistered,omitempty"`
	//IsBlocked bool   `json:"isBlocked,omitempty"`
	Departamento string `json:"departamento,omitempty"`
	Provincia    string `json:"provincia,omitempty"`
	Distrito     string `json:"distrito,omitempty"`
	UpdateBy     string `json:"updatedBy,omitempty"`
}

type UpdateClienteSchema struct {
	Nombre        string `json:"nombre,omitempty"`
	Plan          string `json:"plan,omitempty"`
	AccountId     string `json:"accountId,omitempty"`
	MSISDN        string `json:"msisdn,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Dominio       string `json:"dominio,omitempty"`
	TipoDocumento string `json:"tipoDocumento,omitempty"`
	Documento     string `json:"documento,omitempty"`
	IsRegistered  *bool  `json:"isRegistered,omitempty"`
	IsBlocked     *bool  `json:"isBlocked,omitempty"`
	Departamento  string `json:"departamento,omitempty"`
	Provincia     string `json:"provincia,omitempty"`
	Distrito      string `json:"distrito,omitempty"`
	UpdateBy      string `json:"updatedBy,omitempty"`
}
