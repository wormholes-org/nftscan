package models

import (
	"encoding/json"
	"errors"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	SnftOffset                = 41
	snftCollectionOffset      = 40
	SnftStageOffset           = 39
	SnftCollectionsStageIndex = 30
	WaitIpfsFailTime          = time.Second * 1
)

type SnftInfo struct {
	CreatorAddr          string  `json:"creator_addr"`
	Ownaddr              string  `json:"ownaddr"`
	Contract             string  `json:"nft_contract_addr"`
	Nftaddr              string  `json:"nft_address"`
	Snft                 string  `json:"snft"`
	Snftcollection       string  `json:"snftcollection"`
	Snftstage            string  `json:"snftstage"`
	Name                 string  `json:"name"`
	Desc                 string  `json:"desc"`
	Meta                 string  `json:"meta"`
	Category             string  `json:"category"`
	Royalty              float64 `json:"royalty"`
	SourceUrl            string  `json:"source_url"`
	Md5                  string  `json:"md5"`
	CollectionsName      string  `json:"collections_name"`
	CollectionsCreator   string  `json:"collections_creator"`
	CollectionsExchanger string  `json:"collections_exchanger"`
	CollectionsCategory  string  `json:"collections_category"`
	CollectionsImgUrl    string  `json:"collections_img_url"`
	CollectionsDesc      string  `json:"collections_desc"`
}

type SnftInfoData struct {
	SnftInfo *SnftInfo
	TimeTag  time.Time
}

type IpfsCatch struct {
	Mux      sync.Mutex
	SnftInfo map[string]*SnftInfoData
}

func (n *IpfsCatch) GetByHash(hash string) *SnftInfo {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	fmt.Println("IpfsCatch-GetByHash() GetByHash n.NftInfo catch len=", len(n.SnftInfo))
	if len(n.SnftInfo) == 0 {
		n.SnftInfo = make(map[string]*SnftInfoData)
	}
	if nftinfo := n.SnftInfo[hash]; nftinfo != nil {
		fmt.Println("IpfsCatch-GetByHash() NftFilterCatch hash=", hash)
		s := *nftinfo.SnftInfo
		return &s
	}
	return nil
}

func (n *IpfsCatch) SetByHash(hash string, snftinfo *SnftInfo) *SnftInfo {
	n.Mux.Lock()
	defer n.Mux.Unlock()
	if len(n.SnftInfo) == 0 {
		fmt.Println("IpfsCatch-SetByHash() NftFilterCatch len ==0 ")
		n.SnftInfo = make(map[string]*SnftInfoData)
	}
	s := *snftinfo
	n.SnftInfo[hash] = &SnftInfoData{&s, time.Now().Add(time.Minute * 30)}
	fmt.Println("IpfsCatch-SetByHash() NftFilterCatch", "len=", len(n.SnftInfo), " hash=", hash)
	for s, info := range n.SnftInfo {
		if info.TimeTag.Before(time.Now()) {
			delete(n.SnftInfo, s)
		}
	}
	return nil
}

var ScanIpfsCatch IpfsCatch

func GetSnftInfoFromIPFSWithShell(hash string) (*SnftInfo, error) {
	url := NftIpfsServer
	s := shell.NewShell(url)
	s.SetTimeout(100 * time.Second)
	rc, err := s.Cat(hash)
	if err != nil {
		log.Println("GetSnftInfoFromIPFSWithShell() err=", err)
		return nil, err
	}
	var snft SnftInfo
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Println("GetSnftInfoFromIPFSWithShell() ReadAll() err=", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(b), &snft)
	if err != nil {
		log.Println("GetSnftInfoFromIPFSWithShell() Unmarshal, err=", err)
		return nil, err
	}
	return &snft, nil
}

func copySnftinfo(info *SnftInfo, bn string) *SnftRec {
	var snft SnftRec
	bnu, _ := strconv.ParseUint(bn, 10, 64)
	snft.Blocknumber = bnu
	snft.Creatoraddr = info.CreatorAddr
	snft.Ownaddr = info.Ownaddr
	snft.Contract = info.Contract
	snft.Nftaddr = info.Nftaddr
	snft.Snft = info.Snft
	snft.Snftcollection = info.Snftcollection
	snft.Snftstage = info.Snftstage
	snft.Name = info.Name
	snft.Desc = info.Desc
	snft.Meta = info.Meta
	snft.Category = info.Category
	snft.Royalty = info.Royalty
	snft.Sourceurl = info.SourceUrl
	snft.Md5 = info.Md5
	snft.Collectionsname = info.CollectionsName
	snft.Collectionscreator = info.CollectionsCreator
	snft.Collectionsexchanger = info.CollectionsExchanger
	snft.Collectionscategory = info.CollectionsCategory
	snft.Collectionsimgurl = info.CollectionsImgUrl
	snft.Collectionsdesc = info.CollectionsDesc
	return &snft
}

