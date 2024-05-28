package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"check42/api"
	"check42/store/stores"

	"github.com/go-sql-driver/mysql"
)

var logo string = `
 ██████╗██╗  ██╗███████╗ ██████╗██╗  ██╗ ██╗  ██╗██████╗
██╔════╝██║  ██║██╔════╝██╔════╝██║ ██╔╝ ██║  ██║╚════██╗
██║     ███████║█████╗  ██║     █████╔╝  ███████║ █████╔╝
██║     ██╔══██║██╔══╝  ██║     ██╔═██╗  ╚════██║██╔═══╝
╚██████╗██║  ██║███████╗╚██████╗██║  ██╗      ██║███████╗
 ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝      ╚═╝╚══════╝
`

func main() {
	config := mysql.Config{
		User:      "root",
		Passwd:    "root",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "check42",
		ParseTime: true,
	}

	db, err := connectWithRetries(config, 10)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Connection successful")

	todos := stores.NewMySQLTodoStore(db)
	users := stores.NewMySQLUserStore(db)

	fmt.Println(logo)
	api.RunServer("127.0.0.1:2442", todos, users)
}

func connectWithRetries(config mysql.Config, maxTries int) (*sql.DB, error) {
	tries := 1
	for tries < maxTries {
		db, err := sql.Open("mysql", config.FormatDSN())
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err != nil {
			fmt.Printf("Connection to database failed. Retrying (%d/%d)\n", tries, maxTries)
			time.Sleep(3 * time.Second)
			tries++
			continue
		}
		return db, nil
	}
	log.Fatal("Maximum number of retries exceeded")
	return nil, nil
}
