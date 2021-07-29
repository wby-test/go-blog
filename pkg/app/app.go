package app

import (
	"github.com/fishblog/pkg/errcode"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

type Pager struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	TotalRow int `json:"total_row"`
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Response) ToResponseList(list interface{}, totalRows int) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"list": list,
		"pager": Pager{
			Page:     GetPage(r.Ctx),
			PageSize: GetPageSize(r.Ctx),
			TotalRow: totalRows,
		},
	})
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{
		"code": err.Code(),
		"msg":  err.Msg(),
	}
	detail := err.Details()
	if len(detail) > 0 {
		response["details"] = detail
	}

	r.Ctx.JSON(err.StatusCode(), response)
}
