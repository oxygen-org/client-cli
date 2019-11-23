package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"github.com/oxygen-org/client/clients"
	"github.com/oxygen-org/client/config"
	"github.com/oxygen-org/client/consts"
	"github.com/oxygen-org/client/utils"
	"strconv"

	"strings"

	sjson "github.com/bitly/go-simplejson"
	"github.com/urfave/cli"
	survey "gopkg.in/AlecAivazis/survey.v1"
	//what import
)

const authGroup = "auth"

// VersionCMD desc
var VersionCMD = cli.Command{
	Name:     "version",
	Aliases:  []string{"v"},
	Category: authGroup,
	Usage:    "nitrogen版本",
	Action:   version, //main function of tool
}

// Tool desc
func version(c *cli.Context) error {
	fmt.Println(consts.APPNAME, consts.VERSION)
	return nil
}

// RegisterCMD desc
var RegisterCMD = cli.Command{
	Name:     "register",
	Aliases:  []string{"reg"},
	Category: authGroup,
	Usage:    "新用户注册",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "email, e",
			Usage: "input `EMAIL`",
		},
		cli.StringFlag{
			Name:  "username, u",
			Usage: "input `USERNAME`",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "input `PASSWORD`",
		},

		//flag
	},
	Action: register, //main function of tool
}

// Tool desc
func register(c *cli.Context) error {
	email := c.String("email")
	username := c.String("username")
	password := c.String("password")
	if email == "" && username == "" && password == "" {
		ansmap := make(map[string]interface{})
		hostname, err := os.Hostname()
		if err != nil {
			hostname = ""
		}
		var validationQs = []*survey.Question{
			{
				Name:     "username",
				Prompt:   &survey.Input{Message: "UserName:", Default: hostname},
				Validate: survey.Required,
			},
			{
				Name:   "email",
				Prompt: &survey.Input{Message: "Email:"},
				Validate: func(val interface{}) error {
					str := val.(string)
					err := utils.ValidateFormat(str)
					if err != nil {
						return fmt.Errorf(err.Error())
					}
					vailidSuffix := strings.HasSuffix(str, "xtalpi.com")
					if !vailidSuffix {
						return fmt.Errorf("must be xtalpi.com")
					}
					err = utils.ValidateHost(str)
					if smtpErr, ok := err.(utils.SMTPError); ok && err != nil {
						return fmt.Errorf("Code: %s, Msg: %s", smtpErr.Code(), smtpErr)
					}
					return nil
				},
			},
			{
				Name:   "password",
				Prompt: &survey.Password{Message: "Password:"},
				Validate: func(val interface{}) error {
					str := val.(string)
					if len(str) < 6 {
						return fmt.Errorf("must > 6")
					}
					return nil
				},
			},
		}

		err = survey.Ask(validationQs, &ansmap)
		if err != nil {
			return fmt.Errorf("\n%s", err.Error())
		}
		email = fmt.Sprintf("%v", ansmap["email"])
		username = fmt.Sprintf("%v", ansmap["username"])
		password = fmt.Sprintf("%v", ansmap["password"])
	}
	password = utils.Md5Encrypt(password)
	data := []byte(fmt.Sprintf(`{"email": %s, "password": %s, "username": %s}`, email, password, username))

	resp, err := clients.NewHTTPClient("", "").Post("/auth/registration", nil, data, "json", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	fmt.Printf("register message:%v\n", content.Get("message").MustString())
	return nil
}

// LoginCMD desc
var LoginCMD = cli.Command{
	Name:     "login",
	Aliases:  []string{"logi"},
	Category: authGroup,
	Usage:    "使用前请先登录（登录后8小时内有效）",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "email, e",
			Usage: "input `EMAIL`",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "input `PASSWORD`",
		},

		//flag
	},
	Action: login, //main function of tool
}

