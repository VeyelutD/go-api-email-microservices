package auth

import (
	"encoding/json"
	"errors"
	api "github.com/VeyelutD/go-api-microservice/internal"
	"github.com/VeyelutD/go-api-microservice/internal/tokens"
	"github.com/VeyelutD/go-api-microservice/internal/users"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"time"
)

type VerifyPayload struct {
	AccessToken string `json:"access_token"`
}
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user SendOTPBody
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error("Failed to decode the request body for registration", "error", err)
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	defer r.Body.Close()
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		slog.Warn("Failed to validate the request body for registration", "details", err)
		api.RespondWithError(w, http.StatusBadRequest, "You entered an invalid email")
		return
	}
	newUser, err := h.service.users.Create(r.Context(), user.Email)
	if err != nil {
		slog.Error("Failed to create a user during registration", "details", err)
		api.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	userConfirmationToken, err := h.service.CreateUserConfirmationToken(r.Context(), user.Email)
	if err != nil {
		slog.Error("Failed to create a user confirmation token during registration", "details", err)
		api.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if err := h.service.email.SendConfirmationLink(r.Context(), newUser.Email, userConfirmationToken.Token); err != nil {
		slog.Error("Failed to send the confirmation link during registration", "details", err)
		api.RespondWithError(w, http.StatusInternalServerError, "Couldn't send the confirmation link, try again later")
		return
	}
	api.Respond(w, http.StatusCreated, newUser)
}

func (h *Handler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var user SendOTPBody
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error("Failed to decode the request body during login", "details", err)
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		slog.Warn("Failed to validate the request body during login", "details", err)
		api.RespondWithError(w, http.StatusBadRequest, "You entered an invalid email")
		return
	}
	userInDB, err := h.service.users.GetOneByEmail(r.Context(), user.Email)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			slog.Warn("User not found during login for email", "email", user.Email, "details", err)
			api.RespondWithError(w, http.StatusNotFound, "The user does not exist")
			return
		}
		slog.Error("Failed to find the user for email", "email", user.Email, "details", err)
		api.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if !userInDB.IsConfirmed {
		slog.Warn("User is not confirmed during login", "email", user.Email, "details", err)
		api.RespondWithError(w, http.StatusNotFound, "The user is not confirmed")
		return
	}
	userOTP, err := h.service.GetUserOTP(r.Context(), user.Email)
	if err != nil {
		if !errors.Is(err, ErrUserOTPNotFound) {
			slog.Error("Failed to get the user OTP for email", "email", user.Email, "details", err)
			api.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}
		userOTP, err = h.service.CreateUserOTP(r.Context(), user.Email)
		if err != nil {
			slog.Error("Failed to create the user OTP during login for email", "email", user.Email, "details", err)
			api.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err = h.service.email.SendOTP(r.Context(), user.Email, userOTP.Code); err != nil {
		slog.Error("Failed to send the OTP during login", "details", err)
		api.RespondWithError(w, http.StatusInternalServerError, "Something went wrong sending email")
		return
	}
	api.Respond(w, http.StatusOK, "We sent a one time password to your email, check you inbox")
}
func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var loginInfo VerifyOTPBody
	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		slog.Warn("Failed to decode the request body during login", "details", err)
		api.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	defer r.Body.Close()
	user, err := h.service.VerifyOTPAndGetUser(r.Context(), loginInfo.Email, loginInfo.Code)
	if err != nil {
		slog.Error("Failed to verify the OTP for email", "email", loginInfo.Email, "details", err)
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	accessToken, err := tokens.CreateAccessToken(user.ID)
	if err != nil {
		slog.Error("Failed to create the access token for email", "email", loginInfo.Email, "details", err)
		api.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	api.Respond(w, http.StatusOK, VerifyPayload{AccessToken: accessToken})
}

type VerifyResponse struct {
	Id string `json:"id"`
}

func (h *Handler) VerifyAccessToken(w http.ResponseWriter, r *http.Request) {
	var payload VerifyPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	claims, err := tokens.VerifyAccessToken(payload.AccessToken)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := claims.GetSubject()
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	api.Respond(w, http.StatusOK, VerifyResponse{Id: id})
}

func (h *Handler) ConfirmAccount(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		slog.Error("Empty confirm token provided during account confirmation")
		api.RespondWithError(w, http.StatusBadRequest, "Missing token")
	}
	userConfirmationToken, err := h.service.GetUserConfirmationToken(r.Context(), token)
	if err != nil {
		slog.Error("Failed to get user confirmation token for token during account confirmation", "token", token, "details", err)
		api.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if time.Now().Sub(userConfirmationToken.ExpiresAt.Time) > 0 {
		slog.Warn("Expired token provided during account confirmation for email", "token", token, "email", userConfirmationToken.Email, "details", err)
		api.RespondWithError(w, http.StatusNotFound, "This link has expired")
		return
	}
	if err := h.service.ConfirmUserAndDeleteConfirmationToken(r.Context(), userConfirmationToken.Email, userConfirmationToken.ID); err != nil {
		slog.Error("Failed to confirm the user and delete the confirmation token during account confirmation ", "email", userConfirmationToken.Email, "details", err)
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.Respond(w, http.StatusOK, "User confirmed successfully")
}
