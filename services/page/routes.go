package page

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
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/chapters/{id}/pages", h.handleUploadPage).Methods("POST")
	router.HandleFunc("/chapters/{id}/pages", h.handleGetPages).Methods("GET")
}

func (h *Handler) handleUploadPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chapterID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	pageNumStr := r.FormValue("pageNumber")
	pageNumber, err := strconv.Atoi(pageNumStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid pageNumber"))
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("image is required"))
		return
	}
	defer file.Close()

	os.MkdirAll("uploads/pages", os.ModePerm)

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	filepath := filepath.Join("uploads/pages", filename)

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

	imageURL := fmt.Sprintf("http://localhost:8080/uploads/pages/%s", filename)

	page := &types.Page{
		ChapterID:  chapterID,
		PageNumber: pageNumber,
		ImageURL:   imageURL,
	}

	if err := h.store.CreatePage(chapterID, page); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, page)
}

func (h *Handler) handleGetPages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chapterID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	pages, err := h.store.GetPagesByChapter(chapterID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, pages)
}
