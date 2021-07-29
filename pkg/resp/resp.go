package resp

import "github.com/gin-gonic/gin"

var RespNoData = gin.H{
	"code": 0,
	"msg":  "ok",
}
