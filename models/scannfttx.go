package models

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	ZeroAddr = "0x0000000000000000000000000000000000000000"

	WormHolesNftCount = "1"

	WormHolesMint                      = 0
	WormHolesTransfer                  = 1
	WormHolesExchange                  = 6
	WormHolesPledge                    = 7
	WormHolesUnPledge                  = 8
	WormHolesOpenExchanger             = 11
	WormHolesExToBuyTransfer           = 14
	WormHolesBuyFromSellTransfer       = 15
	WormHolesBuyFromSellMintTransfer   = 16
	WormHolesExToBuyMintToSellTransfer = 17
	WormHolesExAuthToExBuyTransfer     = 18
	WormHolesExAuthToExMintBuyTransfer = 19
	WormHolesExSellNoAuthTransfer      = 20
	WormHolesExSellBatchAuthTransfer   = 27
	WormHolesExForceBuyingAuthTransfer = 28
)

type NftMeta struct {
	Meta    string `json:"meta"`
	TokenId string `json:"token_id"`
}

type WormholesMint struct {
	Version   string `json:"version"`
	Type      uint8  `json:"type"`
	Royalty   uint32 `json:"royalty"`
	MetaUrl   string `json:"meta_url"`
	Exchanger string `json:"exchanger"`
}

type WormholesExchange struct {
	Version    string `json:"version"`
	Type       uint8  `json:"type"`
	NftAddress string `json:"nft_address"`
}

type WormholesPledge struct {
	Version    string `json:"version"`
	Type       uint8  `json:"type"`
	NftAddress string `json:"nft_address"`
}

