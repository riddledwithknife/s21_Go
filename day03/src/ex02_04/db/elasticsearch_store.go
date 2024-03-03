package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"day03/ex02_04/db/types"

	"github.com/elastic/go-elasticsearch/v8"
)

type Store interface {
	GetPlaces(ctx context.Context, limit, offset int) ([]types.Place, int, error)
	GetClosestPlaces(ctx context.Context, lat float64, lon float64, limit int) ([]types.Place, error)
}

type ElasticsearchStore struct {
	client *elasticsearch.Client
}

func NewElasticsearchStore(client *elasticsearch.Client) Store {
	return &ElasticsearchStore{client: client}
}

func (es *ElasticsearchStore) GetPlaces(ctx context.Context, limit int, offset int) ([]types.Place, int, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"from":  offset,
		"size":  limit,
		"query": map[string]interface{}{"match_all": struct{}{}},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("Error encoding query: %s", err)
		return nil, 0, err
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex("places"),
		es.client.Search.WithBody(&buf),
		es.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Printf("Error executing search query: %s", err)
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("Error parsing the response body: %s", err)
			return nil, 0, err
		}
		errMsg := fmt.Sprintf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
		log.Print(errMsg)
		return nil, 0, fmt.Errorf(errMsg)
	}

	var r struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source types.Place `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return nil, 0, err
	}

	places := make([]types.Place, len(r.Hits.Hits))
	for i, hit := range r.Hits.Hits {
		places[i] = hit.Source
	}

	return places, r.Hits.Total.Value, nil
}

func (es *ElasticsearchStore) GetClosestPlaces(ctx context.Context, lat float64, lon float64, limit int) ([]types.Place, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"size": limit,
		"sort": []map[string]interface{}{
			{
				"_geo_distance": map[string]interface{}{
					"location":        map[string]float64{"lat": lat, "lon": lon},
					"order":           "asc",
					"unit":            "km",
					"distance_type":   "arc",
					"ignore_unmapped": true,
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("Error encoding query: %s", err)
		return nil, err
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex("places"),
		es.client.Search.WithBody(&buf),
	)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("Error parsing the response body: %s", err)
			return nil, err
		}
		errMsg := fmt.Sprintf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
		log.Print(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	var r struct {
		Hits struct {
			Hits []struct {
				Source types.Place `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		return nil, err
	}

	places := make([]types.Place, len(r.Hits.Hits))
	for i, hit := range r.Hits.Hits {
		places[i] = hit.Source
	}

	return places, nil
}