// Tool desc
func login(c *cli.Context) error {
	email := c.String("email")
	password := c.String("password")
	if email == "" && password == "" {
		ansmap := make(map[string]interface{})

		var validationQs = []*survey.Question{

			{
				Name:   "email",
				Prompt: &survey.Input{Message: "Email:"},
				Validate: func(val interface{}) error {
					str := val.(string)
					err := utils.ValidateFormat(str)
					if err != nil {
						return fmt.Errorf(err.Error())
					}
					vailidSuffix := strings.HasSuffix(str, "xtalpi.com")
					if !vailidSuffix {
						return fmt.Errorf("must be xtalpi.com")
					}
					err = utils.ValidateHost(str)
					if smtpErr, ok := err.(utils.SMTPError); ok && err != nil {
						return fmt.Errorf("Code: %s, Msg: %s", smtpErr.Code(), smtpErr)
					}
					return nil
				},
			},
			{
				Name:   "password",
				Prompt: &survey.Password{Message: "Password:"},
				Validate: func(val interface{}) error {
					str := val.(string)
					if len(str) < 6 {
						return fmt.Errorf("must > 6")
					}
					return nil
				},
			},
		}
		err := survey.Ask(validationQs, &ansmap)
		if err != nil {
			return fmt.Errorf("\n%s", err.Error())
		}
		email = fmt.Sprintf("%v", ansmap["email"])
		password = fmt.Sprintf("%v", ansmap["password"])
	}

	doLogin(email, password)
	return nil
}

func doLogin(email, password string) error {
	password = utils.Md5Encrypt(password)
	params := map[string]string{
		"email": email, "password": password,
	}
	resp, err := clients.NewHTTPClient(email, password).Get("/api/v1.0/token", params, nil, "json", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("code").MustInt() == 1 {
		token := content.Get("data").Get("token").MustString()
		text := fmt.Sprintf(`{"email": "%s", "token": "%s"}`, email, token)
		ioutil.WriteFile(config.CONFIG.Get("TOKENPATH").MustString(), []byte(text), os.ModePerm)
		fmt.Println("login success")
	} else {
		log.Fatalln("login faild:", content.Get("error").MustString())
	}
	return nil
}

// PasswordCMD desc
var PasswordCMD = cli.Command{
	Name:     "password",
	Aliases:  []string{"pw", "passwd"},
	Category: authGroup,
	Usage:    "修改密码",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "email, e",
			Usage: "input `EMAIL`",
		},
		cli.StringFlag{
			Name:  "old_password, op",
			Usage: "input `OLD_PASSWORD`",
		},
		cli.StringFlag{
			Name:  "new_password, np",
			Usage: "input `NEW_PASSWORD`",
		},
		//flag
	},
	Action: password, //main function of tool
}

