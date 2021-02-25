package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type dollars float32
type database struct {
	DB map[string]dollars
	RW sync.RWMutex
}

func main() {
	db := database{DB: make(map[string]dollars)}
	db.DB["shoes"] = dollars(50)
	db.DB["sock"] = dollars(5)
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/add", db.add)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	db.RW.RLock()
	defer db.RW.RUnlock()
	for item, price := range db.DB {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}

}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	db.RW.RLock()
	defer db.RW.RUnlock()
	item := req.URL.Query().Get("item")
	if price, ok := db.DB[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db *database) add(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if _, ok := db.DB[item]; ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "Item Already Present\n")
		return
	}
	price := req.URL.Query().Get("price")
	if s, err := strconv.ParseFloat(price, 32); err == nil {
		db.RW.Lock()
		db.DB[item] = dollars(s)
		db.RW.Unlock()
		fmt.Fprintf(w, "%s added to list\n", item)
	} else {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "Not a Valid Price\n")
	}
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if _, ok := db.DB[item]; !ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	price := req.URL.Query().Get("price")
	if s, err := strconv.ParseFloat(price, 32); err == nil {
		db.RW.Lock()
		db.DB[item] = dollars(s)
		db.RW.Unlock()
		fmt.Fprintf(w, "%s price updated\n", item)
	} else {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "Not a Valid Price\n")
	}
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if _, ok := db.DB[item]; !ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	db.RW.Lock()
	delete(db.DB, item)
	db.RW.Unlock()
	fmt.Fprintf(w, "%q item deleted from list\n", item)
}
