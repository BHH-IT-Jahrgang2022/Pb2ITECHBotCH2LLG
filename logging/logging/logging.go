package logging

import (
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type LogEntryJSON struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type LogEntryDB struct {
	ID        int    `gorm:"primaryKey"`
	Timestamp int64  `gorm:"column:timestamp"`
	Level     string `gorm:"column:level"`
	Message   string `gorm:"column:message"`
}

func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DATABASE") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&LogEntryDB{}); err != nil {
		return nil, err
	}
	return db, nil
}

// type ConnectionSpecs struct {
// 	username   string
// 	password   string
// 	socketPath string
// 	database   string
// }

// func GetConnStr() string {
// 	connConfig := ConnectionSpecs{
// 		username:   os.Getenv("MYSQL_USER"),
// 		password:   os.Getenv("MYSQL_PASSWORD"),
// 		socketPath: os.Getenv("MYSQL_SOCKET_PATH"),
// 		database:   os.Getenv("MYSQL_DATABASE"),
// 	}

// 	connStr := connConfig.username + ":" + connConfig.password +
// 		"@unix(" + connConfig.socketPath + ")" +
// 		"/" + connConfig.database +
// 		"?charset=utf8"
// 	return connStr
// }

// Creates a database connection and returns the same
// func InitDB() (*sql.DB, error) {
// 	var connStr = GetConnStr()
// 	fmt.Println(connStr)
// 	conn, err := sql.Open("mysql", GetConnStr())
// 	return conn, err
// }

// func InitDB() *sql.DB {
// 	sqlConfig := mysql.Config{
// 		User:   os.Getenv("MYSQL_USER"),
// 		Passwd: os.Getenv("MYSQL_PASSWORD"),
// 		Net:    "tcp",
// 		Addr:   fmt.Sprintf("%s:%s", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT")),
// 		DBName: os.Getenv("Logging_DB"),
// 	}
// 	db, err := sql.Open("mysql", sqlConfig.FormatDSN())
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil
// 	}
// 	if err := db.Ping(); err != nil {
// 		log.Fatal(err)
// 		return nil
// 	}
// 	return db
// }

func StartApi() {
	time.Sleep(60 * time.Second)
	db, err := InitDB()
	if err != nil {
		db_connection_timeout := 12
		for i := 0; i < db_connection_timeout; i++ {
			time.Sleep(10 * time.Second)
			db, err = InitDB()
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Fatal("Could not connect to the database")
		}
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	// Endpoint for services to send logs to
	r.POST("/log", func(c *gin.Context) {
		var log LogEntryJSON
		if err := c.ShouldBindJSON(&log); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}
		var dbLog LogEntryDB
		dbLog.Timestamp = log.Timestamp
		dbLog.Level = log.Level
		dbLog.Message = log.Message
		if err := db.Save(dbLog); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Log saved"})
	})
	// Endpoint for services to get logs
	r.GET("")
	r.Run(os.Getenv("LOGGING_API_PORT"))
}
