package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	CustomerId   int       `gorm:"type:int;primary_key;not_null" json:"customerId,omitempty"`
	AccountId    string    `gorm:"type:char(25);primary_key" json:"accountId,omitempty"`
	CompanyName  string    `gorm:"varchar(36)" json:"companyName,omitempty"`
	ContactName  string    `gorm:"varchar(150)" json:"contactName,omitempty"`
	StateCodeId  int       `gorm:"type:int" json:"stateCodeId,omitempty"`
	Phone        string    `gorm:"varchar(30)" json:"phone,omitempty"`
	EmailAddress string    `gorm:"column:emailaddress;varchar(100)" json:"emailAddress,omitempty"`
	Address      string    `gorm:"varchar(128)" json:"address,omitempty"`
	CountryId    int       `gorm:"type:int" json:"countryId,omitempty"`
	BillingType  string    `gorm:"varchar(16);not null;default:'prepaid'" json:"billingType,omitempty"`
	BillingCycle string    `gorm:"varchar(16);not null;default:'monthly'" json:"billingCycle,omitempty"`
	CreatedBy    string    `gorm:"varchar(36)" json:"createdBy,omitempty"`
	UpdatedBy    string    `gorm:"varchar(36)" json:"updatedBy,omitempty"`
	CreatedDt    time.Time `gorm:"not null;default:'1970-01-01'" json:"createdDt,omitempty"`
	UpdatedDt    time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedDt,omitempty"`
}

func (customer *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

type CreateCustomerSchema struct {
	AccountId    string `json:"accountId" validate:"required"`
	CompanyName  string `json:"companyName" validate:"required"`
	Address      string `json:"address" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	EmailAddress string `json:"emailAddress" validate:"required"`
	ContactName  string `json:"contactName" validate:"required"`
	StateCodeId  int    `json:"stateCodeId" validate:"required"`
	CountryId    string `json:"countryId" validate:"required"`
	BillingType  string `json:"billingType" validate:"required"`
	BillingCycle string `json:"billingCycle" validate:"required"`
}

type UpdateCustomerSchema struct {
	AccountId    string `json:"accountId,omitempty"`
	CompanyName  string `json:"companyName,omitempty"`
	ContactName  string `json:"contactName,omitempty"`
	StateCodeId  int    `json:"stateCodeId,omitempty"`
	Address      string `json:"address,omitempty"`
	Phone        string `json:"phone,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	CountryId    string `json:"countryId,omitempty"`
	BillingType  string `json:"billingType,omitempty"`
	BillingCycle string `json:"billingCycle,omitempty"`
}
