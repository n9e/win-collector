package funcs

import (
	"strings"
	"sync"
	"time"

	"github.com/didi/nightingale/src/dataobj"
	"github.com/toolkits/pkg/logger"
)

var (
	diskStatsMap = make(map[string][2]*diskIOCounter)
	dsLock       = new(sync.RWMutex)
)

func PrepareDiskStats() {
	d := time.Duration(3) * time.Second
	for {
		err := UpdateDiskStats()
		if err != nil {
			logger.Error("update disk stats fail", err)
		}
		time.Sleep(d)
	}
}

func UpdateDiskStats() error {
	dsList, err := IOCounters()
	if err != nil {
		return err
	}

	dsLock.Lock()
	defer dsLock.Unlock()
	for i := 0; i < len(dsList); i++ {
		device := dsList[i].Device
		diskStatsMap[device] = [2]*diskIOCounter{&dsList[i], diskStatsMap[device][0]}
	}
	return nil
}

func IOReadRequests(arr [2]*diskIOCounter) uint64 {
	return arr[0].ReadRequests - arr[1].ReadRequests
}

func IOWriteRequests(arr [2]*diskIOCounter) uint64 {
	return arr[0].WriteRequests - arr[1].WriteRequests
}

func IOMsecRead(arr [2]*diskIOCounter) uint64 {
	return arr[0].MsecRead - arr[1].MsecRead
}

func IOMsecWrite(arr [2]*diskIOCounter) uint64 {
	return arr[0].MsecWrite - arr[1].MsecWrite
}

func IOReadBytes(arr [2]*diskIOCounter) uint64 {
	return arr[0].ReadBytes - arr[1].ReadBytes
}

func IOWriteBytes(arr [2]*diskIOCounter) uint64 {
	return arr[0].WriteBytes - arr[1].WriteBytes
}

func IOUtilTotal(arr [2]*diskIOCounter) uint64 {
	return arr[0].Util + arr[1].Util
}

func IODelta(device string, f func([2]*diskIOCounter) uint64) uint64 {
	val, ok := diskStatsMap[device]
	if !ok {
		return 0
	}

	if val[1] == nil {
		return 0
	}
	return f(val)
}

func IOStatsMetrics() (L []*dataobj.MetricValue) {
	dsLock.RLock()
	defer dsLock.RUnlock()

	for device := range diskStatsMap {
		tags := "device=" + strings.Replace(device, " ", "", -1)
		rio := IODelta(device, IOReadRequests)
		wio := IODelta(device, IOWriteRequests)
		rbytes := IODelta(device, IOReadBytes)
		wbytes := IODelta(device, IOWriteBytes)
		ruse := IODelta(device, IOMsecRead)
		wuse := IODelta(device, IOMsecWrite)
		utilTotal := IODelta(device, IOUtilTotal)
		n_io := rio + wio
		await := 0.0
		if n_io != 0 {
			await = float64(ruse+wuse) / float64(n_io)
		}

		L = append(L, GaugeValue("disk.io.read.request", float64(rio), tags))
		L = append(L, GaugeValue("disk.io.write.request", float64(wio), tags))
		L = append(L, GaugeValue("disk.io.read.bytes", float64(rbytes), tags))
		L = append(L, GaugeValue("disk.io.write.bytes", float64(wbytes), tags))
		L = append(L, GaugeValue("disk.io.read.msec", float64(ruse), tags))
		L = append(L, GaugeValue("disk.io.write.msec", float64(wuse), tags))
		L = append(L, GaugeValue("disk.io.await", await, tags))
		L = append(L, GaugeValue("disk.io.util", float64(utilTotal)/2, tags))
	}

	return
}
