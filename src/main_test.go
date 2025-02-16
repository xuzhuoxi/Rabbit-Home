// Create on 2023/6/5
// @author xuzhuoxi
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/client"
	"github.com/xuzhuoxi/infra-go/netx/httpx"
	"testing"
)

var (
	homeUrl = "http://127.0.0.1:9000"
	server1 = core.LinkEntity{
		Id:         "1001",
		PlatformId: "Test01",
		Name:       "TestServer",
		Network:    "http",
		Addr:       "http://192.168.1.1",
	}
	weight1 = 1.33
	status  = core.EntityStatus{
		Id:     server1.Id,
		Weight: 2.06,
	}
)

func TestLink(t *testing.T) {
	bs, _ := json.Marshal(server1)
	data := base64.StdEncoding.EncodeToString(bs)
	url := fmt.Sprintf("%s/link?data=%s&w=%v", homeUrl, data, weight1)
	httpx.HttpGet(url, nil)
}

func TestClientLinkGet(t *testing.T) {
	client.LinkWithGet(homeUrl, server1, weight1, nil)
}

func TestClientUpdateGet(t *testing.T) {
	client.UpdateWithGet(homeUrl, status, nil)
}

func TestClientUnlinkGet(t *testing.T) {
	client.UnlinkWithGet(homeUrl, server1.Id, nil)
}

// --------------------------------

func TestClientLinkPost(t *testing.T) {
	client.LinkWithPost(homeUrl, server1, weight1, nil)
}

func TestClientUpdatePost(t *testing.T) {
	client.UpdateWithPost(homeUrl, status, nil)
}

func TestClientUnlinkPost(t *testing.T) {
	client.UnlinkWithPost(homeUrl, server1.Id, nil)
}
