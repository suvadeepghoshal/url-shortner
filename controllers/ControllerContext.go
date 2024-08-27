package controllers

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/markbates/goth/gothic"
	uuid "github.com/satori/go.uuid"
	"log/slog"
	"net/http"
)

// ControllerContext Passing the validator instance to all the controllers, so that only one reference of the validator is created and can be used all the time
type ControllerContext struct {
	Validator *validator.Validate
}

type authenticationResponse struct {
	status bool
	err    error
}

type AuthenticationService interface {
	validate(r *http.Request) authenticationResponse
}

type UserAuthenticationService struct {
	reqID string
}

func (service *UserAuthenticationService) validate(r *http.Request) authenticationResponse {

	s, sErr := gothic.Store.Get(r, "auth-session")
	slog.Debug("validate", "session", s)
	if sErr != nil {
		slog.Error("Unable to get auth session", "req_id", service.reqID, "sErr", sErr)
		return authenticationResponse{false, errors.New("user is not authenticated")}
	}

	ua := s.Values["user_authed"]
	if ua == nil || !ua.(bool) {
		slog.Error("User is not authenticated", "req_id", service.reqID, "user_session", ua)
		return authenticationResponse{false, errors.New("user is not authenticated")}
	}

	return authenticationResponse{true, nil}
}

func newUserAuthenticationService(reqId string) *UserAuthenticationService {
	slog.Info("NewUserAuthenticationService", "req_id", reqId)
	return &UserAuthenticationService{
		reqID: reqId,
	}
}

func ProtectedResourceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("ProtectedResourceMiddleware", "req_method", r.Method)
		ctx := context.WithValue(r.Context(), "req_id", uuid.NewV4().String())
		authResp := newUserAuthenticationService(ctx.Value("req_id").(string)).validate(r)
		if !authResp.status {
			http.Error(w, authResp.err.Error(), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
