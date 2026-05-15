package requestid

import "github.com/gin-gonic/gin"

const ContextKey = "requestID"
const Header = "X-Request-Id"

func Get(c *gin.Context) string {
	value, exists := c.Get(ContextKey)
	if !exists {
		return ""
	}

	requestID, ok := value.(string)
	if !ok {
		return ""
	}

	return requestID
}
