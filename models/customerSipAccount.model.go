package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomerSipAccount struct {
	Id                 int       `gorm:"type:int;primary_key;not_null" json:"id,omitempty"`
	Username           string    `gorm:"varchar(30)" json:"username,omitempty"`
	Secret             string    `gorm:"varchar(30)" json:"secret,omitempty"`
	Ipaddress          string    `gorm:"column:ipaddress;varchar(100)" json:"ipaddress,omitempty"`
	Status             string    `gorm:"type:int;default:'0'" json:"status,omitempty"`
	AccountId          string    `gorm:"type:char(30)" json:"accountId,omitempty"`
	SipCc              int       `gorm:"type:int;default:2" json:"sipCc,omitempty"`
	SipCps             int       `gorm:"type:int;default:2" json:"sipCps,omitempty"`
	Ipauthfrom         string    `gorm:"type:char(5);default:'NO'" json:"ipauthfrom,omitempty"`
	ExtensionNo        int       `gorm:"type:int" json:"extensionNo,omitempty"`
	ExtensionId        string    `gorm:"type:char(30)" json:"extensionId,omitempty"`
	CreatedDt          time.Time `gorm:"not null;default:'1970-01-01'" json:"createdDt,omitempty"`
	CreatedBy          string    `gorm:"varchar(30)" json:"createdBy,omitempty"`
	CreatedByAccountId string    `gorm:"varchar(30)" json:"createdByAccountId,omitempty"`
	UpdatedBy          string    `gorm:"varchar(36)" json:"updatedBy,omitempty"`
	UserType           string    `gorm:"type:char(6);default:'PBX'" json:"userType,omitempty"`
	Name               string    `gorm:"varchar(100)" json:"name,omitempty"`
	EmailAddress       string    `gorm:"varchar(150)" json:"emailAddress,omitempty"`
	PhoneNumber        string    `gorm:"varchar(20)" json:"phoneNumber,omitempty"`
}

func (customerSipAccount *CustomerSipAccount) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (CustomerSipAccount) TableName() string {
	return "customer_sip_account"
}

type CreateCustomerSipAccountSchema struct {
	Id                 string    `json:"id" validate:"required"`
	Username           string    `json:"username" validate:"required"`
	Secret             string    `json:"secret" validate:"required"`
	Ipaddress          string    `json:"ipaddress" validate:"required"`
	Status             string    `json:"status" validate:"required"`
	AccountId          string    `json:"accountId" validate:"required"`
	Ipauthfrom         string    `json:"ipauthfrom" validate:"required"`
	CreatedBy          string    `json:"createdBy" validate:"required"`
	CreatedByAccountId string    `json:"createdByAccountId" validate:"required"`
	CreatedDt          time.Time `json:"createdDt" validate:"required"`
	UpdatedBy          string    `json:"updatedBy" validate:"required"`
	UserType           string    `json:"userType" validate:"required"`
	ExtensionNo        int       `json:"extensionNo" validate:"required"`
	ExtensionId        string    `json:"extensiontId" validate:"required"`
	Name               string    `json:"name" validate:"required"`
	EmailAddress       string    `json:"emailAddress" validate:"required"`
	PhoneNumber        string    `json:"phoneNumber" validate:"required"`
}

type UpdateCustomerSipAccountSchema struct {
	Id                 string    `json:"id,omitempty"`
	Username           string    `json:"username,omitempty"`
	Secret             string    `json:"secret,omitempty"`
	Ipaddress          string    `json:"ipaddress,omitempty"`
	Status             string    `json:"status,omitempty"`
	AccountId          string    `json:"accountId,omitempty"`
	Ipauthfrom         string    `json:"ipauthfrom,omitempty"`
	CreatedBy          string    `json:"createdBy,omitempty"`
	CreatedByAccountId string    `json:"createdByAccountId,omitempty"`
	CreatedDt          time.Time `json:"createdDt,omitempty"`
	UpdatedBy          string    `json:"updatedBy,omitempty"`
	UserType           string    `json:"userType,omitempty"`
	ExtensionNo        int       `json:"extensionNo,omitempty"`
	ExtensionId        string    `json:"extensionId,omitempty"`
	Name               string    `json:"name,omitempty"`
	EmailAddress       string    `json:"emailAddress,omitempty"`
	PhoneNumber        string    `json:"phoneNumber,omitempty"`
}
