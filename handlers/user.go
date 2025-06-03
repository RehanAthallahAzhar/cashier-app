package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/rehanazhar/cashier-app/models"
// )

// // func (api *API) Register(c echo.Context) error {
// // 	var creds models.User

// // 	if err := c.Bind(&creds); err != nil {
// // 		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
// // 			Code:  http.StatusBadRequest,
// // 			Error: "Invalid request payload",
// // 		})
// // 	}

// // 	if creds.Username == "" || creds.Password == "" {
// // 		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
// // 			Code:  http.StatusBadRequest,
// // 			Error: "Username and password must not be empty",
// // 		})
// // 	}

// // 	err := api.UserRepo.AddUser(creds)
// // 	if err != nil {
// // 		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
// // 			Code:  http.StatusInternalServerError,
// // 			Error: "Failed to register user",
// // 		})
// // 	}

// // 	return c.JSON(http.StatusOK, models.SuccessResponse{
// // 		Username: creds.Username,
// // 		Message:  "User registered successfully",
// // 	})
// // }

// func (api *API) Login(w http.ResponseWriter, r *http.Request) {
// 	var creds models.User
// 	err := json.NewDecoder(r.Body).Decode(&creds)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Internal Server Error"})
// 		return
// 	}

// 	if creds.Username == "" || creds.Password == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Username or Password empty"})
// 		return
// 	}

// 	if api.UserRepo.CheckPassLength(creds.Password) {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Please provide a password of more than 5 characters"})
// 		return
// 	}

// 	if api.UserRepo.CheckPassAlphabet(creds.Password) {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Please use Password with Contains non Alphabetic Characters"})
// 		return
// 	}

// 	err = api.UserRepo.UserAvail(creds)
// 	if err != nil {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
// 		return
// 	}

// 	sessionToken := uuid.NewString()
// 	expiresAt := time.Now().Add(5 * time.Hour)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Internal Server Error"})
// 		return
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "session_token",
// 		Path:    "/",
// 		Value:   sessionToken,
// 		Expires: expiresAt,
// 	})

// 	session := models.Session{Token: sessionToken, Username: creds.Username, Expiry: expiresAt}
// 	err = api.SessionRepo.AddSessions(session)

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(models.SuccessResponse{Username: creds.Username, Message: "Login Success"})
// }

// func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
// 	username := fmt.Sprintf("%s", r.Context().Value("username"))
// 	c, err := r.Cookie("session_token")
// 	if err != nil {
// 		if err == http.ErrNoCookie {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Internal Server Error"})
// 			return
// 		}
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Internal Server Error"})
// 		return
// 	}
// 	sessionToken := c.Value

// 	api.SessionRepo.DeleteSessions(sessionToken)
// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "session_token",
// 		Value:   "",
// 		Expires: time.Now(),
// 	})

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(models.SuccessResponse{Username: username, Message: "Logout Success"})
// }

// func (api *API) SessionValid(w http.ResponseWriter, r *http.Request) {
// 	username := fmt.Sprintf("%s", r.Context().Value("username"))
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(models.SuccessResponse{Username: username, Message: "Token Valid"})
// }
