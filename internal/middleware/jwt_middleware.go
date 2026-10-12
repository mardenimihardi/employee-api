package middleware

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mardenimihardi/employee-api/internal/config"
)

func JWTAuth(cfg config.Config, db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 ||
				!strings.EqualFold(parts[0], "Bearer") ||
				strings.TrimSpace(parts[1]) == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Missing or invalid authorization header",
				})
			}

			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, echo.NewHTTPError(
						http.StatusUnauthorized,
						"Invalid token signing method",
					)
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Invalid or expired token",
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Invalid token claims",
				})
			}

			employeeID, ok := claims["sub"].(string)
			if !ok || employeeID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Invalid token subject",
				})
			}

			var accountStatus string
			err = db.QueryRowContext(
				c.Request().Context(),
				`SELECT account_status
				 FROM employees
				 WHERE id = $1 AND deleted_at IS NULL`,
				employeeID,
			).Scan(&accountStatus)

			if err == sql.ErrNoRows {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Employee not found or token is no longer valid",
				})
			}

			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"message": "Failed to validate employee",
				})
			}

			if accountStatus != "active" {
				return c.JSON(http.StatusForbidden, map[string]string{
					"message": "Employee account is inactive",
				})
			}

			c.Set("employee_id", employeeID)
			c.Set("claims", claims)

			return next(c)
		}
	}
}
