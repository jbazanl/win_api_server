package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomerCallerid struct {
	Id            int       `gorm:"type:int;primary_key;not_null" json:"id,omitempty"`
	MachingString string    `gorm:"type:char(30);default:'%'" json:"machingString,omitempty"`
	AccountId     string    `gorm:"type:char(30)" json:"accountId,omitempty"`
	AddString     string    `gorm:"type:char(15);default:'%'" json:"addString,omitempty"`
	DisplayString string    `gorm:"type:char(60);default:'%=>%'" json:"displayString,omitempty"`
	ActionType    string    `gorm:"type:char(1);default:'1'" json:"actionType,omitempty"`
	ServiceType   string    `gorm:"type:char(16);default:'SWITCH'" json:"serviceType,omitempty"`
	Route         string    `gorm:"type:char(16);default:'OUTBOUND'" json:"route,omitempty"`
	CreatedDt     time.Time `gorm:"not null;default:'1970-01-01'" json:"createdDt,omitempty"`
	CreatedBy     string    `gorm:"varchar(30)" json:"createdBy,omitempty"`
	UpdatedBy     string    `gorm:"varchar(36)" json:"updatedBy,omitempty"`
}

func (customerCallerid *CustomerCallerid) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (CustomerCallerid) TableName() string {
	return "customer_callerid"
}

type CreateCustomerCalleridSchema struct {
	Id            string    `json:"id" validate:"required"`
	AccountId     string    `json:"accountId" validate:"required"`
	ServiceType   string    `json:"serviceType" validate:"required"`
	MachingString string    `json:"machingString" validate:"required"`
	AddString     string    `json:"addString" validate:"required"`
	DisplayString string    `json:"displayString" validate:"required"`
	ActionType    string    `json:"actionType" validate:"required"`
	Route         string    `json:"route" validate:"required"`
	CreatedBy     string    `json:"createdBy" validate:"required"`
	CreatedDt     time.Time `json:"createdDt" validate:"required"`
	UpdatedBy     string    `json:"updatedBy" validate:"required"`
}

type UpdateCustomerCalleridSchema struct {
	Id            string    `json:"id,omitempty"`
	AccountId     string    `json:"accountId,omitempty"`
	ActionType    string    `json:"actionType,omitempty"`
	ServiceType   string    `json:"serviceType,omitempty"`
	MachingString string    `json:"machingString,omitempty"`
	AddString     string    `json:"addString,omitempty"`
	DisplayString string    `json:"displayString,omitempty"`
	Route         string    `json:"route,omitempty"`
	CreatedBy     string    `json:"createdBy,omitempty"`
	CreatedDt     time.Time `json:"createdDt,omitempty"`
	UpdatedBy     string    `json:"updatedBy,omitempty"`
}
