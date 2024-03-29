package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type order struct {
	OrderUID    string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
	Entry       string `json:"entry"`
	Delivery    struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	} `json:"delivery"`
	Payment struct {
		Transaction  string `json:"transaction"`
		RequestID    string `json:"request_id"`
		Currency     string `json:"currency"`
		Provider     string `json:"provider"`
		Amount       int    `json:"amount"`
		PaymentDt    int    `json:"payment_dt"`
		Bank         string `json:"bank"`
		DeliveryCost int    `json:"delivery_cost"`
		GoodsTotal   int    `json:"goods_total"`
		CustomFee    int    `json:"custom_fee"`
	} `json:"payment"`
	Items []struct {
		ChrtID      int    `json:"chrt_id"`
		TrackNumber string `json:"track_number"`
		Price       int    `json:"price"`
		Rid         string `json:"rid"`
		Name        string `json:"name"`
		Sale        int    `json:"sale"`
		Size        string `json:"size"`
		TotalPrice  int    `json:"total_price"`
		NmID        int    `json:"nm_id"`
		Brand       string `json:"brand"`
		Status      int    `json:"status"`
	} `json:"items"`
	Locale            string `json:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerID        string `json:"customer_id"`
	DeliveryService   string `json:"delivery_service"`
	Shardkey          string `json:"shardkey"`
	SmID              int    `json:"sm_id"`
	DateCreated       string `json:"date_created"`
	OofShard          string `json:"oof_shard"`
}

var cached_data map[string]order

func nats_streaming_connection(conn *pgx.Conn, mutex *sync.Mutex) (*nats.Conn, stan.Conn, stan.Subscription) {
	// connecting to NATS Streaming
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		fmt.Printf("Failed to connect to NATS Streaming server: %v\n", err)
		return nil, nil, nil
	}

	// connecting to NATS Streaming channel
	sc, err := stan.Connect("test-cluster", "client-123", stan.NatsConn(nc))
	if err != nil {
		fmt.Printf("Failed to connect to NATS Streaming channel: %v\n", err)
		return nil, nil, nil
	}

	// subscribing to NATS Streaming channel
	sub, err := sc.Subscribe("test-channel", func(msg *stan.Msg) {
		// anon function for saving data to cache and DB from this channel

		var new_order order
		err := json.Unmarshal(msg.Data, &new_order)
		if err != nil {
			fmt.Printf("Failed to parse JSON: %v\n", err)
			fmt.Printf("CATCHING INVALID DATA AS IT'S SUPPOSED ACCORDING TO THE TASK\n")
		} else {
			// locking mutex and appending new order to the global slice in cache
			mutex.Lock()
			cached_data[new_order.OrderUID] = new_order
			mutex.Unlock()

			// inserting new order to Postgres
			_, err = conn.Exec(context.Background(), "insert into wb.order values ($1, $2)", new_order.OrderUID, string(msg.Data))
			if err != nil {
				fmt.Printf("Failed to insert new order: %v\n", err)
			}
		}
	})
	if err != nil {
		fmt.Printf("Failed to subscribe to NATS Streaming channel: %v\n", err)
		return nil, nil, nil
	}

	return nc, sc, sub
}

func postgres_connection() *pgx.Conn {
	// connecting to Postgres
	var psql_port = "10000/postgres"
	conn, err := pgx.Connect(context.Background(), "postgres://user:passw0rd@localhost:"+psql_port)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
	}

	// creating procedures (you can find more in DBhandler.sql file)
	file_name := "DB/DBhandler.sql"
	commands, err := os.ReadFile(file_name)
	if err != nil {
		fmt.Printf("Failed to read file %v:  %v\n", file_name, err)
		return nil
	}
	_, err = conn.Exec(context.Background(), string(commands))
	if err != nil {
		fmt.Printf("An error in file %v:  %v\n", file_name, err)
	}

	return conn
}

func http_server_start() *http.Server {
	// starting http server
	server := &http.Server{
		Addr: "localhost:8111",
	}
	serverIsRunning := make(chan bool)

	// starting server in goroutine
	go func(serverIsRunning chan bool) {
		fmt.Printf("\nTrying to start server on port 8111...\n")
		serverIsRunning <- true
		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("Server is shutted down. %v\n", err)
		}
	}(serverIsRunning)

	// waiting for goroutine to start execution
	<-serverIsRunning
	close(serverIsRunning)
	if server != nil {
		fmt.Printf("Server started on port 8111\n")
	}
	return server
}

func server_handler(conn *pgx.Conn, sc stan.Conn, mutex *sync.Mutex) {
	http.HandleFunc("/", main_page(conn, sc, mutex))
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
}

func main_page(conn *pgx.Conn, sc stan.Conn, mutex *sync.Mutex) http.HandlerFunc {
	// wrapping main page handler
	return func(w http.ResponseWriter, r *http.Request) {

		got_uid := "got no or invalid uid"
		got_order := order{}
		got_order.OrderUID = "got no or invalid uid"

		// checking for new file in POST request (to send it to NATS Streaming channel)
		if r.Method == http.MethodPost {
			new_file, _, err := r.FormFile("new_file")
			if err != nil {
				fmt.Printf("Failed to read file: %v\n", err)
			} else {
				data, err := io.ReadAll(new_file)
				if err != nil {
					fmt.Printf("Failed to read file: %v\n", err)
				} else {
					var new_order order
					err := json.Unmarshal(data, &new_order)
					if err != nil {
						fmt.Printf("Failed to parse JSON: %v\n", err)
					}

					err = sc.Publish("test-channel", data)
					if err != nil {
						fmt.Printf("Failed to send message to NATS Streaming channel: %v\n", err)
					}

				}
				new_file.Close()
			}
			if r.FormValue("uid") != "" {
				got_uid = r.FormValue("uid")
				got_order = cached_data[got_uid]
				if got_order.OrderUID == "" {
					got_order.OrderUID = "got no or invalid uid"
				}
			}
		}

		data := map[string]interface{}{
			"JSONs": cached_data,
			"UID":   got_uid,
			"order": got_order,
		}

		t, _ := template.ParseFiles("web/main.html")
		err := t.Execute(w, data)
		if err != nil {
			fmt.Printf("Failed to execute template: %v\n", err)
		}
	}
}

func recover_cache_from_database(conn *pgx.Conn, mutex *sync.Mutex) {
	rows, err := conn.Query(context.Background(), "select * from wb.order")
	if err != nil {
		fmt.Printf("Failed to execute query: %v\n", err)
	}
	for rows.Next() {
		var new_order_data []byte
		var id string
		var new_order order
		err := rows.Scan(&id, &new_order_data)
		if err != nil {
			fmt.Printf("Failed to scan row: %v\n", err)
		}
		err = json.Unmarshal(new_order_data, &new_order)
		if err != nil {
			fmt.Printf("Failed to parse JSON: %v\n", err)
		}
		mutex.Lock()
		cached_data[id] = new_order
		mutex.Unlock()
	}
}

func main() {
	cached_data = make(map[string]order)
	var mu sync.Mutex
	conn := postgres_connection()
	if conn == nil {
		fmt.Printf("[MAIN] Failed to connect to Postgres server\n")
		return
	}
	recover_cache_from_database(conn, &mu)
	nc, sc, sub := nats_streaming_connection(conn, &mu)
	if nc == nil || sc == nil || sub == nil {
		fmt.Printf("[MAIN] Failed to connect to NATS Streaming server\n")
		return
	}
	server := http_server_start()
	if server == nil {
		fmt.Printf("[MAIN] Failed to start server\n")
		return
	}

	// creating DB in Postgres if it doesn't exist
	/*_, err := conn.Exec(context.Background(), "call create_DB()")
	if err != nil {
		fmt.Printf("An error calling CREATE_DB:  %v\n", err)
	}*/

	go server_handler(conn, sc, &mu)

	// stop server

	fmt.Printf("\n\nType anything to stop server\n")
	var stop string
	fmt.Scan(&stop)

	server.Shutdown(context.Background())

	// deleting DB in Postgres
	/*_, err := conn.Exec(context.Background(), "call drop_DB()")
	if err != nil {
		fmt.Printf("An error calling DROP_DB:  %v\n", err)
	}*/

	sub.Unsubscribe()
	sc.Close()
	nc.Close()
	conn.Close(context.Background())
	fmt.Printf("[MAIN] Program is shutted down\n")

	time.Sleep(1 * time.Second)
}
