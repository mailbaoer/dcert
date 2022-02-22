package main

import (
	"github.com/DuKanghub/dcert/cmd"
)

func main() {
	domains := []string{"example.com", "*.example.com"} // 通配符证书这么传入域名
	challenge := "dns"                                  // 验证方式，可选值：dns, http
	cmd.GetSSLCerts(challenge, domains)
	// cmd.RenewCert(challenge, domains) // 更新证书
}
