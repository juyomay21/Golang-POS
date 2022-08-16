// database/connect.go
package models

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ittechman101/go-pos/config"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	// Get from .env file
	db_host := config.Config("DB_HOST")
	db_name := config.Config("DB_NAME")
	db_user := config.Config("DB_USER")
	db_passowrd := config.Config("DB_PASSWORD")

	// p := config.Config("DB_PORT")
	// port, err := strconv.ParseUint(p, 10, 32)

	// Get from Docker environment variables
	if len(os.Getenv("MYSQL_HOST")) > 0 {
		db_host = os.Getenv("MYSQL_HOST")
	}

	// p = os.Getenv("MYSQL_PORT")
	// port, err = strconv.ParseUint(p, 10, 32)

	if len(os.Getenv("MYSQL_USER")) > 0 {
		db_user = os.Getenv("MYSQL_USER")
	}

	if len(os.Getenv("MYSQL_PASSWORD")) > 0 {
		db_passowrd = os.Getenv("MYSQL_PASSWORD")
	}

	if len(os.Getenv("MYSQL_DBNAME")) > 0 {
		db_name = os.Getenv("MYSQL_DBNAME")
	}

	// fmt.Println("---env test---")
	// fmt.Println(db_host)
	// fmt.Println(port)
	// fmt.Println(db_user)
	// fmt.Println(db_passowrd)
	// fmt.Println(db_name)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db_user,
		db_passowrd,
		db_host,
		db_name)

	DB, err = gorm.Open("mysql", dsn)

	DB.DB().SetConnMaxLifetime(100)
	DB.DB().SetMaxIdleConns(10)

	// set max connection
	//	DB.SetConnMaxLifetime(100)
	// set max idle connections
	//	DB.SetMaxIdleConns(10)

	// err = DB.DB().Ping()
	// if err != nil {
	// 	log.WithError(err).Fatal("error while pinging DB")
	// }

	if err != nil {
		fmt.Println(
			err.Error(),
		)
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	DB.AutoMigrate(&Cashiers{})
	DB.AutoMigrate(&Categories{})
	DB.AutoMigrate(&Products{})
	DB.AutoMigrate(&Payments{})
	DB.AutoMigrate(&Order{})
	DB.AutoMigrate(&OrderProducts{})
}

func GetDB() *gorm.DB {

	return DB
}
