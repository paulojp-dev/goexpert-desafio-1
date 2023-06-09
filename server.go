package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"
)

func main() {
	http.HandleFunc("/cotacao", GetDollarQuotationHandler)
	http.ListenAndServe(":8080", nil)
}

func GetDollarQuotationHandler(response http.ResponseWriter, request *http.Request) {
	quotation, error := GetDollarQuotation()
	if error != nil {
		response.WriteHeader(http.StatusInternalServerError)
		result := map[string]string{"error": error.Error()}
		json.NewEncoder(response).Encode(result)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	result := map[string]float64{"value": quotation}
	json.NewEncoder(response).Encode(result)
}

func GetDollarQuotation() (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	request, error := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if error != nil {
		return 0, error
	}
	response, error := http.DefaultClient.Do(request)
	if error != nil {
		return 0, error
	}
	defer response.Body.Close()
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		return 0, error
	}
	result := gjson.Get(string(body), "USDBRL.bid")
	if error != nil {
		return 0, error
	}
	error = persistQuotation(result.Float())
	if error != nil {
		return 0, error
	}
	return result.Float(), nil
}

func persistQuotation(quotation float64) error {
	db, error := prepareDB()
	if error != nil {
		return error
	}
	defer db.Close()
	stmt, error := db.Prepare("insert into cotations values(NULL, 'USD', $1, $2)")
	if error != nil {
		return error
	}
	defer stmt.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, error = stmt.ExecContext(ctx, quotation, time.Now())
	if error != nil {
		return error
	}
	return nil
}

func prepareDB() (*sql.DB, error) {
	db, error := sql.Open("sqlite3", "serverDB.db")
	if error != nil {
		return nil, error
	}
	createTable := `
		CREATE TABLE IF NOT EXISTS cotations (
			id INTEGER NOT NULL PRIMARY KEY,
			currency TEXT NOT NULL,
			value REAL NOT NULL,
			time DATETIME NOT NULL
		);
	`
	_, error = db.Exec(createTable)
	if error != nil {
		return nil, error
	}
	return db, nil
}
