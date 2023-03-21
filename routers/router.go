package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"log"
	"nftscan/controllers"
	"nftscan/models"
)

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	err := models.InitDb(models.SqlSvr, models.DbName)
	if err != nil {
		log.Panicln("Init dbase error.")
	}
	registRouterV1()
}

func registRouterV1() {
	//hello world
	beego.Router("/v1/helloWorld", &controllers.NftScanControllerV1{}, "get:HelloWorld")
	beego.Router("/v1/getBlockTrans", &controllers.NftScanControllerV1{}, "post:GetBlockTrans")
	beego.Router("/v1/getBlockSnfts", &controllers.NftScanControllerV1{}, "post:GetBlockSnfts")
}
