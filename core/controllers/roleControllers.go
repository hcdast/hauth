package controllers

import (
	"encoding/json"

	"example-hauth/core/groupcache"
	"example-hauth/core/hrpc"
	"example-hauth/core/models"
	"example-hauth/utils"
	"example-hauth/utils/hret"
	"example-hauth/utils/i18n"
	"example-hauth/utils/jwt"
	"example-hauth/utils/logs"
	"example-hauth/utils/validator"

	"github.com/astaxie/beego/context"
)

type roleController struct {
	models models.RoleModel
}

var RoleCtl = &roleController{
	models.RoleModel{},
}

// swagger:operation GET /v1/auth/role/page StaticFiles roleController
//
// 角色管理页面
//
// 如果用户被授权访问角色管理页面,则系统返回角色管理页面内容,否则返回404错误
//
// ---
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// parameters:
// - name: domain_id
//   in: query
//   description: domain code number
//   required: true
//   type: string
//   format:
// responses:
//   '200':
//     description: success
func (roleController) Page(ctx *context.Context) {
	ctx.Request.ParseForm()
	if !hrpc.BasicAuth(ctx.Request) {
		hret.Error(ctx.ResponseWriter, 403, i18n.NoAuth(ctx.Request))
		return
	}

	rst, err := groupcache.GetStaticFile("AsofdateRolePage")
	if err != nil {
		hret.Error(ctx.ResponseWriter, 404, i18n.PageNotFound(ctx.Request))
		return
	}
	ctx.ResponseWriter.Write(rst)
}

// swagger:operation GET /v1/auth/role/get roleController roleController
//
// 查询角色信息
//
// 查询指定域中的角色信息
//
// ---
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// parameters:
// - name: domain_id
//   in: query
//   description: domain code number
//   required: true
//   type: string
//   format:
// responses:
//   '200':
//     description: success
func (this roleController) Get(ctx *context.Context) {
	ctx.Request.ParseForm()
	if !hrpc.BasicAuth(ctx.Request) {
		hret.Error(ctx.ResponseWriter, 403, i18n.NoAuth(ctx.Request))
		return
	}

	domain_id := ctx.Request.FormValue("domain_id")

	if validator.IsEmpty(domain_id) {
		cookie, _ := ctx.Request.Cookie("Authorization")
		jclaim, err := jwt.ParseJwt(cookie.Value)
		if err != nil {
			logs.Error(err)
			hret.Error(ctx.ResponseWriter, 403, i18n.Disconnect(ctx.Request))
			return
		}
		domain_id = jclaim.DomainId
	}

	if !hrpc.DomainAuth(ctx.Request, domain_id, "r") {
		hret.Error(ctx.ResponseWriter, 403, i18n.Get(ctx.Request, "as_of_date_domain_permission_denied"))
		return
	}

	rst, err := this.models.Get(domain_id)

	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 421, i18n.Get(ctx.Request, "error_role_query"), err)
		return
	}

	hret.Json(ctx.ResponseWriter, rst)
}

// swagger:operation POST /v1/auth/role/post roleController roleController
//
// 新增角色信息
//
// 在某个指定的域中,新增角色信息
//
// ---
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// parameters:
// - name: domain_id
//   in: query
//   description: domain code number
//   required: true
//   type: string
//   format:
// responses:
//   '200':
//     description: success
func (this roleController) Post(ctx *context.Context) {
	ctx.Request.ParseForm()
	if !hrpc.BasicAuth(ctx.Request) {
		hret.Error(ctx.ResponseWriter, 403, i18n.NoAuth(ctx.Request))
		return
	}

	form := ctx.Request.Form
	domainid := form.Get("domain_id")
	if !hrpc.DomainAuth(ctx.Request, domainid, "w") {
		logs.Error("没有权限在这个域中新增角色信息")
		hret.Error(ctx.ResponseWriter, 421, i18n.Get(ctx.Request, "as_of_date_domain_permission_denied"))
		return
	}

	cok, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := jwt.ParseJwt(cok.Value)
	if err != nil {
		hret.Error(ctx.ResponseWriter, 403, i18n.Disconnect(ctx.Request))
		return
	}

	msg, err := this.models.Post(form, jclaim.UserId)
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 421, i18n.Get(ctx.Request, msg), err)
		return
	}
	hret.Success(ctx.ResponseWriter, i18n.Success(ctx.Request))
}

// swagger:operation POST /v1/auth/role/delete roleController roleController
//
// 删除角色信息
//
// 删除某个指定域中的角色信息
//
// ---
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// parameters:
// - name: domain_id
//   in: query
//   description: domain code number
//   required: true
//   type: string
//   format:
// responses:
//   '200':
//     description: success
func (this roleController) Delete(ctx *context.Context) {
	ctx.Request.ParseForm()
	if !hrpc.BasicAuth(ctx.Request) {
		hret.Error(ctx.ResponseWriter, 403, i18n.NoAuth(ctx.Request))
		return
	}

	var allrole []models.RoleInfo
	err := json.Unmarshal([]byte(ctx.Request.FormValue("JSON")), &allrole)
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 421, i18n.Get(ctx.Request, "error_role_json_failed"), err)
		return
	}

	for _, val := range allrole {
		if !hrpc.DomainAuth(ctx.Request, val.Domain_id, "w") {
			hret.Error(ctx.ResponseWriter, 403, i18n.WriteDomain(ctx.Request, val.Domain_id))
			return
		}
	}

	msg, err := this.models.Delete(allrole)
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 418, i18n.Get(ctx.Request, msg))
		return
	}
	hret.Success(ctx.ResponseWriter, i18n.Success(ctx.Request))
}

// swagger:operation PUT /v1/auth/role/put roleController roleController
//
// 更新角色信息
//
// 更新某个域中的角色信息,角色编码不能更新
//
// ---
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// parameters:
// - name: domain_id
//   in: query
//   description: domain code number
//   required: true
//   type: string
//   format:
// responses:
//   '200':
//     description: success
func (this roleController) Update(ctx *context.Context) {
	ctx.Request.ParseForm()
	if !hrpc.BasicAuth(ctx.Request) {
		hret.Error(ctx.ResponseWriter, 403, i18n.NoAuth(ctx.Request))
		return
	}

	form := ctx.Request.Form
	Role_id := form.Get("Role_id")

	did, err := utils.SplitDomain(Role_id)
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 423, i18n.NoSeparator(ctx.Request, Role_id))
	}

	if !hrpc.DomainAuth(ctx.Request, did, "w") {
		hret.Error(ctx.ResponseWriter, 403, i18n.Get(ctx.Request, "as_of_date_domain_permission_denied_modify"))
		return
	}

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := jwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 403, i18n.Disconnect(ctx.Request))
		return
	}

	msg, err := this.models.Update(form, jclaim.UserId)
	if err != nil {
		logs.Error(err.Error())
		hret.Error(ctx.ResponseWriter, 421, i18n.Get(ctx.Request, msg), err)
		return
	}
	hret.Success(ctx.ResponseWriter, i18n.Success(ctx.Request))
}

func init() {
	groupcache.RegisterStaticFile("AsofdateRolePage", "./views/hauth/role_info_page.tpl")
}
