package crawl

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/atomicptr/crab/pkg/crawler"
	"github.com/spf13/cobra"
)

func registerStandardCrawlCommandFlags(cmd *cobra.Command, flagOptions *crawlerFlagOptions) {
	cmd.PersistentFlags().IntVarP(
		&flagOptions.NumberOfWorkers,
		"num-workers",
		"",
		defaultNumberOfWorkers,
		"set number of workers for crawling",
	)
	cmd.PersistentFlags().StringVarP(
		&flagOptions.PrefixUrl,
		"prefix-url",
		"",
		"",
		"prefix/replace all request urls with this one",
	)
	cmd.PersistentFlags().DurationVarP(
		&flagOptions.HttpTimeout,
		"http-timeout",
		"",
		defaultHttpTimeout,
		"set http timeout for requests",
	)
	cmd.PersistentFlags().StringVarP(
		&flagOptions.AuthUsername,
		"auth-user",
		"",
		"",
		"set HTTP basic authentication username",
	)
	cmd.PersistentFlags().StringVarP(
		&flagOptions.AuthPassword,
		"auth-pass",
		"",
		"",
		"set HTTP basic authentication password",
	)
	cmd.PersistentFlags().StringSliceVarP(
		&flagOptions.CookieStrings,
		"cookie",
		"",
		nil,
		"add cookies (as key=value pairs) to each request",
	)
	cmd.PersistentFlags().StringSliceVarP(
		&flagOptions.HeaderStrings,
		"header",
		"",
		nil,
		"add headers (as key=value pairs) to each request",
	)
	cmd.PersistentFlags().StringVarP(
		&flagOptions.FilterStatusQuery,
		"filter-status",
		"",
		"",
		"filter logs by status",
	)
}

func registerStandardCrawlCommandFlagModifiers(modifier *crawler.RequestModifier, flagOptions crawlerFlagOptions) {
	modifier.With(addUserAgentToRequest())

	if isValidUrl(flagOptions.PrefixUrl) {
		modifier.With(addPrefixUrlToRequest(flagOptions.PrefixUrl))
	}

	if len(flagOptions.AuthUsername) > 0 || len(flagOptions.AuthPassword) > 0 {
		modifier.With(addHttpBasicAuthToRequest(flagOptions))
	}

	if len(flagOptions.CookieStrings) > 0 {
		modifier.With(addCookiesToRequest(flagOptions))
	}

	if len(flagOptions.HeaderStrings) > 0 {
		modifier.With(addHeadersToRequest(flagOptions))
	}
}

func crawlUrls(urls []string, agent string, modifier crawler.RequestModifier, flagOptions crawlerFlagOptions, outWriter io.Writer) error {
	requests, err := crawler.CreateRequestsFromUrls(urls, modifier)
	if err != nil {
		return err
	}

	for _, req := range requests {
		requestUrl := req.URL.String()
		if strings.HasPrefix(requestUrl, "/") {
			return fmt.Errorf("you can't use a partial URL paths without specifying a prefix-url: %s", requestUrl)
		}
	}

	crawl := crawler.Crawler{
		NumberOfWorkers: flagOptions.NumberOfWorkers,
		HttpClient: http.Client{
			Timeout: flagOptions.HttpTimeout,
		},
		FilterStatusQuery: flagOptions.FilterStatusQuery,
		OutWriter:         outWriter,
	}
	crawl.Crawl(requests, agent)

	return nil
}

func isValidUrl(rawUrl string) bool {
	if rawUrl == "" {
		return false
	}

	_, err := url.Parse(rawUrl)
	return err == nil
}
