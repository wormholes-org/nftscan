package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type SnftRec struct {
	Blocknumber          uint64  `json:"blocknumber" gorm:"type:bigint unsigned;comment:'block height'"`
	Creatoraddr          string  `json:"creataddr" gorm:"type:char(42);comment:'Create nft address'"`
	Ownaddr              string  `json:"ownaddr" gorm:"type:char(42);comment:'nft owner address'"`
	Contract             string  `json:"contract" gorm:"type:char(42);comment:'contract address'"`
	Nftaddr              string  `json:"nftaddr" gorm:"type:char(42);comment:'Chain of wormholes uniquely identifies the nft flag'"`
	Snft                 string  `json:"snft" gorm:"type:char(42) ;comment:'wormholes chain snft'"`
	Snftcollection       string  `json:"snftcollection" gorm:"type:char(42) DEFAULT NULL;comment:'Wormholes chain snft collection'"`
	Snftstage            string  `json:"snftstage" gorm:"type:char(42) DEFAULT NULL;comment:'wormholes chain snft period'"`
	Name                 string  `json:"name" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft name'"`
	Desc                 string  `json:"desc" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'nft description'"`
	Meta                 string  `json:"meta" gorm:"type:longtext CHARACTER SET utf8mb4  ;comment:'meta information'"`
	Category             string  `json:"categories" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft category'"`
	Royalty              float64 `json:"royalty" gorm:"type:float;COMMENT:'royalty'"`
	Sourceurl            string  `json:"source_url" gorm:"type:varchar(200);comment:'nfc raw data hold address'"`
	Md5                  string  `json:"md5" gorm:"type:longtext ;comment:'Picture md5 value'"`
	Collectionsname      string  `json:"collections" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'NFT collection name'"`
	Collectionscreator   string  `json:"collection_creator_addr" gorm:"type:char(42) ;comment:'Collection creator address'"`
	Collectionsexchanger string  `json:"collections_exchanger" gorm:"type:longtext CHARACTER SET utf8mb4 NOT NULL;comment:'Collection exchange'"`
	Collectionscategory  string  `json:"collections_category" gorm:"type:varchar(200) CHARACTER SET utf8mb4 ;comment:'nft collection category'"`
	Collectionsimgurl    string  `json:"collections_img_url" gorm:"type:longtext;comment:'logo'"`
	Collectionsdesc      string  `json:"collections_desc" gorm:"type:longtext CHARACTER SET utf8mb4;comment:'Collection description'"`
}

type Snfts struct {
	gorm.Model
	SnftRec
}

func (v Snfts) TableName() string {
	return "snfts"
}

func snftRecToSnftinfo(sr *SnftRec) SnftInfo {
	var si SnftInfo
	si.CreatorAddr = sr.Creatoraddr
	si.Ownaddr = sr.Ownaddr
	si.Contract = sr.Contract
	si.Nftaddr = sr.Nftaddr
	si.Name = sr.Name
	si.Desc = sr.Desc
	si.Meta = sr.Meta
	si.Category = sr.Category
	si.Royalty = sr.Royalty
	si.SourceUrl = sr.Sourceurl
	si.Md5 = sr.Md5
	si.CollectionsName = sr.Collectionsname
	si.CollectionsCreator = sr.Collectionscreator
	si.CollectionsExchanger = sr.Collectionsexchanger
	si.CollectionsCategory = sr.Collectionscategory
	si.CollectionsImgUrl = sr.Collectionsimgurl
	si.CollectionsDesc = sr.Collectionsdesc
	return si
}

func GetBlockSnfts(blocks string) ([]SnftInfo, error) {
	fmt.Println("GetBlockSnfts() blocks=", blocks)
	b, _ := strconv.ParseUint(blocks, 10, 64)
	sysp := SysParams{}
	err := GetScanDB().Model(&SysParams{}).Last(&sysp)
	if err.Error != nil {
		return nil, errors.New("snft not sync block number.")
	}
	if b >= sysp.Scannumber {
		return nil, errors.New("not sync block number.")
	}
	fmt.Println("GetBlockSnfts() b=", b, " sysp.Scannumber=", sysp.Scannumber)
	snfts := []Snfts{}
	err = GetScanDB().Model(&Snfts{}).Where("blocknumber = ?", b).Find(&snfts)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		return nil, err.Error
	}
	var snftinfo []SnftInfo
	for _, snft := range snfts {
		snftinfo = append(snftinfo, snftRecToSnftinfo(&snft.SnftRec))
	}
	return snftinfo, nil
}
