package procs

import (
	"time"

	"github.com/n9e/win-collector/stra"
)

func Detect() {
	detect()
	go loopDetect()
}

func loopDetect() {
	for {
		time.Sleep(time.Second * 10)
		detect()
	}
}

func detect() {
	ps := stra.GetProcCollects()
	DelNoPorcCollect(ps)
	AddNewPorcCollect(ps)
}
