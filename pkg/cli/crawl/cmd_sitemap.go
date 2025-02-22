package crawl

import (
	"fmt"
	"github.com/atomicptr/crab/pkg/crawler"
	"github.com/atomicptr/crab/pkg/sitemap"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var sitemapCommandFlags = newCrawlerFlagOptions()

var SitemapCommand = &cobra.Command{
	Use:   "crawl:sitemap [sitemapPath] [agent]",
	Short: "Crawl through a sitemap",
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateCrawlerFlagOptions(sitemapCommandFlags); err != nil {
			fmt.Printf("Flag options are invalid:\n\t%s\n", err)
			os.Exit(1)
		}

		if len(args) != 2 {
			fmt.Println("You have to specify exactly one url or file path to a sitemap xml\n" +
				"\tUsage: crab crawl:sitemap https://domain.com/sitemap.xml 'some agent'")
			os.Exit(1)
		}

		modifier := crawler.RequestModifier{}
		registerStandardCrawlCommandFlagModifiers(&modifier, sitemapCommandFlags)

		sitemapPath := args[0]
		agent := args[1]

		urls, err := sitemap.FetchUrlsFromPath(
			sitemapPath,
			&http.Client{
				Timeout: sitemapCommandFlags.HttpTimeout,

			})
		if err != nil {
			fmt.Printf("Could not read sitemap from %s\n\t%s\n", sitemapPath, err)
			os.Exit(1)
		}

		err = crawlUrls(urls, agent, modifier, sitemapCommandFlags, cmd.OutOrStdout())
		if err != nil {
			fmt.Printf("Could not create crawlable URLs:\n\t%s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	registerStandardCrawlCommandFlags(SitemapCommand, &sitemapCommandFlags)
}
