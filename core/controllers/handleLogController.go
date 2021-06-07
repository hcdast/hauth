package controllers

import (
	"os"
	"path/filepath"

	"example-hauth/core/groupcache"
	"example-hauth/core/hrpc"
	"example-hauth/core/models"
	"example-hauth/utils/hret"
	"example-hauth/utils/i18n"
	"example-hauth/utils/jwt"
	"example-hauth/utils/logs"

	"github.com/astaxie/beego/context"
	"github.com/tealeg/xlsx"
)

type handleLogsController struct {
	model models.HandleLogMode
}

var HandleLogsCtl = &handleLogsController{}

// swagger:operation GET /v1/auth/HandleLogsPage StaticFiles handleLogsController
//
// 操作日志页面
//
// The system will check user permissions.
// So,you must first login system,and then you can send the request.
//
// ---
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// responses:
//   '200':
//     description: success
//   '404':
//     description: page not found
func (this *handleLogsController) Page(ctx *context.Context) {
	ctx.Request.ParseForm()

	if !hrpc.BasicAuth(ctx.Request) {
		hret.Error(ctx.ResponseWriter, 403, i18n.NoAuth(ctx.Request))
		return
	}

	rst, err := groupcache.GetStaticFile("AsofdateHandleLogPage")
	if err != nil {
		hret.Error(ctx.ResponseWriter, 404, i18n.PageNotFound(ctx.Request))
		return
	}
	ctx.ResponseWriter.Write(rst)
}

// swagger:operation GET /v1/auth/handle/logs/download handleLogsController handleLogsController
//
// 下载日志记录,返回excel格式数据
//
// API将会返回用户所属域中的所有操作记录信息.所以,在使用这个API时,必须登录系统.
//
// ---
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// - application/vnd.ms-excel
// responses:
//   '200':
//     description: success
//   '403':
//     description: Insufficient permissions
//   '421':
//     description: query logs information failed.
func (this handleLogsController) Download(ctx *context.Context) {
	ctx.Request.ParseForm()

	if !hrpc.BasicAuth(ctx.Request) {
		hret.Error(ctx.ResponseWriter, 403, i18n.NoAuth(ctx.Request))
		return
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/vnd.ms-excel")

	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := jwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 403, i18n.Disconnect(ctx.Request))
		return
	}
	rst, err := this.model.Download(jclaim.DomainId)
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 421, i18n.Get(ctx.Request, "error_handle_logs_get_failed"))
		return
	}

	file, err := xlsx.OpenFile(filepath.Join(os.Getenv("HBIGDATA_HOME"), "views", "uploadTemplate", "hauthHandleLogsTemplate.xlsx"))
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 421, i18n.Get(ctx.Request, "error_handle_logs_open_error"), err)
		return
	}
	sheet, ok := file.Sheet["handle_logs"]
	if !ok {
		hret.Error(ctx.ResponseWriter, 421, i18n.Get(ctx.Request, "error_handle_logs_sheet_error"))
		return
	}

	for _, v := range rst {
		row := sheet.AddRow()
		cell1 := row.AddCell()
		cell1.Value = v.User_id

		cell2 := row.AddCell()
		cell2.Value = v.Handle_time

		cell3 := row.AddCell()
		cell3.Value = v.Client_ip

		cell4 := row.AddCell()
		cell4.Value = v.Method

		cell5 := row.AddCell()
		cell5.Value = v.Url

		cell6 := row.AddCell()
		cell6.Value = v.Status_code

		cell7 := row.AddCell()
		cell7.Value = v.Data
	}

	file.Write(ctx.ResponseWriter)
}

// swagger:operation GET /v1/auth/handle/logs handleLogsController handleLogsController
//
// 查询用户所属域中的操作日志信息
//
// API只能查询用户所属域的操作日志信息, 数据处理中,采用了分页查询,所以,必须传入2个参数,分别是:
//
// offset: 起始行数
//
// limit : 最大行数
// ---
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// parameters:
// - name: offset
//   in: query
//   description: 起始行数,必须是数字.
//   required: true
//   type: integer
//   format:
// - name: limit
//   in: query
//   description: 最大行数,必须是数字.
//   required: true
//   type: integer
//   format:
// responses:
//   '200':
//      description: success
//   '403':
//      description: Insufficient permissions
//   '421':
//      description: query logs information failed.
func (this handleLogsController) GetHandleLogs(ctx *context.Context) {
	ctx.Request.ParseForm()

	// Check the user permissions
	if !hrpc.BasicAuth(ctx.Request) {
		hret.Error(ctx.ResponseWriter, 403, i18n.NoAuth(ctx.Request))
		return
	}

	// Get form data from client request.
	offset := ctx.Request.FormValue("offset")
	limit := ctx.Request.FormValue("limit")

	// Get user connection information from cookie.
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := jwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 403, i18n.Disconnect(ctx.Request))
		return
	}

	rst, total, err := this.model.Get(jclaim.DomainId, offset, limit)
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 421, i18n.Get(ctx.Request, "error_handle_logs_query_failed"))
		return
	}
	hret.BootstrapTableJson(ctx.ResponseWriter, total, rst)
}

// swagger:operation GET /v1/auth/handle/logs/search handleLogsController handleLogsController
//
// 返回满足用户搜索条件的日志信息
//
// API中会校验用户的权限,如果用户没有登录,将返回权限不足的提示信息
//
// 这个API需要提供3个参数,分别是:
//
// UserId    : 用户账号
//
// StartDate : 日志操作开始日期
//
// EndDate   : 日志操作结束日期
// ---
// produces:
// - application/json
// - application/xml
// - text/xml
// - text/html
// parameters:
// - name: UserId
//   in: query
//   description: domain code number
//   required: true
//   type: string
//   format:
// - name: StartDate
//   in: query
//   description: domain code number
//   required: true
//   type: string
//   format:
// - name: EndDate
//   in: query
//   description: domain code number
//   required: true
//   type: string
//   format:
// responses:
//   '200':
//     description: success
//   '403':
//     description: Insufficient permissions
//   '421':
//     description: query logs information failed.
func (this handleLogsController) SerachLogs(ctx *context.Context) {
	ctx.Request.ParseForm()

	// Check the user permissions
	if !hrpc.BasicAuth(ctx.Request) {
		hret.Error(ctx.ResponseWriter, 403, i18n.NoAuth(ctx.Request))
		return
	}

	// Get form data from request.
	userid := ctx.Request.FormValue("UserId")
	start := ctx.Request.FormValue("StartDate")
	end := ctx.Request.FormValue("EndDate")

	// get user connection information from cookie
	cookie, _ := ctx.Request.Cookie("Authorization")
	jclaim, err := jwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 403, i18n.Disconnect(ctx.Request))
		return
	}

	rst, err := this.model.Search(jclaim.DomainId, userid, start, end)
	if err != nil {
		logs.Error(err)
		hret.Error(ctx.ResponseWriter, 421, i18n.Get(ctx.Request, "error_handle_logs_query_failed"))
		return
	}
	hret.Json(ctx.ResponseWriter, rst)
}

func init() {
	groupcache.RegisterStaticFile("AsofdateHandleLogPage", "./views/hauth/handle_logs_page.tpl")
}
