package resthandlers

import (
	"assetra/pb"
	"assetra/security"
	"assetra/statesvc/restutil"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type AuthHandlers interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type authHandlers struct {
	authSvcClient pb.AuthServiceClient
}

func NewAuthHandlers(authSvcClient pb.AuthServiceClient) AuthHandlers {
	return &authHandlers{authSvcClient: authSvcClient}
}

// SignUp method
func (h *authHandlers) SignUp(w http.ResponseWriter, r *http.Request) {
	// the authetication content must be at the body
	body := r.Body
	if body == nil {
		restutil.WriteError(w, http.StatusBadRequest, restutil.ErrEmptyBody)
		return
	}
	defer body.Close()

	body = http.MaxBytesReader(w, r.Body, 1024*1024) // 1MB max
	content, err := io.ReadAll(body)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user := new(pb.User)
	err = json.Unmarshal(content, user)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// the rest of the parameteres are handled at the service authentication
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()

	resp, err := h.authSvcClient.SignUp(r.Context(), user)
	if err != nil {
		restutil.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Check if we're in production (HTTPS available)
	isSecure := os.Getenv("ENV") == "production"
	
	// Set refresh token in httpOnly cookie
	security.SetRefreshTokenCookie(w, resp.RefreshToken, isSecure)
	
	// Return user data and access token in body
	restutil.WriteAsJson(w, http.StatusOK, map[string]interface{}{
		"user":         resp.User,
		"access_token": resp.Token,
	})
}

func (h *authHandlers) SignIn(w http.ResponseWriter, r *http.Request) {
	// the authetication content must be at the body
	body := r.Body
	if body == nil {
		restutil.WriteError(w, http.StatusBadRequest, restutil.ErrEmptyBody)
		return
	}
	defer body.Close()

	content, err := io.ReadAll(body)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user := new(pb.SignInRequest)
	err = json.Unmarshal(content, user)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	resp, err := h.authSvcClient.SignIn(r.Context(), user)
	if err != nil {
		restutil.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Check if we're in production (HTTPS available)
	isSecure := os.Getenv("ENV") == "production"
	
	// Set refresh token in httpOnly cookie
	security.SetRefreshTokenCookie(w, resp.RefreshToken, isSecure)
	
	// Return user data and access token in body
	restutil.WriteAsJson(w, http.StatusOK, map[string]interface{}{
		"user":         resp.User,
		"access_token": resp.Token,
	})
}

// RefreshToken generates a new access token using the refresh token from cookies
func (h *authHandlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Extract refresh token from cookie
	refreshToken, err := security.ExtractRefreshToken(r)
	if err != nil {
		restutil.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	// Call auth service to validate and generate new tokens
	resp, err := h.authSvcClient.RefreshToken(r.Context(), &pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		restutil.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	// Check if we're in production (HTTPS available)
	isSecure := os.Getenv("ENV") == "production"

	// Set new refresh token in httpOnly cookie (token rotation)
	security.SetRefreshTokenCookie(w, resp.RefreshToken, isSecure)

	// Return new access token in body
	restutil.WriteAsJson(w, http.StatusOK, map[string]interface{}{
		"access_token": resp.Token,
	})
}

// Logout clears authentication cookies and revokes refresh token
func (h *authHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	// Try to extract and revoke the refresh token from cookie
	refreshToken, err := security.ExtractRefreshToken(r)
	if err == nil && refreshToken != "" {
		// Hash it and revoke in database
		// tokenHash := security.HashRefreshToken(refreshToken)
		// Note: We need to add a method to revoke via gRPC
		// For now, just clear the cookie
		// TODO: Call authSvcClient.RevokeRefreshToken(ctx, &pb.RevokeTokenRequest{TokenHash: tokenHash})
	}
	
	// Clear auth cookie
	security.ClearAuthCookies(w)
	restutil.WriteAsJson(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

func (h *authHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	tokenPayload, err := restutil.AuthRequestWithId(r)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	body := r.Body
	if body == nil {
		restutil.WriteError(w, http.StatusBadRequest, restutil.ErrEmptyBody)
		return
	}
	defer body.Close()

	content, err := io.ReadAll(body)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user := new(pb.User)
	err = json.Unmarshal(content, user)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}
	
	// ensure that we are updating the authenticated user
	if !restutil.ContainsRole(tokenPayload.Roles, "admin"){
		user.Id = tokenPayload.UserId
	}

	updatedUser, err := h.authSvcClient.UpdateUser(r.Context(), user)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	restutil.WriteAsJson(w, http.StatusOK, updatedUser)
}

func (h *authHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	tokenPayload, err := restutil.AuthRequestWithId(r)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.authSvcClient.GetUser(r.Context(), &pb.GetUserRequest{Id: tokenPayload.UserId})
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	restutil.WriteAsJson(w, http.StatusOK, user)
}

func (h *authHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	tokenPayload, err := restutil.AuthRequestToken(r)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.authSvcClient.GetUser(r.Context(), &pb.GetUserRequest{Id: tokenPayload.UserId})
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !restutil.ContainsRole(tokenPayload.Roles, "admin") || !restutil.ContainsRole(user.Role, "admin") {
		restutil.WriteError(w, http.StatusBadRequest, restutil.ErrUnauthorized)
		return
	}

	resp, err := h.authSvcClient.ListUsers(r.Context(), &pb.ListUsersRequest{})
	if err != nil {
		restutil.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var users []*pb.User
	for {
		user, err := resp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			restutil.WriteError(w, http.StatusBadRequest, err)
			return
		}
		users = append(users, user)
	}
	restutil.WriteAsJson(w, http.StatusOK, users)
}

func (h *authHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	tokenPayload, err := restutil.AuthRequestWithId(r)
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	targetId := vars["id"]

	user, err := h.authSvcClient.GetUser(r.Context(), &pb.GetUserRequest{Id: tokenPayload.UserId})
	if err != nil {
		restutil.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if tokenPayload.UserId != targetId && !restutil.ContainsRole(user.Role, "admin") {
		restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
		return
	}

	deleteduser, err := h.authSvcClient.DeleteUser(r.Context(), &pb.GetUserRequest{Id: targetId})

	// Clear auth cookies if the user deleted their own account
	if tokenPayload.UserId == targetId {
		security.ClearAuthCookies(w)
	}
	w.Header().Set("Entity", deleteduser.Id)
	restutil.WriteAsJson(w, http.StatusNoContent, nil)
}
