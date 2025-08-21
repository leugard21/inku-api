package page

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/leugard21/inku-api/configs"
	"github.com/leugard21/inku-api/utils"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/chapters/{id}/pages", h.handleGetPages).Methods("GET")
	router.Handle("/chapters/{id}/pages", utils.AuthMiddleware(utils.AdminOnly(http.HandlerFunc(h.handleUploadPages)))).Methods("POST")
}

func (h *Handler) handleUploadPages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chapterID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid chapter id"))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 50<<20)
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	files := r.MultipartForm.File["pages"]
	if len(files) == 0 {
		utils.WriteError(w, http.StatusBadRequest, errors.New("no pages uploaded"))
		return
	}

	var pageURLs []string
	for idx, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		defer file.Close()

		os.MkdirAll("uploads/pages", os.ModePerm)
		filename := fmt.Sprintf("%d_%d_%s", chapterID, idx+1, fileHeader.Filename)
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

		url := fmt.Sprintf("%s/uploads/pages/%s", configs.Envs.PublicHost, filename)
		pageURLs = append(pageURLs, url)

		if err := h.store.CreatePage(chapterID, idx+1, url); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"chapter_id": chapterID,
		"pages":      pageURLs,
	})
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
