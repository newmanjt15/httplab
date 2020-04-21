package main

import (
	"fmt"
	"io/ioutil"
	"log"
	// "net"
	"net/http"
	"crypto/tls"
	// "crypto/x509"
)

type Page struct {
	Title string
	Body []byte
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

func forwardOrder (order string, a *Attendant){
	_, err :=http.Get("https://" + a.domain + ":8080/request")
	if err != nil {
		return
	}

}

func runMaster (orders chan string) {
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
		}
	}

	fmt.Printf("master is finished")
}

func main() {
	order_chan := make(chan string)
	srv := &http.Server{
		Addr: ":443",
		Handler: &handler{Orders: order_chan},
		TLSConfig: &tls.Config{},
	}
	go runMaster(order_chan)
	log.Fatal(srv.ListenAndServeTLS("/etc/letsencrypt/live/testmyprotocol.com/fullchain.pem", "/etc/letsencrypt/live/testmyprotocol.com/privkey.pem"))
}