func ScanWorkerNft(blockS uint64) error {
	snftAddr, err := GetSnftAddressList(big.NewInt(0).SetUint64(blockS), true)
	if err != nil {
		log.Println("ScanWorkerNft() GetSnftAddressList err =", err, "blocks", blockS)
		return err
	}
	//snftInfos := make([]Snfts, len(snftAddr))
	snftInfos := []Snfts{}
	if len(snftAddr) > 0 {
		for _, address := range snftAddr {
			if address.NftAddress.String() == ZeroAddr {
				continue
			}
			spendT := time.Now()
			accountInfo, err := GetAccountInfo(address.NftAddress, big.NewInt(0).SetUint64(blockS))
			if err != nil {
				log.Println("ScanWorkerNft() GetAccountInfo err =", err, "NftAddress= ", address.NftAddress, "blocks", blockS)
				return err
			}
			fmt.Println("ScanWorkerNft() GetAccountInfo spend time=", time.Now().Sub(spendT))
			fmt.Println("ScanWorkerNft() MetaUrl=", accountInfo.MetaURL, "blockS=", blockS)
			index := strings.Index(accountInfo.MetaURL, "/ipfs/")
			if index == -1 {
				log.Printf("ScanWorkerNft() Index ipfs error.\n")
				continue
				return errors.New("ScanWorkerNft(): MetaUrl error.")
			}
			index = strings.LastIndex(accountInfo.MetaURL, "/")
			if index == -1 {
				log.Printf("ScanWorkerNft() LastIndex error.\n")
				continue
				return errors.New("ScanWorkerNft(): MetaUrl error.")
			}
			/*if accountInfo.MetaURL[:index] == "/ipfs/QmYgBEB9CEx356zqJaDd4yjvY92qE276Gh1y2baWeDY3By" ||
				accountInfo.MetaURL[:index] == "/ipfs/QmaiReZpUeWcSRvhWhHwQ4PN2NbggYdZt7hKFAoM8kTVF7" {
				continue
			}*/
			var metaUrl, metaHash string
			if accountInfo.MetaURL[:index] == "/ipfs/QmeCPcX3rYguWqJYDmJ6D4qTQqd5asr8gYpwRcgw44WsS7" ||
				accountInfo.MetaURL[:index] == "/ipfs/QmYgBEB9CEx356zqJaDd4yjvY92qE276Gh1y2baWeDY3By" {
				//metaUrl = "/ipfs/QmVyVJTMQVbHRz8dr8RHrW4c1pgnspcM3Ee1pj9vae2oo8" //1.237
				metaUrl = "/ipfs/QmNbNvhW1StGPQaXhXMQcfT6W7HqEXDY6MfZijuRLf7Roa" //云服务器
				//metaUrl = "/ipfs/QmWpDcyU287P3bgw74nmUmWGDcaRYGud51y8xxQkiK5zDR" //云服务器
				//metaHash = metaUrl + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-4:len(accountInfo.MetaURL)-2])
				metaHash = metaUrl + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-3:len(accountInfo.MetaURL)-1])
			} else {
				//metaHash = accountInfo.MetaURL[:index] + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-4:len(accountInfo.MetaURL)-2])
				metaHash = accountInfo.MetaURL[:index] + "/" + strings.ToLower(accountInfo.MetaURL[len(accountInfo.MetaURL)-3:len(accountInfo.MetaURL)-1])
			}
			fmt.Println("ScanWorkerNft() metaHash=", metaHash)
			var snftinfo *SnftInfo
			spendT = time.Now()
			if snftinfo = ScanIpfsCatch.GetByHash(metaHash); snftinfo == nil {
				retry := 0
				for {
					snftinfo, err = GetSnftInfoFromIPFSWithShell(metaHash)
					if err != nil {
						log.Println("ScanWorkerNft() GetSnftInfoFromIPFS count=", retry, " err =", err, "ipfs hash=", metaHash)
						errflag := strings.Index(err.Error(), "context deadline exceeded")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
						errflag = strings.Index(err.Error(), "connection refused")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
						errflag = strings.Index(err.Error(), "502 Bad Gateway")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
						errflag = strings.Index(err.Error(), "403 Forbidden")
						if errflag != -1 {
							time.Sleep(WaitIpfsFailTime)
							continue
						}
					}
					break
				}
				if err != nil {
					continue
				}
				ScanIpfsCatch.SetByHash(metaHash, snftinfo)
			}
			fmt.Println("ScanWorkerNft() GetSnftInfoFromIPFS spend time=", time.Now().Sub(spendT))
			snftinfo.Ownaddr = strings.ToLower(accountInfo.Owner.String())
			snftinfo.Contract = ""
			snftinfo.Nftaddr = strings.ToLower(address.NftAddress.String())
			snftinfo.Meta = accountInfo.MetaURL
			snftinfo.Snft = snftinfo.Nftaddr[:SnftOffset]
			snftinfo.Snftcollection = snftinfo.Nftaddr[:snftCollectionOffset]
			snftinfo.Snftstage = snftinfo.Nftaddr[:SnftStageOffset]
			b, _ := big.NewInt(0).SetString(snftinfo.Nftaddr[SnftCollectionsStageIndex:SnftStageOffset], 16)
			snftinfo.CollectionsName = b.String() + "-" + snftinfo.CollectionsName
			bn := strconv.FormatUint(blockS, 10)
			snft := Snfts{}
			snft.SnftRec = *copySnftinfo(snftinfo, bn)
			//snftInfos[i].SnftRec = *copySnftinfo(snftinfo, bn)
			snftInfos = append(snftInfos, snft)
		}
	}
	if len(snftInfos) != 0 {
		spendT := time.Now()
		var snft Snfts
		result := GetScanDB().Model(&Snfts{}).Where("blocknumber = ?", blockS).First(&snft)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}
		if result.Error == gorm.ErrRecordNotFound {
			err := GetScanDB().Model(&Snfts{}).Create(&snftInfos)
			if err.Error != nil {
				log.Println("ScanWorkerNft() upload err=", err, " blocknumber=", blockS)
				return err.Error
			}
		}

		fmt.Println("ScanWorkerNft() save snft count=", len(snftInfos), "spend time=", time.Now().Sub(spendT))
	}
	return err
}
