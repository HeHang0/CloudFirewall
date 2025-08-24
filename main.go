package main

import (
	"cloud_firewall/config"
	_ "cloud_firewall/server"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var Version = "dev"

func main() {
	if config.Cfg.Version {
		fmt.Println(Version)
		os.Exit(0)
		return
	}
	addr := config.Cfg.Addr + ":" + strconv.Itoa(config.Cfg.Port)
	log.Println("Listening on " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
