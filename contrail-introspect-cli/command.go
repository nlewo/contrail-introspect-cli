package main

import "fmt"
import "os"
import "strings"

import "github.com/codegangsta/cli"

func GenCommand(descCol DescCollection, name string, usage string) cli.Command {
	return cli.Command{
		Name:      name,
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
				Usage: "Just show the used URL",
			},
			cli.StringFlag{
				Name:  "search, s",
				Usage: fmt.Sprintf("Fuzzy search by %s", descCol.PrimaryField),
				Value: "",
			},
			cli.StringFlag{
				Name:  "strict-search, S",
				Usage: fmt.Sprintf("Strict search by %s", descCol.PrimaryField),
				Value: "",
			},
		},
		Action: func(c *cli.Context) error {
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
				return nil
			}

			var list Shower

			if c.String("s") != "" {
				list = col.SearchFuzzy(c.String("s"))
			} else if c.String("S") != "" {
				list = col.SearchStrict(c.String("S"))
			} else {
				list = col
			}
			
			if c.IsSet("xml") {
				list.Xml()
				return nil
			}
			if c.IsSet("long") {
				list.Long()
				return nil
			}
			list.Short()

			return nil
		},
		BashComplete: func(c *cli.Context) {
			// We only complete the first argument
			if c.NArg() == 0 {
				for _, fqdn := range hosts {
					fmt.Println(fqdn)
				}
			}
		},
	}
}
