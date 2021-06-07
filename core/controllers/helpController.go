/*
 * @Author: hc
 * @Date: 2021-06-01 17:36:46
 * @LastEditors: hc
 * @LastEditTime: 2021-06-03 11:07:20
 * @Description:
 */
package controllers

import (
	"example-hauth/core/groupcache"
	"example-hauth/core/hrpc"
	"example-hauth/utils/hret"
	"example-hauth/utils/i18n"

	"github.com/astaxie/beego/context"
)

type helpController struct {
}

var HelpCtl = &helpController{}

// swagger:operation GET /v1/help/system/help StaticFiles helpController
//
// 系统帮助页面
//
// 将会返回系统帮助首页,其中包含了系统管理操作文档,API文档
//
// ---
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// responses:
//   '200':
//     description: all domain information
func (this helpController) Page(ctx *context.Context) {
	ctx.Request.ParseForm()

	if !hrpc.BasicAuth(ctx.Request) {
		hret.Error(ctx.ResponseWriter, 403, i18n.NoAuth(ctx.Request))
		return
	}

	rst, err := groupcache.GetStaticFile("AsofdateHelpPage")
	if err != nil {
		hret.Error(ctx.ResponseWriter, 404, i18n.PageNotFound(ctx.Request))
		return
	}
	ctx.ResponseWriter.Write(rst)
}

func init() {
	groupcache.RegisterStaticFile("AsofdateHelpPage", "./views/help/auth_help.tpl")
}
