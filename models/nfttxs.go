package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type NftTxRec struct {
	Operator         string `json:"operator" gorm:"type:char(42) ;comment:'nft operator address'"`
	Fromaddr         string `json:"fromaddr" gorm:"type:char(42);comment:'seller address'"`
	Toaddr           string `json:"toaddr" gorm:"type:char(42);comment:'Buyer's address'"`
	Contract         string `json:"contract" gorm:"type:char(42);comment:'contract address'"`
	Tokenid          string `json:"tokenid" gorm:"type:char(42);comment:'Uniquely identifies the nft flag'"`
	Value            string `json:"value" gorm:"type:char(66);comment:'the transaction price'"`
	Price            string `json:"price" gorm:"type:char(66);comment:'the deal price'"`
	Ratio            string `json:"royalty" gorm:"type:char(10) ;COMMENT:'royalty'"`
	Txhash           string `json:"txhash" gorm:"type:char(66);comment:'transaction hash'"`
	Ts               string `json:"transtime" gorm:"type:char(20);comment:'nft creation time'"`
	Blocknumber      uint64 `json:"blocknumber" gorm:"type:bigint unsigned;comment:'block height'"`
	Transactionindex string `json:"transindex" gorm:"type:char(10);comment:'nft trans index'"`
	Metaurl          string `json:"metaurl" gorm:"type:longtext;comment:'meta information'"`
	Nftaddr          string `json:"nft_address" gorm:"type:char(42);comment:'Chain of wormholes uniquely identifies the nft flag'"`
	Nonce            string `json:"nonce" gorm:"type:char(10);comment:'nft trans nonce'"`
	Status           bool   `json:"status" gorm:"type:bool;comment:'nft transaction status'"`
	Transtype        int    `json:"transtype" gorm:"type:int;comment:'nft transaction type'"`
}

type NftTxs struct {
	gorm.Model
	NftTxRec
}

func (v NftTxs) TableName() string {
	return "nfttxs"
}

func NftTxRecToNftTx(nr *NftTxRec) NftTx {
	var n NftTx
	n.Operator = nr.Operator
	n.From = nr.Fromaddr
	n.To = nr.Toaddr
	n.Contract = nr.Contract
	n.TokenId = nr.Tokenid
	n.Value = nr.Value
	n.Price = nr.Price
	n.Ratio = nr.Ratio
	n.TxHash = nr.Txhash
	n.Ts = nr.Ts
	n.BlockNumber = strconv.FormatUint(nr.Blocknumber, 10)
	n.TransactionIndex = nr.Transactionindex
	n.MetaUrl = nr.Metaurl
	n.NftAddr = nr.Nftaddr
	n.Nonce = nr.Nonce
	n.Status = nr.Status
	n.TransType = nr.Transtype
	return n
}

func GetBlockTrans(blocks string) ([]NftTx, error) {
	fmt.Println("GetBlockTrans() blocks=", blocks)
	b, _ := strconv.ParseUint(blocks, 10, 64)
	sysp := SysParams{}
	err := GetScanDB().Model(&SysParams{}).Last(&sysp)
	if err.Error != nil {
		return nil, err.Error
	}
	fmt.Println("GetBlockTrans() b=", b, " sysp.Scannumber=", sysp.Scannumber)
	if b >= sysp.Scannumber {
		return nil, errors.New("nft not sync block number.")
	}
	nftTxs := []NftTxs{}
	err = GetScanDB().Model(&NftTxs{}).Where("blocknumber = ?", b).Find(&nftTxs)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		return nil, err.Error
	}
	var rnft []NftTx
	for _, tx := range nftTxs {
		rnft = append(rnft, NftTxRecToNftTx(&tx.NftTxRec))
	}
	return rnft, nil
}
