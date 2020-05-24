package container

import (
	"fmt"

	"github.com/micro-plat/hydra/registry/conf"
	"github.com/micro-plat/hydra/registry/conf/server"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

//ICloser 关闭
type ICloser interface {
	Close() error
}

//IContainer 组件容器
type IContainer interface {
	GetOrCreate(typ string, name string, creator func(conf *conf.JSONConf) (interface{}, error)) (interface{}, error)
	ICloser
}

//Container 容器用于缓存公共组件
type Container struct {
	cache cmap.ConcurrentMap
}

//NewContainer 构建容器
func NewContainer() *Container {
	return &Container{
		cache: cmap.New(8),
	}

}

//GetOrCreate 获取指定名称的组件，不存在时自动创建
func (c *Container) GetOrCreate(typ string, name string, creator func(conf *conf.JSONConf) (interface{}, error)) (interface{}, error) {
	vc, err := server.Cache.GetVarConf()
	if err != nil {
		return nil, err
	}
	js, err := vc.GetConf(typ, name)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%s_%s_%d", typ, name, js.GetVersion())
	_, obj, err := c.cache.SetIfAbsentCb(key, func(i ...interface{}) (interface{}, error) {
		return creator(js)
	})
	return obj, err

}

//Close 释放组件资源
func (c *Container) Close() error {
	c.cache.RemoveIterCb(func(key string, v interface{}) bool {
		if closer, ok := v.(ICloser); ok {
			closer.Close()
		}
		return true
	})
	return nil
}