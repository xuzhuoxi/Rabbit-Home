// Create on 2023/6/5
// @author xuzhuoxi
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/src/core"
	"github.com/xuzhuoxi/infra-go/netx"
	"testing"
)

var (
	server1 = core.LinkEntity{
		Id:         "1001",
		PlatformId: "Test01",
		Name:       "TestServer",
		Network:    "http",
		Addr:       "http://192.168.1.1",
	}
)

func TestLink(t *testing.T) {
	bs, _ := json.Marshal(server1)
	data := base64.StdEncoding.EncodeToString(bs)
	url := fmt.Sprintf("http://127.0.0.1:9000/link?data=%s", data)
	netx.HttpGet(url, nil)
}
