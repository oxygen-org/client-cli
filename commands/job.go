package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/oxygen-org/client/clients"
	"github.com/oxygen-org/client/utils"

	sjson "github.com/bitly/go-simplejson"

	"github.com/urfave/cli"
	survey "gopkg.in/AlecAivazis/survey.v1"
	//what import
)

const jobGroup = "job"

// CreateCMD desc
var CreateCMD = cli.Command{
	Name:     "create",
	Aliases:  []string{"cre"},
	Category: jobGroup,
	Usage:    "创建job  参数为job.json文件路径",
	Action:   create,
}

// Create desc
func create(c *cli.Context) error {
	// 参数大于 0 直接认定读取配置文件
	// 否则
	// 交互式填写 | 填写配置文件
	confFile := ""
	if c.NArg() < 1 {
		ansmap := make(map[string]interface{})
		var validationQs = []*survey.Question{

			{
				Name:     "conf_file",
				Prompt:   &survey.Input{Message: "Email:"},
				Validate: survey.Required,
			},
		}
		err := survey.Ask(validationQs, &ansmap)
		if err != nil {
			return fmt.Errorf("\n%s", err.Error())
		}
		confFile = fmt.Sprintf("%v", ansmap["conf_file"])

	}
	confFile, err := utils.Expand(c.Args().Get(0))
	if err != nil {
		return fmt.Errorf("\n%s", err.Error())
	}
	if !utils.FileExists(confFile) {
		return fmt.Errorf("\n[%s] not exist", confFile)
	}
	dat, err := ioutil.ReadFile(confFile)
	if err != nil {
		return fmt.Errorf("\n%s", err.Error())
	}
	config, err := sjson.NewJson(dat)
	if err != nil {
		return fmt.Errorf("\nInvalid format: please check your job json")
	}
	gitURL := config.Get("git_url").MustString()
	pattern := `git@(bitbucket.org|github.com|192.168.1.158):.+\.git`
	if regexp.MustCompile(pattern).MatchString(gitURL) {
		log.Fatalln("your git url is invalid.")
	}
	env, ok := config.CheckGet("env")
	if ok {
		envArr, err := env.Map()
		if err != nil {
			log.Fatalln("env is invalid")
		}
		sysEnv := []string{"PYTHONPATH", "DATAPATH", "SAVEDPATH", "GPU"}
		for k := range envArr {
			if !utils.IsUpper(k) {
				log.Fatalln("keys of enviroment must be upper! Please retry after modification.")
			}
			for _, value := range sysEnv {
				if k == value {
					log.Fatalf("Key [%s] is same to system environment variable! Please retry after modification.", k)
				}
			}

		}
	}
	jobName := config.Get("name").MustString()
	if len(jobName) > 32 || len(jobName) < 1 {
		log.Fatalln("job name length limit to 1~32 characters long, more information please fill in description column")
	}
	_, ok = config.CheckGet("docker_image")
	if !ok {
		config.Set("docker_image", "hub.xtalpi.cc/atompai/nitrogen:v6")
	}
	_, ok = config.CheckGet("use_gpu")
	if !ok {
		config.Set("use_gpu", 1)
	}
	data, _ := config.Encode()
	resp, err := clients.NewHTTPClient("", "").Post("/api/v1.0/jobs", nil, data, "json", nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 2 {
		fmt.Printf("create job success, job id is %v", content.Get("job_id").MustInt())
	} else {
		log.Fatalf(content.Get("message").MustString())
	}

	return nil
}

// RunCMD desc
var RunCMD = cli.Command{
	Name:     "run",
	Aliases:  []string{"r"},
	Category: jobGroup,
	Usage:    "运行job",
	Action:   run,
}

// Create desc
func run(c *cli.Context) error {
	// 填写job-id直接后序执行
	// 否则 交互式填写
	fmt.Println("run cmd")
	return nil
}

// LogCMD desc
var LogCMD = cli.Command{
	Name:     "log",
	Aliases:  []string{"lg"},
	Category: jobGroup,
	Usage:    "查看或下载job运行日志",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:   "download, d",
			Usage:  "download log",
			Hidden: true,
		},
		cli.StringFlag{
			Name:  "dir",
			Usage: "the path you want downloaded to",
			Value: ".",
		},
		cli.BoolFlag{
			Name:  "tail,t",
			Usage: "get log of a running job",
		},
		cli.BoolFlag{
			Name:  "realtime,rt",
			Usage: "get real-time log of a running job",
		},
	},
	Action: logIt, //main function of tool
}

// Create desc
func logIt(c *cli.Context) error {
	fmt.Println("log cmd")
	return nil
}

// QueueCMD desc
var QueueCMD = cli.Command{
	Name:     "queue",
	Aliases:  []string{"q"},
	Category: jobGroup,
	Usage:    "查看和取消任务排队",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:   "view, v",
			Usage:  "view job ranking",
			Hidden: true,
		},
		//flag
	},
	Action: queueAbout, //main function of tool
}

// Create desc
func queueAbout(c *cli.Context) error {
	fmt.Println("queue cmd")
	return nil
}

// KillCMD desc
var KillCMD = cli.Command{
	Name:     "kill",
	Aliases:  []string{"k"},
	Category: jobGroup,
	Usage:    "结束正在运行的job",
	Action:   kill,
}

// Create desc
func kill(c *cli.Context) error {
	fmt.Println("kill cmd")
	return nil
}
