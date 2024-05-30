package main

import (
	"check42/api"
	"check42/store/stores"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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
	godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	config := mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      dbHost + ":" + dbPort,
		DBName:    "check42",
		ParseTime: true,
	}

	maxTries, err := strconv.ParseInt(os.Getenv("DB_RETRIES"), 10, 10)
	if err != nil {
		log.Fatal("missing .env file")
	}
	db, err := connectWithRetries(config, int(maxTries))

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Connection successful")

	todos := stores.NewMySQLTodoStore(db)
	users := stores.NewMySQLUserStore(db)

	fmt.Println(logo)

	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")

	api.RunServer(host+":"+port, todos, users)
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
