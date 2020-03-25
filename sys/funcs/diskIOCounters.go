package funcs

import (
	"github.com/StackExchange/wmi"
	"github.com/toolkits/pkg/logger"
)

type Win32_PerfFormattedData struct {
	AvgDiskSecPerRead_Base  uint32
	AvgDiskSecPerWrite_Base uint32
	DiskReadBytesPerSec     uint64
	DiskReadsPerSec         uint32
	DiskWriteBytesPerSec    uint64
	DiskWritesPerSec        uint32
	Name                    string
}

type Win32_PerfFormattedData_IDLE struct {
	Name            string
	PercentIdleTime uint64
}

type diskIOCounter struct {
	Device        string
	MsecRead      uint64
	MsecWrite     uint64
	ReadBytes     uint64
	ReadRequests  uint64
	WriteBytes    uint64
	WriteRequests uint64
	Util          uint64
}

func PerfFormattedData() ([]Win32_PerfFormattedData, error) {

	var dst []Win32_PerfFormattedData
	Query := `SELECT 
				AvgDiskSecPerRead_Base,
				AvgDiskSecPerWrite_Base,
				DiskReadBytesPerSec,
				DiskReadsPerSec,
				DiskWriteBytesPerSec,
				DiskWritesPerSec,
				Name
				FROM Win32_PerfRawData_PerfDisk_PhysicalDisk`
	err := wmi.Query(Query, &dst)

	return dst, err
}

func PerfFormattedDataIDLE(name string) ([]Win32_PerfFormattedData_IDLE, error) {

	var dst []Win32_PerfFormattedData_IDLE
	query := `SELECT PercentIdleTime FROM Win32_PerfFormattedData_PerfDisk_PhysicalDisk WHERE Name = '` + name + `'`
	err := wmi.Query(query, &dst)

	return dst, err
}

func IOCounters() ([]diskIOCounter, error) {
	ret := []diskIOCounter{}
	dst, err := PerfFormattedData()
	if err != nil {
		return nil, err
	}
	for _, d := range dst {

		if d.Name == "_Total" { // not get _Total
			continue
		}

		dstat := diskIOCounter{
			Device:        d.Name,
			MsecRead:      uint64(d.AvgDiskSecPerRead_Base),
			MsecWrite:     uint64(d.AvgDiskSecPerWrite_Base),
			ReadBytes:     d.DiskReadBytesPerSec,
			ReadRequests:  uint64(d.DiskReadsPerSec),
			WriteBytes:    d.DiskWriteBytesPerSec,
			WriteRequests: uint64(d.DiskWritesPerSec),
		}
		idle, err := PerfFormattedDataIDLE(d.Name)
		if err != nil || len(idle) == 0 {
			logger.Error("get disk idle failed: ", err)
		} else {
			dstat.Util = uint64(100) - idle[0].PercentIdleTime
		}
		ret = append(ret, dstat)
	}

	return ret, nil
}
