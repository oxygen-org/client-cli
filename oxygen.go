package main

import (
	"log"
	"os"
	cmd "github.com/oxygen-org/client/commands"
	"github.com/oxygen-org/client/consts"
	"sort"
	"time"

	cli "github.com/urfave/cli"
)

func github.com/oxygen-org() {

	app := cli.NewApp()
	app.Name = consts.APPNAME
	app.Version = consts.VERSION
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Tacey Wong",
			Email: "xinyong.wang@xtalpi.com",
		},
		cli.Author{
			Name: "All Contributors",
		},
	}
	app.Copyright = "(c) 2019 - Forever Tacey Wong"
	app.Usage = "Fusion Computing Management Client"
	app.EnableBashCompletion = true
	app.Before = cmd.HandleBefore
	app.After = cmd.HandleAfter
	app.CommandNotFound = cmd.HandleNotFound
	app.OnUsageError = cmd.HandleUsageError
	app.Action = cmd.BaseAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "language",
			EnvVar: "lang",
			Value:  "english",
			Usage:  "language for the greeting",
		},
	}

	app.Commands = []cli.Command{}
	// mount common
	app.Commands = append(app.Commands, cmd.LsCMD)
	app.Commands = append(app.Commands, cmd.RmCMD)
	app.Commands = append(app.Commands, cmd.PackCMD)
	app.Commands = append(app.Commands, cmd.CheckCMD)
	app.Commands = append(app.Commands, cmd.LookCMD)

	// mount data
	app.Commands = append(app.Commands, cmd.ShareCMD)
	app.Commands = append(app.Commands, cmd.ScpCMD)
	app.Commands = append(app.Commands, cmd.FreezeCMD)
	app.Commands = append(app.Commands, cmd.UnFreezeCMD)
	app.Commands = append(app.Commands, cmd.DownloadCMD)

	// mount job
	app.Commands = append(app.Commands, cmd.CreateCMD)
	app.Commands = append(app.Commands, cmd.RunCMD)
	app.Commands = append(app.Commands, cmd.LogCMD)
	app.Commands = append(app.Commands, cmd.QueueCMD)
	app.Commands = append(app.Commands, cmd.KillCMD)



	// mount auth
	app.Commands = append(app.Commands, cmd.RegisterCMD)
	app.Commands = append(app.Commands, cmd.LoginCMD)
	app.Commands =  append(app.Commands, cmd.PasswordCMD)
	app.Commands = append(app.Commands, cmd.KeyCMD)
	app.Commands = append(app.Commands, cmd.UpdateCMD)
	app.Commands = append(app.Commands, cmd.VersionCMD)

	// mount admin
	app.Commands = append(app.Commands, cmd.DeployCMD)
	app.Commands = append(app.Commands, cmd.ClearCMD)
	app.Commands = append(app.Commands, cmd.NodeCMD)
	app.Commands = append(app.Commands, cmd.DispatchCMD)
	app.Commands = append(app.Commands, cmd.MonitorCMD)

	// 对
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	github.com/oxygen-org()
}
