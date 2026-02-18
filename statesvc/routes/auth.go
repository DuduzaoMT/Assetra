package routes

import (
	"assetra/statesvc/resthandlers"
	"net/http"
)

func NewAuthRoutes(authHandlers resthandlers.AuthHandlers) []*Route {
	return []*Route{
		{
			Path:    "/signup",
			Method:  "POST",
			Handler: authHandlers.SignUp,
		},
		{
			Path:    "/signin",
			Method:  http.MethodPost,
			Handler: authHandlers.SignIn,
		},
		{
			Path:    "/refresh-token",
			Method:  http.MethodPost,
			Handler: authHandlers.RefreshToken,
		},
		{
			Path:    "/logout",
			Method:  http.MethodPost,
			Handler: authHandlers.Logout,
		},
		{
			Path:         "/users",
			Method:       http.MethodGet,
			Handler:      authHandlers.GetUsers,
			AuthRequired: true,
		},
		{
			Path:         "/users/{id}",
			Method:       http.MethodGet,
			Handler:      authHandlers.GetUser,
			AuthRequired: true,
		},
		{
			Path:         "/users/{id}",
			Method:       http.MethodPut,
			Handler:      authHandlers.UpdateUser,
			AuthRequired: true,
		},
		{
			Path:         "/users/{id}",
			Method:       http.MethodDelete,
			Handler:      authHandlers.DeleteUser,
			AuthRequired: true,
		},
	}
}
