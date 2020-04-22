package main

import (
	"fmt"
	"crypto/tls"
	"log"
	"net/http"
	"os/exec"
//	"strings"
	"bytes"
)

type handler struct{
	Orders chan string
	URLs   chan string
	Locations chan string
}

type Order struct{
	URL string
	finished bool
	location string
}

//func handler(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "hello")
//}

func (h *handler) ServeHTTP (w http.ResponseWriter, req *http.Request) {
	if (req.URL.Path[1:] == "heartbeat") {
		w.Write([]byte("alive\n"))
	}else if (req.URL.Path[1:] == "request"){
		fmt.Println("received request")
		//h.Orders<-"orderup"	
		sites, ok := req.URL.Query()["site"]
		fmt.Println(req.URL)
		fmt.Println(sites)
		if !ok || len(sites[0]) < 1 {
			w.Write([]byte("error getting request"))
		}else{
			h.URLs <- sites[0] 
			loc := <- h.Locations
			fmt.Println("location: " + loc)
			w.Write([]byte(loc))
		}
	}else{
		//we just got a url, try to clone it
		h.URLs <- req.URL.Path[1:]
		status := <- h.Locations
		fmt.Println(status)
		if status == "no order" {
			h.Orders<-req.URL.Path[1:]
			//url := <-h.URLs
			//fmt.Println(url)
			w.Write([]byte("submitted"))
		}else{
			w.Write([]byte(status))
		}
	}
}

func run_attendant (orders chan string, urls chan string, order_list [10]Order, locations chan string) {
	x := 0
	for{
		select {
		case o := <-orders:
			fmt.Println("received order for: " + o)
			if x < 10 {
				order_list[x] = Order{URL: o, finished: false, location: "unknown"}		
				cmd := exec.Command("sudo", "bash", "httplab/clone_website.sh", "https://" + o)
				//cmd.Stdin = strings.NewReader("https://" + o)
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					fmt.Println(err)
				}	
				url, err2 := out.ReadString('\n')
				if err2 != nil {
					fmt.Println(err2)
				}	
				fmt.Println(url)
				order_list[x].finished = true
				order_list[x].location = url
				_, err =http.Get("http://54.91.253.3:8080/incoming-" + o + "-" + url)
				if err != nil {				
					fmt.Println("error sending data back")
					fmt.Println(err)
				}
				x += 1
			}
		case r := <-urls:
			for _, n := range order_list{
				if n.URL == r {
					locations <- n.location
				}
			}					
			locations <- "no order"
		}
	}
}

func main() {
	order_chan := make(chan string)
	url_chan := make(chan string)
	location_chan := make(chan string)
	var orders [10]Order
	srv := &http.Server{
                Addr: ":8080",
                Handler: &handler{Orders: order_chan, URLs: url_chan, Locations: location_chan},
                TLSConfig: &tls.Config{},
        }
	go run_attendant(order_chan, url_chan, orders, location_chan)
        log.Fatal(srv.ListenAndServeTLS("/etc/letsencrypt/live/h3.testmyprotocol.com/fullchain.pem", "/etc/letsencrypt/live/h3.testmyprotocol.com/privkey.pem"))
}
