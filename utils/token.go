package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(id string, expireTime time.Time) (string, error) {
	loggedInToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  id,
			"exp": expireTime.Unix(),
		})

	tokenString, err := loggedInToken.SignedString([]byte(os.Getenv("SECRET")))

	fmt.Println(tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func VerifyJWTToken(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

// otherwise it will set token and userId to request's context and command will go to controller section.
func VerifyToken(flag int) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Token Not Found",
			})
			c.Abort()
			return
		}
		err := VerifyJWTToken(token)
		if err != nil {
			if flag == 0 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Access token expired",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Refresh token expired",
				})
			}
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Next()
	}
}

// // otherwise it will set token and userId to request's context and command will go to controller section.
// func VerifyToken(flag int) func(handler http.Handler) http.Handler {
// 	return func(handler http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			token := r.Header.Get("Authorization")
// 			if token == constant.EMPTY_STRING {
// 				errorhandling.SendErrorResponse(r, w, errorhandling.TokenNotFound, constant.EMPTY_STRING)
// 				return
// 			}
// 			token = token[7:]
// 			userId, err := utils.VerifyJWTToken(token)
// 			if err != nil {
// 				if flag == 0 {
// 					errorhandling.SendErrorResponse(r, w, errorhandling.AccessTokenExpired, constant.EMPTY_STRING)
// 				} else {
// 					errorhandling.SendErrorResponse(r, w, errorhandling.RefreshTokenExpired, constant.EMPTY_STRING)
// 				}
// 				return
// 			}
// 			ctx := context.WithValue(r.Context(), constant.TokenKey, token)
// 			ctx = context.WithValue(ctx, constant.UserIdKey, userId)
// 			handler.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }
