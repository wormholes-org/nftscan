package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

const (
	ScanBlockTime  = time.Second * 1
	AddSysRecTime  = time.Minute * 10
	ErrorsWaitTime = time.Second * 5
)

func getStartBlock() (uint64, error) {
	params := SysParams{}
	result := GetScanDB().Model(&SysParams{}).Last(&params)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		time.Sleep(ErrorsWaitTime)
		return 0, result.Error
	}
	if params.ID == 0 {
		params = SysParams{}
		params.Scannumber = 1
		params.Blocknumber = 1
		params.Scansnftnumber = 1
		err := GetScanDB().Model(&SysParams{}).Create(&params)
		if err.Error != nil {
			log.Println("SyncBlock() update params.Scannumber err=", err)
			return params.Scannumber, nil
		}
	}
	return params.Scannumber, nil
}

func SyncBlock() {
	blockS, err := getStartBlock()
	if err != nil {
		return
	}
	for blockE := GetCurrentBlockNumber(); blockS <= blockE; {
		if ScanSnft == "true" {
			err = ScanWorkerNft(blockS)
			if err != nil {
				fmt.Println("SyncProc() call ScanWorkerNft() err=", err)
				break
			}
		}
		err = ScanBlockTxs(blockS)
		if err != nil {
			fmt.Println("SyncProc() SyncBlockTxs err=", err)
			break
		}
		fmt.Println("SyncProc() scan blocknumber=", blockS)
		params := SysParams{}
		result := GetScanDB().Model(&SysParams{}).Last(&params)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			time.Sleep(ErrorsWaitTime)
			return
		}
		fmt.Println("SyncBlock() params.Id=", params.ID, "params.createat=", params.CreatedAt, "time.Now()=", time.Now())
		blockS += 1
		if params.CreatedAt.Add(AddSysRecTime).Before(time.Now()) {
			p := SysParams{}
			p.SysParamsRec = params.SysParamsRec
			p.Scannumber = blockS
			p.Blocknumber = GetCurrentBlockNumber()
			err := GetScanDB().Model(&SysParams{}).Create(&p)
			if err.Error != nil {
				log.Println("SyncBlock() create params.Scannumber err=", err)
				return
			}
			fmt.Println("SyncBlock() create new record upload block number=", blockS)
		} else {
			err := GetScanDB().Model(&SysParams{}).Where("id = ?", params.ID).Update("Scannumber", blockS)
			if err.Error != nil {
				log.Println("SyncBlock() update params.Scannumber err=", err)
				return
			}
			fmt.Println("SyncBlock() update record upload block number=", blockS)
		}
		if blockS >= blockE {
			blockE = GetCurrentBlockNumber()
		}
	}
}

func SyncChain() {
	ticker := time.NewTicker(ScanBlockTime)
	for {
		select {
		case <-ticker.C:
			SyncBlock()
		}
	}
}
