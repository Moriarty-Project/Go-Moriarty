package sherlock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSiteCheckWithStatus(t *testing.T) {
	ubse := &StringBasedSiteElement{
		Name:                   "academia.edu",
		UrlHome:                "https://academia.edu",
		UrlUsername:            "https://independent.academia.edu/%s",
		UnclaimedIfResponseHas: "", //uses the status code, so doesn't matter.
		IsNSFW:                 false,
	}
	sut := NewSiteUserTester(ubse)
	assert.True(t, sut.TestSiteHas(false, "Name")[0])
	assert.False(t, sut.TestSiteHas(false, "reallyreallyrarename")[0])
}

func TestSiteCheckWithResponse(t *testing.T) {
	ubse := &StringBasedSiteElement{
		Name:                   "Apple Discussions",
		UrlHome:                "https://discussions.apple.com",
		UrlUsername:            "https://discussions.apple.com/profile/%s",
		UnclaimedIfResponseHas: "The page you tried was not found. You may have used an outdated link or may have typed the address (URL) incorrectly.",
		IsNSFW:                 false,
	}
	sut := NewSiteUserTester(ubse)
	assert.True(t, sut.TestSiteHas(false, "helloWorld")[0])
	assert.False(t, sut.TestSiteHas(false, "reallyasjnhnsujkdhnukdskjahjn")[0])
}

func TestLoadingElements(t *testing.T) {
	jsonFile := "resources/data.json"
	ubses, err := LoadAllSiteElements(jsonFile)
	assert.NoError(t, err)
	ubse := &StringBasedSiteElement{
		Name:                   "Apple Discussions",
		UrlHome:                "https://discussions.apple.com",
		UrlUsername:            "https://discussions.apple.com/profile/%s",
		UnclaimedIfResponseHas: "The page you tried was not found. You may have used an outdated link or may have typed the address (URL) incorrectly.",
		IsNSFW:                 false,
	}
	assert.Equal(t, ubse, ubses[0])
	siteTesters, err := LoadAllSiteUserTesters(jsonFile)
	assert.NoError(t, err)
	assert.Equal(t, ubse, siteTesters["Apple Discussions"].SiteElement)

	sut := siteTesters["Apple Discussions"]
	assert.True(t, sut.TestSiteHas(false, "helloWorld")[0])
	assert.False(t, sut.TestSiteHas(false, "reallyasjnhnsujkdhnukdskjahjn")[0])

}

func BenchmarkWithCache(b *testing.B) {
	// all cached functions should not change with the benchmark!
	ubse := &StringBasedSiteElement{
		Name:                   "academia.edu",
		UrlHome:                "https://academia.edu",
		UrlUsername:            "https://independent.academia.edu/%s",
		UnclaimedIfResponseHas: "", //uses the status code, so doesn't matter.
		IsNSFW:                 false,
	}
	sut := NewSiteUserTester(ubse)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sut.TestSiteHas(true, "Name")
	}

}
