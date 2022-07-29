package wallet

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (wl *Wallet) GetFunds() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["user"]
		data, err := json.Marshal(wl.userFundsCache[user])
		if err != nil {
			log.Println("[ERROR] GetFunds Marshaling failed")
		}
		w.Write(data)
	}
}

func (wl *Wallet) AddFunds() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["user"]
		var p struct {
			Funds int `json:"funds"`
		}
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("[INFO] New funds: ", p.Funds)
		wl.Lock()
		defer wl.Unlock()
		currentFunds := wl.userFundsCache[user]
		currentFunds += p.Funds
		if currentFunds < 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Not enough funds"))
			return
		}
		err = wl.AddDBFunds(user, currentFunds)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		data, err := json.Marshal(currentFunds)
		if err != nil {
			log.Println("[ERROR] GetFunds Marshaling failed")
		}

		w.Write(data)
	}
}
