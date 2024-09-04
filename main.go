package main

import (
	Translator "ginni-ai-task/src/Translator"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/translate", func(c *gin.Context) {
		var callTranscription []Translator.CallTranscription

		if err := c.ShouldBindJSON(&callTranscription); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if translateError := Translator.Translate(&callTranscription); translateError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": translateError.Error()})
		}
		c.JSON(http.StatusOK, callTranscription)
	})

	router.Run(":8080")
}
