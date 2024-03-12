package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id           int       `gorm:"type:int;primary_key;not_null" json:"id,omitempty"`
	UserId       string    `gorm:"type:char(30)" json:"userId,omitempty"`
	AccountId    string    `gorm:"type:char(30)" json:"accountId,omitempty"`
	UserType     string    `gorm:"varchar(30);default:'CUSTOMERADMIN'" json:"userType,omitempty"`
	Username     string    `gorm:"varchar(30)" json:"username,omitempty"`
	Secret       string    `gorm:"varchar(30)" json:"secret,omitempty"`
	Name         string    `gorm:"varchar(100)" json:"name,omitempty"`
	EmailAddress string    `gorm:"column:emailaddress;varchar(100)" json:"emailAddress,omitempty"`
	Phone        string    `gorm:"varchar(50)" json:"phone,omitempty"`
	Address      string    `gorm:"column:address;varchar(256)" json:"address,omitempty"`
	CountryId    int       `gorm:"type:int;default:0" json:"countryId,omitempty"`
	StatusId     int       `gorm:"type:int;default:1" json:"statusId,omitempty"`
	CreateDt     time.Time `gorm:"not null;default:'1970-01-01'" json:"createDt,omitempty"`
	CreateBy     string    `gorm:"varchar(36)" json:"createBy,omitempty"`
	UpdateBy     string    `gorm:"varchar(36)" json:"updatedBy,omitempty"`
}

func (account *User) BeforeCreate(tx *gorm.DB) (err error) {
	//account.AccountId = uuid.New().String()
	return nil
}

type CreateUserSchema struct {
	UserId       string `json:"userId" validate:"required"`
	AccountId    string `json:"accountId" validate:"required"`
	Username     string `json:"username" validate:"required"`
	Secret       string `json:"secret" validate:"required"`
	Name         string `json:"name" validate:"required"`
	EmailAddress string `json:"emailAddress" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	Address      string `json:"address" validate:"required"`
	CountryId    int    `json:"countryId" validate:"required"`
	StatusId     int    `json:"statusId" validate:"required"`
	UpdateBy     string `json:"updatedBy" validate:"required"`
}

type UpdateUserSchema struct {
	UserId       string `json:"userId,omitempty"`
	AccountId    string `json:"accountId,omitempty"`
	Username     string `json:"username,omitempty"`
	Secret       string `json:"secret,omitempty"`
	Name         string `json:"name,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Address      string `json:"address,omitempty"`
	CountryId    int    `json:"countryId,omitempty"`
	StatusId     int    `json:"statusId,omitempty"`
	UpdateBy     string `json:"updatedBy,omitempty"`
}
