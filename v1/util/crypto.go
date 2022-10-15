package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func AesEncrypt(orig []byte, key []byte) []byte {
	// 转成字节数组
	origData := orig

	// 分组秘钥
	block, _ := aes.NewCipher(key)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// fmt.Println(blockSize)
	// 补全码
	origData = PKCS7Padding(origData, blockSize)

	// fmt.Println("orig len", len(origData))
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	// fmt.Println("En Text", cryted)
	return cryted

}

func AesDecrypt(cryted []byte, key []byte) ([]byte, error) {
	// fmt.Println("Un Text", cryted)
	// 转成字节数组
	crytedByte := cryted

	// 分组秘钥
	block, _ := aes.NewCipher(key)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig, err := PKCS7UnPadding(orig)
	return orig, err
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	res := append(ciphertext, padtext...)
	return res
}

//去码
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	// print(length)
	unpadding := int(origData[length-1])

	if length <= unpadding {
		return nil, fmt.Errorf("error Unpadding")
	}

	return origData[:(length - unpadding)], nil
}
