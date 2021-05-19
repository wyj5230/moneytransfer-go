package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func Init() {
	router := mux.NewRouter()
	router.HandleFunc("/money-side/quotation", createQuotationHandler).Methods("POST")
	router.HandleFunc("/money-side/transaction/{quotationId}", createTransactionHandler).Methods("POST")
	router.HandleFunc("/money-side/transaction/{id}/confirm", confirmTransactionHandler).Methods("POST")
	router.HandleFunc("/money-side/transaction/{id}", getTransactionHandler).Methods("GET")
	router.HandleFunc("/money-side/callback", callbackHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func createQuotationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestBody, requestBodyErr := ioutil.ReadAll(r.Body)
	if requestBodyErr != nil {
		log.Printf("invalid requestBody %s", requestBodyErr)
	}
	var quotationRequest QuotationRequest
	jsonErr := json.Unmarshal(requestBody, &quotationRequest)
	if jsonErr != nil {
		log.Printf("QuotationRequest err: %s", jsonErr)
	}
	encodeErr := json.NewEncoder(w).Encode(CreateQuotation(quotationRequest))
	if encodeErr != nil {
		log.Printf("CreateQuotation encode err: %s", encodeErr)
	}
}

func createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	requestBody, requestBodyErr := ioutil.ReadAll(r.Body)
	if requestBodyErr != nil {
		log.Printf("invalid requestBody %s", requestBodyErr)
	}
	params := mux.Vars(r)
	var transactionRequest TransactionRequest
	jsonErr := json.Unmarshal(requestBody, &transactionRequest)
	if jsonErr != nil {
		log.Printf("TransactionRequest err: %s", jsonErr)
	}
	encodeErr := json.NewEncoder(w).Encode(createTransaction(params["quotationId"], transactionRequest))
	if encodeErr != nil {
		log.Printf("CreateTransaction encode err: %s", encodeErr)
	}
}

func confirmTransactionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	encodeErr := json.NewEncoder(w).Encode(confirmTransaction(params["id"]))
	if encodeErr != nil {
		log.Printf("ConfirmTransaction encode err: %s", encodeErr)
	}
}

func getTransactionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	encodeErr := json.NewEncoder(w).Encode(getTransaction(params["id"]))
	if encodeErr != nil {
		log.Printf("GetTransaction encode err: %s", encodeErr)
	}
}

func callbackHandler(responseWriter http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%s\n", string(body))
		HandlerCallback(body)
	}
	io.WriteString(responseWriter, "Got it\n")
}
