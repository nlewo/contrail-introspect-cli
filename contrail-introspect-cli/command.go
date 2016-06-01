package main

import "fmt"
import "os"
import "strings"

import "github.com/codegangsta/cli"

func GenCommand(descCol DescCollection, name string, usage string) cli.Command {
	return cli.Command{
		Name:      name,
		Aliases:   []string{"a"},
		Usage:     usage,
		ArgsUsage: fmt.Sprintf("%s\n", strings.Join(descCol.PageArgs, " ")),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "long, l",
				Usage: "Long format",
			},
			cli.BoolFlag{
				Name:  "xml, x",
				Usage: "XML output format",
			},
			cli.BoolFlag{
				Name:  "from-file",
				Usage: "Load file instead of URL (for debugging)",
			},
			cli.BoolFlag{
				Name:  "url, u",
				Usage: "Just show URL",
			},
			cli.StringFlag{
				Name:  "search, s",
				Usage: fmt.Sprintf("Search by %s", descCol.PrimaryField),
				Value: "",
			},
		},
		Action: func(c *cli.Context) {
			var page Sourcer
			if c.IsSet("from-file") {
				page = File{Path: c.Args()[0]}
			} else {
				if c.NArg() < len(descCol.PageArgs) {
					cli.ShowSubcommandHelp(c)
					os.Exit(1)
				}
				page = descCol.PageBuilder(c.Args())
			}
			col := page.Load(descCol)
			if c.IsSet("url") {
				fmt.Println(col.url)
				return
			}

			var list Shower

			if c.String("s") != "" {
				list = col.SearchFuzzy(c.String("s"))
			} else {
				list = col
			}

			if c.IsSet("xml") {
				list.Xml()
				return
			}
			if c.IsSet("long") {
				list.Long()
				return
			}
			list.Short()
		},
	}
}