// Tool desc
func password(c *cli.Context) error {
	email := c.String("email")
	oldPassword := c.String("old_password")
	newPassword := c.String("new_password")
	if email == "" && oldPassword == "" && newPassword == "" {
		ansmap := make(map[string]interface{})
		var validationQs = []*survey.Question{
			{
				Name:   "email",
				Prompt: &survey.Input{Message: "Email:"},
				Validate: func(val interface{}) error {
					str := val.(string)
					err := utils.ValidateFormat(str)
					if err != nil {
						return fmt.Errorf(err.Error())
					}
					vailidSuffix := strings.HasSuffix(str, "xtalpi.com")
					if !vailidSuffix {
						return fmt.Errorf("must be xtalpi.com")
					}
					err = utils.ValidateHost(str)
					if smtpErr, ok := err.(utils.SMTPError); ok && err != nil {
						return fmt.Errorf("Code: %s, Msg: %s", smtpErr.Code(), smtpErr)
					}
					return nil
				},
			},
			{
				Name:   "old_password",
				Prompt: &survey.Password{Message: "OldPassword:"},
				Validate: func(val interface{}) error {
					str := val.(string)
					if len(str) < 6 {
						return fmt.Errorf("must > 6")
					}
					return nil
				},
			},
			{
				Name:   "new_password",
				Prompt: &survey.Password{Message: "NewPassword:"},
				Validate: func(val interface{}) error {
					str := val.(string)
					if len(str) < 6 {
						return fmt.Errorf("must > 6")
					}
					return nil
				},
			},
		}
		err := survey.Ask(validationQs, &ansmap)
		if err != nil {
			return fmt.Errorf("\n%s", err.Error())
		}
		email = fmt.Sprintf("%v", ansmap["email"])
		oldPassword = fmt.Sprintf("%v", ansmap["old_password"])
		newPassword = fmt.Sprintf("%v", ansmap["new_password"])
	}
	data := []byte(fmt.Sprintf(`{"email":%s,"old_password":%s ,"new_password": %s
	}`, email, utils.Md5Encrypt(oldPassword), utils.Md5Encrypt(newPassword)))
	resp, err := clients.NewHTTPClient("", "").Put("/auth/password", nil, data, "json", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("code").MustInt() == 1 {
		fmt.Println("reset success")
	} else {
		log.Fatalln("reset faild:", content.Get("message").MustString())
	}
	return nil
}

// KeyCMD desc
var KeyCMD = cli.Command{
	Name:     "key",
	Aliases:  []string{"k"},
	Category: authGroup,
	Usage:    "nitrogen公钥",

	Action: key, //main function of tool
}

// Tool desc
func key(c *cli.Context) error {
	key := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCrH9eJsOVaPcnIG8egJr3sy+Q29wu+fE7w8ybMLr55ePpQA9ko0KmQKWEQp79uwuyk8arCajhnH7/xwLbtFsrpoZ/Du24W04QtQN78GWht4d1PZb8jFV5tAkWJXOTHxzz/K2sQRorWMk1mzcXappLgHS+paBQ9GZs/8COfhjZLMH0uuXPXU+4ewWa2OoBqcLeYw0MDoVJ3FqBJ6h0EWn81lq0maExFAozTPkxEePE7WXV6dEiCpPax0kposle9g7by+RfNCe7LLNlWiFa2J80pEbKT+S7KDBtr4wbx3XAeGUbVPzVEaG9IOXywAXCJIKPmAzOb2/KltZEEE6ylejbR root@912d046f69c3`
	fmt.Println(key)
	return nil
}

// UpdateCMD desc
var UpdateCMD = cli.Command{
	Name:     "update",
	Aliases:  []string{"up"},
	Category: authGroup,
	Usage:    "客户端自动升级",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "download, d",
			Usage: "download install-package only",
		},
		cli.StringFlag{
			Name:  "email,e",
			Usage: "input `EMAIL`",
		},
		cli.StringFlag{
			Name:  "password,p",
			Usage: "input `PASSWORD`",
		},
		//flag
	},
	Action: update, //main function of tool
}

// Tool desc
func update(c *cli.Context) error {
	email := c.String("email")
	password := c.String("password")
	content, ok := checkVersion(email, password)
	if ok {
		fmt.Println("It is the latest version, no need to update")
	} else {
		fmt.Println("更新:", content.Get("message").Get("newest_version").MustString())
	}
	return nil
}

func checkVersion(email, password string) (*sjson.Json, bool) {
	resp, err := clients.NewHTTPClient(email, password).Get("/api/v1.0/version", nil,nil,"",nil)
	if err != nil {
		log.Fatal("Version checking failed! Please retry.")
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1 {
		newestVersion := content.Get("message").Get("newest_version").MustString()
		newestInt := 0
		curInt := 0
		newS := strings.Split(newestVersion, ".")
		curS := strings.Split(consts.VERSION, ".")
		for index, value := range newS {
			nValue, _ := strconv.Atoi(value)
			cValue, _ := strconv.Atoi(curS[index])
			newestInt += int(math.Pow(10, float64((2-index)*2))) * nValue
			curInt += int(math.Pow(10, float64((2-index)*2))) * cValue
		}
		if curInt < newestInt {
			return content, false
		}

	} else {
		log.Fatalln("账户信息有误，请检查后重新输入")
	}
	return nil, true
}
