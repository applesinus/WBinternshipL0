package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

var conn *pgx.Conn
var connection_err error

var psql_port = "10000/postgres"

func main() {
	fmt.Printf("Hello, World!\n")

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		fmt.Printf("Failed to connect to NATS Streaming server: %v\n", err)
	}
	defer nc.Close()

	sc, err := stan.Connect("cluster", "client-123", stan.NatsConn(nc))
	if err != nil {
		fmt.Printf("Failed to connect to NATS Streaming channel: %v\n", err)
	}
	defer sc.Close()

	sub, err := sc.Subscribe("test-channel", func(msg *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(msg.Data))
	})
	if err != nil {
		fmt.Printf("Failed to subscribe to NATS Streaming channel: %v\n", err)
	}
	defer sub.Unsubscribe()

	err = sc.Publish("test-channel", []byte("Hello, NATS Streaming!"))
	if err != nil {
		fmt.Printf("Failed to send message to NATS Streaming channel: %v\n", err)
	}

	// splitter

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
