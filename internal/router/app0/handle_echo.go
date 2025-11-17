package app0

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitee.com/huajinet/go-example/internal/model"
)

func Echo(app ApplicationContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		c.JSON(http.StatusOK, gin.H{
			"headers": c.Request.Header,
		})
	}
}

func CreateBook(app ApplicationContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var book model.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := app.BookService.Create(&book); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusCreated, book)
	}
}

func ListBooks(app ApplicationContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		total, list, err := app.BookService.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"total": total,
			"list":  list,
		})
	}
}
