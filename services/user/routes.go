package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/leugard21/inku-api/types"
	"github.com/leugard21/inku-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")

	router.HandleFunc("/refresh", h.handleRefresh).Methods("POST")
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	user := &types.User{
		Email:    payload.Email,
		Username: payload.Username,
		Password: string(hashed),
	}

	if err := h.store.CreateUser(user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	rt, rtExp, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.store.SaveRefreshToken(user.ID, rt, rtExp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"access_token":  accessToken,
		"refresh_token": rt,
	})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.store.GetUserByIdentifier(payload.Login)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid username/email or password"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid username/email or password"))
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	rt, rtExp, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.store.SaveRefreshToken(user.ID, rt, rtExp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp := map[string]any{
		"access_token":  accessToken,
		"refresh_token": rt,
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := utils.ParseJSON(r, &body); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	rt, err := h.store.GetRefreshToken(body.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if rt == nil || time.Now().After(rt.ExpiresAt) {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("invalid or expired refresh token"))
		return
	}

	_ = h.store.DeleteRefreshToken(body.RefreshToken)
	newRT, newExp, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.store.SaveRefreshToken(rt.UserID, newRT, newExp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Issue new access JWT
	accessToken, err := utils.GenerateAccessToken(rt.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"access_token":  accessToken,
		"refresh_token": newRT,
	})
}
