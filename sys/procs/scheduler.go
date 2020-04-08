package procs

import (
	"strings"
	"time"

	"github.com/shirou/gopsutil/process"
	"github.com/toolkits/pkg/logger"

	"github.com/n9e/win-collector/sys/funcs"
	"github.com/n9e/win-collector/sys/identity"

	"github.com/didi/nightingale/src/dataobj"
	"github.com/didi/nightingale/src/model"
)

type ProcScheduler struct {
	Ticker *time.Ticker
	Proc   *model.ProcCollect
	Quit   chan struct{}
}

func NewProcScheduler(p *model.ProcCollect) *ProcScheduler {
	scheduler := ProcScheduler{Proc: p}
	scheduler.Ticker = time.NewTicker(time.Duration(p.Step) * time.Second)
	scheduler.Quit = make(chan struct{})
	return &scheduler
}

func (this *ProcScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-this.Ticker.C:
				ProcCollect(this.Proc)
			case <-this.Quit:
				this.Ticker.Stop()
				return
			}
		}
	}()
}

func (this *ProcScheduler) Stop() {
	close(this.Quit)
}

type P struct {
	Pid     int32
	Name    string
	Cmdline string
}

func ProcCollect(p *model.ProcCollect) {
	startTime := time.Now()
	ps, err := procs()
	if err != nil {
		logger.Error(err)
		return
	}
	endTime := time.Now()
	logger.Debugf("collect procs complete. Process time %s. Number of procs is %d", endTime.Sub(startTime), len(ps))

	pslen := len(ps)
	cnt := 0
	for i := 0; i < pslen; i++ {
		if isProc(ps[i], p.CollectMethod, p.Target) {
			cnt++
		}
	}

	item := funcs.GaugeValue("proc.num", cnt, p.Tags)
	item.Step = int64(p.Step)
	item.Timestamp = time.Now().Unix()
	item.Endpoint = identity.Identity

	funcs.Push([]*dataobj.MetricValue{item})
}

func isProc(p P, method, target string) bool {
	if method == "name" && target == p.Name {
		return true
	} else if (method == "cmdline" || method == "cmd") && strings.Contains(p.Cmdline, target) {
		return true
	}
	return false
}

func procs() ([]P, error) {
	var processes = []P{}
	var PROCESS P
	pids, err := process.Pids()
	if err != nil {
		return processes, err
	}
	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err == nil {
			pname, err := p.Name()
			pcmdline, err := p.Cmdline()
			if err == nil {
				PROCESS.Name = pname
				PROCESS.Cmdline = pcmdline
				PROCESS.Pid = pid
				processes = append(processes, PROCESS)
			}
		}
	}
	return processes, err
}
