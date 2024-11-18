package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/forgoer/openssl"
	"os"
)

// WritePrivateKeyToFile 将PEM编码的私钥写入文件
func WritePrivateKeyToFile(fileName string, priv *rsa.PrivateKey) error {
	privPEM := PrivateKeyToPEM(priv)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(privPEM)
	return err
}

// WritePublicKeyToFile 将PEM编码的公钥写入文件
func WritePublicKeyToFile(fileName string, pub *rsa.PublicKey) error {
	pubPEM, err := PublicKeyToPEM(pub)
	if err != nil {
		return err
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(pubPEM)
	return err
}

// ReadPrivateKeyFromFile 从文件中读取PEM编码的私钥
func ReadPrivateKeyFromFile(fileName string) (*rsa.PrivateKey, error) {
	privPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return PEMToPrivateKey(privPEM)
}

// ReadPublicKeyFromFile 从文件中读取PEM编码的公钥
func ReadPublicKeyFromFile(fileName string) (*rsa.PublicKey, error) {
	pubPEM, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return PEMToPublicKey(pubPEM)
}

// GenerateKeyPair 生成给定比特大小的RSA密钥对
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// EncryptWithPublicKey 使用给定的公钥对数据进行加密
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return encryptedBytes, nil
}

// DecryptWithPrivateKey 使用给定的私钥解密数据
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return decryptedBytes, nil
}

// PrivateKeyToPEM 将私钥转换为PEM格式
func PrivateKeyToPEM(priv *rsa.PrivateKey) []byte {
	privDER := x509.MarshalPKCS1PrivateKey(priv)
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}
	return pem.EncodeToMemory(&privBlock)
}

// PublicKeyToPEM 将公钥转换为PEM格式
func PublicKeyToPEM(pub *rsa.PublicKey) ([]byte, error) {
	pubDER, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	pubBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDER,
	}
	return pem.EncodeToMemory(&pubBlock), nil
}

// PEMToPrivateKey 解析PEM编码的私钥
func PEMToPrivateKey(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

// PEMToPrivatePKCS8Key 解析PEM编码的私钥
func PEMToPrivatePKCS8Key(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pk := priv.(*rsa.PrivateKey)
	return pk, nil
}

// PEMToPublicKey 解析PEM编码的公钥
func PEMToPublicKey(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("not RSA public key")
	}
}

// opensslRsa签名
func OpensslRsaSign(privateStr string, input string, hash crypto.Hash) string {
	pem2, err2 := PEMToPrivatePKCS8Key([]byte(privateStr))
	if err2 != nil {
		return ""
	}
	private := PrivateKeyToPEM(pem2)
	sign, err := openssl.RSASign([]byte(input), private, hash)
	if err != nil {
		return ""
	}
	return string(sign)
}
