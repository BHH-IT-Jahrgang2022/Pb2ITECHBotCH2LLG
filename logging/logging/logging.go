package logging

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

func (l *LogEntry) Save() error {
	db, err := sql.Open("sqlite3", "logs.db")
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func StartApi() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	// Endpoint for services to send logs to
	r.POST("/log", func(c *gin.Context) {
		var log LogEntry
		if err := c.ShouldBindJSON(&log); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := log.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Log saved"})
	})
	r.GET("")
	r.Run(os.Getenv("LOGGING_API_PORT"))
}
