package main

import (
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/vcraescu/go-xrandr"
	"github.com/vcraescu/rescreen/config"
	"github.com/vcraescu/rescreen/layout"
	"gopkg.in/urfave/cli.v1"
	"os"
	"strings"
)

const version = "0.0.2"

var dryRun bool

func main() {
	app := cli.NewApp()
	app.Name = "rescreen"
	app.Usage = "Configure screen layout"
	app.UsageText = "rescreen [command options] config-file"
	app.ArgsUsage = "Config file"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "dry-run",
			Usage:       "Just print the xrandr commands",
			Destination: &dryRun,
		},
	}
	app.Action = doAction

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(aurora.Red(err))
		os.Exit(1)
	}
}

func doAction(c *cli.Context) error {
	var configPath string
	if c.NArg() == 0 {
		return errors.New("config file path is mandatory")
	}

	configPath = c.Args().First()
	cfg, err := config.LoadFile(configPath)
	if err != nil {
		return err
	}

	screens, err := xrandr.GetScreens()
	if err != nil {
		return err
	}

	lt, err := layout.New(*cfg, screens)
	if err != nil {
		return err
	}

	cmd := xrandr.
		Command().
		DPI(lt.DPI).
		ScreenSize(lt.Resolution)

	for _, node := range lt.Nodes {
		ocmd := cmd.
			Output(node.Monitor).
			Scale(node.Scale).
			SetPrimary(node.Primary).
			Position(node.Position)

		if node.Right != nil {
			ocmd = ocmd.LeftOf(node.Right.Monitor)
		}

		cmd = ocmd.EndOutput()
	}

	if !dryRun {
		return cmd.Run()
	}

	return execDryRun(cmd)
}

func execDryRun(cmd xrandr.CommandBuilder) error {
	fmt.Println(aurora.Bold(aurora.Black(aurora.BgGreen(" DRY RUN "))))
	cmds, err := cmd.RunnableCommands()
	if err != nil {
		return err
	}

	for _, cmd := range cmds {
		fmt.Println(aurora.Green(strings.Join(cmd.Args, " ")))
	}

	return nil
}
