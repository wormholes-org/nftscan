package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
	"sync"
)

type ScanDb struct {
	db     *gorm.DB
	sqlstr string
}

var (
	scandb *ScanDb
	scanmu sync.Mutex
)

func NewNftDb(sqldsn string) (*ScanDb, error) {
	scanmu.Lock()
	defer scanmu.Unlock()
	if scandb != nil {
		return scandb, nil
	}
	d := new(ScanDb)
	var err error
	d.db, err = gorm.Open(mysql.Open(sqldsn), &gorm.Config{})
	if err != nil {
		log.Println("NewNftDb() failed to connect database", err.Error())
		return nil, err
	}
	log.Println("NewNftDb() Open connect database Ok.")
	scandb = d
	return d, err
}

func GetScanDB(sqldsn ...string) *gorm.DB {
	if scandb != nil {
		return scandb.db
	}
	sqlstr := sqldsn[0]
	_, err := NewNftDb(sqlstr)
	if err != nil {
		return nil
	}
	return scandb.db
}

func (d ScanDb) Close() {
	fmt.Println("Close nftdb.")
}

func (d *ScanDb) GetDB() *gorm.DB {
	return d.db
}

func createDb(d *gorm.DB, dbName string) error {
	strOrder := "create database if not exists " + dbName + ";"
	db := d.Exec(strOrder)
	if db.Error != nil {
		fmt.Printf("CreateDataBase err=%s\n", db.Error)
		return db.Error
	}
	strOrder = "use " + dbName
	db = d.Exec(strOrder)
	if db.Error != nil {
		fmt.Printf("use database err=%s\n", db.Error)
	}
	return db.Error
}

func getCreateIndexOrder() []string {
	return []string{
		"CREATE INDEX indexNftsContractTokenidDeleted ON nfts (contract, tokenid, deleted_at);",
		"CREATE INDEX indexNftsTokenidDeletedat ON nfts ( tokenid, deleted_at );",
	}
}

func (d ScanDb) CreateIndexs() error {
	/*for _, s := range getCreateIndexOrder() {
		db := nft.db.Exec(s)
		if db.Error != nil {
			if !strings.Contains(db.Error.Error(), "1061") {
				fmt.Println("CreateIndexs() ",s[len("CREATE INDEX"):strings.Index(s, "ON nfts")],  "err=", db.Error)
				return db.Error
			}
		}
	}*/
	strOrder := "CREATE INDEX indexNftsCreateaddrTokenid ON nfts ( createaddr, tokenid );"
	db := d.db.Exec(strOrder)
	if db.Error != nil {
		if !strings.Contains(db.Error.Error(), "1061") {
			fmt.Printf("CreateIndexs() indexNftsCreateaddrTokenid  err=%s\n", db.Error)
			return db.Error
		}
	}
	return nil
}

/*func getCreateTableObject() []interface{} {
	return []interface{}{
		Nfts{},
		SysParams{},
		SysInfos{},
	}
}*/

/*func (nft NftDb) CreateTables() error {
	for _, s := range getCreateTableObject() {
		err := nft.db.AutoMigrate(s)
		if err != nil {
			t := reflect.TypeOf(s)
			fmt.Println("create table ", t.Name(), "err=", err)
			return err
		}
	}
	return nil
}
*/

