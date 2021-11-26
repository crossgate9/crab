package crawler

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/atomicptr/crab/pkg/filter"
	"github.com/pkg/errors"
)

//Crawler crawls urls in parallel
type Crawler struct {
	HttpClient        http.Client
	NumberOfWorkers   int
	FilterStatusQuery string
	OutWriter         io.Writer
	statusFilter      *filter.Filter
	printMutex        sync.Mutex
}

//Crawl crawls a list of HTTP requests with a set number of workers
func (c *Crawler) Crawl(requests []*http.Request) {
	requestNum := len(requests)
	fmt.Printf("Number of requests: %d\n", requestNum)

	queue := make(chan *http.Request, requestNum)
	for _, req := range requests {
		queue <- req
	}

	wg := sync.WaitGroup{}
	wg.Add(requestNum)

	numberOfWorkers := 1
	if c.NumberOfWorkers > numberOfWorkers {
		numberOfWorkers = c.NumberOfWorkers
	}
	fmt.Printf("Number of workers: %d\n", numberOfWorkers)

	for i := 0; i < numberOfWorkers; i++ {
		go func() {
			for req := range queue {
				c.crawlRequest(req)
				wg.Done()

				fmt.Printf("Sleep for five second.\n")
				time.Sleep(time.Second * 5)
			}
		}()
	}

	wg.Wait()
	close(queue)
}

func (c *Crawler) crawlRequest(req *http.Request) {
	requestStartTime := time.Now()
	res, err := c.HttpClient.Do(req)
	duration := time.Since(requestStartTime)

	if err != nil {
		c.logError(errors.Wrapf(err, "error with url %s", req.URL), req.URL.String(), duration)
		return
	}

	c.log(res.StatusCode, req.URL.String(), duration)
}
