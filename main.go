package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	"time"
)

var mySignKey = "Menyasosali"

type Funds struct {
	Id               int             `json:"id"`
	Name             string          `json:"name"`
	Ticket           string          `json:"ticket"`
	Amount           decimal.Decimal `json:"amount"`
	PricePerItem     decimal.Decimal `json:"priceperitem"`
	PurchasePrice    decimal.Decimal `json:"purchaseprice"`
	PriceCurrent     decimal.Decimal `json:"pricecurrent"`
	PercentChanges   decimal.Decimal `json:"percentchanges"`
	YearlyInvestment decimal.Decimal `json:"yearlyinvestment"`
	ClearMoney       decimal.Decimal `json:"clearmoney"`
	DataPurchase     time.Time       `json:"datapurchase"`
	DataLastUpdate   time.Time       `json:"datalastupdate"`
	Type             string          `json:"type"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var usertest = User{
	Username: "1",
	Password: "1",
}

func main() {
	fmt.Println("My REST Server")
	r := mux.NewRouter()
	r.Handle("/funds/usd/shares", checkAuth(getUSDFundsShares)).Methods("GET")

	r.HandleFunc("/login", login).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	// fmt.Println("user:", u)
	checkLogin(u)
}

func checkLogin(u User) string {
	if usertest.Username != u.Username || usertest.Password != u.Password {
		fmt.Println("Not currect")
		err := "error"
		return err
	}

	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Println(err)
	}

	return validToken
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour * 1000)
	claims["user"] = "Menyasosali"
	claims["authorized"] = true

	tokenString, err := token.SignedString(mySignKey)

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func checkAuth(endpoint func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Toxen"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySignKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func getUSDFundsShares(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var AllShares = myCurrentFunds("share")
	json.NewEncoder(w).Encode(AllShares)
}

func myCurrentFunds(fundType string) []Funds {
	var amountShares []Funds

	db, err := sql.Open("postgres", "postgres://postgres:parol@localhost/fin?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * from fundsusd WHERE type = $1", fundType)

	for rows.Next() {
		f := Funds{}
		err := rows.Scan(&f.Id, &f.Name, &f.Ticket, &f.Amount, &f.PricePerItem, &f.PurchasePrice, &f.PriceCurrent, &f.PercentChanges, &f.YearlyInvestment, &f.ClearMoney, &f.DataPurchase, &f.DataLastUpdate, &f.Type)

		if err != nil {
			log.Fatal(err)
		}
		amountShares = append(amountShares, f)
	}
	return amountShares
}
