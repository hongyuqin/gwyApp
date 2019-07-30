package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gwyApp/pkg/app"
	"gwyApp/pkg/e"
)

func GetArticle(c *gin.Context) {
	appG := app.Gin{C: c}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "success",
	})
}
