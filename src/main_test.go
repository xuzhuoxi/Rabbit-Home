// Create on 2023/6/5
// @author xuzhuoxi
package main

import (
	"encoding/base64"
	"fmt"
	"github.com/xuzhuoxi/Rabbit-Home/core"
	"github.com/xuzhuoxi/Rabbit-Home/core/client"
	"github.com/xuzhuoxi/Rabbit-Home/core/utils"
	"github.com/xuzhuoxi/infra-go/cryptox/asymmetric"
	"github.com/xuzhuoxi/infra-go/cryptox/symmetric"
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/mathx"
	"net/http"
	"testing"
	"time"
)

var (
	homeUrl        = "http://127.0.0.1:9000"
	weight0        = 1.00
	base64Encoding = base64.RawURLEncoding
)

var (
	funcLink         = client.LinkWithGet
	funcUnlink       = client.UnlinkWithGet
	funcUpdate       = client.UpdateWithGet
	funcUpdateDetail = client.UpdateDetailWithGet
	funcQueryRoute   = client.QueryRouteWithGet
)

type extLinkInfo struct {
	core.LinkInfo
	InternalSK        []byte
	InternalAesCipher symmetric.IAESCipher
	OpenSK            []byte
}

var (
	linkInfos = []*extLinkInfo{
		{LinkInfo: core.LinkInfo{Id: "1000", PlatformId: "P-01", TypeName: "LogicalServer", OpenNetwork: "tcp", OpenAddr: "127.0.0.1:41000", OpenKeyOn: false}},
		{LinkInfo: core.LinkInfo{Id: "1001", PlatformId: "P-01", TypeName: "LogicalServer", OpenNetwork: "tcp", OpenAddr: "127.0.0.1:41001", OpenKeyOn: true}},
		{LinkInfo: core.LinkInfo{Id: "1002", PlatformId: "P-01", TypeName: "LogicalServer", OpenNetwork: "tcp", OpenAddr: "127.0.0.1:41002", OpenKeyOn: false}},
		{LinkInfo: core.LinkInfo{Id: "2003", PlatformId: "P-02", TypeName: "LogicalServer", OpenNetwork: "tcp", OpenAddr: "127.0.0.1:41003", OpenKeyOn: true}},
		{LinkInfo: core.LinkInfo{Id: "2004", PlatformId: "P-02", TypeName: "LogicalServer", OpenNetwork: "tcp", OpenAddr: "127.0.0.1:41004", OpenKeyOn: false}},
		{LinkInfo: core.LinkInfo{Id: "2005", PlatformId: "P-02", TypeName: "LogicalServer", OpenNetwork: "tcp", OpenAddr: "127.0.0.1:41005", OpenKeyOn: true}},
	}
	detailInfos = []*core.UpdateDetailInfo{
		{Id: "1000", Links: 1, MaxLinks: 100, TotalReqCount: int64(time.Second), TotalRespTime: int64(time.Minute), EnableKeys: "links,max-links,total-reg,total-resp"},
		{Id: "1001", Links: 11, MaxLinks: 150, TotalReqCount: int64(time.Second), TotalRespTime: int64(time.Minute), EnableKeys: "links,max-links,total-reg,total-resp"},
		{Id: "1002", Links: 1, MaxLinks: 510, TotalReqCount: int64(time.Second), TotalRespTime: int64(time.Minute), EnableKeys: "links,max-links,total-reg,total-resp"},
		{Id: "2003", Links: 61, MaxLinks: 160, TotalReqCount: int64(time.Second), TotalRespTime: int64(time.Minute), EnableKeys: "links,max-links,total-reg,total-resp"},
		{Id: "2004", Links: 15, MaxLinks: 108, TotalReqCount: int64(time.Second), TotalRespTime: int64(time.Minute), EnableKeys: "links,max-links,total-reg,total-resp"},
		{Id: "2005", Links: 91, MaxLinks: 106, TotalReqCount: int64(time.Second), TotalRespTime: int64(time.Minute), EnableKeys: "links,max-links,total-reg,total-resp"},
	}
	updateInfos     []*core.UpdateInfo
	unlinkInfos     []*core.UnlinkInfo
	queryRouteInfos = []*core.QueryRouteInfo{
		{PlatformId: "P-01", TypeName: "LogicalServer"},
		{PlatformId: "P-02", TypeName: "LogicalServer", TempAesKey: utils.DeriveRandomKey32("1234567890123456", "a")},
	}
)
var (
	done             = make(chan bool)
	privatePemPath   = "keys/private/pkcs8_private.pem"
	publicPemPath    = "keys/public/x509_public.pem"
	privateRsaCipher asymmetric.IRSAPrivateCipher
	publicRsaCipher  asymmetric.IRSAPublicCipher
)

var (
	post              = false
	internalKeyEnable = true
	autoUnlink        = false
)

