package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func ServeCallback() {
	handler := func(responseWriter http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Printf("%s\n", string(body))
			HandlerCallback(body)
		}
		io.WriteString(responseWriter, "Got it\n")
	}
	http.HandleFunc("/moneytransfer/callback", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
