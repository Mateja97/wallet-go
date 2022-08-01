package wallet

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//GetFunds return user funds from cache
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

//AddFunds function that manipulate with funds (adding/removing)
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
		log.Println("[INFO] New arrived funds: ", p.Funds)
		wl.cacheLock.Lock()
		defer wl.cacheLock.Unlock()
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
		log.Println("[INFO] New current funds: ", currentFunds, " user: ", user)

		w.Write(data)
	}
}
