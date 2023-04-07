package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func wrapHTTPResult(result interface{}) gin.H {
	return gin.H{"code": "OK", "result": result}
}

func wrapHTTPError(err interface{}, errIdentifier ...string) gin.H {
	if len(errIdentifier) == 0 || errIdentifier[0] == "" {
		errIdentifier = []string{"NOT_SPECIFY_YET"}
	}
	return gin.H{"code": errIdentifier[0], "error": err}
}

func wrapHTTPRespond(c *gin.Context, respond interface{}, err error, errIdentifier ...string) {
	if err != nil {
		c.JSON(http.StatusBadRequest, wrapHTTPError(err.Error(), errIdentifier...))
		return
	}

	c.JSON(http.StatusOK, wrapHTTPResult(respond))
	return
}
