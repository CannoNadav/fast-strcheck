package main

import (
	"flag"
	"fmt"
	"log"
	"main/programargs"
	"main/server"
)

const WORDS_CATALOG_FILE = "words_clean.txt"

func main() {
	flag.Parse()
	_, wordsCatalogFactory, err := programargs.GetWordsCatalogFactory(func(s string) { fmt.Println(s) })
	if err == nil {
		server.InitAndStartServer(wordsCatalogFactory, WORDS_CATALOG_FILE)
	} else {
		log.Println(err.Error())
	}
}
