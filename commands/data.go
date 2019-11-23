package commands

import (
	"fmt"

	"github.com/urfave/cli"
	//what import
)

const dataGroup = "data"

// ShareCMD desc
var ShareCMD = cli.Command{
	Name:     "share",
	Aliases:  []string{"sh"},
	Category: dataGroup,
	Usage:    "dataset文件共享或取消共享",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "cancel, c",
			Usage: "unshared file",
		},
		//flag
	},
	Action: share, //main function of tool
}

// Tool desc
func share(c *cli.Context) error {
	fmt.Println("share cmd")
	return nil
}


// ScpCMD desc
var ScpCMD = cli.Command{
	Name:     "scp",
	Aliases:  []string{"sc"},
	Category: dataGroup,
	Usage:    "上传文件或文件夹",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "file, f",
			Usage: "upload your file or dir",
			Hidden: true,
		},
		cli.StringFlag{
			Name: ",m",
			Usage:"description",
			Hidden: true,
		},
		cli.BoolFlag{
			Name: "share,s",
			Usage: "if you want share this dataset, add this option",
		},
		
	},
	Action: scp, //main function of tool
}

// Tool desc
func scp(c *cli.Context) error {
	fmt.Println("scp cmd")
	return nil
}


// FreezeCMD desc
var FreezeCMD = cli.Command{
	Name:     "freeze",
	Aliases:  []string{"fz"},
	Category: dataGroup,
	Usage:    "上传存档文件",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "resume, r",
			Usage: "upload an archive file from breakpoint",
		},
		cli.BoolFlag{
			Name: "share,s",
			Usage: "if you want share this dataset, add this option",
		},
		cli.StringFlag{
			Name: "dataset,d",
			Value: ".",
			Usage: "freeze a dataset use dataset id",
		},
		//flag
	},
	Action: freeze, //main function of tool
}

// Tool desc
func freeze(c *cli.Context) error {
	fmt.Println("freeze cmd")
	return nil
}


// UnFreezeCMD desc
var UnFreezeCMD = cli.Command{
	Name:     "unfreeze",
	Aliases:  []string{"ufz"},
	Category: dataGroup,
	Usage:    "存档文件下载初始化",
	
	Action: unFreeze, //main function of tool
}

// Tool desc
func unFreeze(c *cli.Context) error {
	fmt.Println("unfreeze cmd")
	return nil
}


// DownloadCMD desc
var DownloadCMD = cli.Command{
	Name:     "down",
	Aliases:  []string{"dw"},
	Category: dataGroup,
	Usage:    "下载文件",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "dir, d",
			Usage: "download remote dir by remote dir name which must be under your name",
		},
		cli.StringFlag{
			Name: "file,f",
			Usage: "download remote file by remote file name which must be under your name",
		},
		cli.StringFlag{
			Name: "path,p",
			Usage: "download path, default is current directory",
		},
		cli.StringFlag{
			Name: "archive,a",
			Usage: "archive id",
		},
	},
	Action: download, //main function of tool
}

// Tool desc
func download(c *cli.Context) error {
	fmt.Println("down cmd")
	return nil
}

