package cryptor

import "testing"

var aesCryp Cryptor

func init() {
	key := "123456781234567812345678"
	aesCryp = NewAesEncrpytor(key)
}

func TestAES(t *testing.T) {
	plainText := "hello world"
	cryptedText, err := aesCryp.Encrypt(plainText)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("cryptedText is %s\n", cryptedText)

	result, err := aesCryp.Decrypt(cryptedText)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("result is %s\n", result)

	if result != plainText {
		t.Fatal("AES encrypt & decrypt fail")
	}
}
