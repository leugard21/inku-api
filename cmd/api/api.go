package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leugard21/inku-api/services/chapter"
	"github.com/leugard21/inku-api/services/comic"
	"github.com/leugard21/inku-api/services/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	comicStore := comic.NewStore(s.db)
	comicHandler := comic.NewHandler(comicStore)
	comicHandler.RegisterRoutes(subrouter)

	chapterStore := chapter.NewStore(s.db)
	chapterHandler := chapter.NewHandler(chapterStore)
	chapterHandler.RegisterRoutes(subrouter)

	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
