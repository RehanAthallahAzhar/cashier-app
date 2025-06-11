package middlewares

import (
	"net/http"

	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/pkg/authclient" // Impor AuthClient
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthMiddlewareOptions berisi dependensi untuk middleware autentikasi
type AuthMiddlewareOptions struct {
	AuthClient *authclient.AuthClient
}

// AuthMiddleware adalah fungsi middleware untuk memvalidasi token
func AuthMiddleware(opts AuthMiddlewareOptions) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token otentikasi tidak ditemukan atau format salah"})
			}
			token := authHeader[7:]

			// Signature fungsi diubah untuk mendapatkan 'userRole'
			isValid, userID, username, userRole, errMsg, err := opts.AuthClient.ValidateToken(token)
			if err != nil {
				c.Logger().Errorf("Kesalahan validasi token gRPC: %v", err)

				// Cek error dari gRPC
				st, ok := status.FromError(err)
				if ok {
					if st.Code() == codes.Unauthenticated {
						// Ini akan menangani pesan error seperti "Token kedaluwarsa" atau "Token tidak valid" dari server gRPC
						return c.JSON(http.StatusUnauthorized, map[string]string{"message": st.Message()})
					}
				}
				// Generic internal server error for other unexpected errors
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Kesalahan server saat memvalidasi token: " + err.Error()})
			}

			if !isValid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token tidak valid: " + errMsg})
			}

			// Jika token valid, Anda bisa menyimpan informasi pengguna di Echo Context
			c.Set("user_id", userID)
			c.Set("username", username)
			c.Set("role", userRole)                                                                   // Set 'role' di context
			c.Logger().Debugf("User %s (ID: %s, Role: %s) authenticated", username, userID, userRole) // Opsional: log info user

			return next(c)
		}
	}
}
