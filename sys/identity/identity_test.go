package identity

import (
	"testing"
)

func Test_myip(t *testing.T) {
	ips, err := MyIp4List()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ips)

}

func Benchmark_myip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MyIp4List()
	}
}
