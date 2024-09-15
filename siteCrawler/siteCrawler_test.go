package siteCrawler

import (
	"fmt"
	"testing"
)


func TestBasicSiteCrawlerFromCache(t *testing.T) {
	scc, err := NewBasicSiteCrawler("testsite.domain/%v", "", "/siteCrawler/crawlTesting/val%v.json")
	if err != nil {
		t.Fatal(err)
	}

	scc.AddAuthenticator(NewBearerAuth("foo"))
	fmt.Println(scc.GetAll(NewDefaultIterator(5000, 6000, 1), 50))
	fmt.Println(scc.cache.Save())
}
