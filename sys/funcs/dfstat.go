package funcs

import (
	"fmt"
	"strings"

	"github.com/n9e/win-collector/sys"

	"github.com/didi/nightingale/src/dataobj"
	"github.com/shirou/gopsutil/disk"
	"github.com/toolkits/pkg/logger"
)

func diskUsage(path string) (*disk.UsageStat, error) {
	diskUsage, err := disk.Usage(path)
	return diskUsage, err
}

func diskPartitions() ([]disk.PartitionStat, error) {
	diskPartitions, err := disk.Partitions(true)
	return diskPartitions, err
}

func deviceMetrics(ignoreMountPointsPrefix []string) (L []*dataobj.MetricValue) {
	diskPartitions, err := diskPartitions()

	if err != nil {
		logger.Error("get disk stat failed", err)
		return
	}

	var diskTotal uint64 = 0
	var diskUsed uint64 = 0
	for _, device := range diskPartitions {
		if hasIgnorePrefix(device.Mountpoint, ignoreMountPointsPrefix) {
			continue
		}
		du, err := diskUsage(device.Mountpoint)
		if err != nil {
			logger.Error("get disk stat fail", err)
			continue
		}

		diskTotal += du.Total
		diskUsed += du.Used
		tags := fmt.Sprintf("mount=%s,fstype=%s", device.Mountpoint, device.Fstype)
		L = append(L, GaugeValue("disk.bytes.total", du.Total, tags))
		L = append(L, GaugeValue("disk.bytes.used", du.Used, tags))
		L = append(L, GaugeValue("disk.bytes.free", du.Free, tags))
		L = append(L, GaugeValue("disk.bytes.used.percent", du.UsedPercent, tags))
	}

	if len(L) > 0 && diskTotal > 0 {
		L = append(L, GaugeValue("disk.cap.bytes.total", float64(diskTotal)))
		L = append(L, GaugeValue("disk.cap.bytes.used", float64(diskUsed)))
		L = append(L, GaugeValue("disk.cap.bytes.free", float64(diskTotal-diskUsed)))
		L = append(L, GaugeValue("disk.cap.bytes.used.percent", float64(diskUsed)*100.0/float64(diskTotal)))
	}
	return
}

func DeviceMetrics() (L []*dataobj.MetricValue) {
	return deviceMetrics(sys.Config.MountIgnorePrefix)
}

func hasIgnorePrefix(fsFile string, ignoreMountPointsPrefix []string) bool {
	hasPrefix := false
	if len(ignoreMountPointsPrefix) > 0 {
		for _, ignorePrefix := range ignoreMountPointsPrefix {
			if strings.HasPrefix(fsFile, ignorePrefix) {
				hasPrefix = true
				logger.Debugf("mount point %s has ignored prefix %s", fsFile, ignorePrefix)
				break
			}
		}
	}
	return hasPrefix
}
