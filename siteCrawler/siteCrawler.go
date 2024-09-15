package siteCrawler

import (
	"GoMoriarty/utils"
	"context"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/carlmjohnson/requests"
)

// site crawler is in charge of going to an endpoint that will allow it, and just downloading all of the items in a row.
// really, only works on poorly implemented backends.

type SiteCrawler struct {
	urlSetter      func(string) *requests.Builder
	savePath       string                         //must contain one replaceable field. This is where items are saved.
	authenticators []SiteAuthenticator            //here should be things like, adding cookies to a request if needed, ect.
	validators     []func(r *http.Response) error //check to see if it is a valid response
	cache          *siteCrawlerCache
}

func NewSiteCrawler(urlSetter func(string) *requests.Builder, savePath string) (*SiteCrawler, error) {
	abs, _ := utils.GetAbsolutePath(path.Dir(savePath))

	sc := &SiteCrawler{
		urlSetter:      urlSetter,
		savePath:       path.Join(abs, path.Base(savePath)),
		authenticators: []SiteAuthenticator{},
		validators:     []func(r *http.Response) error{},
		cache:          LoadOrNewSiteCrawlerCache(savePath),
	}

	return sc, nil
}

// pass "" for method to use GET, url must contain exactly one replacement case. eg, %v
func NewBasicSiteCrawler(url, method, savePath string) (*SiteCrawler, error) {
	return NewSiteCrawler(
		func(val string) *requests.Builder {
			req := requests.URL(fmt.Sprintf(url, val))
			return req
		},
		savePath,
	)
}
func (sc *SiteCrawler) AddAuthenticator(authenticators ...SiteAuthenticator) {
	sc.authenticators = append(sc.authenticators, authenticators...)
}
func (sc *SiteCrawler) AddValidator(validators ...func(r *http.Response) error) {
	sc.validators = append(sc.validators, validators...)
}

// to ignore maxFailInRow, pass -1
// returns true if worked, false if maxFailInRow was surpassed
func (sc *SiteCrawler) GetAll(iterator SiteIterator, maxFailInRow int) bool {
	failStreak := 0
	for val, exists := iterator.Next(); exists && failStreak != maxFailInRow; val, exists = iterator.Next() {
		// check if it's in the cache
		if !sc.cache.IsNew(val) {
			// we've already seen it.
			if sc.cache.IsFound(val) {
				// it was found correctly already.
				failStreak = 0
			} else {
				// we couldn't find it correctly.
				failStreak += 1
			}
			continue
		}
		// get the URL formatting.
		err := sc.Get(val, fmt.Sprintf(sc.savePath, val))
		sc.cache.AddFound(val, err == nil)
		if err != nil {
			failStreak += 1
			fmt.Println(err)
			fmt.Println("")
		} else {
			failStreak = 0
		}

	}
	return true
}

// Get assumes all of the work for incrementing the URL has already been done.
func (sc *SiteCrawler) Get(val string, savedTo string) (err error) {
	ctx, cancelFunc := context.WithDeadline(context.TODO(), time.Now().Add(time.Second*60))
	defer cancelFunc()

	req := sc.urlSetter(val).
		ToFile(savedTo)

	// add all the auths.
	for _, auth := range sc.authenticators {
		auth.AddAuthentication(req)
	}
	// add any extra validators
	for _, validator := range sc.validators {
		req.AddValidator(validator)
	}

	return req.Fetch(ctx)
}
