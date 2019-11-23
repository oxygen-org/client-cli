package utils



import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"time"
)

var defaultCipher = "Private key"

var (
	h = md5.New()
)

func cipherEncode(sourceText, cipher string) string {
	if cipher == "" {
		cipher = defaultCipher
	}
	h.Write([]byte(cipher))
	cipherHash := fmt.Sprintf("%x", h.Sum(nil))
	h.Reset()
	inputData := []byte(sourceText)
	loopCount := len(inputData)
	outData := make([]byte, loopCount)
	for i := 0; i < loopCount; i++ {
		outData[i] = inputData[i] ^ cipherHash[i%32]
	}
	return fmt.Sprintf("%s", outData)
}

// AzDGEncrypt DESC
func AzDGEncrypt(sourceText, cipher string) string {
	h.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
	noise := fmt.Sprintf("%x", h.Sum(nil))
	h.Reset()
	inputData := []byte(sourceText)
	loopCount := len(inputData)
	outData := make([]byte, loopCount*2)
	for i, j := 0, 0; i < loopCount; i, j = i+1, j+1 {
		outData[j] = noise[i%32]
		j++
		outData[j] = inputData[i] ^ noise[i%32]
	}
	return base64.StdEncoding.EncodeToString([]byte(cipherEncode(fmt.Sprintf("%s", outData), cipher)))
}

// AzDGDecrypt DESC
func AzDGDecrypt(sourceText, cipher string) string {
	buf, err := base64.StdEncoding.DecodeString(sourceText)
	if err != nil {
		fmt.Printf("Decode(%q) failed: %v", sourceText, err)
		return ""
	}
	inputData := []byte(cipherEncode(fmt.Sprintf("%s", buf), cipher))
	loopCount := len(inputData)
	outData := make([]byte, loopCount)
	for i, j := 0, 0; i < loopCount; i, j = i+2, j+1 {
		outData[j] = inputData[i] ^ inputData[i+1]
	}
	return fmt.Sprintf("%s", outData)
}