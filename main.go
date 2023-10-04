package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	serverPort = 8080
)

func LuhnAlgo(cardNo string) bool {
	nDigits := len(cardNo)
	nSum := 0
	isSecond := false

	for i := nDigits - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(cardNo[i]))
		if isSecond {
			digit = digit * 2
			if digit > 9 {
				digit = digit - 9
			}
		}

		nSum += digit
		isSecond = !isSecond
	}

	return nSum%10 == 0
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("server: %s/\n", r.Method)
	fmt.Printf("server: content-type: %s\n", r.Header.Get("content-type"))
	fmt.Printf("server: headers:\n")
	for headerName, headerValue := range r.Header {
		fmt.Printf("\t%s = %s\n", headerName, strings.Join(headerValue, ", "))
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("server: could not read request body: %s\n", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		fmt.Printf("server: could not unmarshal request body: %s\n", err)
	}

	ccNumber := data["credit_card"]

	isValid := LuhnAlgo(ccNumber.(string))

	if isValid {
		fmt.Fprintf(w, `{"message": "credit card is valid!"}`)
	} else {
		fmt.Fprintf(w, `{"message": "credit card is invalid!"}`)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)

	ctx, cancelCtx := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server is closed")
		} else if err != nil {
			fmt.Printf("error listening server: %s", err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}
