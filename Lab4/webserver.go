package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/add", db.add)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

//RW : Read Write lock for DB
var RW sync.RWMutex

func (db database) list(w http.ResponseWriter, req *http.Request) {
	RW.RLock()
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
	RW.RUnlock()
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	RW.RLock()
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
	RW.RUnlock()
}

func (db database) add(w http.ResponseWriter, req *http.Request) {
	RW.Lock()
	item := req.URL.Query().Get("item")
	if _, ok := db[item]; ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "Item Already Present\n")
		RW.Unlock()
		return
	}
	price := req.URL.Query().Get("price")
	if s, err := strconv.ParseFloat(price, 32); err == nil {
		db[item] = dollars(s)
		fmt.Fprintf(w, "%s added to list\n", item)
	} else {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "Not a Valid Price\n")
	}
	RW.Unlock()
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	RW.Lock()
	item := req.URL.Query().Get("item")
	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusBadRequest) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		RW.Unlock()
		return
	}
	price := req.URL.Query().Get("price")
	if s, err := strconv.ParseFloat(price, 32); err == nil {
		db[item] = dollars(s)
		fmt.Fprintf(w, "%s price updated\n", item)
	} else {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "Not a Valid Price\n")
	}
	RW.Unlock()
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	RW.Lock()
	item := req.URL.Query().Get("item")
	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusBadRequest) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		RW.Unlock()
		return
	}
	delete(db, item)
	fmt.Fprintf(w, "%q item deleted from list\n", item)
	RW.Unlock()
}
