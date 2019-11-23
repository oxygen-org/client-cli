package commands

import (
	"fmt"
	"log"
	"github.com/oxygen-org/client/clients"
	"strconv"
	"strings"

	"github.com/bndr/gotabulate"

	"github.com/urfave/cli"
	survey "gopkg.in/AlecAivazis/survey.v1"
	//what import
)

const adminGroup = "admin"

// DeployCMD desc
var DeployCMD = cli.Command{
	Name:     "deploy",
	Aliases:  []string{"dp"},
	Category: adminGroup,
	Usage:    "管理员功能，部署镜像至各节点",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "nodes, n",
			Usage: `deploy image, multiple nodes would be separated by ","`,
		},
		//flag
	},
	Action: deploy, //main function of tool
}

func deploy(c *cli.Context) error {
	nodes := c.String("nodes")
	image := ""
	if c.NArg() > 0 {
		image = c.Args().Get(0)
	}
	if nodes == "" && c.NArg() < 1 {
		ansmap := make(map[string]interface{})
		var validationQs = []*survey.Question{
			{
				Name:   "nodes",
				Prompt: &survey.Input{Message: "Please Input Nodes(split with ,):", Default: ""},
			},
			{
				Name:     "image",
				Prompt:   &survey.Input{Message: "Please Input Image:", Default: ""},
				Validate: survey.Required,
			},
		}
		err := survey.Ask(validationQs, &ansmap)
		if err != nil {
			return fmt.Errorf("\n%s", err.Error())
		}
		nodes = ansmap["nodes"].(string)
		image = ansmap["image"].(string)
	}
	data := []byte(fmt.Sprintf(`{"image":"%s","nodes":"%s"}`, nodes, image))
	resp, err := clients.NewHTTPClient("", "").Post("/api/v1.0/image/deploy", nil, data, "json", nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1 {
		fmt.Printf(`image [%s] will pulling, please check result after.`, image)
	} else {
		log.Fatalln(content.Get("message").MustString())
	}
	return nil
}

// ClearCMD desc
var ClearCMD = cli.Command{
	Name:     "clear",
	Aliases:  []string{"cls"},
	Category: adminGroup,
	Usage:    "管理员功能，清除日志缓存",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "job, j",
			Usage: `delete log of the job, multiple ids would be separated by ","`,
		},
		//flag
	},
	Action: clear, //main function of tool
}

func clear(c *cli.Context) error {
	jobs := c.String("job")
	if jobs == "" {
		ansmap := make(map[string]interface{})
		var validationQs = []*survey.Question{
			{
				Name:   "job",
				Prompt: &survey.Input{Message: "Please Input Job(multi split with ,):", Default: ""},
				Validate: func(val interface{}) error {
					err := survey.Required(val)
					if err != nil {
						return err
					}
					jobList := strings.Split(val.(string), ",")
					if len(jobList) < 1 {
						jobs = "all"
					}
					for _, job := range jobList {
						_, err := strconv.Atoi(job)
						if err != nil {
							return fmt.Errorf("\n%s is not numeric", job)
						}
					}
					return nil
				},
			},
		}
		err := survey.Ask(validationQs, &ansmap)
		if err != nil {
			return fmt.Errorf("\n%s", err.Error())
		}
		jobs = ansmap["job"].(string)

	}
	data := []byte(fmt.Sprintf(`{"job_id":"%s"}`, jobs))
	resp, err := clients.NewHTTPClient("", "").Delete("/api/v1.0/job/log", nil, data, "json", nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1 {
		fmt.Printf(`log of job [%s] has been deleted`, jobs)
	} else {
		log.Fatalln(content.Get("message").MustString())
	}
	return nil
}

// NodeCMD desc
var NodeCMD = cli.Command{
	Name:     "node",
	Aliases:  []string{"nd"},
	Category: adminGroup,
	Usage:    "管理员功能，查看/改变节点状态",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "on, o",
			Usage: "turn node status available",
		},
		cli.StringFlag{
			Name:  "node, nd",
			Usage: "chose an node",
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "chose all nodes",
		},

		//flag
	},
	Action: node, //main function of tool
}

func node(c *cli.Context) error {
	ansmap := make(map[string]interface{})
	var validationQs = []*survey.Question{
		{
			Name: "action",
			Prompt: &survey.Select{
				Message: "Please Select One Action:",
				Options: []string{"ON", "OFF", "SHOW"},
			},
		},
	}
	err := survey.Ask(validationQs, &ansmap)
	if err != nil {
		return fmt.Errorf("\n%s", err.Error())
	}
	if ansmap["action"] != "SHOW" {
		var nodeAsk = []*survey.Question{
			{
				Name: "node",
				Prompt: &survey.Input{
					Message: "Please Input Node:",
					Default: "all",
				},
				Validate: survey.Required,
			},
		}
		err = survey.Ask(nodeAsk, &ansmap)
		if err != nil {
			return fmt.Errorf("\n%s", err.Error())
		}
		node := ansmap["node"].(string)
		status := 1
		if ansmap["action"].(string) == "OFF" {
			status = 1
		}
		data := []byte(fmt.Sprintf(`{"node_name":"%s","status":%s}`, node, status))
		resp, err := clients.NewHTTPClient("", "").Patch("/api/v1.0/station", nil, data, "json", nil)
		if err != nil {
			log.Fatalln(err.Error())
		}
		content := resp.GetJSON()
		if content.Get("result").MustInt() == 1 {
			fmt.Println("operate successfully")
		} else {
			fmt.Println(content.Get("message").MustString())
		}
	} else {
		resp, err := clients.NewHTTPClient("", "").Get("/api/v1.0/station", nil, nil, "", nil)
		if err != nil {
			log.Fatalln(err.Error())
		}
		content := resp.GetJSON()
		if content.Get("result").MustInt() == 1 {
			data := content.Get("data")
			dataList := [][]string{}
			for index := range data.MustArray() {
				one := data.GetIndex(index)
				node := one.Get("node").MustString()
				statusByte, _ := one.Get("status").Encode()
				status := string(statusByte)
				dataList = append(dataList, []string{node, status})
			}
			headers := []string{"node", "status"}
			tabulate := gotabulate.Create(dataList)
			tabulate.SetHeaders(headers)
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple
		}
	}

	return nil
}

