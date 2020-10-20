package funcs

import (
	"github.com/didi/nightingale/src/common/dataobj"
)

func CollectorMetrics() []*dataobj.MetricValue {
	return []*dataobj.MetricValue{
		GaugeValue("proc.agent.alive", 1),
	}
}
