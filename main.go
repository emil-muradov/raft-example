package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func main() {
	port := flag.String("port", "8080", "Specifies a tcp port that server listens to")
	flag.Parse()
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	nodeIpAddr := GetOutboundIP()
	nodeFullAddr := fmt.Sprintf("%s:%s", nodeIpAddr, *port)
	rdb.SAdd(ctx, "cluster:node:ip", nodeFullAddr)
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	})
	log.Println("App started")
	http.ListenAndServe(fmt.Sprintf(":%s", *port), router)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80") // Connect to a known external address (Google DNS)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddress := conn.LocalAddr().(*net.UDPAddr)
	return localAddress.IP
}
