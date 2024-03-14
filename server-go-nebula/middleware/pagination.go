package middleware

import (
	"net/http"
	res "server/shares/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kataras/i18n"
)

func Pagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetString("lang")
		_limit := c.DefaultQuery("limit", "15")
		limit, err := strconv.ParseInt(_limit, 10, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
			return
		}
		_page := c.DefaultQuery("page", "1")
		page, err := strconv.ParseInt(_page, 10, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, res.Err(i18n.Tr(lang, "format_args")))
			return
		}
		c.Set("limit", limit)
		c.Set("page", page)
		c.Set("offset", page-1)
		c.Next()
	}
}
