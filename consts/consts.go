package consts

import (
	sjson "github.com/bitly/go-simplejson"
)

// APPNAME 程序名称
const APPNAME = "client"

// VERSION 版本号
const VERSION = "0.0.1"


// CONFNAME 配置文件名称
const CONFNAME = ".nitrogen.json"

// TOKENNAME TOKEN文件名称
const TOKENNAME = ".token.json"

// CONFPATH 配置文件搜索路径
// 没办法用const
var CONFPATH = [3]string{"./", "$HOME/.nitrogen", "/etc/nitrogen"}



// S3LOCATIONS S3区域
var S3LOCATIONS = []string{"cn-north-1", "cn-northwest-1"}

// DEFAULTCONF 缺省配置
var DEFAULTCONF = `{"server_host": "http://192.168.1.182:7000", "bucket_url": "https://s3.cn-north-1.amazonaws.com.cn/nitrogen/", "s3_bucket": "s3://nitrogen/", "bucket": "nitrogen", "RD_HOST": "52.82.105.38", "RD_PORT": "6379", "RD_DB": 5, "RD_PWD": "xtalpi_redis", "hdfs_url": "http://192.168.1.120:50070", "hdfs_user": "jp"}`
// DEFAULTCONFJSON 缺省配置
var DEFAULTCONFJSON,_ = sjson.NewJson([]byte(DEFAULTCONF))