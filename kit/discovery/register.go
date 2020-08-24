/*
@Time : 2020/8/24 17:18
@Author : ZhouHui2
@File : register
@Software: GoLand
*/
package discovery

import "time"

type RegisterInfo struct {
	Host           string
	Port           int
	ServiceName    string
	UpdateInterval time.Duration
}

type Register interface {
	Register(info RegisterInfo) error
	DeRegister(info RegisterInfo) error
}
