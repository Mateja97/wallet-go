package wallet

import (
	"database/sql"
	"log"
)

type DbUser struct {
	Wid      int    `json:"id"`
	Funds    int    `json:"funds"`
	Username string `json:"username"`
}

func (wl *Wallet) GetDBFunds(user string) (int, error) {

	query := `select funds from wallet where username = $1;`
	var funds int
	err := wl.db.QueryRow(query, user).Scan(&funds)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("[INFO] No player found with username : ", user)
			return 0, err
		}
		log.Println("[ERROR] Query error", err.Error())
		return 0, err
	}
	return funds, nil
}

func (wl *Wallet) AddDBFunds(user string, funds int) error {

	query := `update wallet set funds = $1 where username = $2;`
	err := wl.db.QueryRow(query, funds, user).Err()
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("[INFO] No player found with username : ", user)
			return err
		}
		log.Println("[ERROR] Query error", err.Error())
		return err
	}
	return nil
}
