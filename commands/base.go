package commands

import (
	"fmt"
	"github.com/oxygen-org/client/consts"
	"github.com/urfave/cli"
)



// HandleBefore Before钩子
func HandleBefore(c *cli.Context) error {
	fmt.Fprintf(c.App.Writer, "这是前置钩子\n")
	return nil
}

// HandleAfter After钩子
func HandleAfter(c *cli.Context) error {
	fmt.Fprintf(c.App.Writer, "这是后置钩子\n")
	return nil
}

// HandleNotFound NotFound钩子
func HandleNotFound(c *cli.Context, cmd string) {
	fmt.Fprintf(c.App.Writer, "这是not found\n")
	fmt.Fprintf(c.App.Writer, "Thar be no %q here.\n", cmd)
}

// HandleUsageError UsageError钩子
func HandleUsageError(c *cli.Context, err error, isSubcommand bool) error {
	if isSubcommand {
		return err
	}

	fmt.Fprintf(c.App.Writer, "这是WRONG: %#v\n", err)
	return nil
}

func showLogo() {
	logo := `                                                   
	 @@@@@@   @@@  @@@  @@@ @@@   @@@@@@@@  @@@@@@@@  @@@  @@@  
	@@@@@@@@  @@@  @@@  @@@ @@@  @@@@@@@@@  @@@@@@@@  @@@@ @@@  
	@@!  @@@  @@!  !@@  @@! !@@  !@@        @@!       @@!@!@@@  
	!@!  @!@  !@!  @!!  !@! @!!  !@!        !@!       !@!!@!@!  
	@!@  !@!   !@@!@!    !@!@!   !@! @!@!@  @!!!:!    @!@ !!@!  
	!@!  !!!    @!!!      @!!!   !!! !!@!!  !!!!!:    !@!  !!!  
	!!:  !!!   !: :!!     !!:    :!!   !!:  !!:       !!:  !!!  
	:!:  !:!  :!:  !:!    :!:    :!:   !::  :!:       :!:  !:!  
	::::: ::   ::  :::     ::     ::: ::::   :: ::::   ::   ::  
	 : :  :    :   ::      :      :: :: :   : :: ::   ::    :                                 
  `
	fmt.Println(logo)
	fmt.Println("\t", consts.APPNAME, consts.VERSION)
	fmt.Println("\tPlease type `github.com/oxygen-org -h/--help` for the help of usage")
}

// BaseAction x
func BaseAction(c *cli.Context) error {
	showLogo()
	return nil
}


