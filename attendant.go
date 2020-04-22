package main

import (
	"fmt"
	"crypto/tls"
	"log"
	"net/http"
)

type handler struct{
	Orders chan string
}

//func handler(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "hello")
//}

func (h *handler) ServeHTTP (w http.ResponseWriter, req *http.Request) {
	if (req.URL.Path[1:] == "heartbeat") {
		w.Write([]byte("alive\n"))
	}else if (req.URL.Path[1:] == "request"){
		fmt.Println("received request")
		h.Orders<-"orderup"	
	}else{
		w.Write([]byte("wut\n"))
	}
}

func run_attendant (orders chan string) {
	for{
		select {
		case o := <-orders:
			fmt.Println("received order for: " + o)
		}
	}
}

func main() {
	order_chan := make(chan string)
	srv := &http.Server{
                Addr: ":8080",
                Handler: &handler{Orders: order_chan},
                TLSConfig: &tls.Config{},
        }
	go run_attendant(order_chan)
        log.Fatal(srv.ListenAndServeTLS("/etc/letsencrypt/live/h3.testmyprotocol.com/fullchain.pem", "/etc/letsencrypt/live/h3.testmyprotocol.com/privkey.pem"))
}
