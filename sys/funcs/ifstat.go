package funcs

import (
	"fmt"
	"strings"
	"time"

	"github.com/n9e/win-collector/sys"

	"github.com/StackExchange/wmi"
	"github.com/didi/nightingale/src/dataobj"
	"github.com/shirou/gopsutil/net"
	"github.com/toolkits/pkg/logger"
)

func netIoCountersStat(ifacePrefix []string) ([]net.IOCountersStat, error) {
	netIOCounter, err := net.IOCounters(true)
	netIfs := []net.IOCountersStat{}
	for _, iface := range ifacePrefix {
		for _, netIf := range netIOCounter {
			logger.Debug("discover netif names: ", netIf.Name)
			if strings.Contains(netIf.Name, iface) {
				logger.Debug("match netif names: ", netIf.Name)
				netIfs = append(netIfs, netIf)
			}
		}
	}
	return netIfs, err
}

type Win32_NetworkAdapter struct {
	Name            string
	NetConnectionID string
	Speed           uint64
	NetEnabled      bool
}

func NetworkAdapterData(NetConnectionID string) ([]Win32_NetworkAdapter, error) {

	var dst []Win32_NetworkAdapter
	Query := `SELECT 
				Name,
				NetConnectionID,
				Speed,
				NetEnabled
			 	FROM  Win32_NetworkAdapter
				WHERE NetConnectionID = '` + NetConnectionID + `'`
	err := wmi.Query(Query, &dst)

	return dst, err
}

type CumIfStat struct {
	inBytes    uint64
	outBytes   uint64
	inPackets  uint64
	outPackets uint64
	inDrop     uint64
	outDrop    uint64
	inErr      uint64
	outErr     uint64
	speed      uint64
}

var (
	historyIfStat map[string]CumIfStat
	lastTime      time.Time
)

func netMetrics(ifacePrefix []string) (ret []*dataobj.MetricValue) {
	netIfs, err := netIoCountersStat(ifacePrefix)
	if err != nil {
		logger.Error("get network iocounters failed: ", err)
		return []*dataobj.MetricValue{}
	}
	now := time.Now()
	newIfStat := make(map[string]CumIfStat)
	for _, netIf := range netIfs {
		ifstat := CumIfStat{
			inBytes:    netIf.BytesRecv,
			outBytes:   netIf.BytesSent,
			inPackets:  netIf.PacketsRecv,
			outPackets: netIf.PacketsSent,
			inDrop:     netIf.Dropin,
			outDrop:    netIf.Dropout,
			inErr:      netIf.Errin,
			outErr:     netIf.Errout,
		}
		networkAdapter, err := NetworkAdapterData(netIf.Name)
		if err != nil || len(networkAdapter) == 0 {
			logger.Error("get network adapter speed failed: ", netIf.Name, err)
		} else {
			if !networkAdapter[0].NetEnabled {
				logger.Debug("network adapter is not enbaled, ignore: ", netIf.Name)
				continue
			}
			ifstat.speed = networkAdapter[0].Speed
		}

		newIfStat[netIf.Name] = ifstat
	}
	interval := now.Unix() - lastTime.Unix()
	lastTime = now

	var totalBandwidth uint64 = 0
	inTotalUsed := 0.0
	outTotalUsed := 0.0

	if historyIfStat == nil {
		historyIfStat = newIfStat
		return []*dataobj.MetricValue{}
	}
	for iface, stat := range newIfStat {
		tags := fmt.Sprintf("iface=%s", strings.Replace(iface, " ", "", -1))
		oldStat := historyIfStat[iface]
		inbytes := float64(stat.inBytes-oldStat.inBytes) / float64(interval)
		if inbytes < 0 {
			inbytes = 0
		}

		inbits := inbytes * 8
		ret = append(ret, GaugeValue("net.in.bits", inbits, tags))

		outbytes := float64(stat.outBytes-oldStat.outBytes) / float64(interval)
		if outbytes < 0 {
			outbytes = 0
		}
		outbits := outbytes * 8
		ret = append(ret, GaugeValue("net.out.bits", outbits, tags))

		v := float64(stat.inDrop-oldStat.inDrop) / float64(interval)
		if v < 0 {
			v = 0
		}
		ret = append(ret, GaugeValue("net.in.dropped", v, tags))

		v = float64(stat.outDrop-oldStat.outDrop) / float64(interval)
		if v < 0 {
			v = 0
		}
		ret = append(ret, GaugeValue("net.out.dropped", v, tags))

		v = float64(stat.inPackets-oldStat.inPackets) / float64(interval)
		if v < 0 {
			v = 0
		}
		ret = append(ret, GaugeValue("net.in.pps", v, tags))

		v = float64(stat.outPackets-oldStat.outPackets) / float64(interval)
		if v < 0 {
			v = 0
		}
		ret = append(ret, GaugeValue("net.out.pps", v, tags))

		v = float64(stat.inErr-oldStat.inErr) / float64(interval)
		if v < 0 {
			v = 0
		}
		ret = append(ret, GaugeValue("net.in.errs", v, tags))

		v = float64(stat.outErr-oldStat.outErr) / float64(interval)
		if v < 0 {
			v = 0
		}
		ret = append(ret, GaugeValue("net.out.errs", v, tags))

		if strings.HasPrefix(iface, "vnet") { //vnet采集到的stat.speed不准确，不计算percent
			continue
		}

		inTotalUsed += inbits

		inPercent := float64(inbits) * 100 / float64(stat.speed)

		if inPercent < 0 || stat.speed <= 0 {
			ret = append(ret, GaugeValue("net.in.percent", 0, tags))
		} else {
			ret = append(ret, GaugeValue("net.in.percent", inPercent, tags))
		}

		outTotalUsed += outbits
		outPercent := float64(outbits) * 100 / float64(stat.speed)
		if outPercent < 0 || stat.speed <= 0 {
			ret = append(ret, GaugeValue("net.out.percent", 0, tags))
		} else {
			ret = append(ret, GaugeValue("net.out.percent", outPercent, tags))
		}

		ret = append(ret, GaugeValue("net.bandwidth.mbits", stat.speed, tags))
		totalBandwidth += stat.speed
	}

	ret = append(ret, GaugeValue("net.bandwidth.mbits.total", totalBandwidth))
	ret = append(ret, GaugeValue("net.in.bits.total", inTotalUsed))
	ret = append(ret, GaugeValue("net.out.bits.total", outTotalUsed))

	if totalBandwidth <= 0 {
		ret = append(ret, GaugeValue("net.in.bits.total.percent", 0))
		ret = append(ret, GaugeValue("net.out.bits.total.percent", 0))
	} else {
		inTotalPercent := float64(inTotalUsed) / float64(totalBandwidth) * 100
		ret = append(ret, GaugeValue("net.in.bits.total.percent", inTotalPercent))

		outTotalPercent := float64(outTotalUsed) / float64(totalBandwidth) * 100
		ret = append(ret, GaugeValue("net.out.bits.total.percent", outTotalPercent))
	}

	historyIfStat = newIfStat
	return ret
}

func NetMetrics() (ret []*dataobj.MetricValue) {
	return netMetrics(sys.Config.IfacePrefix)
}
