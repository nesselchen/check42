package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// var logo string = `
//  ██████╗██╗  ██╗███████╗ ██████╗██╗  ██╗ ██╗  ██╗██████╗
// ██╔════╝██║  ██║██╔════╝██╔════╝██║ ██╔╝ ██║  ██║╚════██╗
// ██║     ███████║█████╗  ██║     █████╔╝  ███████║ █████╔╝
// ██║     ██╔══██║██╔══╝  ██║     ██╔═██╗  ╚════██║██╔═══╝
// ╚██████╗██║  ██║███████╗╚██████╗██║  ██╗      ██║███████╗
//  ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝      ╚═╝╚══════╝
// `

func main() {
	// store := store.NewJsonTodoStore("db.json")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("root:%s@tcp(mysql:%s)/", password, port)
	time.Sleep(2 * time.Second)

	db, _ := sql.Open("mysql", dsn)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(`CREATE DATABASE IF NOT EXISTS check42`); err != nil {
		log.Fatal("Err in line 35 ", err)
	}

	db, _ = sql.Open("mysql", dsn+"check42")
	// r, err := db.Query("SELECT * FROM user")
	// if err != nil {
	// 	log.Fatal("Err in line 43 ", err)
	// }
	// for r.Next() {
	// 	var id int
	// 	var name string
	// 	r.Scan(&id, &name)
	// 	fmt.Println(id, name)
	// }

	// fmt.Println(logo)
	// api.RunServer("0.0.0.0:2442", store)
}
