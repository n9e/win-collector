# windows 服务器基础资源采集 agent

## 功能
系统指标采集

### 指标
指标与 linux 版本基本保持一致，部分指标上有差异。

以下仅列出差异部分

#### cpu
|metric|linux|windows|
|--|--|--|
|cpu.softirq|支持|不支持|
|cpu.steal|支持|不支持|
|cpu.iowait|支持|不支持|
|cpu.nice|支持|不支持|
|cpu.guest|支持|不支持|
|cpu.core.softirq|支持|不支持|
|cpu.core.steal|支持|不支持|
|cpu.core.iowait|支持|不支持|
|cpu.core.nice|支持|不支持|
|cpu.core.guest|支持|不支持|
|cpu.loadavg.1|支持|不支持|
|cpu.loadavg.5|支持|不支持|
|cpu.loadavg.15|支持|不支持|

#### mem
|metric|linux|windows|
|--|--|--|
|mem.bytes.buffers|支持|不支持|
|mem.bytes.cached|支持|不支持|
|mem.swap.*|支持|不支持|


#### disk
|metric|linux|windows|
|--|--|--|
|disk.inodes.*|支持|不支持|
|disk.io.avgrq_sz|支持|不支持|
|disk.io.avgqu_sz|支持|不支持|
|disk.io.svctm|支持|不支持|
|disk.io.read.msec|不支持|支持|
|disk.io.write.msec|不支持|支持|
|disk.rw.error|支持|不支持|


#### net
|metric|linux|windows|
|--|--|--|
|net.sockets.*|支持|不支持|

#### sys
|metric|linux|windows|
|--|--|--|
|sys.net.netfilter.*|支持|不支持|
|sys.net.tcp.ip4.confailures|不支持|支持|
|sys.net.tcp.ip4.confailures|不支持|支持|
|sys.net.tcp.ip4.conpassive|不支持|支持|
|sys.net.tcp.ip4.conestablished|不支持|支持|
|sys.net.tcp.ip4.conreset|不支持|支持|
|sys.ps.entity.total|支持|不支持|

#### log
暂不支持，Todo

#### port 
无差异

#### proc
无差异

注意由于 windows 获取 process 列表性能较差，启用进程监控时记得把采集周期调大一点。

可以开启 debug 来查看获取 process 的开销耗时，务必将采集周期调到采集开销的2倍以上。否则长时间会导致 oom


## 修改配置文件
win-collector.yml
```
logger:
  dir: logs/collector
  level: WARNING
  keepHours: 2
identity:
  specify: "backup.010607001.default.192.168.21.108.xxf"
sys:
  # timeout in ms
  # interval in second
  timeout: 1000
  interval: 20
  ifacePrefix:
    - 以太网
    - 本地连接
  mountPoint: []
  mountIgnorePrefix:
    - /var/lib
  ntpServers:
    - ntp.aliyun.com
  plugin: plugin/
  ignoreMetrics:
    - cpu.core.idle
    - cpu.core.util
    - cpu.core.sys
    - cpu.core.user
    - cpu.core.irq

```
注意：
根据自身需求，如需固定主机名，则可参照此案例。
如specify 留空的话就自动取第一块网卡的 ip，有东西的话用 specify 的值做 endpoint

address.yml
```
---
monapi:
  http: 0.0.0.0:5801
  addresses:
    - 192.168.21.215

transfer:
  http: 0.0.0.0:5810
  rpc: 0.0.0.0:5811
  addresses:
    - 192.168.21.215

collector:
  http: 127.0.0.1:2058
```
注意：monapi的http地址如需要远程调用也就是跨公网跨云调用，需要做IP或者域名映射

## 测试验证
```
打开CMD窗口将win-collecor.exe程序鼠标左键拖到该CMD窗口后，添加启动参数-f指定配置文件路径，如下：
C:\Users\Administrator>D:\win-collector-win64-0.2.1\win-collector.exe -f D:\win-collector-win64-0.2.1\etc\win-collector.yml
2020/09/03 21:31:11 maxprocs: Leaving GOMAXPROCS=4: CPU quota undefined 
collector start, use configuration file: D:\win-collector-win64-0.2.1\etc\win-collector.ymlrunner.cwd: D:\win-collector-win64-0.2.1
2020/09/03 21:31:11 collector.go:64: endpoint: backup.010607001.default.192.168.21.108.xxf
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/collector/ping       --> github.com/n9e/win-collector/http/routes.ping (1 handlers)
[GIN-debug] GET    /api/collector/pid        --> github.com/n9e/win-collector/http/routes.pid (1 handlers)
[GIN-debug] GET    /api/collector/addr       --> github.com/n9e/win-collector/http/routes.addr (1 handlers)
[GIN-debug] GET    /api/collector/stra       --> github.com/n9e/win-collector/http/routes.getStrategy (1 handlers)
[GIN-debug] POST   /api/collector/push       --> github.com/n9e/win-collector/http/routes.pushData (1 handlers)
[GIN-debug] GET    /api/collector/debug/pprof/ --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /api/collector/debug/pprof/cmdline --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /api/collector/debug/pprof/profile --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] POST   /api/collector/debug/pprof/symbol --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /api/collector/debug/pprof/symbol --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /api/collector/debug/pprof/trace --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /api/collector/debug/pprof/allocs --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /api/collector/debug/pprof/block --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /api/collector/debug/pprof/goroutine --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /api/collector/debug/pprof/heap --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /api/collector/debug/pprof/mutex --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
[GIN-debug] GET    /api/collector/debug/pprof/threadcreate --> github.com/gin-contrib/pprof.pprofHandler.func1 (1 handlers)
2020/09/03 21:31:12 http.go:40: starting http server, listening on: 127.0.0.1:2058
```
表示启动正常，再检查服务端是否上报

## 运行
管理员权限直接运行 `win-collector.exe` ([下载win-collector](https://github.com/n9e/win-collector/releases)) 即可。配置文件在 `etc/address.yml` 和 `etc/win-collector.yml` 内

## 注册为服务
可以使用 [nssm](https://nssm.cc/) 将其注册为服务
