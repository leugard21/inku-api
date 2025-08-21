package comic

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/leugard21/inku-api/types"
	"github.com/leugard21/inku-api/utils"
)

type Handler struct {
	store types.ComicStore
}

func NewHandler(store types.ComicStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/comics", h.handleCreateComic).Methods("POST")
	router.HandleFunc("/comics", h.handleGetComics).Methods("GET")
	router.HandleFunc("/comic/{id}", h.handleGetComicByID).Methods("GET")
	router.HandleFunc("/comics/search", h.handleSearchComics).Methods("GET")
}

func (h *Handler) handleCreateComic(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	author := r.FormValue("author")
	status := r.FormValue("status")
	genres := r.MultipartForm.Value["genres"]

	if title == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.New("title is required"))
		return
	}

	file, header, err := r.FormFile("cover")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	os.MkdirAll("uploads", os.ModePerm)

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	filepath := filepath.Join("uploads", filename)

	out, err := os.Create(filepath)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	coverURL := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)

	comic := &types.Comic{
		Title:       title,
		Description: description,
		Author:      author,
		CoverURL:    coverURL,
		Status:      status,
		Genres:      genres,
	}

	if err := h.store.CreateComic(comic); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, comic)
}

func (h *Handler) handleGetComics(w http.ResponseWriter, r *http.Request) {
	comics, err := h.store.GetAllComics()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, comics)
}

func (h *Handler) handleGetComicByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid comic id"))
		return
	}

	comic, err := h.store.GetComicByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if comic == nil {
		utils.WriteError(w, http.StatusNotFound, errors.New("comic not found"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, comic)
}

func (h *Handler) handleSearchComics(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	genre := r.URL.Query().Get("genre")
	status := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")

	results, err := h.store.SearchComicsAdvanced(q, genre, status, sort)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, results)
}
