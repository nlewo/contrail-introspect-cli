package utils

import "sort"
import "strings"

import "github.com/moovweb/gokogiri/xml"

func Pretty(nodes []xml.Node) string {
	ret := make([]string, len(nodes))
	for i, n := range nodes {
		ret[i] = ResolveIp(n.Content())
	}
	sort.Strings(ret)
	return strings.Join(ret, " ; ")
}
