package clients

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"github.com/oxygen-org/client/config"
	"github.com/oxygen-org/client/consts"
	"github.com/oxygen-org/client/utils"
	"path"
	"sort"
	"strings"
	"time"

	sjson "github.com/bitly/go-simplejson"
	"golang.org/x/net/publicsuffix"
)

var totalTimeLimit = 10 * time.Second
var dialTimeLimti = 5 * time.Second
var shakeTimeLimit = 5 * time.Second

// HTTPClient custom http-client
type HTTPClient struct {
	Client   http.Client
	BaseURL  string
	User     string
	Password string
	Header   map[string]string
}

// HTTPResponse custom respond
type HTTPResponse struct {
	http.Response
}

// GetBody 获取Respond的Body原始结果
func (resp *HTTPResponse) GetBody() ([]byte, error) {
	return ioutil.ReadAll(resp.Body)
}

// CheckGetJSON 获取Respond Body的JSON格式
func (resp *HTTPResponse) CheckGetJSON() (*sjson.Json, error) {
	content, err := resp.GetBody()
	if err != nil {
		return nil, err
	}
	return sjson.NewJson(content)
}

// GetJSON 获取Respond Body的JSON格式
func (resp *HTTPResponse) GetJSON() *sjson.Json {
	content, _ := resp.CheckGetJSON()
	return content
}

// GetText 获取Respond的文本字符串格式
func (resp *HTTPResponse) GetText() string {
	content, _ := resp.CheckGetText()
	return string(content)

}

// CheckGetText 获取Respond的文本字符串格式
func (resp *HTTPResponse) CheckGetText() (string, error) {
	content, err := resp.GetBody()
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: dialTimeLimti,
	}).Dial,
	TLSHandshakeTimeout: shakeTimeLimit,
}

// NewHTTPClient 创建custom http-client
func NewHTTPClient(user string, password string) *HTTPClient {
	c := new(HTTPClient)
	c.Client = http.Client{
		Timeout:   totalTimeLimit,
		Transport: netTransport}
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	c.Client.Jar = jar
	c.Client.CheckRedirect = func(request *http.Request, via []*http.Request) error {
		request.SetBasicAuth(c.User, c.Password)
		return nil
	}
	c.BaseURL = "https://httpbin.org/" //consts.DEFAULTCONFJSON.Get("server_host").MustString()
	c.User = user
	c.Password = password
	return c
}

// UploadFile http-post-file-request
func (c *HTTPClient) UploadFile(url string, fieldName string, filePath string) (*HTTPResponse, error) {

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)

	fileWriter, err := multiPartWriter.CreateFormFile(fieldName, "name.txt")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Fatalln(err)
	}

	fieldWriter, err := multiPartWriter.CreateFormField("normal_field")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = fieldWriter.Write([]byte("Value"))
	if err != nil {
		log.Fatalln(err)
	}

	multiPartWriter.Close()

	request, err := http.NewRequest("POST", "https://httpbin.org/post", &requestBody)
	if err != nil {
		log.Fatalln(err)
	}
	c.setHeader(request, map[string]string{"Content-Type": multiPartWriter.FormDataContentType()})

	resp, err := c.Client.Do(request)
	if err != nil && resp != nil {
		log.Fatalln(err)
	}
	return nil, nil
}

func (c *HTTPClient) setHeader(r *http.Request, extra map[string]string) {
	for key, value := range extra {
		r.Header.Add(key, value)
	}
	tokenFile := config.CONFIG.Get("TOKENPATH").MustString()
	r.Header.Add("User-Agent", consts.APPNAME+"-"+consts.VERSION)
	if c.User != "" || c.Password != "" {
		r.SetBasicAuth(c.User, c.Password)
	} else if utils.FileExists(tokenFile) {
		dat, _ := ioutil.ReadFile(tokenFile)
		tokeInfo, err := sjson.NewJson(dat)
		if err != nil {
			log.Fatalln("Token file error:", err)
		}
		token := tokeInfo.Get("token").MustString()
		r.SetBasicAuth(token, "")
	}

}

func (c *HTTPClient) doRequest(method, urlPath string, query map[string]string, data []byte, bodyType string, header map[string]string) (*HTTPResponse, error) {
	u, err := url.Parse(c.BaseURL)
	u.Path = path.Join(u.Path, urlPath)
	if query != nil {
		params := url.Values{}
		for key, value := range query {
			params.Set(key, value)
		}
		u.RawQuery = params.Encode()
	}
	var body io.Reader
	var headers = map[string]string{}
	{
		switch bodyType {
		case "json":
			headers["Content-Type"] = "application/json"
			body = bytes.NewBuffer(data)

		case "form":
			headers["Content-Type"] = "application/x-www-form-urlencoded"
			if data == nil {
				break
			}
			jsonData, err := sjson.NewJson(data)
			if err != nil {
				log.Fatalln(err)
			}
			mapData, err := jsonData.Map()
			if err != nil {
				log.Fatalln(err)
			}
			params := url.Values{}
			for k, v := range mapData {
				params.Set(k, v.(string))
			}
			body = strings.NewReader(params.Encode())
		default:
			if data == nil {
				break
			}
			jsonData, err := sjson.NewJson(data)
			if err != nil {
				log.Fatalln(err)
			}
			mapData, err := jsonData.Map()
			if err != nil {
				log.Fatalln(err)
			}
			params := url.Values{}
			for k, v := range mapData {
				params.Set(k, v.(string))
			}
			var buf strings.Builder
			keys := make([]string, 0, len(params))
			for k := range params {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				vs := params[k]

				for _, v := range vs {
					if buf.Len() > 0 {
						buf.WriteByte('&')
					}
					buf.WriteString(k)
					buf.WriteByte('=')
					buf.WriteString(v)
				}
			}
			body = strings.NewReader(buf.String())

		}
	}
	request, err := http.NewRequest(strings.ToUpper(method), u.String(), body)
	c.setHeader(request, headers)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	tresp := HTTPResponse{*resp}
	return &tresp, nil
}

// Get http-get-method-request
func (c *HTTPClient) Get(urlPath string, query map[string]string, data []byte, bodyType string, header map[string]string) (*HTTPResponse, error) {
	return c.doRequest("GET", urlPath, query, data, bodyType, header)
}

// Post http-post-method-request
func (c *HTTPClient) Post(urlPath string, query map[string]string, data []byte, bodyType string, header map[string]string) (*HTTPResponse, error) {
	return nil, nil
}

// Put http-put-method-request
func (c *HTTPClient) Put(urlPath string, query map[string]string, data []byte, bodyType string, header map[string]string) (*HTTPResponse, error) {
	return nil, nil
}

// Option http-option-method-request
func (c *HTTPClient) Option(urlPath string, query map[string]string, data []byte, bodyType string, header map[string]string) (*HTTPResponse, error) {
	return nil, nil
}

// Patch http-patch-method-request
func (c *HTTPClient) Patch(urlPath string, query map[string]string, data []byte, bodyType string, header map[string]string) (*HTTPResponse, error) {
	return nil, nil

}

// Delete http-delete-method-request
func (c *HTTPClient) Delete(urlPath string, query map[string]string, data []byte, bodyType string, header map[string]string) (*HTTPResponse, error) {
	return nil, nil
}

// Head http-head-method-request
func (c *HTTPClient) Head(urlPath string, query map[string]string, data []byte, bodyType string, header map[string]string) (*HTTPResponse, error) {
	return nil, nil

}