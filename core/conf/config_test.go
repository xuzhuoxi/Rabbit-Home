// Create on 2025/2/12
// @author xuzhuoxi
package conf

import (
	"gopkg.in/yaml.v2"
	"testing"
)

const cfgStr = `allows_on: false
allows:
blocks_on: true
blocks:
  - "192.168.0.1"
  - "10.0.0.1-20"
  - "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
  - "2001:db8:85a3::8a2e:370:7334"
  - "2001:0db8:85a3:0000:0000:8a2e:0370:1-7334"
  - "2001:0db8:85a3:0000:0000:8a2e:0370:1-7334"`

var (
	ipAdders = []string{
		"192.168.0.1",
		"192.168.0.2",
		"10.0.0.1",
		"10.0.0.20",
		"2001:0db8:85a3:0000:0000:8a2e:0370:0",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"2001:db8:85a3::8a2e:370:0",
		"2001:db8:85a3::8a2e:370:7334",
	}
)

func TestIpVerifier(t *testing.T) {
	verifier := &IpVerifier{}
	err := yaml.Unmarshal([]byte(cfgStr), verifier)
	if nil != err {
		t.Fatal(err)
	}
	verifier.PreProcess()
	for i, addr := range ipAdders {
		t.Logf("Verify Ok(%v)! index:%d, addr:%s", verifier.CheckIpAddr(addr), i, addr)
	}
}
