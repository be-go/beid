package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func test_router(s *Server) {
	var r *gin.Engine = s.router
	// var db *pgxpool.Pool = s.db

	// curl http://localhost:8080/jwt-get
	r.GET("/jwt-get", func(c *gin.Context) {

		type MyCustomClaims struct {
			Foo string `json:"foo"`
			jwt.RegisteredClaims
		}
		var claims MyCustomClaims

		// // Create claims with multiple fields populated
		claims = MyCustomClaims{
			"bar",
			jwt.RegisteredClaims{
				// A usual scenario is to set the expiration time relative to the current time
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "test",
				Subject:   "somebody",
				ID:        "1",
				Audience:  []string{"somebody_else"},
			},
		}

		// // Create claims while leaving out some of the optional fields
		// claims = MyCustomClaims{
		// 	"bar",
		// 	jwt.RegisteredClaims{
		// 		// Also fixed dates can be used for the NumericDate
		// 		ExpiresAt: jwt.NewNumericDate(time.Unix(1916239022, 0)),
		// 		Issuer:    "test",
		// 	},
		// }

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString("kHS256Secret")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"token": ss})
		}
	})

	// curl http://localhost:8080/jwt-verify?tok=xxxxxx
	r.GET("/jwt-verify", func(c *gin.Context) {
		tokenString := c.Query("tok")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return "kHS256Secret", nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.JSON(http.StatusOK, claims)
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not ok"})
			return
		}

	})

}