func InitDb(sqlsvr string, dbName string) error {
	nft, err := gorm.Open(mysql.Open(sqlsvr), &gorm.Config{})
	if err != nil {
		log.Println("NewNftDb() failed to connect database", err.Error())
		return err
	}
	if err != nil {
		fmt.Printf("InitDb()->connectDb() err=%s\n", err)
		return err
	}
	err = createDb(nft, dbName)
	if err != nil {
		fmt.Printf("Create Db err=%s\n", err)
		return err
	}
	err = nft.AutoMigrate(&SysParams{})
	if err != nil {
		fmt.Println("create table SysParams{} err=", err)
		return err
	}
	/*err = nft.AutoMigrate(&SysInfos{})
	if err != nil {
		fmt.Println("create table SysInfos{} err=", err)
		return err
	}*/
	err = nft.AutoMigrate(&NftTxs{})
	if err != nil {
		fmt.Println("create table NftTxs{} err=", err)
		return err
	}
	err = nft.AutoMigrate(&Snfts{})
	if err != nil {
		fmt.Println("create table Snfts{} err=", err)
		return err
	}

	strOrder := "CREATE INDEX indexNftsBlocknumberNftaddr ON snfts (blocknumber, nftaddr);"
	db := nft.Exec(strOrder)
	if db.Error != nil {
		if !strings.Contains(db.Error.Error(), "1061") {
			fmt.Printf("CreateIndexs() indexNftsBlocknumberNftaddr  err=%s\n", db.Error)
			return db.Error
		}
	}
	strOrder = "CREATE INDEX indexSNftsOwnAddr ON snfts ( ownaddr);"
	db = nft.Exec(strOrder)
	if db.Error != nil {
		if !strings.Contains(db.Error.Error(), "1061") {
			fmt.Printf("CreateIndexs() indexNftsBlocknumberNftaddr  err=%s\n", db.Error)
			return db.Error
		}
	}
	strOrder = "CREATE INDEX indexNfttxsBlockDelete ON nfttxs (Blocknumber, deleted_at);"
	db = nft.Exec(strOrder)
	if db.Error != nil {
		if !strings.Contains(db.Error.Error(), "1061") {
			fmt.Printf("CreateIndexs() indexNftsBlocknumberNftaddr  err=%s\n", db.Error)
			return db.Error
		}
	}
	strOrder = "CREATE INDEX indexSnftsBlockDelete ON snfts (Blocknumber, deleted_at);"
	db = nft.Exec(strOrder)
	if db.Error != nil {
		if !strings.Contains(db.Error.Error(), "1061") {
			fmt.Printf("CreateIndexs() indexNftsBlocknumberNftaddr  err=%s\n", db.Error)
			return db.Error
		}
	}
	strOrder = "CREATE INDEX indexSNftsNftaddrDeleted ON snfts (nftaddr, deleted_at);"
	db = nft.Exec(strOrder)
	if db.Error != nil {
		if !strings.Contains(db.Error.Error(), "1061") {
			fmt.Printf("CreateIndexs() indexSNftsNftaddrDeleted  err=%s\n", db.Error)
			return db.Error
		}
	}
	strOrder = "CREATE INDEX indexSNftsBlocknumberDeleted ON snfts (blocknumber, deleted_at);"
	db = nft.Exec(strOrder)
	if db.Error != nil {
		if !strings.Contains(db.Error.Error(), "1061") {
			fmt.Printf("CreateIndexs() indexSNftsBlocknumberDeleted  err=%s\n", db.Error)
			return db.Error
		}
	}
	strOrder = "CREATE INDEX indexSNftsOwnaddrDeleted ON snfts (ownaddr, deleted_at);"
	db = nft.Exec(strOrder)
	if db.Error != nil {
		if !strings.Contains(db.Error.Error(), "1061") {
			fmt.Printf("CreateIndexs() indexSNftsOwnaddrDeleted  err=%s\n", db.Error)
			return db.Error
		}
	}
	strOrder = "CREATE INDEX indexNfttxsTxhasDeleted ON nfttxs ( txhash, deleted_at );"
	db = nft.Exec(strOrder)
	if db.Error != nil {
		if !strings.Contains(db.Error.Error(), "1061") {
			fmt.Printf("CreateIndexs() indexNfttxsTxhasDeleted  err=%s\n", db.Error)
			return db.Error
		}
	}
	//nft.Close()
	_, err = NewNftDb(sqlsvr + dbName + localtime)
	return err
}
