package model

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var mySqlDB MySqlDB

func GetMySqlDB() MySqlDB {
	return mySqlDB
}

type MySqlDB struct {
	*gorm.DB
}

// var dbUserName = *flag.String("dbUserName", "s9_master", "db user name")
// var dbPassword = *flag.String("dbPassword", "CwrcxxTB72zDBZEw", "db password")
// var dbHost = *flag.String("dbHost", "127.0.0.1", "db host")
// var dbPort = *flag.String("dbPort", "3306", "db port")
// var dbName = *flag.String("dbName", "s9_platform", "db name")
var dbUserName = "s9_master"
var dbPassword = "CwrcxxTB72zDBZEw"
var dbHost = "127.0.0.1"
var dbPort = "3306"
var dbName = "s9_platform"

func init() {
	// flag.Parse()
	DBInit()
}

func DBInit() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUserName, dbPassword, dbHost, dbPort, dbName)
	// log.Fatal(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
		Logger:                                   logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal(fmt.Sprintf("%+v", errors.WithStack(err)))
	}

	mySqlDB = MySqlDB{db}

	initTableAndProcedure()
}

func initTableAndProcedure() {

	// create table
	mySqlDB.AutoMigrate(&PlinkoBalls{})
	mySqlDB.AutoMigrate(&PlinkoOdds{})

	// create initial
	initialData("plinko_balls", "resources/initialdata/plinko_balls.sql")
	initialData("plinko_odds", "resources/initialdata/plinko_odds.sql")
}

func initialData(tableName, sqlFilePath string) {
	if mySqlDB.Migrator().HasTable(tableName) {
		var value int64
		mySqlDB.Table(tableName).Count(&value)
		if value == 0 {
			mydir, _ := os.Getwd()
			// fmt.Println(mydir + "/" + sqlFilePath)
			query, err := ioutil.ReadFile(mydir + "/" + sqlFilePath)
			if err != nil {
				log.Fatal(fmt.Sprintf("%+v", errors.WithStack(err)))
			}

			sqlAll := string(query)
			for _, v := range strings.Split(sqlAll, ";") {
				if v != "" {
					mySqlDB.Exec(v)
				}
			}

		}
	}
}
