package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func main() {
	http.HandleFunc("/cotacao", GetDolarCotationHandler)
	http.ListenAndServe(":8080", nil)

	// TODO: 10ms máximo para salvar no banco de dados SQLite
}

func GetDolarCotationHandler(response http.ResponseWriter, request *http.Request) {
	cotation, error := GetDolarCotation()
	if error != nil {
		response.WriteHeader(http.StatusInternalServerError)
		result := map[string]string{"error": error.Error()}
		json.NewEncoder(response).Encode(result)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	result := map[string]float64{"value": cotation}
	json.NewEncoder(response).Encode(result)
}

func GetDolarCotation() (float64, error) {
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
	return result.Float(), nil
}
