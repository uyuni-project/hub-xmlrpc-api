package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func main() {
	color.Set(color.FgBlue, color.Italic)
	defer color.Unset()
	app := &cli.App{
		Name:  "Hub XMLRPC API",
		Usage: "helps to operate SUSE Manager or Uyuni infrastructures with one Server, called a Hub, managing several Servers.",
	}
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "example",
			Aliases: []string{"ex"},
			Usage:   "list examples to help user start  with",
			Subcommands: []*cli.Command{
				{
					Name:    "manual_login_mode",
					Aliases: []string{"manual"},
					Usage:   "example to show usage of manual mode",
					Action: func(c *cli.Context) error {
						readExampleContents("manual")
						return nil
					},
				},
				{
					Name:    "auto_realy_login_mode",
					Aliases: []string{"relay"},
					Usage:   "example to show usage of relay mode",
					Action: func(c *cli.Context) error {
						readExampleContents("relay")
						return nil
					},
				},
				{
					Name:    "auto_autoconnect_login_mode",
					Aliases: []string{"auto"},
					Usage:   "example to show usage of autoconnect mode",
					Action: func(c *cli.Context) error {
						readExampleContents("auto")
						return nil
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func readExampleContents(mode string) {
	filePath := fmt.Sprintf("../example_scripts/auth_%s_mode.py", mode)
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Could read the file %s", err)
	}
	defer file.Close()
	b, _ := ioutil.ReadAll(file)
	fmt.Println(string(b))
}
