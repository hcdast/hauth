/*
 * @Author: hc
 * @Date: 2021-06-01 17:36:46
 * @LastEditors: hc
 * @LastEditTime: 2021-06-03 11:10:24
 * @Description:
 */
package service

import (
	"net/http"

	"example-hauth/utils/jwt"
)

const redirect = `
<script type="text/javascript">
    $.Hconfirm({
		cancelBtn:false,
        header:"连接已断开",
        body:"用户连接已断开，请重新登录",
        callback:function () {
            window.location.href="/"
        }
    })
</script>
`

func CheckConnection(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil || !jwt.CheckToken(cookie.Value) {
		w.Write([]byte(redirect))
	}
}
