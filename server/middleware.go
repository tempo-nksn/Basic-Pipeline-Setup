package server

import (
	"github.com/tempo-nksn/Tempo-Backend/constants"
	"github.com/tempo-nksn/Tempo-Backend/models"
	"log"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type loginVariables struct {
	Username string `form:"username" json:"u_name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

//DB middleware attaches a database connection to gin's Context
func DB(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constants.ContextDB, db)
		c.Next()
	}
}

//JWT is
func JWT() *jwt.GinJWTMiddleware {
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm: "Tempo",
		Key:   []byte("Car Pool"),
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.Rider); ok {
				return jwt.MapClaims{
					"id": v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals loginVariables
			if err := c.Bind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			db := getDB(c)
			var user models.Rider
			log.Println(loginVals)
			if err := db.Where(&models.Rider{UName: loginVals.Username}).First(&user).Error; err != nil {
				return "", jwt.ErrFailedAuthentication
			}
			if err := db.Where(&models.Rider{UName: loginVals.Username, Password: loginVals.Password}).First(&user).Error; err != nil {
				return "", jwt.ErrFailedAuthentication
			}

			return &user, nil

		},
	}
	authMiddleware.MiddlewareInit()

	return authMiddleware
}
