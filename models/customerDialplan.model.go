package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomerDialplan struct {
	Id            int       `gorm:"type:int;primary_key;not_null" json:"id,omitempty"`
	AccountId     string    `gorm:"type:char(30)" json:"accountId,omitempty"`
	DialplanId    string    `gorm:"type:char(30)" json:"dialplanId,omitempty"`
	MachingString string    `gorm:"type:char(30);default:'%'" json:"machingString,omitempty"`
	DisplayString string    `gorm:"type:char(30)" json:"displayString,omitempty"`
	AddString     string    `gorm:"type:char(30)" json:"addString,omitempty"`
	RemoveString  string    `gorm:"type:char(30)" json:"removeString,omitempty"`
	CreatedDt     time.Time `gorm:"not null;default:'1970-01-01'" json:"createdDt,omitempty"`
	CreatedBy     string    `gorm:"varchar(30)" json:"createdBy,omitempty"`
	UpdatedBy     string    `gorm:"varchar(36)" json:"updatedBy,omitempty"`
}

func (customerDialplan *CustomerDialplan) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (CustomerDialplan) TableName() string {
	return "customer_dialplan"
}

type CreateCustomerDialplanSchema struct {
	Id            string    `json:"id" validate:"required"`
	AccountId     string    `json:"accountId" validate:"required"`
	DialplanId    string    `json:"dialplanId" validate:"required"`
	MachingString string    `json:"machingString" validate:"required"`
	AddString     string    `json:"addString" validate:"required"`
	RemoveString  string    `json:"removeString" validate:"required"`
	DisplayString string    `json:"displayString" validate:"required"`
	CreatedBy     string    `json:"createdBy" validate:"required"`
	CreatedDt     time.Time `json:"createdDt" validate:"required"`
	UpdatedBy     string    `json:"updatedBy" validate:"required"`
}

type UpdateCustomerDialplanSchema struct {
	Id            string    `json:"id,omitempty"`
	AccountId     string    `json:"accountId,omitempty"`
	DialplanId    string    `json:"dialplanId,omitempty"`
	MachingString string    `json:"machingString,omitempty"`
	AddString     string    `json:"addString,omitempty"`
	RemoveString  string    `json:"removeString,omitempty"`
	DisplayString string    `json:"displayString,omitempty"`
	CreatedBy     string    `json:"createdBy,omitempty"`
	CreatedDt     time.Time `json:"createdDt,omitempty"`
	UpdatedBy     string    `json:"updatedBy,omitempty"`
}
