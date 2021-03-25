package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dollars float32

type Item struct {
	Name  string  `bson:"name"`
	Price dollars `bson:"price"`
}

const (
	//mongodbEndpoint = "mongodb://172.17.0.2:27017" // Find this from the Mongo container
        mongodbEndpoint = "mongodb://mongo:27017"
)

func main() {
	http.HandleFunc("/list", list)
	http.HandleFunc("/price", price)
	http.HandleFunc("/add", add)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

func list(w http.ResponseWriter, req *http.Request) {

	client, col := dbconnection()
	defer disconnect(client)

	var results []Item
	//findOptions := options.Find()
	cur, _ := col.Find(context.TODO(), bson.D{{}})
	for cur.Next(context.TODO()) {
		//fmt.Println("hi")
		// create a value into which the single document can be decoded
		var elem Item
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(elem)
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())
	for _, item := range results {
		fmt.Println(item)
		fmt.Fprintf(w, "%s: %s\n", item.Name, item.Price)
	}

}

func price(w http.ResponseWriter, req *http.Request) {

	client, col := dbconnection()
	defer disconnect(client)

	item := req.URL.Query().Get("item")
	var result Item
	filter := bson.D{{"name", item}}
	err := col.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "no such item: %q %s\n", item, err)
	} else {
		fmt.Fprintf(w, "%s\n", result.Price)
	}
}

func add(w http.ResponseWriter, req *http.Request) {
	client, col := dbconnection()
	defer disconnect(client)

	item := req.URL.Query().Get("item")
	var result Item
	filter := bson.D{{"name", item}}
	err := col.FindOne(context.TODO(), filter).Decode(&result)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "Item Already Present\n")
		return
	}

	price := req.URL.Query().Get("price")
	if s, err := strconv.ParseFloat(price, 32); err == nil {
		product := Item{Name: item, Price: dollars(s)}
		fmt.Println(product)
		_, err = col.InsertOne(context.TODO(), product)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(w, "Database operation failed with error %s\n", err)
		} else {
			fmt.Fprintf(w, "%s added to list\n", item)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "Not a Valid Price\n")
	}
}

func update(w http.ResponseWriter, req *http.Request) {
	client, col := dbconnection()
	defer disconnect(client)

	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	if s, err := strconv.ParseFloat(price, 32); err == nil {
		filter := bson.D{{"name", item}}
		update := bson.D{
			{"$set", bson.D{
				{"price", s},
			}},
		}
		updateResult, err := col.UpdateOne(context.TODO(), filter, update)
		if updateResult.MatchedCount == 0 {
			w.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(w, "no such item: %q\n", item)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(w, "db error: %q\n", err)
			return
		}
		fmt.Fprintf(w, "%s price updated\n", item)
	} else {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "Not a Valid Price\n")
	}

}

func delete(w http.ResponseWriter, req *http.Request) {

	client, col := dbconnection()
	defer disconnect(client)
	item := req.URL.Query().Get("item")
	filter := bson.D{{"name", item}}
	deleteResult, err := col.DeleteMany(context.TODO(), filter)
	if deleteResult.DeletedCount == 0 {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "db error: %q\n", err)
		return
	}

	fmt.Fprintf(w, "%q item deleted from list\n", item)
}

func dbconnection() (*mongo.Client, *mongo.Collection) {
	// Set client options
	clientOptions := options.Client().ApplyURI(mongodbEndpoint)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for your collection
	col := client.Database("store").Collection("items")
	return client, col
}

func disconnect(client *mongo.Client) {
	client.Disconnect(context.TODO())
	fmt.Println("Disconnected to MongoDB!")
}
