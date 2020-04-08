package funcs

import (
	"github.com/n9e/win-collector/sys"

	"github.com/didi/nightingale/src/dataobj"
)

type FuncsAndInterval struct {
	Fs       []func() []*dataobj.MetricValue
	Interval int
}

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := sys.Config.Interval
	Mappers = []FuncsAndInterval{
		{
			Fs: []func() []*dataobj.MetricValue{
				CollectorMetrics,
				CpuMetrics,
				MemMetrics,
				NetMetrics,
				IOStatsMetrics,
				NtpOffsetMetrics,
				TcpipMetrics,
			},
			Interval: interval,
		},
		{
			Fs: []func() []*dataobj.MetricValue{
				DeviceMetrics,
			},
			Interval: interval,
		},
	}
}
