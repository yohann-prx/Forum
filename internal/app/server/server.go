package server

import (
	"SPORTALK/internal/store"
	"SPORTALK/internal/store/sqlite"
	"log"
	"net/http"
)

type server struct {
	store  store.Store
	router *http.ServeMux
	logger *log.Logger
}

// NewServer creates a new server
func NewServer(store store.Store) *server {
	return &server{
		store:  store,
		router: &http.ServeMux{},
		logger: log.Default(),
	}
}

// ServeHTTP implements the http.Handler interface
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// HandlePaths handles the paths for the server
func Start(con Config) error {
	db, err := InitDB(con)
	if err != nil {
		return err
	}

	store := sqlite.NewSQL(db)

	server := NewServer(store)

	server.HandlePaths()

	log.Println("Starting server: http://localhost:8080")

	return http.ListenAndServe(con.Port, server)
}
