package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	sjson "github.com/bitly/go-simplejson"
	"golang.org/x/net/publicsuffix"
)



func main() {

	resp ,_:= NewHTTPClient("", "").Get("/get",nil,nil,"form",nil)
	fmt.Println(resp.GetText())
}
