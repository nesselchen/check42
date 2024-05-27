package main

import (
	"fmt"
	"log"
	"time"

	"check42/api"
	"check42/model/todos"
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
	fmt.Print("Connecting...")
	for {
		var err error
		store, err = stores.NewMySQLTodoStore(config)
		if err == nil {
			fmt.Println()
			break
		}
		fmt.Print(".")
		time.Sleep(3 * time.Second)
	}
	defer store.Close()

	err := store.CreateTodo(todos.Todo{
		Owner: 1,
		Text:  "Test",
		Done:  false,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(logo)
	api.RunServer("0.0.0.0:2442", store)
}
