package basic

//Option 配置选项
type Option func(*BasicAuth)

//WithUP 添加用户名密码
func WithUP(userName string, pwd string) Option {
	return func(b *BasicAuth) {
		b.Members[userName] = pwd
	}
}

//WithExcludes 排除的服务或请求
func WithExcludes(p ...string) Option {
	return func(b *BasicAuth) {
		b.Excludes = p
	}
}