package funcs

import (
	"encoding/json"
	"testing"
	"time"
)

func Test_cpu(t *testing.T) {
	d := time.Duration(3) * time.Second
	if err := UpdateCpuStat(); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(d)
	if err := UpdateCpuStat(); err != nil {
		t.Error(err)
		return
	}
	for _, m := range CpuMetrics() {
		bs, _ := json.Marshal(m)
		t.Log(string(bs))
	}
}

func Test_dfstat(t *testing.T) {
	ignore := []string{"C"}
	for _, m := range deviceMetrics(ignore) {
		bs, _ := json.Marshal(m)
		t.Log(string(bs))
	}
}

func Test_DiskIOStat(t *testing.T) {
	d := time.Duration(3) * time.Second
	if err := UpdateDiskStats(); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(d)
	if err := UpdateDiskStats(); err != nil {
		t.Error(err)
		return
	}
	for _, m := range IOStatsMetrics() {
		bs, _ := json.Marshal(m)
		t.Log(string(bs))
	}
}

func Test_IfStat(t *testing.T) {
	ifacePrefix := []string{"以太网"}
	d := time.Duration(3) * time.Second
	time.Sleep(d)
	netMetrics(ifacePrefix)
	for _, m := range netMetrics(ifacePrefix) {
		bs, _ := json.Marshal(m)
		t.Log(string(bs))
	}
}

func Test_MemInfo(t *testing.T) {
	for _, m := range MemMetrics() {
		bs, _ := json.Marshal(m)
		t.Log(string(bs))
	}
}

func Test_Ntp(t *testing.T) {
	ntpServers := []string{"ntp.aliyun.com"}
	for _, m := range ntpOffsetMetrics(ntpServers) {
		bs, _ := json.Marshal(m)
		t.Log(string(bs))
	}
}