// DispatchCMD desc
var DispatchCMD = cli.Command{
	Name:     "dispatch",
	Aliases:  []string{"dis"},
	Category: adminGroup,
	Usage:    "管理员功能，查看/更新调度参数",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "update, u",
			Usage: "update dispatcher params",
		},
		//flag
	},
	Action: dispatch, //main function of tool
}

func dispatch(c *cli.Context) error {
	update := c.Bool("update")
	if !update {
		ansmap := make(map[string]interface{})
		var validationQs = []*survey.Question{
			{
				Name: "update",
				Prompt: &survey.Select{
					Message: "Will Update?:",
					Options: []string{"Yes", "No"},
				},
			},
		}
		err := survey.Ask(validationQs, &ansmap)
		if err != nil {
			return fmt.Errorf("\n%s", err.Error())
		}
		if ansmap["update"] == "YES"{
			update = true
		}
	}

	if update{
		resp,err := clients.NewHTTPClient("","").Post("/api/v1.0/dispatcher",nil,nil,"json",nil)
		if err != nil{
			log.Fatalln(err.Error())
		}
		content := resp.GetJSON()
		if content.Get("result").MustInt() == 1{
			fmt.Println("Operation sucess!")
		}else{
			log.Fatalln(content.Get("message").MustString())
		}
	}else{
		resp,err := clients.NewHTTPClient("","").Get("/api/v1.0/dispatcher",nil,nil,"",nil)
		if err != nil{
			log.Fatalln(err.Error())
		}
		content := resp.GetJSON()
		if content.Get("result").MustInt() ==1{
			data:= content.Get("data")
			fmt.Println("Threshold:")
			headerT := []string{"cpu_core", "gpu_memory", "memory", "disk"}
			thresholdList := []string{}
			for _,v := range headerT{
				threshold := data.Get("Threshold")
				thresholdList = append(thresholdList,threshold.Get(v).MustString())
			}
			tabulate := gotabulate.Create(thresholdList)
			tabulate.SetHeaders(headerT)
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple
			fmt.Println("CPU Job:")
			headerC := []string{"cpu", "memory"}
			cpuList := []string{}
			for _,v := range headerC{
				cpu := data.Get("cpu_job")
				cpuList = append(cpuList,cpu.Get(v).MustString())
			}
			tabulate = gotabulate.Create(cpuList)
			tabulate.SetHeaders(headerT)
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple


			fmt.Println("GPU Job:")
			headerG := []string{"cpu", "memory","gpu_memory"}
			gpuList := []string{}
			for _,v := range headerG{
				gpu := data.Get("gpu_job")
				gpuList = append(gpuList,gpu.Get(v).MustString())
			}
			tabulate = gotabulate.Create(gpuList)
			tabulate.SetHeaders(headerG)
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple


			fmt.Println("Sort Parameter:")
			headerS := []string{"cpu", "memory","dispatched_job"}
			sortList := []string{}
			for _,v := range headerS{
				s := data.Get("gpu_job")
				sortList = append(sortList,s.Get(v).MustString())
			}
			tabulate = gotabulate.Create(sortList)
			tabulate.SetHeaders(headerS)
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple


			fmt.Println("Loop Interval:")
			tabulate = gotabulate.Create([]string{data.Get("loop_interval").MustString()})
			tabulate.SetHeaders([]string{"interval",})
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple


		}else{
			log.Fatalln(content.Get("message").MustString())
		}
	}

	fmt.Println("dispatch cmd")
	return nil
}

// MonitorCMD desc
var MonitorCMD = cli.Command{
	Name:     "monitor",
	Aliases:  []string{"mo"},
	Category: adminGroup,
	Usage:    "管理员功能，查看集群资源使用情况",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "now, n",
			Usage: "view resource usage now",
		},
		cli.IntFlag{
			Name:  "day, d",
			Usage: "view data within some days, default today",
			Value: 0,
		},
		cli.StringFlag{
			Name:  "hour, h",
			Usage: "view data for a period of per day, format like 0,24, default is worktime 9:00-20:00",
			Value:"9,20",
		},
		//flag
	},
	Action: monitor, //main function of tool
}

func monitor(c *cli.Context) error {
	fmt.Println("monitor cmd")
	return nil
}
