package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomerVoipminute struct {
	Id                   int       `gorm:"type:int;primary_key;not_null" json:"id,omitempty"`
	CustomerVoipminuteId string    `gorm:"type:char(30)" json:"customerVoipminuteId,omitempty"`
	AccountId            string    `gorm:"type:char(30)" json:"accountId,omitempty"`
	AccountType          string    `gorm:"type:char(30)" json:"accountType,omitempty"`
	TariffId             string    `gorm:"type:char(30)" json:"tariffId,omitempty"`
	Status               string    `gorm:"type:char(1);default:'1'" json:"status,omitempty"`
	CreatedDt            time.Time `gorm:"not null;default:'1970-01-01'" json:"createdDt,omitempty"`
	CreatedBy            string    `gorm:"varchar(30)" json:"createdBy,omitempty"`
	UpdatedBy            string    `gorm:"varchar(36)" json:"updatedBy,omitempty"`
}

func (customerVoipminute *CustomerVoipminute) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (CustomerVoipminute) TableName() string {
	return "customer_voipminuts"
}

type CreateCustomerVoipminuteSchema struct {
	Id                   string    `json:"id" validate:"required"`
	CustomerVoipminuteId string    `json:"customerVoipminuteId" validate:"required"`
	AccountId            string    `json:"accountId" validate:"required"`
	AccountType          string    `json:"accountType" validate:"required"`
	TariffId             string    `json:"tariffId" validate:"required"`
	Status               string    `json:"status" validate:"required"`
	CreatedBy            string    `json:"createdBy" validate:"required"`
	CreatedDt            time.Time `json:"createdDt" validate:"required"`
	UpdatedBy            string    `json:"updatedBy" validate:"required"`
}

type UpdateCustomerVoipminuteSchema struct {
	Id                   string    `json:"id,omitempty"`
	CustomerVoipminuteId string    `json:"customerVoipminuteId,omitempty"`
	AccountId            string    `json:"accountId,omitempty"`
	AccountType          string    `json:"accountType,omitempty"`
	TariffId             string    `json:"tariffId,omitempty"`
	Status               string    `json:"status,omitempty"`
	CreatedBy            string    `json:"createdBy,omitempty"`
	CreatedDt            time.Time `json:"createdDt,omitempty"`
	UpdatedBy            string    `json:"updatedBy,omitempty"`
}
