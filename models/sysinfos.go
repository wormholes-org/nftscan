package models

import "gorm.io/gorm"

type SysInfoRec struct {
	Nfttotal     uint64 `json:"nfttotal" gorm:"type:bigint unsigned DEFAULT 0;comment:'Total number of nfts'"`
	Snfttotal    uint64 `json:"Snfttotal" gorm:"type:bigint unsigned DEFAULT 0;comment:'total number of snfts'"`
	Nfthour      uint64 `json:"nfthour" gorm:"type:int unsigned DEFAULT 0 ;comment:'nft additions per hour'"`
	Snfthour     uint64 `json:"snfthour" gorm:"type:bigint unsigned DEFAULT 0;comment:'snft increments per hour'"`
	Transcnthour uint64 `json:"transcnthour" gorm:"type:int unsigned DEFAULT 0 ;comment:'Transactions per hour'"`
	Transamthour uint64 `json:"transamthour" gorm:"type:bigint unsigned DEFAULT 0;comment:'Total transaction amount per hour'"`
	Transcntday  uint64 `json:"transcntday" gorm:"type:int unsigned DEFAULT 0 ;comment:'daily transactions'"`
	Transamtday  uint64 `json:"transamtday" gorm:"type:bigint unsigned DEFAULT 0;comment:'Total daily transaction amount'"`
}

type SysInfos struct {
	gorm.Model
	SysInfoRec
}

func (v SysInfos) TableName() string {
	return "sysinfos"
}
