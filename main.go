package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "secret"
	hostname = "127.0.0.1:3306"
	dbname   = "test"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

var DB *sql.DB

func init() {
	_db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		panic(err)
	}
	DB = _db
}

func createEC2(servedId int) {
	fmt.Println("Creating server")
	_, err := DB.Exec("UPDATE Server SET status = 'TODO' WHERE id = ?", servedId)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
	_, err = DB.Exec("UPDATE Server SET status = 'IN_PROGRESS' WHERE id = ?", servedId)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server creation in progress")

	time.Sleep(5 * time.Second)
	_, err = DB.Exec("UPDATE Server SET status = 'DONE' WHERE id = ?", servedId)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server Created")
}

func main() {
	ge := gin.Default()

	ge.POST("/servers", func(ctx *gin.Context) {
		go createEC2(1)
		ctx.JSON(200, map[string]interface{}{"message": "ok"})
	})

	// short polling
	ge.GET("/short/status/:server_id", func(ctx *gin.Context) {
		serverId := ctx.Param("server_id")
		var status string
		row := DB.QueryRow("SELECT status FROM Server WHERE id = ?;", serverId)
		if row.Err() != nil {
			panic(row.Err())
		}

		row.Scan(&status)
		ctx.JSON(200, map[string]interface{}{"status": status})
	})

	// long polling
	ge.GET("/long/status/:server_id", func(ctx *gin.Context) {
		serverId := ctx.Param("server_id")
		currentStatus := ctx.Query("status")

		var status string
		for {
			row := DB.QueryRow("SELECT status FROM Server WHERE id = ?;", serverId)
			if row.Err() != nil {
				panic(row.Err())
			}
			row.Scan(&status)

			if currentStatus != status {
				break
			}

			time.Sleep(1 * time.Second)
		}

		ctx.JSON(200, map[string]interface{}{"status": status})
	})

	ge.Run(":9000")
}
