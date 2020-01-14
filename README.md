### 运行
```
git clone https://github.com/swirldawn/go_probe.git
cd go_probe
cp .config.example .config
vi .config
# 编辑你需要监控的服务器
[
	{
		"name":"本地",//显示的名称
		"ip":"127.0.0.1",//服务器ip
		"port":"8084"//请求端口
	},
]
#启动监控
nohup ./main >> /tmp/go.log 2>&1 &
```

##### 主服务器需要配置 .config文件 配置每个服务器
##### 主服务器的域名也需要配置到配置文件里面 不然看不到主服务器的,不想看可以不配置
##### 然后访问主服务器的 /index 路由即可