package requests

import "fmt"
import "github.com/codegangsta/cli"
import "log"

import "github.com/nlewo/contrail-introspect-cli/utils"

func Path() cli.Command {
	return cli.Command{
		Name:      "controller-path",
		Usage:     "Find the path to go from a prefix to another on in a routing instance",
		ArgsUsage: "controller-fqdn ri source-prefix dest-prefix",
		Action: func(c *cli.Context) error {

			controller := c.Args()[0]
			ri := c.Args()[1]
			srcIp := c.Args()[2]
			dstIp := c.Args()[3]

			page := DescCtrlRoute().PageBuilder([]string{controller, ri})
			col := page.Load(DescCtrlRoute())
			elt := col.SearchStrict(srcIp)
			if len(elt) < 1 {
				log.Fatal(fmt.Sprintf("Prefix %s not found in RI %s", srcIp, ri))
			}
			srcNode := elt[0].GetField("paths/list/ShowRoutePath/next_hop")

			page = DescCtrlRoute().PageBuilder([]string{controller, ri})
			col = page.Load(DescCtrlRoute())
			elt = col.SearchStrict(dstIp)
			if len(elt) < 1 {
				log.Fatal(fmt.Sprintf("Prefix %s not found in RI %s", dstIp, ri))
			}
			dstNode := elt[0].GetField("paths/list/ShowRoutePath/next_hop")
			label := elt[0].GetField("paths/list/ShowRoutePath/label")

			fmt.Printf("From prefix %s on %s to dst %s on %s with label %s\n", srcIp, utils.ResolveIp(srcNode), dstIp, utils.ResolveIp(dstNode), label)

			return nil
		},
	}
}
