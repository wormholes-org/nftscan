package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/crypto/sha3"
	"log"
	"math"
	"math/big"
	"time"
)

const (
	ReDialDelyTime = 5
)

type BeneficiaryAddress struct {
	Address    common.Address
	NftAddress common.Address
}
type BeneficiaryAddressList []*BeneficiaryAddress

type Account struct {
	Nonce   uint64
	Balance *big.Int
	// *** modify to support nft transaction 20211220 begin ***
	//NFTCount uint64		// number of nft who account have
	// *** modify to support nft transaction 20211220 end ***
	Root           common.Hash // merkle root of the storage trie
	CodeHash       []byte
	PledgedBalance *big.Int
	// *** modify to support nft transaction 20211215 ***
	//Owner common.Address
	// whether the account has a NFT exchanger
	ExchangerFlag bool
	BlockNumber   *big.Int
	// The ratio that exchanger get.
	FeeRate       uint32
	ExchangerName string
	ExchangerURL  string
	// ApproveAddress have the right to handle all nfts of the account
	ApproveAddressList []common.Address
	// NFTBalance is the nft number that the account have
	NFTBalance uint64
	AccountNFT
}
type AccountNFT struct {
	//Account
	Name                  string
	Symbol                string
	Price                 *big.Int
	Direction             uint8 // 0:未交易,1:买入,2:卖出
	Owner                 common.Address
	NFTApproveAddressList common.Address
	//Auctions map[string][]common.Address
	// MergeLevel is the level of NFT merged
	MergeLevel uint8

	Creator   common.Address
	Royalty   uint32
	Exchanger common.Address
	MetaURL   string
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

func TextAndHash(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}

func GetEthAddr(msg string, sigStr string) (common.Address, error) {
	sigData, _ := hexutil.Decode(sigStr)
	if len(sigData) != 65 {
		return common.Address{}, fmt.Errorf("552,signature must be 65 bytes long")
	}
	if sigData[64] != 27 && sigData[64] != 28 {
		return common.Address{}, fmt.Errorf("552,invalid Ethereum signature (V is not 27 or 28)")
	}
	sigData[64] -= 27
	hash, _ := TextAndHash([]byte(msg))
	fmt.Println("sigdebug hash=", hexutil.Encode(hash))
	rpk, err := crypto.SigToPub(hash, sigData)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*rpk), nil
}

func GetSnftAddressList(blockNumber *big.Int, fulltx bool) ([]*BeneficiaryAddress, error) {
	client, err := rpc.Dial(EthNode)
	if err != nil {
		fmt.Println("GetSnftAddressList() err=", err)
		return nil, err
	}
	var result BeneficiaryAddressList
	err = client.CallContext(context.Background(), &result, "eth_getBlockBeneficiaryAddressByNumber", toBlockNumArg(blockNumber), fulltx)
	if err != nil {
		fmt.Println("GetSnftAddressList() err=", err)
		return nil, err
	}
	return result, err
}

func GetAccountInfo(nftaddr common.Address, blockNumber *big.Int) (*Account, error) {
	client, err := rpc.Dial(EthNode)
	if err != nil {
		log.Println("GetAccountInfo() err=", err)
		return nil, err
	}
	var result Account
	err = client.CallContext(context.Background(), &result, "eth_getAccountInfo", nftaddr, toBlockNumArg(blockNumber))
	if err != nil {
		log.Println("GetAccountInfo() err=", err)
		return nil, err
	}
	return &result, err
}

func GetForcedSaleSNFTAddresses(nftaddr string, buyAddr common.Address, blockNumber *big.Int) ([]common.Address, error) {
	client, err := rpc.Dial(EthNode)
	if err != nil {
		log.Println("GetForcedSaleSNFTAddresses() err=", err)
		return nil, err
	}
	var result []common.Address
	err = client.CallContext(context.Background(), &result, "eth_getForcedSaleSNFTAddresses", nftaddr, buyAddr, toBlockNumArg(blockNumber))
	if err != nil {
		log.Println("GetForcedSaleSNFTAddresses() err=", err)
		return nil, err
	}
	return result, err
}

