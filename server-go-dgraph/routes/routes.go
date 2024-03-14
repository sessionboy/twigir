package routes

import (
	"server/routes/message"
	"server/routes/notify"
	"server/routes/status"
	"server/routes/users"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	users.Route(r)
	status.Route(r)
	message.Route(r)
	notify.Route(r)
}
