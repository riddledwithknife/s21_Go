package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"day03/ex02_04/db"
	"day03/ex02_04/db/types"

	"github.com/dgrijalva/jwt-go"
	"github.com/elastic/go-elasticsearch/v8"
)

const pageSize = 10

var jwtKey = []byte("your_secret_key")

func createClient() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	return es
}

func main() {
	esClient := createClient()
	store := db.NewElasticsearchStore(esClient)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		htmlEndpointHandler(w, r, store)
	})
	http.HandleFunc("/api/places", func(w http.ResponseWriter, r *http.Request) {
		apiPlacesHandler(w, r, store)
	})
	http.HandleFunc("/api/get_token", getTokenHandler)
	http.HandleFunc("/api/recommend", jwtMiddleware(func(w http.ResponseWriter, r *http.Request) {
		recommendHandler(w, r, store)
	}))

	log.Fatal(http.ListenAndServe(":8888", nil))
}

func htmlEndpointHandler(w http.ResponseWriter, r *http.Request, store db.Store) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		http.Error(w, fmt.Sprintf("Invalid 'page' value: 'foo'"), http.StatusBadRequest)
	}

	offset := (page - 1) * pageSize
	places, total, err := store.GetPlaces(ctx, pageSize, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %s", err), http.StatusInternalServerError)
		return
	}

	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	data := struct {
		Places       []types.Place
		Total        int
		CurrentPage  int
		PreviousPage int
		NextPage     int
		LastPage     int
	}{
		Places:       places,
		Total:        total,
		CurrentPage:  page,
		PreviousPage: int(math.Max(float64(page-1), 1)),
		NextPage:     int(math.Min(float64(page+1), float64(lastPage))),
		LastPage:     lastPage,
	}

	tmpl, err := template.New("places").Parse(`
		<!doctype html>
		<html>
		<head>
			<meta charset="utf-8">
			<title>Places</title>
		</head>
		<body>
		<h5>Total: {{.Total}}</h5>
		<ul>
			{{range .Places}}
				<li>
					<div>{{.Name}}</div>
					<div>{{.Address}}</div>
					<div>{{.Phone}}</div>
				</li>
			{{end}}
		</ul>
		{{if gt .CurrentPage 1}}<a href="/?page={{.PreviousPage}}">Previous</a>{{end}}
		{{if lt .CurrentPage .LastPage}}<a href="/?page={{.NextPage}}">Next</a>{{end}}
		</body>
		</html>
	`)

	if err != nil {
		http.Error(w, "Error creating template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %s", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func apiPlacesHandler(w http.ResponseWriter, r *http.Request, store db.Store) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		respondWithError(w, 400, "Invalid 'page' value: 'foo'")
	}

	offset := (page - 1) * pageSize
	places, total, err := store.GetPlaces(ctx, pageSize, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching data: %s", err))
		return
	}

	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	response := struct {
		Places       []types.Place `json:"places"`
		Total        int           `json:"total"`
		CurrentPage  int           `json:"current_page"`
		PreviousPage int           `json:"previous_page"`
		NextPage     int           `json:"next_page"`
		LastPage     int           `json:"last_page"`
	}{
		Places:       places,
		Total:        total,
		CurrentPage:  page,
		PreviousPage: int(math.Max(float64(page-1), 1)),
		NextPage:     int(math.Min(float64(page+1), float64(lastPage))),
		LastPage:     lastPage,
	}

	w.Header().Set("Content-Type", "application/json")

	prettyJSON, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		log.Printf("Error encoding response to pretty JSON: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error encoding response")
		return
	}

	if _, err := w.Write(prettyJSON); err != nil {
		log.Printf("Error writing pretty JSON to response: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error writing response")
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	prettyJSON, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON response: %s", err)
		return
	}

	_, err = w.Write(prettyJSON)
	if err != nil {
		log.Printf("Error writing JSON response: %s", err)
	}
}

func recommendHandler(w http.ResponseWriter, r *http.Request, store db.Store) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	latStr, lonStr := r.URL.Query().Get("lat"), r.URL.Query().Get("lon")
	lat, latErr := strconv.ParseFloat(latStr, 64)
	lon, lonErr := strconv.ParseFloat(lonStr, 64)

	if latErr != nil || lonErr != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid latitude or longitude")
		return
	}

	const limit = 3
	places, err := store.GetClosestPlaces(ctx, lat, lon, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching recommendations: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, places)
}

func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &types.Claims{
		Username: "User",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in generating token"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			log.Printf("Error parsing token: %v", err)
			http.Error(w, "Unauthorized - Error parsing token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Printf("Invalid token")
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
