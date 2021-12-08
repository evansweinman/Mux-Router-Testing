package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Item struct {
	UID   string  `json:"UID"`
	Name  string  `json:"Name"`
	Desc  string  `json:"Desc"`
	Price float64 `json:"Price"`
}

var inventory []Item

//write a response with the HTTP library that we imported into our main.go file
func homePage(w http.ResponseWriter, r *http.Request) {
	//I believe Fprintf writes to browser
	fmt.Fprintf(w, "Endpoints called: homePage()")
}

func getInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("FunctionCalled: getInventory()")

	//takes writer ojbect that we passed from our mux router and encodes it to pop up in browser
	json.NewEncoder(w).Encode(inventory)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Item
	//send a request for the data
	_ = json.NewDecoder(r.Body).Decode(&item)

	//appending item to our slice data inventory
	inventory = append(inventory, item)

	//whatever item datat we have is going to render in teh data we return back as a response
	json.NewEncoder(w).Encode(item)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	_deleteItemAtUid(params["uid"])

	json.NewEncoder(w).Encode(inventory)
}

func _deleteItemAtUid(uid string) {
	for index, item := range inventory {
		if item.UID == uid {
			// Delete item from slice
			inventory = append(inventory[:index], inventory[index+1:]...)
			break
		}
	}
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)

	params := mux.Vars(r)

	// Delete item at UID
	_deleteItemAtUid(params["uid"])
	// Create it with new data
	inventory = append(inventory, item)

	json.NewEncoder(w).Encode(inventory)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage).Methods("GET")

	// Create a route to view all items in inventory
	// localhost:9000/inventory
	router.HandleFunc("/inventory", getInventory).Methods("GET")

	//Posting typically allows us to add data on an HTTP server
	router.HandleFunc("/inventory", createItem).Methods("POST")

	//this route will get an ID proeprty of uid
	router.HandleFunc("/inventory/{uid}", deleteItem).Methods("DELETE")

	//Use PUT to update items
	router.HandleFunc("/inventory/{uid}", updateItem).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	inventory = append(inventory, Item{
		UID:   "0",
		Name:  "Cheese",
		Desc:  "A fine block of cheese",
		Price: 4.99,
	})
	inventory = append(inventory, Item{
		UID:   "1",
		Name:  "Milk",
		Desc:  "A jug of milk",
		Price: 3.25,
	})
	handleRequests()
}
