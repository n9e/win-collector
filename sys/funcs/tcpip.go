package funcs

import (
	"time"

	"github.com/StackExchange/wmi"
	"github.com/didi/nightingale/src/dataobj"
	"github.com/toolkits/pkg/logger"
)

var (
	historyTcpipStat *Win32_TCPPerfFormattedData
	lasttcpipTime    time.Time
)

func TcpipMetrics() (L []*dataobj.MetricValue) {
	tcpipStat, err := TcpipCounters()
	if err != nil || len(tcpipStat) == 0 {
		logger.Error("Get tcpip data fail: ", err)
		return
	}
	newTcpipStat := &tcpipStat[0]
	now := time.Now()
	interval := now.Unix() - lastTime.Unix()
	lasttcpipTime = now

	if historyTcpipStat == nil {
		historyTcpipStat = newTcpipStat
		return []*dataobj.MetricValue{}
	}
	v := float64(newTcpipStat.ConnectionFailures-historyTcpipStat.ConnectionFailures) / float64(interval)
	if v < 0 {
		v = 0
	}
	L = append(L, GaugeValue("sys.net.tcp.ip4.con.failures", v))

	v = float64(newTcpipStat.ConnectionsPassive-historyTcpipStat.ConnectionsPassive) / float64(interval)
	if v < 0 {
		v = 0
	}
	L = append(L, GaugeValue("sys.net.tcp.ip4.con.passive", v))

	v = float64(newTcpipStat.ConnectionsReset-historyTcpipStat.ConnectionsReset) / float64(interval)
	if v < 0 {
		v = 0
	}
	L = append(L, GaugeValue("sys.net.tcp.ip4.con.reset", v))

	v = float64(newTcpipStat.ConnectionsActive-historyTcpipStat.ConnectionsActive) / float64(interval)
	if v < 0 {
		v = 0
	}
	L = append(L, GaugeValue("sys.net.tcp.ip4.con.active", v))

	v = float64(newTcpipStat.ConnectionsEstablished+historyTcpipStat.ConnectionsEstablished) / 2
	L = append(L, GaugeValue("sys.net.tcp.ip4.con.established", v))

	return
}

type Win32_TCPPerfFormattedData struct {
	ConnectionFailures     uint64
	ConnectionsActive      uint64
	ConnectionsPassive     uint64
	ConnectionsEstablished uint64
	ConnectionsReset       uint64
}

func TcpipCounters() ([]Win32_TCPPerfFormattedData, error) {
	var dst []Win32_TCPPerfFormattedData
	err := wmi.Query("SELECT ConnectionFailures,ConnectionsActive,ConnectionsPassive,ConnectionsEstablished,ConnectionsReset FROM Win32_PerfRawData_Tcpip_TCPv4", &dst)
	return dst, err
}
