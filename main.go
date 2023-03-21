package main

import (
	beego "github.com/beego/beego/v2/server/web"
	"log"
	"nftscan/models"
	_ "nftscan/routers"
)

func main() {
	log.Println("nftscan start.")
	go models.SyncChain()
	beego.Run()
}
