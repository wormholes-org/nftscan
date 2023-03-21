package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
	"log"
	"sync"
	"time"
)

type CatchInfo struct {
	Date    interface{}
	TimeTag time.Time
}

type NftQueryCatch struct {
	Mux       sync.Mutex
	CatchData map[string]*CatchInfo
}

func (n *NftQueryCatch) GetByHash(hash string, catchData interface{}) error {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	fmt.Println("NftQueryCatch() GetByHash() n.CatchData catch len=", len(n.CatchData))
	if len(n.CatchData) == 0 {
		n.CatchData = make(map[string]*CatchInfo)
	}
	if catchinfo, ok := n.CatchData[hash]; ok {
		fmt.Println("NftQueryCatch() GetByHash() hash=", hash)
		err := json.Unmarshal(catchinfo.Date.([]byte), catchData)
		if err == nil {
			return nil
		}
	}
	return errors.New("no catch data.")
}

func (n *NftQueryCatch) SetByHash(hash string, data interface{}) {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	if len(n.CatchData) == 0 {
		log.Println("NftQueryCatch() SetByHash() len ==0 ")
		n.CatchData = make(map[string]*CatchInfo)
	}
	deleteKeys := make([]string, 0)
	if len(n.CatchData) > 1000 {
		for key, info := range n.CatchData {
			if info.TimeTag.Add(time.Second * 60).Before(time.Now()) {
				//if info.TimeTag.Add(time.Minute * 60).Before(time.Now()) {
				deleteKeys = append(deleteKeys, key)
			}
		}
	}
	if len(deleteKeys) != 0 {
		for _, key := range deleteKeys {
			delete(n.CatchData, key)
		}
		if len(n.CatchData) == 0 {
			log.Println("NftQueryCatch() SetByHash() len ==0 ")
			n.CatchData = make(map[string]*CatchInfo)
		}
	}
	catch := CatchInfo{}
	mdata, _ := json.Marshal(data)
	catch.Date = mdata
	catch.TimeTag = time.Now()
	n.CatchData[hash] = &catch
	fmt.Println("NftQueryCatch() SetByHash()", "len=", len(n.CatchData), " hash=", hash)
}

func (n *NftQueryCatch) ClearCatch( /*flag NftFlushType*/ ) {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	fmt.Println("ClearCatch() clear catch  hash")
	n.CatchData = make(map[string]*CatchInfo)
}

func (n NftQueryCatch) NftCatchHash(data string) string {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)
	return hexutil.Encode(hash)
}

var NftCatch NftQueryCatch

func GetQueryCatch() *NftQueryCatch {
	return &NftCatch
}
