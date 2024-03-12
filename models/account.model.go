package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	Id               int       `gorm:"type:int;primary_key;not_null" json:"id,omitempty"`
	AccountId        string    `gorm:"type:char(25)" json:"accountId,omitempty"`
	ParentAccountId  string    `gorm:"type:char(25)" json:"parentAccountId,omitempty"`
	StatusId         string    `gorm:"type:char(1)" json:"statusId,omitempty"`
	AccountType      string    `gorm:"varchar(16);default:'CUSTOMER'" json:"accountType,omitempty"`
	AccountLevel     int       `gorm:"type:int" json:"accountLevel,omitempty"`
	Dp               int       `gorm:"type:int;default:4" json:"dp,omitempty"`
	CurrencyId       int       `gorm:"type:int;default:4" json:"currencyId,omitempty"`
	AccountCc        int       `gorm:"type:int;default:2" json:"accountCc,omitempty"`
	AccountCps       int       `gorm:"type:int;default:2" json:"accountCps,omitempty"`
	TaxNumber        string    `gorm:"varchar(30)" json:"taxNumber,omitempty"`
	TaxType          string    `gorm:"varchar(16);default:'exclusive'" json:"taxType,omitempty"`
	AccountCodecs    string    `gorm:"varchar(255);default:'PCMU,PCMA'" json:"accountCodecs,omitempty"`
	MediaTranscoding string    `gorm:"char(1;default:'0'" json:"mediaTranscoding,omitempty"`
	MediaRtpproxy    string    `gorm:"char(1;default:'1'" json:"mediaRtpproxy,omitempty"`
	MaxCallDuration  int       `gorm:"column:max_callduration;type:int;default:120" json:"maxCallDuration,omitempty"`
	CreateDt         time.Time `gorm:"not null;default:'1970-01-01'" json:"createDt,omitempty"`
	CreateBy         string    `gorm:"varchar(36)" json:"createBy,omitempty"`
	UpdateBy         string    `gorm:"varchar(36)" json:"updateBy,omitempty"`
	//UpdatedDt    time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedDt,omitempty"`
}

func (account *Account) BeforeCreate(tx *gorm.DB) (err error) {
	//account.AccountId = uuid.New().String()
	return nil
}

func (Account) TableName() string {
	return "account"
}

type CreateAccountSchema struct {
	AccountId       string    `json:"accountId" validate:"required"`
	ParentAccountId string    `json:"companyName" validate:"required"`
	StatusId        string    `json:"statusId" validate:"required"`
	AccountType     string    `json:"accountType" validate:"required"`
	AccountLevel    string    `json:"accountLevel" validate:"required"`
	CreateDt        time.Time `json:"createDt" validate:"required"`
	CreateBy        string    `json:"createBy" validate:"required"`
	UpdateBy        string    `json:"updateBy" validate:"required"`
}

type UpdateAccountSchema struct {
	AccountId       string    `json:"accountId,omitempty"`
	ParentAccountId string    `json:"parentAccountId,omitempty"`
	StatusId        string    `json:"statusId,omitempty"`
	AccountType     string    `json:"accountType,omitempty"`
	AccountLevel    string    `json:"accountLevel,omitempty"`
	CreateDt        time.Time `json:"createDt,omitempty"`
	CreateBy        string    `json:"createBy,omitempty"`
	UpdateBy        string    `json:"updateBy,omitempty"`
}
