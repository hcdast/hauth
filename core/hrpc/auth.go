/*
 * @Author: hc
 * @Date: 2021-06-01 17:36:46
 * @LastEditors: hc
 * @LastEditTime: 2021-06-03 11:09:02
 * @Description:
 */
package hrpc

// hrpc package
// this package provide permissions related function
import (
	"net/http"

	"example-hauth/utils/jwt"
	"example-hauth/utils/logs"
	"example-hauth/utils/validator"

	"github.com/hzwy23/dbobj"
)

// 校验用户是否有权限访问当前API
func BasicAuth(r *http.Request) bool {
	cookie, _ := r.Cookie("Authorization")
	jclaim, err := jwt.ParseJwt(cookie.Value)
	if err != nil {
		logs.Error(err)
		return false
	}
	if jclaim.UserId == "admin" {
		return true
	}
	cnt := 0
	err = dbobj.QueryRow(sys_rdbms_hrpc_006, jclaim.UserId, r.URL.Path).Scan(&cnt)
	if err != nil {
		logs.Error(err)
		return false
	}
	if cnt == 0 {
		logs.Error("insufficient privileges", "user id is :", jclaim.UserId, "api is :", r.URL.Path)
		return false
	}
	return true
}

// 检查用户对指定的域的权限
// 第一个参数中,http.Request,包含了用户的连接信息,cookie中.
// 第二个参数中,domain_id,是用户想要访问的域
// 第三个参数是访问模式,r 表示 只读, w 表示 读写.
// 如果返回true,表示用户有权限
// 返回false,表示用户没有权限
func DomainAuth(req *http.Request, domain_id string, pattern string) bool {
	if validator.IsEmpty(domain_id) {
		return false
	}

	level := checkDomainAuthLevel(req, domain_id)
	switch pattern {
	case "r":
		if level != -1 {
			return true
		} else {
			return false
		}
	case "w":
		if level != 2 {
			return false
		} else {
			return true
		}
	default:
		return false
	}
}

func IsRoot(domainId string) bool {
	if domainId == "vertex_root" {
		return true
	}
	return false
}
