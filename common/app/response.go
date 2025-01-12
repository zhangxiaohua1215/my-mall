package app

import (
	"my-mall/common/errcode"

	"github.com/gin-gonic/gin"
)

type response struct {
	ctx        *gin.Context
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	RequestId  string      `json:"request_id"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *pagination `json:"pagination,omitempty"`
}

func NewResponse(c *gin.Context) *response {
	return &response{ctx: c}
}

// SetPagination 设置Response的分页信息
func (r *response) SetPagination(pagination *pagination) *response {
	r.Pagination = pagination
	return r
}

func (r *response) Success(data interface{}) {
	r.Code = errcode.Success.Code()
	r.Msg = errcode.Success.Msg()
	if id, exists := r.ctx.Get("traceid"); exists {
		r.RequestId = id.(string)
	}
	r.Data = data
	r.ctx.JSON(errcode.Success.HttpStatusCode(), r)
}

func (r *response) SuccessOk() {
	r.Success("")
}

func (r *response) Error(err *errcode.AppError) {
	r.Code = err.Code()
	r.Msg = err.Msg()
	if id, exists := r.ctx.Get("traceid"); exists {
		r.RequestId = id.(string)
	}
	// 兜底记一条响应错误, 项目自定义的AppError中有错误链条, 方便出错后排查问题
	// logger.Ctx(r.ctx).Errorw("api_response_error", "err", err)
	r.ctx.JSON(err.HttpStatusCode(), r)
}
