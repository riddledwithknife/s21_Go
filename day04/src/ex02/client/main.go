package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getClient() *http.Client {
	data, _ := os.ReadFile("../ca/minica.pem")
	cp, _ := x509.SystemCertPool()
	cp.AppendCertsFromPEM(data)

	cert, _ := tls.LoadX509KeyPair("../ca/client/cert.pem", "../ca/client/key.pem")

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      cp,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
		Timeout: 30 * time.Second,
	}

	return client
}

func main() {
	fk := flag.String("k", "", "Two letter abbreviation for candy to buy (supports only this ones: CE, AA, NT, DE, YR)")
	fc := flag.String("c", "", "Count of candy to buy")
	fm := flag.String("m", "", "Amount of money to pass to \"machine\"")
	flag.Parse()

	if *fk == "" || *fc == "" || *fm == "" {
		log.Fatalln("Can't use not every flag")
	}

	money, _ := strconv.Atoi(*fm)
	count, _ := strconv.Atoi(*fc)

	if money <= 0 {
		log.Fatalln("Amount of money must be positive")
	}

	switch *fk {
	case "CE", "AA", "NT", "DE", "YR":

	default:
		log.Fatalln("Invalid candy type")
	}

	if count <= 0 {
		log.Fatalln("Amount of candy must be positive")
	}

	jsonBody := []byte(fmt.Sprintf(`{"money": %d, "candyType": "%s", "candyCount": %d}`, money, *fk, count))
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, "https://localhost:3333/buy_candy", bodyReader)
	if err != nil {
		log.Fatalf("client: could not create request: %s\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := getClient()
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request error: %s", err)
	}

	var response map[string]any
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Fatalf("Error decoding response: %s\n", err)
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %s\n", err)
	}
	fmt.Println(string(jsonResponse))
}
