package wallet

import (
	"context"
	"database/sql"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"wallet-go/cors"
	"wallet-go/debezium"
	"wallet-go/kafka"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Wallet struct {
	server         *http.Server
	db             *sql.DB
	cacheLock      sync.Mutex
	dbzMess        chan *sarama.ConsumerMessage
	consumer       kafka.KafkaConsumer
	userFundsCache map[string]int
}

//Init wallet service providing necessary parameters for database and kafka consumer
func (w *Wallet) Init(port, dbUsr, dbPw, dbHost, dbPort, dbName, kafkaTopic string, brokers []string, ch chan *sarama.ConsumerMessage) error {
	w.dbzMess = ch
	w.userFundsCache = make(map[string]int)
	connStr := fmt.Sprintf("postgresql://%s:%s@%s%s/%s?sslmode=disable", dbUsr, dbPw, dbHost, dbPort, dbName)

	var err error
	w.db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	r := mux.NewRouter()
	r.HandleFunc("/funds/{user}", w.GetFunds()).Methods("GET")
	r.HandleFunc("/funds/{user}", w.AddFunds()).Methods("POST")

	w.server = &http.Server{
		Addr:    port,
		Handler: cors.CORSEnabled(r),
	}
	err = w.consumer.Init(brokers, kafkaTopic)
	if err != nil {
		log.Println("[ERROR]Consumer init failed")
		return err
	}
	return nil
}

//Run - runs wallet to listen rest api request, consume kafka topics and store data to the userFundsCache
func (w *Wallet) Run() {
	go func() {
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("[ERROR] Wallet server run failed")
		}
	}()
	go w.consumer.ConsumeMessage(w.dbzMess)
	for {
		msg := <-w.dbzMess
		var dbz debezium.DebeziumMsg
		err := json.Unmarshal(msg.Value, &dbz)
		if err != nil {
			log.Println("[ERROR] Unmarshalling debezium message failed", err.Error())
			continue
		}
		if dbz.Payload.Op == "d" {
			var user DbUser
			err := json.Unmarshal(dbz.Payload.Before, &user)
			if err != nil {
				log.Println("[ERROR] Unmarshalling debezium player failed", err.Error())
				continue
			}
			delete(w.userFundsCache, user.Username)
		}
		var user DbUser
		err = json.Unmarshal(dbz.Payload.After, &user)
		if err != nil {
			log.Println("[ERROR] Unmarshalling debezium player failed", err.Error())
			continue
		}
		w.userFundsCache[user.Username] = user.Funds

	}

}

//Stop gracefully stop http server, kafka consumer and database connection
func (w *Wallet) Stop() {
	if err := w.server.Shutdown(context.Background()); err != nil {
		log.Println("[ERROR] Wallet server shutdown failed")
	}
	err := w.consumer.Stop()
	if err != nil {
		log.Fatalln("[ERROR] Could not close consumer gracefully:", err.Error())
	}
	err = w.db.Close()
	if err != nil {
		log.Fatalln("[ERROR] Could not close db gracefully:", err.Error())
	}
	log.Println("[INFO] Wallet stopped gracefully")
}
