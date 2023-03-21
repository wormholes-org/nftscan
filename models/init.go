package models

import (
	"flag"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"log"
	"os"
	"testing"
)

var (
	SqlSvr        string
	Sqldsndb      string
	DbName        string
	EthNode       string
	NftIpfsServer string
	ScanSnft      string
)

const localtime = "?parseTime=true&loc=Local"
const version = "0.8.7"

func DisplayVersion() {
	v := flag.Bool("version", false, "display version")
	testing.Init()
	flag.Parse()
	if *v {
		fmt.Println("version =", version)
		os.Exit(0)
	}
}

func init() {
	DisplayVersion()
	dbname, err := beego.AppConfig.String("dbname")
	if err != nil {
		log.Println("models init params dbname err=", err)
	}
	DbName = dbname
	fmt.Println("models init param DbName=", DbName)
	dbUserName, _ := beego.AppConfig.String("dbusername")
	if err != nil {
		log.Println("models init params dbusername err=", err)
	}
	fmt.Println("models init param dbUserName=", dbUserName)
	dbUserPassword, _ := beego.AppConfig.String("dbuserpassword")
	if err != nil {
		log.Println("models init params dbuserpassword err=", err)
	}
	dbServerIP, _ := beego.AppConfig.String("dbserverip")
	fmt.Println("models init param dbServerIP=", dbServerIP)
	if err != nil {
		log.Println("models init params dbserverip err=", err)
	}
	dbServerPort, _ := beego.AppConfig.String("dbserverport")
	if err != nil {
		log.Println("models init params dbserverport err=", err)
	}
	fmt.Println("models init param dbServerPort=", dbServerPort)
	SqlSvr = dbUserName + ":" + dbUserPassword + "@tcp(" + dbServerIP + ":" + dbServerPort + ")/"
	//fmt.Println("SqlSvr=", SqlSvr)
	Sqldsndb = SqlSvr + dbname + localtime
	EthNode, _ = beego.AppConfig.String("EthNode")
	fmt.Println("models init param EthersNode=", EthNode)
	NftIpfsServer, _ = beego.AppConfig.String("nftIpfsServer")
	scansnft, _ := beego.AppConfig.String("ScanSnft")
	ScanSnft = "true"
	if scansnft == "false" {
		ScanSnft = "false"
	}
}
