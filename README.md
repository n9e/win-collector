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


## 运行
管理员权限直接运行 `win-collector.exe` ([下载win-collector](https://github.com/n9e/win-collector/releases)) 即可。配置文件在 `etc/address.yml` 和 `etc/win-collector.yml` 内

## 注册为服务
可以使用 [nssm](https://nssm.cc/) 将其注册为服务
