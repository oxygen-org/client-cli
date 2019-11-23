package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"

	sjson "github.com/bitly/go-simplejson"
	"github.com/bndr/gotabulate"
)

func pprint() {
	// Some Strings
	str1 := []string{"TV", "1000$", "Sold"}
	str2 := []string{"PC", "50%", "on Hold"}

	// Create Object
	tabulate := gotabulate.Create([][]string{str1, str2})

	// Set Headers
	tabulate.SetHeaders([]string{"Type", "Cost", "Status"})

	// Render
	fmt.Println(tabulate.Render("grid"))

}

func loadConfig(configPath string) (*sjson.Json, error) {
	dat, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config, err := sjson.NewJson(dat)
	return config, err

}

func getEmail() {

}

func hashFileGen(filePath string) (string, error) {
	const filechunk = 65536
	var returnSHA1String string

	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA1String, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return "Cannot access file", err
	}
	filesize := info.Size()

	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

	hasher := sha1.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)
		file.Read(buf)
		io.WriteString(hasher, string(buf))
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func hashFileCopy(filePath string) (string, error) {
	var returnSHA1String string

	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA1String, err
	}

	defer file.Close()

	hasher := sha1.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return returnSHA1String, err
	}

	hashInBytes := hasher.Sum(nil)[:20]
	returnSHA1String = hex.EncodeToString(hashInBytes)

	return returnSHA1String, nil
}

// Md5Encrypt md5 encrypt
func Md5Encrypt(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
