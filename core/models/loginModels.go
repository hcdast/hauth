/*
 * @Author: hc
 * @Date: 2021-06-01 17:36:46
 * @LastEditors: hc
 * @LastEditTime: 2021-06-07 10:56:03
 * @Description:
 */
package models

import (
	"example-hauth/utils/logs"
	"fmt"

	"example-hauth/dbobj"
)

type LoginModels struct {
}

func (this LoginModels) GetDefaultPage(user_id string) string {
	row := dbobj.QueryRow(sys_rdbms_078, user_id)
	var url = "./views/hauth/theme/default/index.tpl"
	err := row.Scan(&url)
	if err != nil {
		logs.Debug("get default theme.")
		url = "./views/hauth/theme/default/index.tpl"
	}
	return url
}

func (this LoginModels) GetDefaultOrgId(user_id string) (org_id string, err error) {
	err = dbobj.QueryRow(sys_rdbms_080, user_id).Scan(&org_id)
	return
}
