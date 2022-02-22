# 介绍

演示如何在代码中使用lego，Fork from [Dcert](https://github.com/DuKanghub/dcert) 代码。去除aliyun相关代码，改为 manual 校验方式。

# 使用

## 修改配置

将代码config 目录的`config.yaml.sample`改名为`config.yaml`，并在对应字段填上自己的信息。

- Email: 邮箱，用于注册账号。
- SavePath：保存目录。用于保存用户信息和签发的证书。
- CADirURL：CA目录地址，可以理解成申请证书的apiserver。(可选)
- DNSNameservers：用于校验 TXT Record 的DNS服务器地址。

# 申请证书

打开main.go

```go//  
定义你要申请证书的域名//定义你要申请证书的域名

domains := []string{"dev.dukanghub.com"}	
// 如果要申请通配符证书，可如下定义
// domains := []string{"dukanghub.com", "*.dukanghub.com"}
// 使用dns验证
challenge := "dns"	//验证方式，可选值：dns, http
// 证书签发
cmd.GetSSLCerts(challenge, domains)
// 证书续期
cmd.RenewSSLCerts(challenge, domains)
```

如果使用http验证，需要将域名解析到公网可访问的服务器上，然后将`.well-known`反向代理至运行本程序的ip 18888 端口。

证书如果签发成功。会保存在`${SavePath}/certificates`目录下。

假设申请域名为`dev.dukanghub.com`，则证书是在这个目录下的`dev.dukanghub.com.crt`或`dev.dukanghub.com.pem`(两个都可以，二者选其一，看自己喜欢)，私钥是`dev.dukanghub.com.key`，这两个文件就是web服务器要用到的ssl证书。



