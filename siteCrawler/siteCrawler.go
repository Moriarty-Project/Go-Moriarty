package siteCrawler

import (
	"GoMoriarty/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
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
	val, exists := iterator.Next()
	for exists {
		// check if it's in the cache
		if !sc.cache.IsNew(val) {
			// we've already seen it.
			if sc.cache.IsFound(val) {
				// it was found correctly already.
				failStreak = 0
			} else {
				// we couldn't find it correctly.
				failStreak += 1
				if failStreak == maxFailInRow {
					return false
				}
			}
		} else {
			// get the URL formatting.
			err := sc.Get(val, fmt.Sprintf(sc.savePath, val))
			sc.cache.AddFound(val, err == nil)
			if err != nil {
				failStreak += 1
				if failStreak == maxFailInRow {
					return false
				}
				fmt.Println(err)
				fmt.Println("")
			} else {
				failStreak = 0
			}
		}
		val, exists = iterator.Next()
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

type SiteAuthenticator interface {
	AddAuthentication(*requests.Builder)
}
type CookieAuth struct {
	cookies map[string]string
}

func NewCookieAuth(from string) *CookieAuth {
	ca := &CookieAuth{
		cookies: map[string]string{},
	}
	parts := strings.Split(from, "; ")
	for _, part := range parts {
		keyVal := strings.Split(part, "=")
		ca.cookies[keyVal[0]] = keyVal[1]
	}
	return ca
}
func (ca CookieAuth) AddAuthentication(req *requests.Builder) {
	for key, val := range ca.cookies {
		req.Cookie(key, val)
	}
}

type BearerAuth struct {
	val string
}

func NewBearerAuth(val string) BearerAuth {
	val, _ = strings.CutPrefix(val, "Bearer ")
	return BearerAuth{
		val: val,
	}
}
func (ba BearerAuth) AddAuthentication(req *requests.Builder) {
	req.Bearer(ba.val)
}

type SiteIterator interface {
	Next() (val string, exists bool) //get the next value and increment the next iteration
	Peak(n int) (string, bool)       //peak ahead, but doesn't actually skip ahead.
	Skip(n int)                      //skip ahead by n points in the iteration.
}
type DefaultIterator struct {
	current  int
	max      int
	stepSize int
}

func NewDefaultIterator(start, stop, step int) *DefaultIterator {
	return &DefaultIterator{
		current:  start,
		max:      stop,
		stepSize: step,
	}
}
func (di *DefaultIterator) Next() (found string, exists bool) {
	if di.current > di.max {
		return "", false
	}
	found = fmt.Sprintf("%v", di.current)
	di.current += di.stepSize
	return found, true
}
func (di *DefaultIterator) Peak(n int) (found string, exists bool) {
	val := di.current + n
	if n > di.max {
		return "", false
	}
	return fmt.Sprintf("%v", val), true
}
func (di *DefaultIterator) Skip(n int) {
	di.current += n
}

type siteCrawlerCache struct {
	savePath string //where this file is saved.
	Found    map[string]bool
}

func LoadOrNewSiteCrawlerCache(savePath string) *siteCrawlerCache {
	absPath, _ := utils.GetAbsolutePath(path.Dir(savePath))
	scc := &siteCrawlerCache{
		savePath: absPath,
		Found:    map[string]bool{},
	}
	file, err := os.Open(path.Join(scc.savePath, "cache.json"))
	if err == nil {
		// file exists
		defer file.Close()
		data, err := io.ReadAll(file)
		if err == nil {
			json.Unmarshal(data, scc)
		}

	}
	return scc
}
func (scc *siteCrawlerCache) Save() error {
	file, err := os.Create(path.Join(scc.savePath, "cache.json"))
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(scc)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}
func (scc *siteCrawlerCache) IsNew(val string) bool {
	_, exists := scc.Found[val]
	return !exists
}
func (scc *siteCrawlerCache) IsFound(val string) bool {
	return scc.Found[val]
}
func (scc *siteCrawlerCache) AddFound(val string, correctly bool) {
	scc.Found[val] = correctly
}
