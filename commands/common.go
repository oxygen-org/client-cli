package commands

import (
	"fmt"
	"log"
	"github.com/oxygen-org/client/clients"
	"strconv"
	"strings"

	sjson "github.com/bitly/go-simplejson"
	"github.com/bndr/gotabulate"
	"github.com/urfave/cli"
	survey "gopkg.in/AlecAivazis/survey.v1"
	//what import
)

const commonGroup = "common"

// LsCMD desc
var LsCMD = cli.Command{
	Name:     "ls",
	Aliases:  []string{"l"},
	Category: commonGroup,
	Usage:    "显示数据集、公开数据集、存档文件或任务",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "job, j",
			Usage: "list your jobs",
		},
		cli.BoolFlag{
			Name:  "dataset, d",
			Usage: "list your datasets",
		},
		cli.BoolFlag{
			Name:  "archive, a",
			Usage: "list your archive files",
		},
		cli.BoolFlag{
			Name:  "public, p",
			Usage: "list public datasets",
		},
	},
	Action: ls, //main function of tool
}

func ls(c *cli.Context) error {
	jobS := "(job)list your jobs"
	datasetS := "(dataset)list your datasets"
	archiveS := "(archive)list your archive files"
	publicS := "(public)list public datasets"
	ansmap := make(map[string]interface{})
	var validationQs = []*survey.Question{
		{
			Name: "type",
			Prompt: &survey.Select{
				Message: "Choose a ls-type:",
				Options: []string{jobS, datasetS, archiveS, publicS},
			},
		},
	}
	err := survey.Ask(validationQs, &ansmap)
	if err != nil {
		return fmt.Errorf("\n%s", err.Error())
	}
	switch ansmap["type"].(string) {
	case jobS:
		return lsJob()
	case datasetS:
		return lsDataset()
	case archiveS:
		return lsArchive()
	case publicS:
		lsPublic()
	}
	return nil
}

