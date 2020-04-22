package main

import (
	"fmt"
	"io/ioutil"
	"log"
	// "net"
	"net/http"
	// "crypto/tls"
	// "crypto/x509"
	"strings"
)

type Page struct {
	Title string
	Body []byte
}

type Order struct {
	URL string
	status bool
	finished bool
	loc string
	result string
}


type Attendant struct {
	IP				string
	domain			string
	http_version	string
	busy			bool
	active			bool
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	str, err := loadPage(r.URL.Path[1:])
// 	if err != nil {
// 		fmt.Fprint(w, "Error")
// 	}
// 	fmt.Fprint(w, "%s", str)
// }

type handler struct{
	Orders chan string
}

func (h *handler) ServeHTTP (w http.ResponseWriter, req *http.Request) {
	// path := req.URL.Path[1:]
	
	w.Write([]byte("hello\n"))
	h.Orders<-"hello"
}

func pingAttendant (a *Attendant) {
	resp, err :=http.Get("https://" + a.domain + ":8080/heartbeat")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(responseData)
	fmt.Println(responseString)
	if (responseString == "alive\n"){
		a.active = true
	}

}

func checkOrder (order string, a *Attendant){
result, err :=http.Get("https://" + a.domain + ":8080/" + order)
	if err != nil {
		return
	}
	fmt.Println("check complete")
	fmt.Println(result)
}

func forwardOrder (order string, a *Attendant){
	result, err :=http.Get("https://" + a.domain + ":8080/" + order)
	if err != nil {
		return
	}
	fmt.Println("order complete")
	fmt.Println(result)
}

func runMaster (orders chan string, check_chan chan string) {
	// list of attendants
	var attendants [2]*Attendant
	// fill ones we know of now
	attendants[0] = &Attendant{IP: "3.85.192.243", domain: "h1.testmyprotocol.com", http_version: "h1", busy: false, active: false}
	attendants[1] = &Attendant{IP: "3.86.81.22", domain: "h3.testmyprotocol.com", http_version: "h3", busy: false, active: false}
	pingAttendant(attendants[0])
	pingAttendant(attendants[1])
	for _, attendant := range attendants {
		if attendant.active {
			fmt.Println(attendant.domain + " is active!")
		}
	}

	for {
		select {
		case c:= <-orders:
				fmt.Println("Received order for: " + c)
				forwardOrder(c, attendants[1])
		case c:= <-check_chan:
				fmt.Println("received check for: " + c)
				checkOrder(c, attendants[1])
		}
	}

	fmt.Printf("master is finished")
}

func main() {
	order_chan := make(chan string)
	check_chan := make(chan string)
	var orders [10]Order
	x := 0
	// srv := &http.Server{
	// 	Addr: ":8080",
	// 	Handler: &handler{Orders: order_chan},
	// 	TLSConfig: &tls.Config{},
	// }
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "incoming") {
			fmt.Println("received data")
			fmt.Println(r.URL.Path)
			// got result, store it for now
			parts := strings.Split(r.URL.Path, "-")
			fmt.Println(parts)
			for i, n := range orders {
				if parts[1] == n.URL {
					fmt.Println("updating order")
					fmt.Println(parts[2])
					orders[i].loc = parts[2]
				}
			}
		}else{
			url := r.URL.Path[1:]
			for _, n := range orders {
				if url == n.URL {
					w.Write([]byte(n.loc[1:]))
					return
				}
			}
			orders[x] = Order{URL: url, loc: "nowhere"}
			x += 1
			order_chan <- url
		}
	})
	go runMaster(order_chan, check_chan)
	// log.Fatal(srv.ListenAndServeTLS("/etc/letsencrypt/live/testmyprotocol.com/fullchain.pem", "/etc/letsencrypt/live/testmyprotocol.com/privkey.pem"))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
