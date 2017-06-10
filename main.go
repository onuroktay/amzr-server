package main

import (
	"fmt"
	"github.com/onuroktay/amazon-reader/amzr-server/elasticsearch"
	"log"
	"net/http"
)

const (
	//USER ROLE
	USER = 1
	// EDITOR ROLE
	EDITOR = 2
	// ADMIN ROLE
	ADMIN = 3
)

var (
	// path is a string
	path     = "/Users/onur/Go/src/github.com/onuroktay/amazon-reader/AmzR-Client/dist"
	// certPath is a string
	certPath  = "/Users/onur/Go/src/github.com/onuroktay/amazon-reader/AmzR-Server/"
	database *DATABASE
)

func main() {
	// Connect to ElasticSearch
	es, err := OnurTPIES.NewElasticSearch("amazonreader")

	if err != nil {
		log.Fatalln("ElasticSearch connection error:", err.Error())
	}

	// Set Database
	database = &DATABASE{accesser: es}

	fmt.Print("Server started on port :8080\n")

	routes()

	// err = http.ListenAndServe(":8080", nil)nprivate
	err = http.ListenAndServeTLS(":8080", certPath + "server.pem", certPath + "server.key", nil)
	if err != nil {
		fmt.Println(err.Error())
	}


}
