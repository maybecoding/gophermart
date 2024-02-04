package http

import "github.com/gin-gonic/gin"

type response struct {
	Error string `json:"error"`
}

func errorResponse(c *gin.Context, err error, code int) {
	c.AbortWithStatusJSON(code,
		response{
			Error: err.Error(),
		})
}