type WormholesOpenExchanger struct {
	Version     string `json:"version"`
	Type        uint8  `json:"type"`
	Feerate     uint32 `json:"fee_rate"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Blocknumber string `json:"block_number"`
}

type WormholesTransfer struct {
	Version    string `json:"version"`
	Type       uint8  `json:"type"`
	NftAddress string `json:"nft_address"`
}

type Buyer struct {
	Price       string `json:"price"`
	Nftaddress  string `json:"nft_address"`
	Exchanger   string `json:"exchanger"`
	Blocknumber string `json:"block_number"`
	Seller      string `json:"seller"`
	Sig         string `json:"sig"`
}

type WormholesFixTrans struct {
	Version string `json:"version"`
	Type    uint8  `json:"type"` //14
	Buyer   `json:"buyer"`
}

type Seller1 struct {
	Price       string `json:"price"`
	Nftaddress  string `json:"nft_address"`
	Exchanger   string `json:"exchanger"`
	Blocknumber string `json:"block_number"`
	Sig         string `json:"sig"`
}

type WormholesBuyFromSellTrans struct {
	Version string `json:"version"`
	Type    uint8  `json:"type"` //15
	Seller1 `json:"seller1"`
}

type Seller2 struct {
	Price         string `json:"price"`
	Royalty       string `json:"royalty"`
	Metaurl       string `json:"meta_url"`
	Exclusiveflag string `json:"exclusive_flag"`
	Exchanger     string `json:"exchanger"`
	Blocknumber   string `json:"block_number"`
	Sig           string `json:"sig"`
}

type WormholesBuyFromSellMintTrans struct {
	Version string `json:"version"`
	Type    uint8  `json:"type"`
	Seller2 `json:"seller2"`
}

type Buyer1 struct {
	Price       string `json:"price"`
	Exchanger   string `json:"exchanger"`
	Blocknumber string `json:"block_number"`
	Seller      string `json:"seller"`
	Sig         string `json:"sig"`
}

type ExchangerMintTrans struct {
	Version string  `json:"version"`
	Type    uint8   `json:"type"`
	Seller  Seller2 `json:"seller2"`
	Buyer   Buyer1  `json:"buyer"`
}

type ExchangerAuth struct {
	Exchangerowner string `json:"exchanger_owner"`
	To             string `json:"to"`
	Blocknumber    string `json:"block_number"`
	Sig            string `json:"sig"`
}

type WormholesFixTransAuth struct {
	Version       string `json:"version"`
	Type          uint8  `json:"type"`
	Buyer         `json:"buyer"`
	Seller1       `json:"seller1"`
	Exchangerauth ExchangerAuth `json:"exchanger_auth"`
}

type ExchangerAuthMintTrans struct {
	Version       string        `json:"version"`
	Type          uint8         `json:"type"`
	Seller        Seller2       `json:"seller2"`
	Buyer         Buyer1        `json:"buyer"`
	Exchangerauth ExchangerAuth `json:"exchanger_auth"`
}

type WormholesExchangeTransNoAuth struct {
	Version string `json:"version"`
	Type    uint8  `json:"type"`
	Buyer   `json:"buyer"`
	Seller1 `json:"seller1"`
}

type Buyauth struct {
	Exchanger   string `json:"exchanger"`
	Blocknumber string `json:"block_number"`
	Sig         string `json:"sig"`
}

type Sellerauth struct {
	Exchanger   string `json:"exchanger"`
	Blocknumber string `json:"block_number"`
	Sig         string `json:"sig"`
}

type WormholesBatchAuthFixTrans struct {
	Version       string `json:"version"`
	Type          uint8  `json:"type"`
	Buyauth       `json:"buyer_auth"`
	Buyer         `json:"buyer"`
	Sellerauth    `json:"seller_auth"`
	Seller1       `json:"seller1"`
	Exchangerauth ExchangerAuth `json:"exchanger_auth"`
}

type ExchangerBatchAuthTrans struct {
	Worm WormholesBatchAuthFixTrans `json:"wormholes"`
}

//type Buyer2 struct {
//	Nftaddress  string `json:"nft_address"`
//	Exchanger   string `json:"exchanger"`
//	Blocknumber string `json:"block_number"`
//	Sig         string `json:"sig"`
//}

type WormholesForceBuyingTrans struct {
	Version       string `json:"version"`
	Type          uint8  `json:"type"`
	Buyauth       `json:"buyer_auth"`
	Buyer         `json:"buyer"`
	Exchangerauth ExchangerAuth `json:"exchanger_auth"`
}

type ExchangerForceBuyingAuthTrans struct {
	Worm WormholesForceBuyingTrans `json:"wormholes"`
}

type Wormholes struct {
	Version string `json:"version"`
	Type    uint8  `json:"type"`
}

type NftTx struct {
	Operator         string
	From             string
	To               string
	Contract         string
	TokenId          string
	Value            string
	Price            string
	Ratio            string
	TxHash           string
	Ts               string
	BlockNumber      string
	TransactionIndex string
	MetaUrl          string
	NftAddr          string
	Nonce            string
	Status           bool
	TransType        int
}

func GenNftAddr(UserMintDeep *big.Int) error {
	UserMintDeep = UserMintDeep.Add(UserMintDeep, big.NewInt(1))
	return nil
}

func hashMsg(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}

func recoverAddress(msg string, sigStr string) (*common.Address, error) {
	sigData, err := hexutil.Decode(sigStr)
	if err != nil {
		log.Println("recoverAddress() err=", err)
		return nil, err
	}
	if len(sigData) != 65 {
		return nil, fmt.Errorf("signature must be 65 bytes long")
	}
	if sigData[64] != 27 && sigData[64] != 28 {
		return nil, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sigData[64] -= 27
	hash, _ := hashMsg([]byte(msg))
	rpk, err := crypto.SigToPub(hash, sigData)
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(*rpk)
	return &addr, nil
}

func copyNftTx(wnfttx *NftTx) *NftTxRec {
	var nftTx NftTxRec
	bn, _ := strconv.ParseUint(wnfttx.BlockNumber, 10, 64)
	nftTx.Operator = wnfttx.Operator
	nftTx.Fromaddr = wnfttx.From
	nftTx.Toaddr = wnfttx.To
	nftTx.Contract = wnfttx.Contract
	nftTx.Tokenid = wnfttx.TokenId
	nftTx.Value = wnfttx.Value
	nftTx.Price = wnfttx.Price
	nftTx.Ratio = wnfttx.Ratio
	nftTx.Txhash = wnfttx.TxHash
	nftTx.Ts = wnfttx.Ts
	nftTx.Blocknumber = bn
	nftTx.Transactionindex = wnfttx.TransactionIndex
	nftTx.Metaurl = wnfttx.MetaUrl
	nftTx.Nftaddr = wnfttx.NftAddr
	nftTx.Nonce = wnfttx.Nonce
	nftTx.Status = wnfttx.Status
	nftTx.Transtype = wnfttx.TransType
	return &nftTx
}

func ScanBlockTxs(blockNum uint64) error {
	spendT := time.Now()
	client, err := ethclient.Dial(EthNode)
	if err != nil {
		log.Println("ScanBlockTxs() err=", err)
		return err
	}
	defer client.Close()
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
	if err != nil {
		log.Println("ScanBlockTxs() BlockByNumber() err=", err)
		return err
	}
	fmt.Println("ScanBlockTxs() blocknumber=", blockNum)
	/*	chainId, err := client.ChainID(context.Background())
		if err != nil {
			log.Println("ScanBlockTxs() ChainID() err=", err)
			return err
		}
		fmt.Println("ScanBlockTxs() chainId=", chainId)*/
	UserMintDeep := big.NewInt(0)
	if blockNum > 1 {
		mintdeep, err := GetUserMintDeep(blockNum - 1)
		if err != nil {
			log.Println("ScanBlockTxs() GetUserMintDeep() err=", err)
			return err
		}
		UserMintDeep, ok := UserMintDeep.SetString(mintdeep, 16)
		if !ok {
			log.Println("ScanBlockTxs() UserMintDeep.SetString() errors.")
			return err
		}
		log.Println("ScanBlockTxs() UserMintDeep= ", UserMintDeep)
	}
	transT := block.Time()
	log.Println(time.Unix(int64(transT), 0))
	wnfttxs := make([]*NftTx, 0, 20)
	wminttxs := make([]*NftTx, 0, 20)
	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			continue
		}
		nonce := tx.Nonce()
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Println("ScanBlockTxs() TransactionReceipt() err=", err)
			return err
		}
		transFlag := true
		if receipt.Status != 1 {
			log.Println("ScanBlockTxs() receipt.Status != 1")
			transFlag = false
			//continue
		}
		data := tx.Data()
		if len(data) > 10 && string(data[:10]) == "wormholes:" {
			var wormholes Wormholes
			jsonErr := json.Unmarshal(data[10:], &wormholes)
			if jsonErr != nil {
				log.Println("ScanBlockTxs() wormholes type err=", err)
				continue
			}
			switch wormholes.Type {
			case WormHolesMint:
				if !transFlag {
					continue
				}
				wormMint := WormholesMint{}
				jsonErr := json.Unmarshal(data[10:], &wormMint)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() wormholes mint type err=", err)
					continue
				}
				var nftmeta NftMeta
				metabyte, _ := hex.DecodeString(wormMint.MetaUrl)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesMint
				nftx.To = strings.ToLower(tx.To().String())
				//nftx.Contract = WormHolesContract
				nftx.Contract = strings.ToLower(wormMint.Exchanger)
				//nftx.TokenId = wormMint.NftAddress
				nftx.TokenId = nftmeta.TokenId
				nftx.Value = WormHolesNftCount
				nftx.Ratio = strconv.Itoa(int(wormMint.Royalty))
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftx.MetaUrl = nftmeta.Meta
				nftx.NftAddr = common.BytesToAddress(UserMintDeep.Bytes()).String()
				wminttxs = append(wminttxs, &nftx)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("ScanBlockTxs() wormholes mint type err=", err)
				}
			case WormHolesExchange:
				if !transFlag {
					continue
				}
				wormtrans := WormholesExchange{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}
				nftx := NftTx{}
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.Status = transFlag
				nftx.TransType = WormHolesExchange
				nftx.From = strings.ToLower(tx.To().String())
				nftx.To = ZeroAddr
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				nftx.NftAddr = strings.ToLower(wormtrans.NftAddress)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesPledge:
				if !transFlag {
					continue
				}
				wormtrans := WormholesPledge{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() Unmarshal err=", err)
					continue
				}
				from := strings.ToLower(tx.To().String())
				msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesPledge() err=", err)
					msg, err = tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
					if err != nil {
						log.Println("ScanBlockTxs() WormHolesPledge() err=", err)
					} else {
						from = strings.ToLower(msg.From().String())
					}
				} else {
					from = strings.ToLower(msg.From().String())
				}
				nftx := NftTx{}
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.Status = transFlag
				nftx.TransType = WormHolesPledge
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				nftx.From = from
				nftx.To = strings.ToLower(tx.To().String())
				nftx.NftAddr = strings.ToLower(wormtrans.NftAddress)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesUnPledge:
				if !transFlag {
					continue
				}
				wormtrans := WormholesPledge{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}
				from := strings.ToLower(tx.To().String())
				msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesUnPledge() err=", err)
					msg, err = tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
					if err != nil {
						log.Println("ScanBlockTxs() WormHolesUnPledge() err=", err)
					} else {
						from = strings.ToLower(msg.From().String())
					}
				} else {
					from = strings.ToLower(msg.From().String())
				}
				nftx := NftTx{}
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.Status = transFlag
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				nftx.TransType = WormHolesUnPledge
				nftx.From = from
				nftx.To = strings.ToLower(tx.To().String())
				nftx.NftAddr = strings.ToLower(wormtrans.NftAddress)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesOpenExchanger:
				if !transFlag {
					continue
				}
				wormOpen := WormholesOpenExchanger{}
				jsonErr := json.Unmarshal(data[10:], &wormOpen)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() Unmarshal err=", err)
					continue
				}
				wormOpen.Blocknumber = strconv.FormatUint(block.NumberU64(), 10)
				exInfo, err := json.Marshal(&wormOpen)
				if err != nil {
					log.Println("ScanBlockTxs() Marshal err=", err)
					continue
				}
				if wormOpen.Name == "exchanger test." {
					exchangerInfo := string(exInfo)
					log.Println("ScanBlockTxs() find open exchanger trans.", exchangerInfo)
				}
				from := strings.ToLower(tx.To().String())
				msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesOpenExchanger() err=", err)
					msg, err = tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
					if err != nil {
						log.Println("ScanBlockTxs() WormHolesOpenExchanger() err=", err)
					} else {
						from = strings.ToLower(msg.From().String())
					}
				} else {
					from = strings.ToLower(msg.From().String())
				}
				fmt.Println("ScanBlockTxs() WormHolesOpenExchanger() from=", from)

				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesOpenExchanger
				nftx.From = from
				nftx.To = strings.ToLower(tx.To().String())
				nftx.Price = tx.Value().String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				nftx.MetaUrl = string(exInfo)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesTransfer:
				if !transFlag {
					continue
				}
				wormtrans := WormholesTransfer{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() WormHolesTransfer unmarshal type err=", err)
					continue
				}
				if wormtrans.NftAddress == "" {
					log.Println("ScanBlockTxs() WormHolesTransfer nftaddress equal null.")
					continue
				}
				from := strings.ToLower(tx.To().String())
				msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesTransfer() err=", err)
					msg, err = tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
					if err != nil {
						log.Println("ScanBlockTxs() WormHolesTransfer() err=", err)
					} else {
						from = strings.ToLower(msg.From().String())
					}
				} else {
					from = strings.ToLower(msg.From().String())
				}
				//from = msg.From()
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesTransfer
				nftx.From = from
				nftx.To = strings.ToLower(tx.To().String())
				//nftx.Contract = WormHolesContract
				//nftx.Contract = wormtrans.Exchanger
				nftx.NftAddr = wormtrans.NftAddress
				nftx.Value = WormHolesNftCount
				nftx.Price = tx.Value().String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesExToBuyTransfer:
				if !transFlag {
					continue
				}
				wormtrans := WormholesFixTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() wormholes mint type err=", err)
					continue
				}
				/*if ExchangeOwer != wormtrans.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}*/

				msg := wormtrans.Price + wormtrans.Nftaddress + wormtrans.Exchanger + wormtrans.Blocknumber + wormtrans.Seller
				toaddr, err := recoverAddress(msg, wormtrans.Sig)
				if err != nil {
					log.Println("ScanBlockTxs() recoverAddress() err=", err)
					//return err
					continue
				}
				if toaddr.String() != tx.To().String() {
					log.Println("ScanBlockTxs() PubkeyToAddress() buyer address error.")
					//return err
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExToBuyTransfer
				nftx.From = strings.ToLower(wormtrans.Seller)
				msgas, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesBuyFromSellMintTransfer() err=", err)
					msgas, err = tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
					if err != nil {
						log.Println("ScanBlockTxs() WormHolesBuyFromSellMintTransfer() err=", err)
					} else {
						nftx.From = strings.ToLower(msgas.From().String())
					}
				} else {
					nftx.From = strings.ToLower(msgas.From().String())
				}
				nftx.To = strings.ToLower(tx.To().String())
				//nftx.Contract = WormHolesContract
				nftx.Contract = strings.ToLower(wormtrans.Exchanger)
				nftx.NftAddr = wormtrans.Nftaddress
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesBuyFromSellTransfer:
				if !transFlag {
					continue
				}
				wormtrans := WormholesBuyFromSellTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() wormholes mint type err=", err)
					continue
				}
				/*if ExchangeOwer != wormtrans.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}*/

				msg := wormtrans.Price + wormtrans.Nftaddress + wormtrans.Exchanger + wormtrans.Blocknumber
				fromaddr, err := recoverAddress(msg, wormtrans.Sig)
				if err != nil {
					log.Println("ScanBlockTxs() recoverAddress() err=", err)
					//return nil, err
					continue
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesBuyFromSellTransfer
				nftx.From = strings.ToLower(fromaddr.String())
				nftx.To = strings.ToLower(tx.To().String())
				nftx.Contract = strings.ToLower(wormtrans.Exchanger)
				nftx.NftAddr = wormtrans.Nftaddress
				nftx.Value = WormHolesNftCount
				//price, _ := hexutil.DecodeUint64(wormtrans.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesBuyFromSellMintTransfer:
				if !transFlag {
					continue
				}
				wormtrans := WormholesBuyFromSellMintTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() wormholes mint type err=", err)
					continue
				}
				/*if ExchangeOwer != wormtrans.Exchanger {
					log.Println("GetBlockTxs() ExchangeOwer err=")
					continue
				}*/
				from := strings.ToLower(tx.To().String())
				msgas, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesBuyFromSellMintTransfer() err=", err)
					msgas, err = tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
					if err != nil {
						log.Println("ScanBlockTxs() WormHolesBuyFromSellMintTransfer() err=", err)
					} else {
						from = strings.ToLower(msgas.From().String())
					}
				} else {
					from = strings.ToLower(msgas.From().String())
				}
				msg := wormtrans.Price + wormtrans.Royalty + wormtrans.Metaurl + wormtrans.Exclusiveflag +
					wormtrans.Exchanger + wormtrans.Blocknumber
				toaddr, err := recoverAddress(msg, wormtrans.Sig)
				if err != nil {
					log.Println("ScanBlockTxs() recoverAddress() err=", err)
					continue
				}

				var nftmeta NftMeta
				metabyte, _ := hex.DecodeString(wormtrans.Metaurl)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftxm := NftTx{}
				nftxm.Status = transFlag
				nftxm.TransType = WormHolesBuyFromSellMintTransfer
				nftxm.To = strings.ToLower(toaddr.String())
				nftxm.Contract = strings.ToLower(wormtrans.Exchanger)
				nftxm.TokenId = nftmeta.TokenId
				nftxm.Value = WormHolesNftCount
				royalty, _ := hexutil.DecodeUint64(wormtrans.Royalty)
				nftxm.Ratio = strconv.FormatUint(royalty, 10)
				nftxm.TxHash = strings.ToLower(tx.Hash().String())
				nftxm.Ts = strconv.FormatUint(transT, 10)
				nftxm.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftxm.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftxm.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftxm.MetaUrl = nftmeta.Meta
				NftAddr := common.BytesToAddress(UserMintDeep.Bytes()).String()
				nftxm.NftAddr = NftAddr
				wminttxs = append(wminttxs, &nftxm)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("ScanBlockTxs() wormholes mint type err=", err)
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesBuyFromSellMintTransfer
				nftx.To = from
				nftx.From = strings.ToLower(toaddr.String())
				nftx.Contract = strings.ToLower(wormtrans.Exchanger)
				nftx.TokenId = nftmeta.TokenId
				nftx.NftAddr = NftAddr
				nftx.Value = WormHolesNftCount

				price, _ := hexutil.DecodeBig(wormtrans.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesExToBuyMintToSellTransfer:
				if !transFlag {
					continue
				}
				wormtrans := ExchangerMintTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}

				msg := wormtrans.Buyer.Price + wormtrans.Buyer.Exchanger + wormtrans.Buyer.Blocknumber + wormtrans.Buyer.Seller

				toaddr, err := recoverAddress(msg, wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("ScanBlockTxs() recoverAddress() err=", err)
					continue
					//return err
				}
				if toaddr.String() != tx.To().String() {
					log.Println("ScanBlockTxs() PubkeyToAddress() buyer address error.")
					//return nil, err
				}
				msg = wormtrans.Seller.Price + wormtrans.Seller.Royalty + wormtrans.Seller.Metaurl + wormtrans.Seller.Exclusiveflag +
					wormtrans.Seller.Exchanger + wormtrans.Seller.Blocknumber

				fromAddr, err := recoverAddress(msg, wormtrans.Seller.Sig)
				if err != nil {
					log.Println("ScanBlockTxs() recoverAddress() err=", err)
					//return err
					continue
				}

				var nftmeta NftMeta
				metabyte, jsonErr := hex.DecodeString(wormtrans.Seller.Metaurl)
				if jsonErr != nil {
					log.Println("GetBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftxm := NftTx{}
				nftxm.Status = transFlag
				nftxm.TransType = WormHolesExToBuyMintToSellTransfer
				nftxm.To = strings.ToLower(fromAddr.String())
				nftxm.Contract = strings.ToLower(wormtrans.Seller.Exchanger)
				nftxm.TokenId = nftmeta.TokenId
				nftxm.Value = WormHolesNftCount
				royalty, _ := strconv.Atoi(wormtrans.Seller.Royalty)
				nftxm.Ratio = strconv.Itoa(royalty)
				nftxm.TxHash = strings.ToLower(tx.Hash().String())
				nftxm.Ts = strconv.FormatUint(transT, 10)
				nftxm.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftxm.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftxm.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftxm.MetaUrl = nftmeta.Meta
				NftAddr := common.BytesToAddress(UserMintDeep.Bytes()).String()
				nftxm.NftAddr = NftAddr
				wminttxs = append(wminttxs, &nftxm)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("ScanBlockTxs() wormholes mint type err=", err)
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExToBuyMintToSellTransfer
				nftx.From = strings.ToLower(fromAddr.String())
				nftx.To = strings.ToLower(toaddr.String())
				nftx.NftAddr = NftAddr
				nftx.Value = WormHolesNftCount
				price, _ := hexutil.DecodeUint64(wormtrans.Buyer.Price)
				nftx.Price = strconv.FormatUint(price, 10)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesExAuthToExBuyTransfer:
				if !transFlag {
					continue
				}
				wormtrans := WormholesFixTransAuth{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() wormholes mint type err=", err)
					continue
				}
				msg := wormtrans.Buyer.Price + wormtrans.Buyer.Nftaddress + wormtrans.Buyer.Exchanger + wormtrans.Buyer.Blocknumber + wormtrans.Buyer.Seller

				toaddr, err := recoverAddress(msg, wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("ScanBlockTxs() recoverAddress() err=", err)
					//return err
					continue
				}
				if toaddr.String() != tx.To().String() {
					log.Println("ScanBlockTxs() PubkeyToAddress() buyer address error.")
					//return errors.New("buyer address error.")
					continue
				}
				msg = wormtrans.Seller1.Price + wormtrans.Seller1.Nftaddress + wormtrans.Seller1.Exchanger + wormtrans.Seller1.Blocknumber

				fromAddr, err := recoverAddress(msg, wormtrans.Seller1.Sig)
				if err != nil {
					log.Println("ScanBlockTxs() recoverAddress() err=", err)
					//return err
					continue
				}

				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExAuthToExBuyTransfer
				nftx.Contract = strings.ToLower(wormtrans.Exchangerauth.Exchangerowner)
				//nftx.From = strings.ToLower(wormtrans.Seller)
				nftx.From = strings.ToLower(fromAddr.String())
				nftx.To = strings.ToLower(tx.To().String())
				nftx.NftAddr = wormtrans.Buyer.Nftaddress
				nftx.Value = WormHolesNftCount
				price, _ := hexutil.DecodeBig(wormtrans.Buyer.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesExAuthToExMintBuyTransfer:
				if !transFlag {
					continue
				}
				wormtrans := ExchangerAuthMintTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}

				msg := wormtrans.Buyer.Price + wormtrans.Buyer.Exchanger + wormtrans.Buyer.Blocknumber + wormtrans.Buyer.Seller

				toaddr, err := recoverAddress(msg, wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("ScanBlockTxs() recoverAddress() err=", err)
					//return err
					continue
				}
				if toaddr.String() != tx.To().String() {
					log.Println("ScanBlockTxs() PubkeyToAddress() buyer address error.")
					//return nil, err
				}
				msg = wormtrans.Seller.Price + wormtrans.Seller.Royalty + wormtrans.Seller.Metaurl + wormtrans.Seller.Exclusiveflag +
					wormtrans.Seller.Exchanger + wormtrans.Seller.Blocknumber

				fromAddr, err := recoverAddress(msg, wormtrans.Seller.Sig)
				if err != nil {
					log.Println("ScanBlockTxs() recoverAddress() err=", err)
					//return err
					continue
				}
				//if fromAddr.String() != tx.To().String() {
				//	log.Println("GetBlockTxs() PubkeyToAddress() buyer address error.")
				//	//return nil, err
				//}
				var nftmeta NftMeta
				metabyte, jsonErr := hex.DecodeString(wormtrans.Seller.Metaurl)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() hex.DecodeString err=", err)
					continue
				}
				jsonErr = json.Unmarshal(metabyte, &nftmeta)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() NftMeta unmarshal type err=", err)
					continue
				}
				nftxm := NftTx{}
				nftxm.Status = transFlag
				nftxm.TransType = WormHolesExAuthToExMintBuyTransfer
				nftxm.To = strings.ToLower(fromAddr.String())
				nftxm.Contract = strings.ToLower(wormtrans.Exchangerauth.Exchangerowner)
				nftxm.TokenId = nftmeta.TokenId
				nftxm.Value = WormHolesNftCount
				royalty, _ := strconv.Atoi(wormtrans.Seller.Royalty)
				nftxm.Ratio = strconv.Itoa(royalty)
				nftxm.TxHash = strings.ToLower(tx.Hash().String())
				nftxm.Ts = strconv.FormatUint(transT, 10)
				nftxm.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftxm.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftxm.Nonce = strconv.FormatUint(uint64(nonce), 10)
				nftxm.MetaUrl = nftmeta.Meta
				NftAddr := common.BytesToAddress(UserMintDeep.Bytes()).String()
				nftxm.NftAddr = NftAddr
				wminttxs = append(wminttxs, &nftxm)
				err = GenNftAddr(UserMintDeep)
				if err != nil {
					log.Println("ScanBlockTxs() wormholes mint type err=", err)
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExAuthToExMintBuyTransfer
				nftx.From = strings.ToLower(fromAddr.String())
				nftx.To = strings.ToLower(toaddr.String())
				nftx.NftAddr = NftAddr
				nftx.Value = WormHolesNftCount
				nftx.Contract = strings.ToLower(wormtrans.Exchangerauth.Exchangerowner)
				//price, _ := hexutil.DecodeUint64(wormtrans.Buyer.Price)
				//nftx.Price = strconv.FormatUint(price, 10)
				price, _ := hexutil.DecodeBig(wormtrans.Buyer.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesExSellNoAuthTransfer:
				if !transFlag {
					continue
				}
				wormtrans := WormholesExchangeTransNoAuth{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() wormholes mint type err=", err)
					continue
				}
				msg := wormtrans.Buyer.Price + wormtrans.Buyer.Nftaddress + wormtrans.Buyer.Exchanger + wormtrans.Buyer.Blocknumber + wormtrans.Buyer.Seller + wormtrans.Buyer.Seller

				toaddr, err := recoverAddress(msg, wormtrans.Buyer.Sig)
				if err != nil {
					log.Println("ScanBlockTxs() recoverAddress() err=", err)
					//return err
					continue
				}
				if toaddr.String() != tx.To().String() {
					log.Println("ScanBlockTxs() PubkeyToAddress() buyer address error.")
					//return errors.New("buyer address error.")
					continue
				}
				msg = wormtrans.Seller1.Price + wormtrans.Seller1.Nftaddress + wormtrans.Seller1.Exchanger + wormtrans.Seller1.Blocknumber

				fromAddr, err := recoverAddress(msg, wormtrans.Seller1.Sig)
				if err != nil {
					log.Println("ScanBlockTxs() recoverAddress() err=", err)
					//return err
					continue
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExSellNoAuthTransfer
				nftx.Contract = strings.ToLower(wormtrans.Seller1.Exchanger)
				//nftx.From = strings.ToLower(wormtrans.Seller)
				nftx.From = strings.ToLower(fromAddr.String())
				nftx.To = strings.ToLower(tx.To().String())
				nftx.NftAddr = wormtrans.Buyer.Nftaddress
				nftx.Value = WormHolesNftCount
				price, _ := hexutil.DecodeBig(wormtrans.Buyer.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesExSellBatchAuthTransfer:
				if !transFlag {
					continue
				}
				wormtrans := WormholesBatchAuthFixTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() wormholes mint type err=", err)
					continue
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExSellBatchAuthTransfer
				nftx.Contract = strings.ToLower(wormtrans.Seller1.Exchanger)
				//nftx.From = strings.ToLower(wormtrans.Seller)
				nftx.From = strings.ToLower(wormtrans.Buyer.Seller)
				nftx.To = strings.ToLower(tx.To().String())
				nftx.NftAddr = wormtrans.Buyer.Nftaddress
				nftx.Value = WormHolesNftCount
				price, _ := hexutil.DecodeBig(wormtrans.Buyer.Price)
				nftx.Price = price.String()
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			case WormHolesExForceBuyingAuthTransfer:
				if !transFlag {
					continue
				}
				wormtrans := WormholesForceBuyingTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() wormholes force buy type err=", err)
					continue
				}
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExForceBuyingAuthTransfer
				nftx.Contract = strings.ToLower(wormtrans.Buyauth.Exchanger)
				if wormtrans.Buyer.Seller != "" {
					nftx.From = strings.ToLower(wormtrans.Buyer.Seller)
				} else {
					nftx.From = strings.ToLower(wormtrans.Buyer.Exchanger)
				}
				nftx.To = strings.ToLower(tx.To().String())
				nftx.NftAddr = wormtrans.Buyer.Nftaddress
				nftx.Value = WormHolesNftCount
				price, err := GetForcedSaleAmountByAddress(common.HexToAddress(nftx.NftAddr + "0"))
				if err != nil {
					log.Println("ScanBlockTxs() GetForcedSaleAmountByAddress err=", err)
					continue
				}
				buyAddr := common.HexToAddress(nftx.To)
				nftList, err := GetForcedSaleSNFTAddresses(nftx.NftAddr, buyAddr, new(big.Int).SetUint64(block.NumberU64()-1))
				if err != nil {
					log.Println("ScanBlockTxs() GetForcedSaleSNFTAddresses err=", err)
					return err
				}
				nftCnt := len(nftList)
				fmt.Println("ScanBlockTxs() nftList=", nftCnt)
				price = price.Mul(price, big.NewInt(int64(nftCnt)))
				nftx.Price = price.Text(10)
				nftx.TxHash = strings.ToLower(tx.Hash().String())
				nftx.Ts = strconv.FormatUint(transT, 10)
				nftx.BlockNumber = strconv.FormatUint(block.NumberU64(), 10)
				nftx.TransactionIndex = strconv.FormatUint(uint64(receipt.TransactionIndex), 10)
				nftx.Nonce = strconv.FormatUint(nonce, 10)
				wnfttxs = append(wnfttxs, &nftx)
			}
		}
	}

	var wnfts []NftTxs
	if len(wnfttxs) != 0 {
		fmt.Println("ScanBlockTxs() create tx count", len(wnfttxs), " block num=", blockNum)
		for _, wnfttx := range wnfttxs {
			var wnft NftTxs
			wnft = NftTxs{}
			wnft.NftTxRec = *copyNftTx(wnfttx)
			wnfts = append(wnfts, wnft)
		}
	}
	if len(wminttxs) != 0 {
		fmt.Println("ScanBlockTxs() create mint count", len(wminttxs), " block num=", blockNum)
		for _, wnfttx := range wminttxs {
			var wnft NftTxs
			wnft = NftTxs{}
			wnft.NftTxRec = *copyNftTx(wnfttx)
			wnfts = append(wnfts, wnft)
		}
	}
	if len(wnfts) != 0 {
		spendT := time.Now()
		wnft := NftTxs{}
		result := GetScanDB().Model(&NftTxs{}).Where("blocknumber = ?", blockNum).First(&wnft)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}
		if result.Error == gorm.ErrRecordNotFound {
			terr := GetScanDB().Transaction(func(tx *gorm.DB) error {
				for _, txs := range wnfts {
					err := tx.Model(&NftTxs{}).Where("Txhash = ?", txs.Txhash).Delete(&NftTxs{})
					if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
						log.Println("ScanWorkerNft() upload err=", err)
						return err.Error
					}
				}
				err := tx.Model(&NftTxs{}).Create(&wnfts)
				if err.Error != nil {
					log.Println("ScanWorkerNft() upload err=", err)
					return err.Error
				}
				return nil
			})
			if terr != nil {
				log.Println("ScanWorkerNft() upload err=", err)
				return terr
			}
		} else {
			log.Println("ScanWorkerNft() upload exit block=", blockNum)
		}
		fmt.Println("ScanBlockTxs() create txs spend time=", time.Now().Sub(spendT), " block number=", blockNum, " tx count", len(wnfttxs))
	}
	fmt.Println("ScanBlockTxs() end spend time=", time.Now().Sub(spendT), " block number=", blockNum, " tx count=", len(wnfts))
	return nil
}
