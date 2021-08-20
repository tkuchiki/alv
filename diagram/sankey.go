package diagram

import (
	"fmt"
	"io"
	"os"

	"github.com/tkuchiki/alv/parser"

	"github.com/go-echarts/go-echarts/v2/components"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	sankeyLinks       = []opts.SankeyLink{}
	sankeyLinksValues = map[opts.SankeyLink]float32{}
)

type Sankey struct {
}

func NewSankey() *Sankey {
	return &Sankey{}
}

func (s *Sankey) Render(f *os.File, data map[string][]parser.ParsedLog) error {
	page := components.NewPage()
	page.PageTitle = "alv"
	page.AddCharts(
		createSankey(data),
	)

	return page.Render(io.MultiWriter(f))
}

func createNodes(names []string) []opts.SankeyNode {
	nodes := make([]opts.SankeyNode, 0, len(names))
	for _, name := range names {
		nodes = append(nodes, opts.SankeyNode{
			Name: name,
		})
	}

	return nodes
}

type accessLog struct {
	method string
	uri    string
}

func (a accessLog) uriWithMethod() string {
	return fmt.Sprintf("%s %s", a.method, a.uri)
}

// ユーザIDを map のキーに、URL を []string で持つ
// それを元に opts.SankeyLink を生成する
func createLinks(urls []parser.ParsedLog) {
	switch l := len(urls); {
	case l == 0:
		return
	case l == 1:
		link := opts.SankeyLink{
			Source: urls[0].Request(),
		}
		sankeyLinksValues[link]++
		if sankeyLinksValues[link] == 1 {
			sankeyLinks = append(sankeyLinks, link)
		}
	case l > 1:
		for i := 0; i < l-1; i++ {
			link := opts.SankeyLink{
				Source: urls[i].Request(), Target: urls[i+1].Request(),
			}
			sankeyLinksValues[link]++
			if sankeyLinksValues[link] == 1 {
				sankeyLinks = append(sankeyLinks, link)
			}
		}
	}
}

func createSankey(logs map[string][]parser.ParsedLog) *charts.Sankey {
	sankey := charts.NewSankey()
	sankey.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "access pattern",
		}),
	)

	uriMap := map[string]struct{}{}
	var links []opts.SankeyLink
	for _, accessLogs := range logs {
		createLinks(accessLogs)
		for _, accessLog := range accessLogs {
			uriMap[accessLog.Request()] = struct{}{}
		}
	}
	var uris []string
	for uri, _ := range uriMap {
		uris = append(uris, uri)
	}

	nodes := createNodes(uris)

	for _, link := range sankeyLinks {
		val := sankeyLinksValues[link]
		link.Value = val
		links = append(links, link)
	}

	sankey.AddSeries("sankey", nodes, links).
		SetSeriesOptions(
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:     "source",
				Curveness: 0.5,
			}),
			charts.WithLabelOpts(opts.Label{
				Show: true,
			}),
		)
	return sankey
}
