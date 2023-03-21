package controllers

import (
	"errors"
	beego "github.com/beego/beego/v2/server/web"
	"log"
	"regexp"
)

var (
	ERRINPUTINVALID = errors.New("549,input data invalid")
	ERRTOKEN        = errors.New("550,token invalid, please relogin!")
	ERRINPUT        = errors.New("551,Incorrect user information entered")
)

const (
	PattenString      = "^[0-9a-zA-Z_]+$"
	PattenNumber      = "^[0-9]+$"
	PattenHex         = "^[0-9a-fA-F]+$"
	PattenOperator    = "^[<>=]+$"
	PattenEmail       = "^[A-Za-z0-9]+([-_.][A-Za-z0-9]+)*@([A-Za-z0-9]+[-.])+[A-Za-z0-9]{2,4}$"
	PattenAddr        = "^0x[0-9a-fA-F]{40}$"
	PattenImageBase64 = "^data:image/[a-zA-Z0-9]+;base64,[a-zA-Z0-9/+]+=?=?$"
)

type NftScanControllerV1 struct {
	beego.Controller
}

func (receiver NftScanControllerV1) name() {

}
func (nft *NftScanControllerV1) Hello() {
	nft.Ctx.ResponseWriter.Write([]byte("Hello World!"))
}

func (nft *NftScanControllerV1) verifyInputData(data map[string]string) error {
	//regString, _ := regexp.Compile(PattenString)
	regNumber, _ := regexp.Compile(PattenNumber)

	if data["blocknumber"] != "" {
		match := regNumber.MatchString(data["blocknumber"])
		if !match {
			log.Println("verifyInputData() data[\"blocknumber\"] error.")
			return ERRINPUTINVALID
		}
	}
	if data["index"] != "" {
		match := regNumber.MatchString(data["index"])
		if !match {
			log.Println("verifyInputData() data[\"index\"] error.")
			return ERRINPUTINVALID
		}
	}
	if data["count"] != "" {
		match := regNumber.MatchString(data["count"])
		if !match {
			log.Println("verifyInputData() data[\"count\"] error.")
			return ERRINPUTINVALID
		}
	}
	return nil
}

func (nft *NftScanControllerV1) HelloWorld() {
	HelloWorld(nft)
}

func (nft *NftScanControllerV1) GetBlockTrans() {
	GetBlockTrans(nft)
}

func (nft *NftScanControllerV1) GetBlockSnfts() {
	GetBlockSnfts(nft)
}
