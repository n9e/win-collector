package identity

import (
	"testing"
)

func Test_myip(t *testing.T) {
	ips, err := myIp4List()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ips)

}
