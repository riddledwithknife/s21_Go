package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Place struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
}

func createClient() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return es
}

func createIndex(es *elasticsearch.Client) {
	mapping := `{
        "mappings": {
            "properties": {
                "name": { "type": "text" },
                "address": { "type": "text" },
                "phone": { "type": "text" },
                "location": { "type": "geo_point" }
            }
        }
    }`

	req := esapi.IndicesCreateRequest{
		Index: "places",
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error response: %s", res.String())
	} else {
		log.Println("Index created successfully")
	}
}

func uploadData(es *elasticsearch.Client, places []Place) {
	var buf strings.Builder
	for _, r := range places {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "places", "_id" : "%d" } }%s`, r.ID, "\n"))
		data, err := json.Marshal(r)
		if err != nil {
			log.Fatalf("Cannot encode article %d: %s", r.ID, err)
		}
		data = append(data, "\n"...)

		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}

	res, err := es.Bulk(strings.NewReader(buf.String()), es.Bulk.WithIndex("places"))
	if err != nil {
		log.Fatalf("Failure indexing batch: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error response: %s", res.String())
	} else {
		log.Println("Batch indexed successfully")
	}
}

func main() {
	es := createClient()

	jsonData, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
	}

	var places []Place
	err = json.Unmarshal(jsonData, &places)
	if err != nil {
		log.Fatalf("Error parsing JSON data: %s", err)
	}

	createIndex(es)
	uploadData(es, places)
}
