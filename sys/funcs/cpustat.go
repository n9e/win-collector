package funcs

import (
	"fmt"
	"sync"
	"time"

	"github.com/didi/nightingale/src/dataobj"
	"github.com/shirou/gopsutil/cpu"
	"github.com/toolkits/pkg/logger"
)

const (
	historyCount int = 2
)

type CpuStats struct {
	User   float64
	System float64
	Idle   float64
	Irq    float64
	Total  float64
}

type ProcStat struct {
	Cpu  CpuStats
	Cpus []CpuStats
}

var (
	psHistory [historyCount]*ProcStat
	psLock    = new(sync.RWMutex)
)

func CurrentProcStat() (procStat ProcStat, err error) {
	psPer, err := cpu.Times(true)
	if err != nil {
		return
	}
	var cpuTotal CpuStats
	var cpus []CpuStats
	for _, ps := range psPer {
		var c CpuStats
		c.User = ps.User
		cpuTotal.User += c.User

		c.System = ps.System
		cpuTotal.System += c.System

		c.Idle = ps.Idle
		cpuTotal.Idle += c.Idle

		c.Irq = ps.Irq
		cpuTotal.Irq += c.Irq

		c.Total = ps.Total()
		cpuTotal.Total += c.Total

		cpus = append(cpus, c)
	}
	procStat.Cpu = cpuTotal
	procStat.Cpus = cpus
	return
}

func PrepareCpuStat() {
	d := time.Duration(3) * time.Second
	for {
		err := UpdateCpuStat()
		if err != nil {
			logger.Error("update cpu stat fail", err)
		}
		time.Sleep(d)
	}
}

func UpdateCpuStat() error {
	ps, err := CurrentProcStat()
	if err != nil {
		return err
	}
	psLock.Lock()
	defer psLock.Unlock()
	for i := historyCount - 1; i > 0; i-- {
		psHistory[i] = psHistory[i-1]
	}
	psHistory[0] = &ps
	return nil
}

func deltaTotal() float64 {
	if len(psHistory) < 2 {
		return 0
	}
	return psHistory[0].Cpu.Total - psHistory[1].Cpu.Total
}

func CpuIdles() (res []*CpuStats) {
	psLock.RLock()
	defer psLock.RUnlock()
	if len(psHistory) < 2 {
		return
	}
	if len(psHistory[0].Cpus) != len(psHistory[1].Cpus) {
		return
	}
	for i, c := range psHistory[0].Cpus {
		stats := new(CpuStats)
		dt := c.Total - psHistory[1].Cpus[i].Total
		if dt == 0 {
			return
		}
		invQuotient := 100.00 / float64(dt)
		stats.Idle = float64(c.Idle-psHistory[1].Cpus[i].Idle) * invQuotient
		stats.User = float64(c.User-psHistory[1].Cpus[i].User) * invQuotient
		stats.System = float64(c.System-psHistory[1].Cpus[i].System) * invQuotient
		stats.Irq = float64(c.Irq-psHistory[1].Cpus[i].Irq) * invQuotient
		res = append(res, stats)
	}
	return
}

func CpuIdle() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(psHistory[0].Cpu.Idle-psHistory[1].Cpu.Idle) * invQuotient
}

func CpuUser() float64 {
	psLock.Lock()
	defer psLock.Unlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(psHistory[0].Cpu.User-psHistory[1].Cpu.User) * invQuotient
}

func CpuSystem() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(psHistory[0].Cpu.System-psHistory[1].Cpu.System) * invQuotient
}

func CpuIrq() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(psHistory[0].Cpu.Irq-psHistory[1].Cpu.Irq) * invQuotient
}

func CpuPrepared() bool {
	psLock.RLock()
	defer psLock.RUnlock()
	return psHistory[1] != nil
}

func CpuMetrics() []*dataobj.MetricValue {
	if !CpuPrepared() {
		return []*dataobj.MetricValue{}
	}

	var ret []*dataobj.MetricValue

	cpuIdleVal := CpuIdle()
	idle := GaugeValue("cpu.idle", cpuIdleVal)
	util := GaugeValue("cpu.util", 100.0-cpuIdleVal)
	user := GaugeValue("cpu.user", CpuUser())
	system := GaugeValue("cpu.sys", CpuSystem())
	irq := GaugeValue("cpu.irq", CpuIrq())
	ret = []*dataobj.MetricValue{idle, util, user, system, irq}

	idles := CpuIdles()
	for i, stats := range idles {
		tags := fmt.Sprintf("core=%d", i)
		ret = append(ret, GaugeValue("cpu.core.idle", stats.Idle, tags))
		ret = append(ret, GaugeValue("cpu.core.util", 100.0-stats.Idle, tags))
		ret = append(ret, GaugeValue("cpu.core.user", stats.User, tags))
		ret = append(ret, GaugeValue("cpu.core.sys", stats.System, tags))
		ret = append(ret, GaugeValue("cpu.core.irq", stats.Irq, tags))
	}

	return ret
}
