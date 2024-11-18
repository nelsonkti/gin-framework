# 打包命令
xxxxxx代表相应的服务名称

make build service=xxxxxx
# 先执行打包命令后才能执行下面命令

# 启动命令分别有

restart: 重启

start:  后台启动

stop:  停止

debug:  调试启动

status：查看运行中的服务

用法：make debug ms-api-X service=xxxxxx

比如1服：make debug ms-api-1 service=xxxxxx

# 状态查看：

make status ms-api-X service=xxxxx

## 例如出现：
```sh
PID   USER     TIME  COMMAND
    1 root      0:50 tail -f /dev/null
    7 root      0:00 /bin/sh
  583 root      0:00 /bin/sh
 1202 root      5:49 /app/main/main
 1220 root      0:00 ps -ef
```

/app/main/main代表正在运行main服务，/app/main/xxxxxx就是运行中的服务名，TIME是运行时间