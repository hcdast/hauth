package panda

const (
	admin = "admin"
)

// 检查是否为超级管理员
func IsAdmin(str string) bool {
	return str == admin
}
