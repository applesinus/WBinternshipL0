package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn
var connection_err error

var psql_port = "10000/postgres"

func main() {
	fmt.Printf("Hello, World!\n")

	conn, connection_err = pgx.Connect(context.Background(), "postgres://user:passw0rd@localhost:"+psql_port)

	if connection_err != nil {
		fmt.Printf("DB CONNECTION IS FAILED: %v\n", connection_err)
	}

	fmt.Printf("DB CONNECTION IS SUCCESS\n")

	file_name := "DB/DBhandler.sql"

	commands, err := os.ReadFile(file_name)
	if err != nil {
		fmt.Printf("NO SUCH FILE! err: %v\n", err)
		return
	}
	_, err = conn.Exec(context.Background(), string(commands))
	if err != nil {
		fmt.Printf("An error in file %v:  %v\n", file_name, err)
	}

	_, err = conn.Exec(context.Background(), "call create_DB()")
	if err != nil {
		fmt.Printf("An error calling CREATE_DB in file %v:  %v\n", file_name, err)
	}
	_, err = conn.Exec(context.Background(), "call drop_DB()")
	if err != nil {
		fmt.Printf("An error calling CREATE_DB in file %v:  %v\n", file_name, err)
	}

	var wait int
	fmt.Scanln(&wait)
}
