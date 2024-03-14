package routes

import (
	"server/routes/status"
	"server/routes/users"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	users.Route(r)
	status.Route(r)
}
