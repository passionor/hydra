package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/context"
)

var _ context.IPath = &rpath{}

//rpath 处理请求的路径信息
type rpath struct {
	*gin.Context
}

//GetMethod 获取服务请求方式
func (c *rpath) GetMethod() string {
	return c.Request.Method
}

//GetService 获取处理的服务名
func (c *rpath) GetService() string {
	return ""
}

//GetPath 获取请求路径
func (c *rpath) GetPath() string {
	return c.Context.Request.URL.Path
}

//GetHeader 获取请求头信息
func (c *rpath) GetHeader(key string) string {
	return c.Context.GetHeader(key)
}

//GetCookie 获取cookie信息
func (c *rpath) GetCookie(name string) (string, bool) {
	if cookie, err := c.Context.Request.Cookie(name); err == nil {
		return cookie.Value, true
	}
	return "", false
}
