/*
 * @Author: hc
 * @Date: 2021-06-01 17:36:47
 * @LastEditors: hc
 * @LastEditTime: 2021-06-08 13:57:59
 * @Description:
 */
package hret

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example-hauth/utils/logs"
)

type HttpOkMsg struct {
	Version    string      `json:"version"`
	Reply_code int         `json:"reply_code"`
	Reply_msg  string      `json:"reply_msg"`
	Data       interface{} `json:"data,omitempty"`
	Total      int64       `json:"total,omitempty"`
	Rows       interface{} `json:"rows,omitempty"`
}

type HttpErrMsg struct {
	Error_code    int         `json:"error_code"`
	Error_msg     string      `json:"error_msg"`
	Error_details interface{} `json:"error_details,omitempty"`
	Version       string      `json:"version"`
}

func Json(w http.ResponseWriter, data interface{}) ([]byte, error) {
	ijs, err := json.Marshal(data)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(`{error_code:765,error_msg:"` + err.Error() + `",error_details:"format json type info failed."}`))
		return ijs, err
	}
	if string(ijs) == "null" {
		w.Write([]byte("[]"))
		return ijs, nil
	}
	_, err = w.Write(ijs)
	return ijs, err
}

// 请求错误处理
func Error(w http.ResponseWriter, code int, msg string, details ...interface{}) {
	e := HttpErrMsg{
		Version:       "v1.0",
		Error_code:    code,
		Error_msg:     msg,
		Error_details: details,
	}
	ijs, err := json.Marshal(e)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(`{error_code:` + strconv.Itoa(e.Error_code) + `,error_msg:"` + e.Error_msg + `",error_details:"format json type info failed."}`))
		return
	}
	w.WriteHeader(e.Error_code)
	w.Write(ijs)
}

// 请求成功处理
func Success(w http.ResponseWriter, v interface{}) {
	ok := HttpOkMsg{
		Version:    "v1.0",
		Reply_code: 200,
		Reply_msg:  "execute successfully.",
		Data:       v,
	}
	ojs, err := json.Marshal(ok)
	if err != nil {
		logs.Error(err.Error())
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(`{error_code:` + strconv.Itoa(http.StatusExpectationFailed) + `,error_msg:"format json type info failed.",error_details:"format json type info failed."}`))
		return
	}
	w.Write(ojs)
}

func BootstrapTableJson(w http.ResponseWriter, total int64, v interface{}) {
	ok := HttpOkMsg{
		Version:    "v1.0",
		Reply_code: 200,
		Reply_msg:  "execute successfully.",
		Rows:       v,
		Total:      total,
	}

	ijs, err := json.Marshal(ok)
	if err != nil {
		logs.Error(err.Error())
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(`{error_code:` + strconv.Itoa(http.StatusExpectationFailed) + `,error_msg:"format json type info failed.",error_details:"format json type info failed."}`))
		return
	}
	w.Write(ijs)
}
