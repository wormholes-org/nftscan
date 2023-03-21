package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nftscan/models"
	"time"
)

func GetBlockSnfts(nft *NftScanControllerV1) {
	timeBegin := time.Now()
	fmt.Println("GetBlockSnfts()>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", timeBegin)
	var httpResponseData HttpResponseData
	var data map[string]string
	bytes, _ := ioutil.ReadAll(nft.Ctx.Request.Body)
	defer nft.Ctx.Request.Body.Close()
	err := json.Unmarshal(bytes, &data)
	if err == nil {
		inputDataErr := nft.verifyInputData(data)
		if inputDataErr != nil {
			httpResponseData.Code = "500"
			httpResponseData.Msg = inputDataErr.Error()
			httpResponseData.Data = []interface{}{}
		} else {
			snfts, inputDatarr := models.GetBlockSnfts(data["blocknumber"])
			if inputDatarr == nil {
				httpResponseData.Code = "200"
				httpResponseData.Data = snfts
			} else {
				httpResponseData.Code = "500"
				httpResponseData.Msg = inputDatarr.Error()
				httpResponseData.Data = []interface{}{}
			}
		}
	} else {
		httpResponseData.Code = "500"
		httpResponseData.Msg = ERRINPUT.Error()
		httpResponseData.Data = []interface{}{}
	}
	responseData, _ := json.Marshal(httpResponseData)
	nft.Ctx.ResponseWriter.Write(responseData)
	fmt.Println("GetBlockSnfts()<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", time.Now().Sub(timeBegin))
}