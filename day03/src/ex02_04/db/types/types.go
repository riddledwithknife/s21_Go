package types

import "github.com/dgrijalva/jwt-go"

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

type JsonResponse struct {
	Name     string  `json:"name"`
	Total    int     `json:"total"`
	Places   []Place `json:"places"`
	PrevPage int     `json:"prev_page,omitempty"`
	NextPage int     `json:"next_page,omitempty"`
	LastPage int     `json:"last_page"`
	Error    string  `json:"error,omitempty"`
}

type RecommendationResponse struct {
	Name   string  `json:"name"`
	Places []Place `json:"places"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
