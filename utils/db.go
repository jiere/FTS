package utils

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var (
	db *gorm.DB
	//DbDialect such as "mysql", "postgres", etc.
	dbDialect string
	dbHost    string
	dbPort    uint32
	dbUser    string
	dbPasswd  string
	// DbName is the project databse name
	dbName string
)

// DB returns the gorm DB handle
func DB() *gorm.DB {
	return db
}

// Host returns the server IP where fts runs
func Host() string {
	return dbHost
}

func parseDBConfig() {
	config := viper.New()
	config.AddConfigPath("cfg/")
	config.SetConfigName("db")
	config.SetConfigType("ini")

	if err := config.ReadInConfig(); err != nil {
		log.Fatalln(err)		
		//panic(err)
	}

	dbDialect = config.GetString("db.dialect")
	dbHost = config.GetString("db.host")
	dbPort = config.GetUint32("db.port")
	dbUser = config.GetString("db.user")
	dbPasswd = config.GetString("db.passwd")
	dbName = config.GetString("db.name")
}

// Initializes the database instance
func setup() {
	var err error
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPasswd, dbHost, dbPort, dbName)
	fmt.Println(connStr)
	db, err = gorm.Open(dbDialect, connStr)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
		fmt.Println("failed to connect database:", err)
		return
	}
	fmt.Println("connect database success")
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

// DBInit :called by main.go only once when system init
func DBInit() {
	parseDBConfig()
	setup()
}
