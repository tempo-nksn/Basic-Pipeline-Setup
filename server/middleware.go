package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/tempo-nksn/Tempo-Backend/constants"
)

//DB middleware attaches a database connection to gin's Context
func DB(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constants.ContextDB, db)
		c.Next()
	}
}
