package main

import (
	"fmt"

	"github.com/javif89/dotenv"
	"github.com/urfave/cli/v2"

	// "github.com/javif89/dotenv"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name: "dotenv",
		Usage: "Create and manipulate .env files in your system",
		UsageText: "dotenv [path to file] command [command options]",
		Authors: []*cli.Author{
			{
				Name: "Javier Feliz",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "set",
				Aliases: []string{"s"},
				Usage: "Set an environment variable",
				UsageText: "dotenv set -f [path to file] -k [key] -v [value]",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "path",Value: "./.env",Usage: "Path to the `file`",Aliases: []string{"p", "f"},},
					&cli.StringFlag{Name: "key",Value: "",Usage: "Key to set",Aliases: []string{"k"},},
					&cli.StringFlag{Name: "value",Value: "",Usage: "Value to set",Aliases: []string{"v"},},
				},
				Action: func(c *cli.Context) error {
					// Validate parameters
					if c.String("key") == "" || c.String("value") == "" {
						return cli.Exit("Key and value are required", 1)
					}

					file := dotenv.LoadOrCreate(c.String("path"))
					file.Add(c.String("key"), c.String("value"))
					file.Save()
					return nil
				},
			},
			{
				Name: "get",
				Aliases: []string{"g"},
				Usage: "Get the value of a key in a file",
				UsageText: "dotenv get -f [path to file] -k [key]",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "path",Value: "./.env",Usage: "Path to the `file`",Aliases: []string{"p", "f"},},
					&cli.StringFlag{Name: "key",Value: "",Usage: "Key to set",Aliases: []string{"k"},},
				},
				Action: func(c *cli.Context) error {
					// Validate parameters
					if c.String("key") == "" {
						return cli.Exit("Key is required", 1)
					}

					file := dotenv.Load(c.String("path"))
					fmt.Println(file.Get(c.String("key")))
					return nil
				},
			},
			{
				Name: "keys",
				Aliases: []string{"k"},
				Usage: "List all the keys in a file",
				UsageText: "dotenv keys -f [path to file]",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "path",Value: "./.env",Usage: "Path to the `file`",Aliases: []string{"p", "f"},},
				},
				Action: func(c *cli.Context) error {
					file := dotenv.Load(c.String("path"))
					for _, k := range file.Keys() {
						fmt.Println(k)
					}
					return nil
				},
			},
		},
	}

	// Run the app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}