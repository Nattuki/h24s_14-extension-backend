package database

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/srinathgs/mysqlstore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	label_mysql      *gorm.DB
	label_mysql_once sync.Once
)

func createMysqlDB(dbname, host, user, pass, port string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	return db
}

func GetDBConnection() *gorm.DB {
	label_mysql_once.Do(func() {
		if label_mysql == nil {
			dbName := os.Getenv("NS_MARIADB_DATABASE")
			host := os.Getenv("NS_MARIADB_HOSTNAME")
			port := os.Getenv("NS_MARIADB_PORT")
			user := os.Getenv("NS_MARIADB_USER")
			pass := os.Getenv("NS_MARIADB_PASSWORD")
			label_mysql = createMysqlDB(dbName, host, user, pass, port)
		}
	})

	return label_mysql
}

func UseSessionStore(e *echo.Echo) {
	db := GetDBConnection()
	db_sql, _ := db.DB()
	store, err := mysqlstore.NewMySQLStoreFromConnection(db_sql, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		log.Fatal(err)
	}
	e.Use(session.Middleware(store))
}
