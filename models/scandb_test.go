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
	"github.com/fsnotify/fsnotify"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

//const sqlsvrLcT = "admin:user123456@tcp(192.168.1.235:3306)/"
const sqlsvrLcT = "admin:user123456@tcp(192.168.56.122:3306)/"

const dbNameT = "scandb"

const localtimeT = "?parseTime=true&loc=Local"

//const localtimeT = "?charset=utf8mb4&parseTime=True&loc=Local"

const sqldsnT = sqlsvrLcT + dbNameT + localtimeT

func TestCreateNfts(t *testing.T) {
	_, err := NewNftDb(sqldsnT)
	if err != nil {
		t.Fatalf("NewNftDb")
	}
	nd := GetScanDB()
	nftTx := NftTxs{}
	nftTx.Status = false
	nftTx.Value = "10000"
	nftTx.Blocknumber = 20000
	nftTx.Ts = "6666666"
	nftTx.Nonce = "99999999"
	dberr := nd.Model(&NftTxs{}).Create(&nftTx)
	if dberr.Error != nil {
		fmt.Println("BuyResult() create trans record err=", dberr.Error)
		return
	}
	nftTx = NftTxs{}
	dberr = nd.Model(&NftTxs{}).Last(&nftTx)
	if dberr.Error != nil {
		fmt.Println("BuyResult() create trans record err=", dberr.Error)
		return
	}
}

func TestSnfts(t *testing.T) {
	/*_, err := NewNftDb(sqldsnT)
	if err != nil {
		t.Fatalf("NewNftDb")
	}*/
	nd := GetScanDB(sqldsnT)
	snftTx := Snfts{}
	snftTx.Blocknumber = 10000
	snftTx.Royalty = 100.00011
	dberr := nd.Model(&Snfts{}).Create(&snftTx)
	if dberr.Error != nil {
		fmt.Println("BuyResult() create trans record err=", dberr.Error)
		return
	}
	snftTx = Snfts{}
	dberr = nd.Model(&Snfts{}).Last(&snftTx)
	if dberr.Error != nil {
		fmt.Println("BuyResult() create trans record err=", dberr.Error)
		return
	}
}

func TestScanSnft(t *testing.T) {
	EthNode = "http://43.129.181.130:8561"
	//EthNode = "http://150.109.149.220:8560"
	NftIpfsServer = "192.168.1.235:5001"
	_ = GetScanDB(sqldsnT)
	for i := 1; i < 100; i++ {
		err := ScanWorkerNft(uint64(i) + 1058)
		if err != nil {
			t.Fatal("err= ", err)
		}
	}
}

func TestInitDB(t *testing.T) {
	err := InitDb(sqlsvrLcT, dbNameT)
	if err != nil {
		log.Panicln("Init dbase error.")
	}
}

func TestRecoverAddr(t *testing.T) {
	msg := "0x53444835ec580000" + "0x0000000000000000000000000000000000000014" +
		"0xdf129ff495cb69b87ba3c65ea4bfb6398b479d56" + "0x14507" + "0x8edb587b9aedd348a76d36b341fffbcf63f2a5a9"

	toaddr, err := recoverAddress(msg, "0x23eb25b582d3cf128bb364dacb5301f0e6f979f21c88b83be7bfc4afa19dbea932b7b7ad5f93dcde1392c8d22106abba8f278275ad930348703c3a4f7a4be5cd1b")
	if err != nil {
		log.Println("ScanBlockTxs() recoverAddress() err=", err)
		//return err
	}
	fmt.Println(toaddr.String())
}

func ScanBlockTxModify(blockNum uint64) error {
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
				fromAddr, err := client.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				if err != nil {
					fmt.Println(fromAddr)
				}
				from := strings.ToLower(tx.To().String())
				msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesOpenExchanger() err=", err)

				} else {
					from = strings.ToLower(msg.From().String())
				}
				msg, err = tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
				if err != nil {
					log.Println("ScanBlockTxs() WormHolesOpenExchanger() err=", err)
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
				wormtrans := ExchangerMintTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}

				msg := wormtrans.Buyer.Price + wormtrans.Buyer.Exchanger + wormtrans.Buyer.Blocknumber

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
				nftx := NftTx{}
				nftx.Status = transFlag
				nftx.TransType = WormHolesExAuthToExBuyTransfer
				nftx.Contract = strings.ToLower(wormtrans.Exchangerauth.Exchangerowner)
				nftx.From = strings.ToLower(wormtrans.Seller)
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
				wormtrans := ExchangerAuthMintTrans{}
				jsonErr := json.Unmarshal(data[10:], &wormtrans)
				if jsonErr != nil {
					log.Println("ScanBlockTxs() WormHolesExMintTransfer mint type err=", err)
					continue
				}

				msg := wormtrans.Buyer.Price + wormtrans.Buyer.Exchanger + wormtrans.Buyer.Blocknumber

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
				//nftx.From = strings.ToLower(wormtrans.Seller)
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
					continue
				}
				buyAddr := common.HexToAddress(nftx.To)
				nftList, err := GetForcedSaleSNFTAddresses(nftx.NftAddr, buyAddr, new(big.Int).SetUint64(block.NumberU64()-1))
				nftCnt := len(nftList)
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
	return nil
}

