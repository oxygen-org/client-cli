package config

import (
	"io/ioutil"
	"os"
	"github.com/oxygen-org/client/consts"
	"github.com/oxygen-org/client/utils"
	"path"

	sj "github.com/bitly/go-simplejson"
)

//CONFIG JSON配置
var CONFIG = &sj.Json{}

func init() {
	confPath := ""
	confFile := ""
	confFlag := false
	for _, value := range consts.CONFPATH {
		confPath, _ = utils.Expand(value)
		confFile = path.Join(confPath, consts.CONFNAME)

		if utils.FileExists(confFile) {
			confFlag = true
			break
		}
	}
	if confFlag {
		dat, _ := ioutil.ReadFile(confFile)
		config, _ := sj.NewJson(dat)
		CONFIG = config

	} else {
		CONFIG = consts.DEFAULTCONFJSON
		confPath, _ = utils.Expand("~/.nitrogen")
		confFile = path.Join(confPath, consts.CONFNAME)
		confContent, _ := consts.DEFAULTCONFJSON.EncodePretty()
		utils.CreateDirIfNotExists(confPath)
		ioutil.WriteFile(confFile, confContent, os.ModePerm)
	}
	CONFIG.Set("CONFIGPATH", confFile)
	CONFIG.Set("TOKENPATH", path.Join(confPath, consts.TOKENNAME))
}