func init() {
	if post {
		funcLink = client.LinkWithPost
		funcUnlink = client.UnlinkWithPost
		funcUpdate = client.UpdateWithPost
		funcUpdateDetail = client.UpdateDetailWithPost
		funcQueryRoute = client.QueryRouteWithPost
	}
	for index := 0; index < len(linkInfos); index++ {
		updateInfos = append(updateInfos, &core.UpdateInfo{Id: linkInfos[index].Id, Weight: mathx.RandFloat64(0.1, 2.0)})
		unlinkInfos = append(unlinkInfos, &core.UnlinkInfo{Id: linkInfos[index].Id})
	}

	if internalKeyEnable {
		privateRsa, err := asymmetric.LoadPrivateCipherPEM(filex.FixFilePath(privatePemPath))
		if nil != err {
			panic(err)
		}
		privateRsaCipher = privateRsa
		publicRsa, err := asymmetric.LoadPublicCipherPEM(filex.FixFilePath(publicPemPath))
		if nil != err {
			panic(err)
		}
		publicRsaCipher = publicRsa

		for index := 0; index < len(linkInfos); index++ {
			linkSignature, err := privateRsa.SignBase64(linkInfos[index].OriginalSignData(), core.Base64Encoding)
			if nil != err {
				panic(err)
			}
			linkInfos[index].Signature = linkSignature
			unlinkSignature, err := privateRsa.SignBase64(unlinkInfos[index].OriginalSignData(), core.Base64Encoding)
			if nil != err {
				panic(err)
			}
			unlinkInfos[index].Signature = unlinkSignature
		}
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

func TestQuery(t *testing.T) {
	for index := 0; index < len(queryRouteInfos); index++ {
		saveIndex := index
		t.Log("QueryRouteInfo:", queryRouteInfos[saveIndex])
		funcQueryRoute(homeUrl, *queryRouteInfos[saveIndex], publicRsaCipher, func(res *http.Response, body *[]byte) {
			suc, fail, err := client.ParseQueryRouteBackInfo(res, body)
			if nil != err {
				t.Fatal(err)
				return
			}
			if nil != fail {
				t.Fatal(fail)
				return
			}
			err = suc.ComputeOpenSK(queryRouteInfos[saveIndex].TempAesKey)
			if nil != err {
				t.Fatal(err)
				return
			}
			t.Log("QueryRouteInfo Suc.", saveIndex, suc)
		})
	}
}

func TestInternal(t *testing.T) {
	for index := 0; index < len(linkInfos); index++ {
		go testLink(t, index)
	}
	<-done // 等待
}
func testLink(t *testing.T, index int) {
	defaultSleep()
	t.Log("LinkInfo:", index, linkInfos[index])
	funcLink(homeUrl, linkInfos[index].LinkInfo, weight0, func(res *http.Response, body *[]byte) {
		suc, fail, err := client.ParseLinkBackInfo(res, body, privateRsaCipher)
		if nil != err {
			t.Fatal(err)
			return
		}
		if nil != fail {
			t.Fatal(fail)
			return
		}
		if internalKeyEnable {
			linkInfos[index].InternalSK = suc.InternalSK
			linkInfos[index].InternalAesCipher = symmetric.NewAESCipher(suc.InternalSK)
		}
		if len(suc.OpenSK) > 0 {
			linkInfos[index].OpenSK = suc.OpenSK
		}
		t.Log(fmt.Sprintf("%d Id(%s) InternalSK(%d)=%v, OpenSK(%d)=%v", index, linkInfos[index].Id, len(suc.InternalSK), suc.InternalSK, len(suc.OpenSK), suc.OpenSK))

		t.Log("Link Suc.")
		defaultSleep()
		testUpdateAndDetail(t, index)
	})
}

func testUpdateAndDetail(t *testing.T, index int) {
	if !autoUnlink {
		for {
			go testUpdateRandom(t, index)
			go testUpdateDetailRandom(t, index)
			defaultSleep()
		}
	} else {
		go testUpdate(t, index)
		go testUpdateDetail(t, index)
		defaultSleep()
		testUnlink(t, index)
	}
}
func testUpdateRandom(t *testing.T, index int) {
	updateInfos[index].Weight = mathx.RandFloat64(0, 100)
	testUpdate(t, index)
}
func testUpdateDetailRandom(t *testing.T, index int) {
	detailInfos[index].Links = mathx.RandUint64(0, 200)
	testUpdateDetail(t, index)
}

func testUpdate(t *testing.T, index int) {
	t.Log("UpdateInfo:", index, updateInfos[index])
	funcUpdate(homeUrl, *updateInfos[index], linkInfos[index].InternalAesCipher, func(res *http.Response, body *[]byte) {
		fail, err := client.ParseUpdateBackInfo(res, body)
		if nil != err {
			t.Fatal(err)
			return
		}
		if nil != fail {
			t.Fatal(fail)
			return
		}
		t.Log("Update Suc.")
	})
}
func testUpdateDetail(t *testing.T, index int) {
	t.Log("DetailInfo:", index, detailInfos[index])
	funcUpdateDetail(homeUrl, *detailInfos[index], linkInfos[index].InternalAesCipher, func(res *http.Response, body *[]byte) {
		fail, err := client.ParseUpdateBackInfo(res, body)
		if nil != err {
			t.Fatal(err)
			return
		}
		if nil != fail {
			t.Fatal(fail)
			return
		}
		t.Log("Update Detail Suc")
	})
}

func testUnlink(t *testing.T, index int) {
	t.Log("UnlinkInfo:", index, unlinkInfos[index])
	funcUnlink(homeUrl, *unlinkInfos[index], func(res *http.Response, body *[]byte) {
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

func defaultSleep() {
	randomSleep(time.Second, time.Second*5)
}
func randomSleep(min, max time.Duration) {
	s := time.Duration(mathx.RandInt64(int64(min), int64(max)))
	time.Sleep(s)
}
