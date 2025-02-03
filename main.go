package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Define the structure data type
type Purchase struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`           // Name of supplier
	Describe       string `json:"describe"`       // Product
	Reference      string `json:"reference"`      // Product reference
	PeriodEntrance string `json:"periodentrance"` // Quantity of products was arrived
	PeriodSale     string `json:"periodsale"`     // Quantity of sales
	StockCurrent   string `json:"stockcurrent"`   // Many products available
	Status         string `json:"status"`         // Available, Critical and Empty
	Quantity       string `json:"quantity"`       // Quantity per purchase
}

type PurchaseService struct { // Define the map of structure
	Purchases map[int]Purchase
}

// Implement math request for JSON
// type MathPurchase struct {
// 	PeriodEntrance int
// 	PeriodSale int
// 	StockCurrent int

// 	StockCurrent := (PeriodEntrance - PeriodSale)

// }

// Implement the CRUD (Create, Read, Update and Delete)
// Create
func (p *PurchaseService) Create(w http.ResponseWriter, r *http.Request) {
	var purchase Purchase
	err := json.NewDecoder(r.Body).Decode(&purchase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := len(p.Purchases) + 1
	purchase.Id = id

	p.Purchases[id] = purchase
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(purchase)
	w.WriteHeader(http.StatusCreated)
}

// Read
func (p *PurchaseService) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var purchases []Purchase

	for _, ps := range p.Purchases {
		purchases = append(purchases, ps)
	}

	json.NewEncoder(w).Encode(purchases)
}

// Read
func (p *PurchaseService) Get(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	if val, ok := p.Purchases[id]; ok {
		json.NewEncoder(w).Encode(val)
	} else {
		http.Error(w, "Purchase not found", http.StatusNotFound)
	}
}

// Delete
func (p *PurchaseService) Delete(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	if _, ok := p.Purchases[id]; ok {
		delete(p.Purchases, id)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Purchase not found", http.StatusNotFound)
	}

}

// Update
func (p *PurchaseService) Update(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	var purchase Purchase
	err := json.NewDecoder(r.Body).Decode(&purchase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := p.Purchases[id]; ok {
		purchase.Id = id
		p.Purchases[id] = purchase
	} else {
		http.Error(w, "Purchase not found", http.StatusNotFound)
	}
}

// Request of CRUD
// Request Update
func handleUpdatePurchase(w http.ResponseWriter, r *http.Request, service *PurchaseService) {
	q := r.URL.Query()
	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Update(w, r, id)
	} else {
		http.Error(w, "Purchase not found", http.StatusNotFound)
	}
}

// Request Delete
func handleDeletePurchase(w http.ResponseWriter, r *http.Request, service *PurchaseService) {
	q := r.URL.Query()
	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Delete(w, r, id)
	} else {
		http.Error(w, "Purchase not found", http.StatusNotFound)
	}
}

// Request Read
func handleGetPurchase(w http.ResponseWriter, r *http.Request, service *PurchaseService) {
	q := r.URL.Query()
	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Get(w, r, id)
	} else {
		service.List(w, r)
	}
}

// Request Create
func handleCreatePurchase(w http.ResponseWriter, r *http.Request, service *PurchaseService) {
	service.Create(w, r)
}

// The function main to run server CRUD
func main() {
	service := &PurchaseService{Purchases: make(map[int]Purchase)}
	mux := http.NewServeMux()

	mux.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetPurchase(w, r, service)
		case http.MethodPost:
			handleCreatePurchase(w, r, service)
		case http.MethodDelete:
			handleDeletePurchase(w, r, service)
		case http.MethodPut:
			handleUpdatePurchase(w, r, service)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	address := "localhost:8080"

	log.Printf("Starting server on %s", address)

	err := http.ListenAndServe(address, mux)
	if err != nil {
		log.Fatal(err)
	}
}