func subAskInput(msg string) string {
	ansmap := make(map[string]interface{})
	var validationQs = []*survey.Question{
		{
			Name: "value",
			Prompt: &survey.Input{
				Message: msg,
			},
		},
	}
	survey.Ask(validationQs, &ansmap)
	return ansmap["value"].(string)
}
func lsJob() error {
	resp, err := clients.NewHTTPClient("", "").Get("/api/v1.0/jobs", nil, nil, "", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1 {
		jobs := content.Get("jobs")
		jobList := [][]string{}
		for index := range jobs.MustArray() {
			id := jobs.GetIndex(index).Get("id").MustString()
			name := jobs.GetIndex(index).Get("name").MustString()
			state := jobs.GetIndex(index).Get("state").MustString()
			jobList = append(jobList, []string{id, name, state})
		}
		if len(jobList) < 1 {
			fmt.Println("no jobs")
		} else {
			tabulate := gotabulate.Create(jobList)
			tabulate.SetHeaders([]string{"ID", "Name", "State"})
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple
		}
	} else {
		log.Fatalln("Error:", content.Get("message").MustString())
	}
	return nil
}
func lsDataset() error {
	resp, err := clients.NewHTTPClient("", "").Get("/api/v1.0/datasets", nil, nil, "", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1 {
		datasets := content.Get("datasets")
		datasetList := [][]string{}
		for index := range datasets.MustArray() {
			if datasets.GetIndex(index).Get("user").MustString() != content.Get("user").MustString() {
				continue
			}
			id := strconv.Itoa(datasets.GetIndex(index).Get("id").MustInt())
			filename := datasets.GetIndex(index).Get("filename").MustString()
			filetype := datasets.GetIndex(index).Get("type").MustString()
			createDate := datasets.GetIndex(index).Get("created_date").MustString()
			descByte, _ := datasets.GetIndex(index).Get("description").EncodePretty()
			desc := string(descByte)
			isShared := datasets.GetIndex(index).Get("shared").MustBool()
			ownership := ""
			if isShared {
				ownership = "public"
			} else {
				ownership = "private"
			}
			datasetList = append(datasetList, []string{id,
				filename, filetype, createDate, desc, ownership})
		}
		if len(datasetList) < 1 {
			fmt.Println("no dataset")
		} else {
			tabulate := gotabulate.Create(datasetList)
			tabulate.SetHeaders([]string{"Id", "File Name", "File Type",
				"Date", "Description", "Ownership"})
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple
		}
	} else {
		log.Fatalln("Error:", content.Get("message").MustString())
	}

	return nil
}
func lsArchive() error {
	resp, err := clients.NewHTTPClient("", "").Get("/api/v1.0/archive", nil, nil, "", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1 {
		archives := content.Get("data")
		archiveList := [][]string{}
		for index := range archives.MustArray() {
			// if archives.GetIndex(index).Get("user").MustString() != content.Get("user").MustString(){
			// 	continue
			// }
			id := strconv.Itoa(archives.GetIndex(index).Get("id").MustInt())
			filename := archives.GetIndex(index).Get("filename").MustString()
			filetype := archives.GetIndex(index).Get("type").MustString()
			createDate := archives.GetIndex(index).Get("created_date").MustString()
			isShared := archives.GetIndex(index).Get("shared").MustBool()
			ownership := ""
			if isShared {
				ownership = "public"
			} else {
				ownership = "private"
			}
			archiveList = append(archiveList, []string{id,
				filename, filetype, createDate, ownership})
		}
		if len(archiveList) < 1 {
			fmt.Println("no archive")
		} else {
			tabulate := gotabulate.Create(archiveList)
			tabulate.SetHeaders([]string{"Id", "File Name", "File Type",
				"Date", "Ownership"})
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple
		}
	} else {
		log.Fatalln("Error:", content.Get("message").MustString())
	}
	return nil
}
func lsPublic() error {
	resp, err := clients.NewHTTPClient("", "").Get("/api/v1.0/datasets", nil, nil, "", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1 {
		datasets := content.Get("datasets")
		datasetList := [][]string{}
		for index := range datasets.MustArray() {
			if !datasets.GetIndex(index).Get("shared").MustBool() {
				continue
			}
			id := strconv.Itoa(datasets.GetIndex(index).Get("id").MustInt())
			filename := datasets.GetIndex(index).Get("filename").MustString()
			filetype := datasets.GetIndex(index).Get("type").MustString()
			createDate := datasets.GetIndex(index).Get("created_date").MustString()
			descByte, _ := datasets.GetIndex(index).Get("description").EncodePretty()
			desc := string(descByte)
			ownership := "public"
			datasetList = append(datasetList, []string{id,
				filename, filetype, createDate, desc, ownership})
		}
		if len(datasetList) < 1 {
			fmt.Println("no dataset")
		} else {
			tabulate := gotabulate.Create(datasetList)
			tabulate.SetHeaders([]string{"Id", "File Name", "File Type",
				"Date", "Description", "Ownership"})
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple
		}
	} else {
		log.Fatalln("Error:", content.Get("message").MustString())
	}

	return nil
}

// RmCMD desc
var RmCMD = cli.Command{
	Name:     "rm",
	Aliases:  []string{"r"},
	Category: commonGroup,
	Usage:    "删除数据集、存档文件或任务",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "job, j",
			Usage: "delete dataset by dataset id",
		},
		cli.BoolFlag{
			Name:  "dataset, d",
			Usage: "delete dataset by dataset id",
		},
		cli.BoolFlag{
			Name:  "archive, a",
			Usage: "delete archive by archive id",
		},
		//flag
	},
	Action: rm, //main function of tool
}

func rm(c *cli.Context) error {
	jobS := "(job)delete dataset by dataset id"
	datasetS := "(dataset)delete dataset by dataset id"
	archiveS := "(archive)delete archive by archive id"
	ansmap := make(map[string]interface{})
	var validationQs = []*survey.Question{
		{
			Name: "type",
			Prompt: &survey.Select{
				Message: "Choose a rm-type:",
				Options: []string{jobS, datasetS, archiveS},
			},
		},
	}
	err := survey.Ask(validationQs, &ansmap)
	if err != nil {
		return fmt.Errorf("\n%s", err.Error())
	}
	switch ansmap["type"].(string) {
	case jobS:
		jobID := subAskInput("Please Input Job ID")
		return rmJob(jobID)
	case datasetS:
		datasetID := subAskInput("Please Input Dataset ID")
		return rmDataset(datasetID)
	case archiveS:
		archiveID := subAskInput("Please Input Archive ID")
		return rmArchive(archiveID)
	}
	return nil
}

func rmJob(jobID string) error {
	data := []byte(fmt.Sprintf(`{"job_id": %s}`, jobID))
	resp, err := clients.NewHTTPClient("", "").Delete("/api/v1.0/jobs", nil, data, "json", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	fmt.Println(content)
	if content.Get("result").MustInt() == 1 {
		job := content.Get("job")
		log, have := job.CheckGet("notes")
		if have && log.MustString() != "" {
			// s3 hdfs 删除
		}
		fmt.Printf("Job %s has been deleted\n", jobID)
	} else {
		fmt.Println(content.Get("message").MustString())
	}
	return nil
}
func rmDataset(datasetID string) error {
	data := []byte(fmt.Sprintf(`{"dataset_id:%s"}`, datasetID))
	resp, err := clients.NewHTTPClient("", "").Delete("/api/v1.0/datasets", nil, data, "json", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	fmt.Println(content)
	if content.Get("result").MustInt() == 1 {
		dataset := content.Get("datasets")
		fileType := dataset.Get("type").MustString()
		if fileType == "file" {
			// s3 hdfs 删除文件
		} else {
			// s3 hdfs 删除文件夹
		}
		fmt.Printf("Dataset %s has been deleted\n", datasetID)
	} else {
		fmt.Println(content.Get("message").MustString())
	}
	return nil

}
func rmArchive(archiveID string) error {

	archiveURL := fmt.Sprintf(`/api/v1.0/archive/%s/detail`, archiveID)
	resp, err := clients.NewHTTPClient("", "").Get(archiveURL, nil, nil, "", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	archiveInfo := resp.GetJSON()
	if archiveInfo.Get("result").MustInt() != 1 {
		log.Fatalln(archiveInfo.Get("message").MustString())
	}
	glacierID, ok := archiveInfo.Get("data").CheckGet("glacier_id")
	if ok && glacierID.MustString() != "" {
		fmt.Println("delete from AWS glacier...")
		// 从s3 删除
		// 错误：delete failed! please retry again.
	}
	data := []byte(fmt.Sprintf(`{"archive": %s}`,archiveID))
	resp, err = clients.NewHTTPClient("", "").Delete(archiveURL, nil,data,"json",nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1 {
		fmt.Println("Operation success!")
	} else {
		return fmt.Errorf(content.Get("message").MustString())
	}
	return nil
}

// PackCMD desc
var PackCMD = cli.Command{
	Name:     "pack",
	Aliases:  []string{"p"},
	Category: commonGroup,
	Usage:    "自动打包",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name, n",
			Usage: "set name of your image",
		},
		cli.StringFlag{
			Name:  "version, v",
			Usage: "set version of your image",
			Value: "v1",
		},
		//flag
	},
	Action: pack, //main function of tool
}

func pack(c *cli.Context) error {
	name := c.String("name")
	version := c.String("version")
	jobID := ""
	if len(c.Args()) > 0 {
		jobID = c.Args().Get(0)
	}
	if name == "" && version == "v1" {
		ansmap := make(map[string]interface{})

		var validationQs = []*survey.Question{
			{
				Name:   "name",
				Prompt: &survey.Input{Message: "Please Input Image Name:"},
				Validate: func(val interface{}) error {
					err := survey.Required(val)
					if err != nil {
						return err
					}
					str := val.(string)
					if str != strings.ToLower(str) {
						return fmt.Errorf("image name must be lowercase")
					}
					return nil
				},
			},
			{
				Name:     "version",
				Prompt:   &survey.Input{Message: "Please Input Image Version:", Default: "v1"},
				Validate: survey.Required,
			},
			{
				Name:     "job_id",
				Prompt:   &survey.Input{Message: "Please Input Job ID:"},
				Validate: survey.Required,
			},
		}

		err := survey.Ask(validationQs, &ansmap)
		if err != nil {
			return fmt.Errorf("\n%s", err.Error())
		}
		name = ansmap["name"].(string)
		version = ansmap["version"].(string)
		jobID = ansmap["job_id"].(string)
	} else {
		if name != strings.ToLower(name) {
			return fmt.Errorf("image name must be lowercase")
		}
		if jobID == "" {
			return fmt.Errorf("please input job id")
		}
	}
	data := []byte(fmt.Sprintf(`{
		"auto": 1, "name": "%s", "version": "%s", "job_id": %s,
	}`, name, version, jobID))
	resp, err := clients.NewHTTPClient("", "").Post("/api/v1.0/images", nil, data, "json", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1 {
		fmt.Printf(`create job success, image id is %s`,
			content.Get("image_id").MustString())
	} else {
		return fmt.Errorf(content.Get("message").MustString())
	}

	return nil
}

// CheckCMD desc
var CheckCMD = cli.Command{
	Name:     "check",
	Aliases:  []string{"ch"},
	Category: commonGroup,
	Usage:    "查看镜像在节点的分布情况",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "nodes, n",
			Usage: "multiple nodes would be separated by \",\"",
		},
		//flag
	},
	Action: check, //main function of tool
}

func check(c *cli.Context) error {

	nodes := strings.TrimSpace(c.String("nodes"))
	// nodesList := strings.Split(nodes, ",")
	image := ""
	if len(c.Args()) > 0 {
		image = c.Args().Get(0)
	}
	if nodes == "" {
		ansmap := make(map[string]interface{})

		var validationQs = []*survey.Question{
			{
				Name:     "image",
				Prompt:   &survey.Input{Message: "Please Input Image Name:"},
				Validate: survey.Required,
			},
			{
				Name:     "nodes",
				Prompt:   &survey.Input{Message: `Please Input Node,separated By ","`},
				Validate: survey.Required,
			},
		}

		err := survey.Ask(validationQs, &ansmap)
		if err != nil {
			return fmt.Errorf("\n%s", err.Error())
		}
		image = ansmap["image"].(string)
		nodes = ansmap["nodes"].(string)
	}

	data := map[string]string{"image": image, "nodes": nodes}
	clients.NewHTTPClient("", "").Get("/api/v1.0/image/deploy/confirmation", data, nil, "", nil)
	// shit api respond
	return nil
}

// LookCMD desc
var LookCMD = cli.Command{
	Name:     "look",
	Aliases:  []string{"lk"},
	Category: commonGroup,
	Usage:    "目前支持查看job详情,dataset详情(可选树状结构),存档文件详情,节点剩余资源",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "job, j",
			Usage: "look job detail",
		},
		cli.StringFlag{
			Name:  "archive, a",
			Usage: "look archive detail",
		},
		cli.StringFlag{
			Name:  "dataset, d",
			Usage: "look dataset detail",
		},
		cli.StringFlag{
			Name:  "tree,t",
			Usage: "look dataset tree",
		},
		cli.StringFlag{
			Name:  "machine,m",
			Usage: "look all remained resources of machine in cluster",
		},
		cli.StringFlag{
			Name:  "image,i",
			Usage: "look all docker image",
		},
		//flag
	},
	Action: look, //main function of tool
}

func look(c *cli.Context) error {
	jobS := "(job)look job detail"
	datasetS := "(dataset)look dataset detail"
	treeS := "(dataset)look dataset tree"
	archiveS := "(archive)look archive detail"
	machineS := "(machine)look all remained resources of machine in cluster"
	imageS := "(image)look all docker image"
	ansmap := make(map[string]interface{})
	var validationQs = []*survey.Question{
		{
			Name: "type",
			Prompt: &survey.Select{
				Message: "Choose a ls-type:",
				Options: []string{jobS, datasetS, treeS, archiveS, machineS, imageS,},
			},
		},
	}
	err := survey.Ask(validationQs, &ansmap)
	if err != nil {
		return fmt.Errorf("\n%s", err.Error())
	}
	switch ansmap["type"].(string) {
	case jobS:
		jobID := subAskInput("Please Input Job ID")
		return lookJob(jobID)
	case datasetS:
		datasetID := subAskInput("Please Input Dataset ID")
		return lookDataset(datasetID)
	case treeS:
		datasetID := subAskInput("Please Input Dataset ID")
		return lookDataset(datasetID)
	case archiveS:
		archiveID := subAskInput("Please Input Archive ID")
		return lookArchive(archiveID)
	case machineS:
		return lookMachine()
	case imageS:
		return lookImage()
	}
	return nil
	return nil
}

func lookJob(jobID string) error {
	urlPath := fmt.Sprintf(`api/v1.0/job/%s`, jobID)
	resp, err := clients.NewHTTPClient("", "").Get(urlPath, nil, nil, "", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1 {
		jobInfo := content.Get("data")
		jobInfoList := []string{}
		for index := range jobInfo.MustArray() {
			v, _ := jobInfo.GetIndex(index).Encode()
			jobInfoList = append(jobInfoList, string(v))
		}
		if len(jobInfoList) < 1 {
			fmt.Println("no detail")
		} else {
			tabulate := gotabulate.Create(jobInfoList)
			tabulate.SetHeaders([]string{"name", "description", "datasets", "git_url", "git_branch", "git_commit", "docker_image", "command",
				"environment", "use_gpu", "state", "hostname", "start_time", "end_time"})
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple
		}
	} else {
		log.Fatalln("Error:", content.Get("message").MustString())
	}
	return nil
}

func getDatasetInfo(datasetID string) (*sjson.Json, error) {
	urlPath := fmt.Sprintf(`api/v1.0/datasets/%s/detail`, datasetID)
	resp, err := clients.NewHTTPClient("", "").Get(urlPath, nil, nil, "", nil)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()

	return content, nil
}
func lookDataset(datasetID string) error {
	content, err := getDatasetInfo(datasetID)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if content.Get("result").MustInt() == 1 {
		datasetInfo := content.Get("data")
		datasetList := []string{}
		if datasetInfo.Get("filetype").MustString() == "file"{
			for index := range datasetInfo.MustArray() {
				v, _ := datasetInfo.GetIndex(index).Encode()
				datasetList = append(datasetList, string(v))
			}
			if len(datasetList) < 1 {
				fmt.Println("no content")
			} else {
				tabulate := gotabulate.Create(datasetList)
				tabulate.SetHeaders([]string{"filename", "created_date", "filetype", "filesize", "description"})
				tabulate.SetAlign("center")
				fmt.Println(tabulate.Render("grid")) //simple
			}
		}else if datasetInfo.Get("filetype").MustString() == "dir"{
			fmt.Println("need s3 & hdfs query")
			// S3 Query
			// hdfs Query
			// files = h_files if len(h_files) >= len(s3_files) else s3_files
			// if not files:
            //         print("ERROR: dataset is not exist in file system. please contact administrator")
		}
	
	} else {
		log.Fatalln("Query failed, this data is not exist or you are not owner.please try again")
	}
	return nil

	return nil
}

func lookTree(datasetID string) error {
	return nil
}

func lookArchive(archiveID string) error {
	urlPath := fmt.Sprintf(`/api/v1.0/archive/%s/detail`, archiveID)
	resp, err := clients.NewHTTPClient("", "").Get(urlPath, nil, nil, "", nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1 {
		archiveInfo := content.Get("data")
		archiveList := []string{}
		headers := []string{"filename", "filetype", "filesize", "description"}
		for _,value := range headers {
			itemV, _ := archiveInfo.Get(value).Encode()
			archiveList = append(archiveList, string(itemV))
		}
		if len(archiveList) < 1 {
			fmt.Println("no detail")
		} else {
			tabulate := gotabulate.Create(archiveList)
			tabulate.SetHeaders(headers)
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple
		}
	} else {
		log.Fatalln("Error:", content.Get("message").MustString())
	}
	return nil
}

func lookMachine() error {
	resp ,err := clients.NewHTTPClient("","").Get("/api/v1.0/resources",nil,nil,"",nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1{
		resource := content.Get("data")
		resourceList := [][]string{}
		for index := range resource.MustArray(){
			resourceItem := resource.GetIndex(index)
			node := resourceItem.Get("node").MustString()
			idleGPU := resourceItem.Get("gpu").Get("unused").MustString()
			if idleGPU == ""{
				idleGPU = "all busy"
			}
			idleGPUMemeoryByte,_ := resourceItem.Get("gpu").Get("gpu_memory_remain").Encode()
			idleGPUMemeory := string(idleGPUMemeoryByte)
			idleCPU := resourceItem.Get("cpu").Get("unused").MustString()
			idleDisk := resourceItem.Get("disk_free").MustString()
			idleMem := resourceItem.Get("memory").Get("unused").MustString()
			resourceList = append(resourceList,[]string{
				node,idleGPU,idleGPUMemeory,idleCPU,idleMem,idleDisk,
			})
			headers := []string{
				"node", "idle GPU", "idle GPU memory(G)", "idle CPU(%)", 
				"idle memory(G)", "idle disk(G)",
			}
			tabulate := gotabulate.Create(resourceList)
			tabulate.SetHeaders(headers)
			tabulate.SetAlign("center")
			fmt.Println(tabulate.Render("grid")) //simple
		}
	}
	return nil
}


func lookImage() error{
	resp ,err := clients.NewHTTPClient("","").Get("/api/v1.0/resources",nil,nil,"",nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	content := resp.GetJSON()
	if content.Get("result").MustInt() == 1{
		imageInfo := content.Get("data")
		imageList := [][]string{}
		for index := range imageInfo.MustArray(){
			image := imageInfo.GetIndex(index)
			imageList = append(imageList,[]string{
				image.Get("name").MustString(),
				image.Get("version").MustString(),
				image.Get("description").MustString(),
				image.Get("created_date").MustString(),
			})
		}
		headers := []string{"name", "version", "description", "created_date",}
		tabulate := gotabulate.Create(imageList)
		tabulate.SetHeaders(headers)
		tabulate.SetAlign("center")
		fmt.Println(tabulate.Render("grid")) //simple

	}else{
		log.Fatalln(content.Get("message").MustString())
	}
	return nil
}