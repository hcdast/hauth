/*
 * @Author: hc
 * @Date: 2021-06-01 17:36:46
 * @LastEditors: hc
 * @LastEditTime: 2021-06-08 17:11:58
 * @Description:
 */
package service

import (
	"sync"

	"example-hauth/utils/logs"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type RegisterFunc func()

// key 应用名称
// value 注册路由方法
var regApp = make(map[string]RegisterFunc)
var regLock = new(sync.RWMutex)

// 注册服务
func AppRegister(name string, registerFunc RegisterFunc) {
	regLock.Lock()
	defer regLock.Unlock()
	if _, ok := regApp[name]; ok {
		panic("应用已经被注册，无法再次注册")
	} else {
		regApp[name] = registerFunc
	}
}

func Bootstrap() {
	// 开启消息，
	// 将80端口的请求，重定向到443上
	// go RedictToHtpps()

	// 插件 请求后执行
	beego.InsertFilter("/*", beego.FinishRouter, func(ctx *context.Context) {
		go WriteHandleLogs(ctx)
	}, false)

	// 插件 请求前执行
	beego.InsertFilter("/v1/auth/*", beego.BeforeRouter, func(ctx *context.Context) {
		go CheckConnection(ctx.ResponseWriter, ctx.Request)
	}, false)

	// 注册路由信息
	registerRouter()

	// 遍历服务 并执行其中方法
	for key, fc := range regApp {
		logs.Info("register App, name is:", key)
		fc()
	}
	// 启动beego服务
	beego.Run()
}
