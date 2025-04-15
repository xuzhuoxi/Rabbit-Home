// Create on 2023/6/5
// @author xuzhuoxi
package main

import (
	"encoding/base64"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/client"
	"github.com/xuzhuoxi/infra-go/cryptox/asymmetric"
	"github.com/xuzhuoxi/infra-go/cryptox/symmetric"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/mathx"
	"net/http"
	"testing"
	"time"
)

var (
	homeUrl  = "http://127.0.0.1:9000"
	linkInfo = &core.LinkInfo{
		Id:          "1001",
		PlatformId:  "Test01",
		Name:        "TestServer",
		OpenNetwork: "tcp",
		OpenAddr:    "127.0.0.1:41000",
	}
	updateInfo = &core.UpdateInfo{
		Id:     "1001",
		Weight: 1.33,
	}
	detailInfo = &core.UpdateDetailInfo{
		Id:            "1001",
		Links:         1,
		MaxLinks:      10,
		TotalReqCount: int64(time.Second),
		TotalRespTime: int64(time.Minute),
		EnableKeys:    "links,max-links,total-reg,total-resp",
	}
	unlinkInfo = &core.UnlinkInfo{
		Id: "1001",
	}
	weight0        = 1.00
	base64Encoding = base64.RawURLEncoding
)
var (
	done             = make(chan bool)
	privatePemPath   = "keys/private/pkcs8_private.pem"
	publicPemPath    = "keys/public/x509_public.pem"
	privateRsaCipher asymmetric.IRSAPrivateCipher
	publicRsaCipher  asymmetric.IRSAPublicCipher
	tempKey          []byte
	tempAesCipher    symmetric.IAESCipher
)

var (
	funcLink         = client.LinkWithGet
	funcUnlink       = client.UnlinkWithGet
	funcUpdate       = client.UpdateWithGet
	funcUpdateDetail = client.UpdateDetailWithGet
)

var (
	post       = false
	keyEnable  = true
	autoUnlink = true
)

func init() {
	privateRsa, err := asymmetric.LoadPrivateCipherPEM(filex.FixFilePath(privatePemPath))
	if nil != err {
		panic(err)
	}
	privateRsaCipher = privateRsa
	linkSignature, err := privateRsa.SignBase64(linkInfo.OriginalSignData(), core.Base64Encoding)
	if nil != err {
		panic(err)
	}
	linkInfo.Signature = linkSignature
	unlinkSignature, err := privateRsa.SignBase64(unlinkInfo.OriginalSignData(), core.Base64Encoding)
	if nil != err {
		panic(err)
	}
	unlinkInfo.Signature = unlinkSignature
	publicRsa, err := asymmetric.LoadPublicCipherPEM(filex.FixFilePath(publicPemPath))
	if nil != err {
		panic(err)
	}
	publicRsaCipher = publicRsa
	if post {
		funcLink = client.LinkWithPost
		funcUnlink = client.UnlinkWithPost
		funcUpdate = client.UpdateWithPost
		funcUpdateDetail = client.UpdateDetailWithPost
	}
}

func TestVerifySign(t *testing.T) {
	original := []byte("abcde")
	base64Signature, err := privateRsaCipher.SignBase64(original, core.Base64Encoding)
	if nil != err {
		t.Fatal(err)
		return
	}
	_, err = publicRsaCipher.VerifySignBase64(original, base64Signature, core.Base64Encoding)
	if nil != err {
		t.Fatal(err)
		return
	}
}

func TestInternal(t *testing.T) {
	go testLink(t)
	<-done // 等待
}
func testLink(t *testing.T) {
	t.Log("LinkInfo:", linkInfo)
	funcLink(homeUrl, *linkInfo, weight0, func(res *http.Response, body *[]byte) {
		suc, fail, err := client.ParseLinkBackInfo(res, body)
		if nil != err {
			t.Fatal(err)
			return
		}
		if nil != fail {
			t.Fatal(fail)
			return
		}
		if keyEnable {
			encryptTempKey, err := base64Encoding.DecodeString(suc.TempBase64Key)
			if nil != err {
				t.Fatal(err)
				return
			}
			tempKey, err = privateRsaCipher.Decrypt(encryptTempKey)
			if nil != err {
				t.Fatal(err)
				return
			}
			t.Log(fmt.Sprintf("Id(%s)TempKey(%d):%v", linkInfo.Id, len(tempKey), tempKey))
			tempAesCipher = symmetric.NewAESCipher(tempKey)
		}

		t.Log("Link Suc.")
		time.Sleep(1 * time.Second)
		testUpdate(t)
	})
}

func testUpdate(t *testing.T) {
	t.Log("UpdateInfo:", updateInfo)
	funcUpdate(homeUrl, *updateInfo, tempAesCipher, func(res *http.Response, body *[]byte) {
		fail, err := client.ParseUpdateBackInfo(res, body, tempAesCipher)
		if nil != err {
			t.Fatal(err)
			return
		}
		if nil != fail {
			t.Fatal(fail)
			return
		}
		t.Log("Update Suc.")
		time.Sleep(1 * time.Second)
		if autoUnlink {
			testUpdateDetail(t)
		} else {
			testUpdateRandom(t)
		}
	})
}
func testUpdateDetail(t *testing.T) {
	t.Log("DetailInfo:", detailInfo)
	funcUpdateDetail(homeUrl, *detailInfo, tempAesCipher, func(res *http.Response, body *[]byte) {
		fail, err := client.ParseUpdateBackInfo(res, body, tempAesCipher)
		if nil != err {
			t.Fatal(err)
			return
		}
		if nil != fail {
			t.Fatal(fail)
			return
		}
		t.Log("Update Detail Suc")
		time.Sleep(1 * time.Second)
		testUnlink(t)
	})
}
func testUpdateRandom(t *testing.T) {
	for i := 0; i < 100; i++ {
		updateInfo.Weight = mathx.RandFloat64(0, 100)
		testUpdate(t)
		time.Sleep(time.Second)
	}
	testUnlink(t)
}

func testUnlink(t *testing.T) {
	t.Log("UnlinkInfo:", unlinkInfo)
	funcUnlink(homeUrl, *unlinkInfo, func(res *http.Response, body *[]byte) {
		suc, fail, err := client.ParseUnlinkBackInfo(res, body)
		if nil != err {
			t.Fatal(err)
			return
		}
		if nil != fail {
			t.Fatal(fail)
			return
		}
		t.Log("Unlink Suc.", suc.Id)
		done <- true // 结束
	})
}
