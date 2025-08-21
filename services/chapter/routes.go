package chapter

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/leugard21/inku-api/types"
	"github.com/leugard21/inku-api/utils"
)

type Handler struct {
	store types.ChapterStore
}

func NewHandler(store types.ChapterStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/comics/{id}/chapters", h.handleCreateChapter).Methods("POST")
	router.HandleFunc("/comics/{id}/chapters", h.handleGetChapters).Methods("GET")
}

func (h *Handler) handleCreateChapter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	comicID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var payload types.CreateChapterPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	chapter := &types.Chapter{
		ComicID:       comicID,
		Title:         payload.Title,
		ChapterNumber: payload.ChapterNumber,
	}

	if err := h.store.CreateChapter(comicID, chapter); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, chapter)
}

func (h *Handler) handleGetChapters(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	comicID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	chapters, err := h.store.GetChaptersByComic(comicID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, chapters)
}