func GetForcedSaleAmount(nftaddr common.Address) (string, error) {
	client, err := rpc.Dial(EthNode)
	if err != nil {
		log.Println("GetForcedSaleAmount() err=", err)
		return "", err
	}
	var result hexutil.Big
	err = client.CallContext(context.Background(), &result, "eth_getForcedSaleAmount", nftaddr)
	if err != nil {
		log.Println("GetForcedSaleAmount() err=", err)
		return "", err
	}
	return result.ToInt().String(), err
}

func IsOfficialNFT(nftAddress common.Address) bool {
	maskByte := byte(128)
	nftByte := nftAddress[0]
	result := maskByte & nftByte
	if result == 128 {
		return true
	}
	return false
}

var ExchangePeriod = uint64(6160) // 365 * 720 * 24 * 4 / 4096
func GetExchangAmount(nftaddress common.Address, initamount *big.Int) *big.Int {
	nftInt := new(big.Int).SetBytes(nftaddress.Bytes())
	baseInt, _ := big.NewInt(0).SetString("8000000000000000000000000000000000000000", 16)
	nftInt.Sub(nftInt, baseInt)
	//nftInt.Add(nftInt, big.NewInt(1))
	nftInt.Div(nftInt, big.NewInt(4096))
	times := nftInt.Uint64() / ExchangePeriod
	rewardratio := math.Pow(0.88, float64(times))
	result := big.NewInt(0)
	new(big.Float).Mul(big.NewFloat(rewardratio), new(big.Float).SetInt(initamount)).Int(result)

	return result
}

func CalculateExchangeAmount(level uint8, mergenumber uint32) *big.Int {
	//nftNumber := math.BigPow(16, int64(level))
	nftNumber := big.NewInt(int64(mergenumber))
	switch {
	case level == 0:
		radix, _ := big.NewInt(0).SetString("30000000000000000", 10)
		return big.NewInt(0).Mul(nftNumber, radix)
	case level == 1:
		radix, _ := big.NewInt(0).SetString("143000000000000000", 10)
		return big.NewInt(0).Mul(nftNumber, radix)
	case level == 2:
		radix, _ := big.NewInt(0).SetString("271000000000000000", 10)
		return big.NewInt(0).Mul(nftNumber, radix)
	default:
		radix, _ := big.NewInt(0).SetString("650000000000000000", 10)
		return big.NewInt(0).Mul(nftNumber, radix)
	}
}

func GetForcedSaleAmountByAddress(nftAddress common.Address) (*big.Int, error) {
	if !IsOfficialNFT(nftAddress) {
		return nil, errors.New("not official nft")
	}
	initAmount := CalculateExchangeAmount(1, 1)
	amount := GetExchangAmount(nftAddress, initAmount)
	return amount, nil
}

func GetUserMintDeep(blockNumber uint64) (string, error) {
	client, err := rpc.Dial(EthNode)
	if err != nil {
		log.Println("GetUserMintDeep() err=", err)
		return "", err
	}
	var result string
	blockN := hexutil.EncodeUint64(blockNumber)
	err = client.CallContext(context.Background(), &result, "eth_getUserMintDeep", blockN)
	if err != nil {
		log.Println("GetUserMintDeep() err=", err)
		return "", err
	}
	return result, err
}

func GetCurrentBlockNumber() uint64 {
	var client *ethclient.Client
	var err error
	for {
		for {
			client, err = ethclient.Dial(EthNode)
			if err != nil {
				log.Println("GetCurrentBlockNumber()", "EthNode=", EthNode, " connect err=", err)
				time.Sleep(ReDialDelyTime * time.Second)
			} else {
				//log.Println("GetCurrentBlockNumber() connect OK!")
				//log.Println("GetCurrentBlockNumber() connect OK!")
				break
			}
		}
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Println("GetCurrentBlockNumber() get HeaderByNumber err=", err)
			client.Close()
			time.Sleep(ReDialDelyTime * time.Second)
		} else {
			log.Println("GetCurrentBlockNumber() header.Number=", header.Number.String())
			client.Close()
			return header.Number.Uint64()
		}
	}
}
