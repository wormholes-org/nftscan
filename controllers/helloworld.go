package controllers

import (
	"encoding/json"
	"fmt"
	"nftscan/models"
	"time"
)

func HelloWorld(nft *NftScanControllerV1) {
	fmt.Println("HelloWorld()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", time.Now())
	var httpResponseData HttpResponseData
	rstr, inputDatarr := models.HelloWorld()
	if inputDatarr == nil {
		httpResponseData.Code = "200"
		httpResponseData.Data = rstr
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = inputDatarr.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("HelloWorld()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now())
}
