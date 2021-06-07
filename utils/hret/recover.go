/*
 * @Author: hc
 * @Date: 2021-06-01 17:36:47
 * @LastEditors: hc
 * @LastEditTime: 2021-06-07 10:10:34
 * @Description: 异常处理 —— Go语言没有异常系统，其使用 panic 触发宕机类似于其他语言的抛出异常，recover 的宕机恢复机制就对应其他语言中的 try/catch 机制。
 */
package hret

import "example-hauth/utils/logs"

type httpPanicFunc func()

// HttpPanic user for stop panic up.
func HttpPanic(f ...httpPanicFunc) {
	if r := recover(); r != nil {
		logs.Error("system generator panic.", r)
		for _, val := range f {
			val()
		}
	}
}
