package cryptor

// Cryptor 密码器
type Cryptor interface {
	// Encrypt 加密方法
	Encrypt(plainText string) (string, error)

	// Decrypt 解密方法
	Decrypt(cryptedText string) (string, error)
}
