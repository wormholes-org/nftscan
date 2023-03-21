package models

import "gorm.io/gorm"

type SysParamsRec struct {
	Blocknumber    uint64 `json:"blocknumber" gorm:"type:bigint unsigned DEFAULT 0;comment:'block height'"`
	Scannumber     uint64 `json:"scannumber" gorm:"type:bigint unsigned DEFAULT 1;comment:'Scanned block height'"`
	Scansnftnumber uint64 `json:"scansnftnumber" gorm:"type:bigint unsigned DEFAULT 0;comment:'Scanned snft block height'"`
	Savedsnft      string `json:"snft" gorm:"type:char(42) ;comment:'snft backed up to ipfs'"`
}

type SysParams struct {
	gorm.Model
	SysParamsRec
}

func (v SysParams) TableName() string {
	return "sysparams"
}
