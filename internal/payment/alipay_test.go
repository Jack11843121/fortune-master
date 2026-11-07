package payment

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

// TestRSA2SignVerify 验证支付宝 RSA2 签名/验签闭环与参数排序。
func TestRSA2SignVerify(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	params := map[string]string{
		"app_id":      "202100000000",
		"method":      "alipay.trade.wap.pay",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"version":     "1.0",
		"out_trade_no": "FT202410250001",
	}
	sig, err := rsa2Sign(key, buildSignStr(params))
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}
	ok, err := rsa2Verify(&key.PublicKey, buildSignStr(params), sig)
	if err != nil || !ok {
		t.Fatalf("验签失败 ok=%v err=%v", ok, err)
	}

	// 伪造篡改应验签失败
	bad := map[string]string{"app_id": "202199999999", "method": "alipay.trade.wap.pay"}
	ok, _ = rsa2Verify(&key.PublicKey, buildSignStr(bad), sig)
	if ok {
		t.Fatal("篡改参数后验签仍通过，存在安全漏洞")
	}
}

// TestBuildSignStr 验证支付宝签名串的字典序与编码。
func TestBuildSignStr(t *testing.T) {
	got := buildSignStr(map[string]string{"b": "2", "a": "1"})
	if got != "a=1&b=2" {
		t.Fatalf("签名串顺序错误: %q", got)
	}
	if buildSignStr(map[string]string{"x": "a b*c"}) != "x=a%20b%2Ac" {
		t.Fatalf("URL 编码错误: %q", buildSignStr(map[string]string{"x": "a b*c"}))
	}
}