func TestGetForcedSaleAmountByAddress(t *testing.T) {
	price, err := GetForcedSaleAmountByAddress(common.HexToAddress("0x800000000000000000000000000000000001865e"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(price)
}

func TestGetForcedSaleSNFTAddresses(t *testing.T) {
	EthNode = "http://192.168.4.240:8560"
	buyAddr := common.HexToAddress("0x85d3fda364564c365870233e5ad6b611f2227846")
	nftList, err := GetForcedSaleSNFTAddresses("0x800000000000000000000000000000000000000", buyAddr, big.NewInt(5))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(nftList))
}

func TestScanWnftTrans(t *testing.T) {
	//EthNode = "http://43.129.181.130:8561"
	//EthNode = "http://150.109.149.220:8560"
	EthNode = "http://192.168.4.240:8560"
	_ = GetScanDB(sqldsnT)
	for i := 0; i < 1000; i++ {
		//err := ScanBlockTxs(uint64(i) + 1141)
		//err := ScanBlockTxs(uint64(i) + 34976)
		//err := ScanBlockTxs(uint64(i) + 22191)
		//err := ScanBlockTxs(uint64(i) + 4256)
		//err := ScanBlockTxs(uint64(i) + 21255)
		//err := ScanBlockTxs(uint64(i) + 21239) //jiaoyisuoyinwenti
		//err := ScanBlockTxs(uint64(i) + 14871) //jiaoyisuoyinwenti
		err := ScanBlockTxModify(uint64(i) + 36) //jiaoyisuoyinwenti
		//err := ScanBlockTxs( uint64(i) + 68449)
		//err := ScanBlockTxs( uint64(i) + 69880)
		if err != nil {
			t.Fatal("err= ", err)
		}
	}
}

func TestRecover(t *testing.T) {
	msg := "0x98a7d9b8314c0000" + "0x0000000000000000000000000000000000000001" +
		"0x01842a2cf56400a245a56955dc407c2c4137321e" + "0x604127"
	msg = "0x56bc75e2d63100000" + "0x0000000000000000000000000000000000000008" +
		"0x0dd939bd2f55e052595d556d02fb534404ac1234" + "0x607b7d"
	fromAddr, err := recoverAddress(msg, "0x8b1ae1954ff4c9929af47682bc4a5fad2c8d207888bacff560829e35c68cde58730b40ccc5f9d13f9746c1dff30db66cbaf33e652abcc07dad8ee96328129f691b")
	fmt.Println(fromAddr.String(), err)
}

func TestQueryCatch(t *testing.T) {
	for i := 0; i < 10; i++ {
		GetQueryCatch().SetByHash(strconv.Itoa(i), i)
		var m int
		err := GetQueryCatch().GetByHash(strconv.Itoa(i), &m)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(m)
	}
}

func TestScanTx(t *testing.T) {
	ScanSnft = "false"
	//EthNode = "http://150.109.149.220:8560"
	EthNode = "http://192.168.4.240:8560"
	_ = GetScanDB(sqldsnT)
	SyncBlock()
}

func TestMonitorFileName(t *testing.T) {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()
	err = watch.Add("D:\\temp\\m.txt")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("create file : ", ev.Name)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("write file : ", ev.Name)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err)
					return
				}
			}
		}
	}()

	//循环
	select {}
}

func TestGeneratePrv(t *testing.T) {
	keys := []string{}
	for i := 0; i < 6; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			fmt.Println("failed GenerateKey with.", err)
			return
		}
		prvKey := strings.ToLower(hexutil.Encode(crypto.FromECDSA(key)))
		keys = append(keys, prvKey)
	}
	data, err := json.Marshal(keys)
	if err != nil {
		t.Fatal("marshal err=", err)
	}
	f, err := os.OpenFile("d:/temp/keys.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Fatal("openfile() err=", err)
	}
	defer f.Close()
	n, err := f.Write([]byte(data))
	if err != nil {
		t.Fatal("write file err=", err)
	}
	fmt.Println("write n=", n)
	rf, err := os.OpenFile("d:/temp/keys.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Fatal("openfile() err=", err)
	}
	defer rf.Close()
	if finfo, err := rf.Stat(); err != nil {
		t.Fatal("f.stat() err=", err)
	} else {
		data = make([]byte, finfo.Size())
	}
	n, err = rf.Read(data)
	if err != nil {
		t.Fatal("read file err=", err)
	}
	prvKeys := map[string]string{}
	prvKeys["0xEbE809C70406Fc07e70Bc2c590bC7Dd9Ba9272Ac"] = "0x58291806b354fda6c7a1ef171e2f1dd3ce9bc187677312d81e03df3ee970b308"
	prvKeys["0xEbE809C70406Fc07e70Bc2c590bC7Dd9Ba9272A1"] = "0x58291806b354fda6c7a1ef171e2f1dd3ce9bc187677312d81e03df3ee970b309"
	data, err = json.Marshal(prvKeys)
	err = json.Unmarshal(data, &prvKeys)
	if err != nil {
		t.Fatal("unmarshal() err=", err)
	}
	for i, key := range prvKeys {
		fmt.Println(i, key)
	}
}
