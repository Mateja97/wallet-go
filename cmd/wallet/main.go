package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"wallet-go/wallet"

	"github.com/Shopify/sarama"
)

var port = flag.String("port", ":8080", "")
var dbHost = flag.String("db.host", "", "Database host address")
var dbName = flag.String("db.name", "", "Database table name")
var dbPort = flag.String("db.port", "", "Database port")
var dbUsr = flag.String("db.usr", "", "Database username")
var dbPw = flag.String("db.pw", "", "Database password")
var kafkaBrokers = flag.String("kafka.brokers", "", "Ip address of kafka broker")
var kafkaTopic = flag.String("kafka.topic", "", "Topic of data")

func main() {

	flag.Parse()
	brokers := strings.Split(*kafkaBrokers, ",")
	dbzMess := make(chan *sarama.ConsumerMessage)

	w := wallet.Wallet{}
	err := w.Init(*port, *dbUsr, *dbPw, *dbHost, *dbPort, *dbName, *kafkaTopic, brokers, dbzMess)
	if err != nil {
		log.Fatalln("[ERROR] Wallet init failed, error: ", err.Error())
	}
	go w.Run()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done
	w.Stop()
}
