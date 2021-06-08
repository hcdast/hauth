/*
 * @Author: hc
 * @Date: 2021-06-01 17:36:46
 * @LastEditors: hc
 * @LastEditTime: 2021-06-08 11:38:34
 * @Description:
 */
package models

import (
	"example-hauth/utils/logs"

	"example-hauth/dbobj"
)

type LoginModels struct {
}

func (LoginModels) GetDefaultPage(user_id string) string {
	row := dbobj.QueryRow(sys_rdbms_078, user_id)
	var url = "./views/hauth/theme/default/index.tpl"
	err := row.Scan(&url)
	if err != nil {
		logs.Debug("get default theme.")
		url = "./views/hauth/theme/default/index.tpl"
	}
	return url
}

func (LoginModels) GetDefaultOrgId(user_id string) (org_id string, err error) {
	err = dbobj.QueryRow(sys_rdbms_080, user_id).Scan(&org_id)
	return
}
