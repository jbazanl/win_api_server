package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomerBalance struct {
	Id             int       `gorm:"type:int;primary_key;not_null" json:"id,omitempty"`
	AccountId      string    `gorm:"type:char(30)" json:"accountId,omitempty"`
	ServiceType    string    `gorm:"type:char(30);default:'SWITCH'" json:"serviceType,omitempty"`
	CreditLimit    float64   `gorm:"type:float64;default:0.00" json:"creditLimit,omitempty"`
	MaxcreditLimit float64   `gorm:"type:float64;default:0.00" json:"maxcreditLimit,omitempty"`
	Balance        float64   `gorm:"type:float64;default:0.00" json:"balance,omitempty"`
	UpdateDt       time.Time `gorm:"not null;default:'1970-01-01'" json:"updateDt,omitempty"`
}

func (customerBalance *CustomerBalance) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (CustomerBalance) TableName() string {
	return "customer_balance"
}

type CreateCustomerBalanceSchema struct {
	Id             string  `json:"id" validate:"required"`
	AccountId      string  `json:"accountId" validate:"required"`
	ServiceType    string  `json:"serviceType" validate:"required"`
	CreditLimit    float64 `json:"creditLimit" validate:"required"`
	MaxcreditLimit float64 `json:"maxcreditLimit" validate:"required"`
	Balance        float64 `json:"balance" validate:"required"`
	UpdateDt       string  `json:"updateDt" validate:"required"`
}

type UpdateCustomerBalanceSchema struct {
	Id             string  `json:"id,omitempty"`
	AccountId      string  `json:"accountId,omitempty"`
	ServiceType    string  `json:"serviceType,omitempty"`
	CreditLimit    float64 `json:"creditLimit,omitempty"`
	MaxcreditLimit float64 `json:"maxcreditLimit,omitempty"`
	Balance        float64 `json:"balance,omitempty"`
	UpdateDt       string  `json:"updateDt,omitempty"`
}
