package payment

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"testing"
)

func TestSignVerify(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	msg := []byte("POST\n/v3/pay/transactions/native\n1699999999\nnonce123\n{\"a\":1}\n")
	sig, err := signRSA(key, msg)
	if err != nil {
		t.Fatal(err)
	}
	if !verifyRSA(&key.PublicKey, msg, sig) {
		t.Fatal("验签失败，签名应有效")
	}
	if verifyRSA(&key.PublicKey, []byte("被篡改的内容\n"), sig) {
		t.Fatal("被篡改内容不应验签通过")
	}
}

func TestGCMDecrypt(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, 12)
	_, _ = rand.Read(nonce)

	plaintext := []byte(`{"out_trade_no":"FT123","trade_state":"SUCCESS","amount":{"total":990}}`)
	aad := []byte("transaction")
	ct := gcm.Seal(nil, nonce, plaintext, aad)

	got, err := gcmDecrypt(key, string(nonce), base64.StdEncoding.EncodeToString(ct), string(aad))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(plaintext) {
		t.Fatalf("解密结果不一致: %s", got)
	}
}