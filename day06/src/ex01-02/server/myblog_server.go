package server

import (
	"bufio"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

type BlogPost struct {
	Title       string
	ArticleDate string
	ArticleText string
}

type RateLimiter struct {
	requests map[string]int
	mu       sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]int),
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		rl.mu.Lock()
		defer rl.mu.Unlock()

		rl.requests[clientIP]++

		if rl.requests[clientIP] > 100 {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func GetDBCredentials() string {
	file, _ := os.Open("credentials/db_credentials.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			log.Fatalln("Wrong admin or password format in admin_credentials file:", line)
		}
		user := strings.TrimSpace(parts[0])
		password := strings.TrimSpace(parts[1])

		return fmt.Sprintf("postgresql://%s:%s@localhost:5432/blog_db", user, password)
	}
	return ""
}

func HandleMainPage(conn pgx.Conn) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			page, err := strconv.Atoi(req.URL.Query().Get("page"))
			if err != nil || page < 1 {
				page = 1
			}

			pageSize := 3
			offset := (page - 1) * pageSize

			var totalCount int
			err = conn.QueryRow(context.Background(), "SELECT COUNT(article_id) FROM articles").Scan(&totalCount)
			if err != nil {
				http.Error(res, "Error executing query of getting count of posts from db", http.StatusGatewayTimeout)
				return
			}

			totalPages := totalCount / pageSize
			if totalCount%pageSize != 0 {
				totalPages++
			}

			if page > totalPages {
				http.Redirect(res, req, fmt.Sprintf("?page=%d", totalPages), http.StatusSeeOther)
				return
			}

			rows, err := conn.Query(context.Background(), "SELECT title, article_date::text, article_text FROM articles LIMIT $1 OFFSET $2", pageSize, offset)
			if err != nil {
				http.Error(res, "Error executing query of getting posts from db", http.StatusGatewayTimeout)
				return
			}
			defer rows.Close()

			var blogPosts []BlogPost
			for rows.Next() {
				var post BlogPost
				if err := rows.Scan(&post.Title, &post.ArticleDate, &post.ArticleText); err != nil {
					http.Error(res, "Error scanning SELECT output", http.StatusGatewayTimeout)
					return
				}
				blogPosts = append(blogPosts, post)
			}

			if err := rows.Err(); err != nil {
				http.Error(res, "Error after iteration over SELECT output", http.StatusGatewayTimeout)
				return
			}

			pageNumbers := make([]int, totalPages)
			for i := 0; i < totalPages; i++ {
				pageNumbers[i] = i + 1
			}

			paginationData := struct {
				BlogPosts   []BlogPost
				TotalPages  int
				PageNumbers []int
			}{
				BlogPosts:   blogPosts,
				TotalPages:  totalPages,
				PageNumbers: pageNumbers,
			}

			tmpl := template.Must(template.ParseFiles("html/index.html"))
			_ = tmpl.Execute(res, paginationData)
			return
		}
	}
}

func HandleLoginPage() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			tmpl := template.Must(template.ParseFiles("html/login.html"))
			_ = tmpl.Execute(res, nil)
			return
		}

		var adminCredentials = make(map[string]string)
		getAdminCredentials(&adminCredentials)

		username := req.FormValue("username")
		password := req.FormValue("password")

		if pwd, ok := adminCredentials[username]; ok && pwd == password {
			session, _ := store.Get(req, "admin-session")
			session.Values["authenticated"] = true
			err := session.Save(req, res)
			if err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
			}
			http.Redirect(res, req, "/admin", http.StatusFound)
			return
		}

		http.Error(res, "Invalid username or password", http.StatusUnauthorized)
	}
}

func getAdminCredentials(adminCredentials *map[string]string) {
	file, _ := os.Open("credentials/admin_credentials.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			log.Fatalln("Wrong admin or password format in admin_credentials file:", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		(*adminCredentials)[key] = value
	}
}

func HandleAdminPage(conn pgx.Conn) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "admin-session")
		loggedIn, ok := session.Values["authenticated"].(bool)
		if !loggedIn || !ok {
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}

		if req.Method == "GET" {
			tmpl := template.Must(template.ParseFiles("html/admin.html"))
			_ = tmpl.Execute(res, nil)
		}

		if req.Method == "POST" {
			if err := req.ParseForm(); err != nil {
				http.Error(res, "Error parsing form", http.StatusBadRequest)
				return
			}

			title := req.FormValue("title")
			articleText := req.FormValue("article_text")

			_, err := conn.Exec(context.Background(), "INSERT INTO articles (title, article_date, article_text) VALUES ($1, CURRENT_DATE, $2)", title, articleText)
			if err != nil {
				http.Error(res, "Error executing sql query", http.StatusInternalServerError)
				return
			}

			tmpl := template.Must(template.ParseFiles("html/admin.html"))
			_ = tmpl.Execute(res, nil)
		}
	}
}

func HandleLogout() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "admin-session")
		session.Values["authenticated"] = false
		session.Save(req, res)

		http.Redirect(res, req, "/", http.StatusFound)
	}
}
