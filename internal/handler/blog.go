package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"blog-system/internal/auth"
	"blog-system/internal/service"

	"github.com/gorilla/mux"
)

type BlogHandler struct {
	service *service.BlogService
}

func NewBlogHandler(service *service.BlogService) *BlogHandler {
	return &BlogHandler{service: service}
}

func (h *BlogHandler) RegisterRoutes(r *mux.Router, sessionManager *auth.SessionManager) {
	protected := (*r).PathPrefix("").Subrouter()
	(*protected).Use(auth.AuthMiddleware(sessionManager))

	articleStemPath := "/articles"
	articleSpecificPath := articleStemPath + "/{id:[0-9]+}"

	// public articles
	(*r).HandleFunc(articleStemPath, (*h).GetAllArticles).Methods("GET")
	(*r).HandleFunc(articleSpecificPath, (*h).GetArticle).Methods("GET")

	// public comments
	(*r).HandleFunc(articleSpecificPath+"/comments", (*h).AddComment).Methods("POST")

	// protected articles
	(*protected).HandleFunc(articleStemPath, (*h).CreateArticle).Methods("POST")
	(*protected).HandleFunc(articleSpecificPath, (*h).UpdateArticle).Methods("PUT")
	(*protected).HandleFunc(articleSpecificPath, (*h).DeleteArticle).Methods("DELETE")
}

// -- articles --
func (h *BlogHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Author  string   `json:"author"`
		Tags    []string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	article, err := (*h).service.CreateArticle(
		req.Title,
		req.Content,
		req.Author,
		req.Tags)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content/Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}

func (h *BlogHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	id, err := (*h).getIDFromPath(r)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	article, err := (*h).service.GetArticle(id)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func (h *BlogHandler) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := (*h).service.GetAllArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func (h *BlogHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	id, err := (*h).getIDFromPath(r)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	article, err := (*h).service.UpdateArticle(id, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func (h *BlogHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id, err := (*h).getIDFromPath(r)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	if err := (*h).service.DeleteArticle(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// -- comments --
func (h *BlogHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	articleID, err := (*h).getIDFromPath(r)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Author  string `json:"author"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	comment, err := (*h).service.AddComment(articleID, req.Author, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// -- helpers --
func (h *BlogHandler) getIDFromPath(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["id"])
}
