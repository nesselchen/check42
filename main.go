package main

import (
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
		Addr:      "mysql:3306",
		DBName:    "check42",
		ParseTime: true,
	}

	var store *stores.MySQLTodoStore
	retries := 1
	const maxRetries = 10
	for {
		var err error
		store, err = stores.NewMySQLTodoStore(config)
		if err == nil {
			fmt.Println()
			break
		}
		if retries >= maxRetries {
			log.Fatal("Maximum number of retries exceeded")
		}
		fmt.Printf("Connection to database failed. Retrying (%d/%d)\n", retries, maxRetries)
		time.Sleep(3 * time.Second)
		retries++
	}

	fmt.Println("Connection successful")
	defer store.Close()

	fmt.Println(logo)
	api.RunServer("0.0.0.0:2442", store)
}
