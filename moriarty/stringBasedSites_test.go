package moriarty

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
	sut := NewDataToUserTester(ubse)
	assert.True(t, sut.TestSourceHas("Name")[0])
	assert.False(t, sut.TestSourceHas("reallyreallyrarename")[0])
}
func TestSiteCheckWithResponse(t *testing.T) {
	ubse := &StringBasedSiteElement{
		Name:                   "Apple Discussions",
		UrlHome:                "https://discussions.apple.com",
		UrlUsername:            "https://discussions.apple.com/profile/%s",
		UnclaimedIfResponseHas: "The page you tried was not found. You may have used an outdated link or may have typed the address (URL) incorrectly.",
		IsNSFW:                 false,
	}
	sut := NewDataToUserTester(ubse)
	assert.True(t, sut.TestSourceHas("helloWorld")[0])
	assert.False(t, sut.TestSourceHas("reallyasjnhnsujkdhnukdskjahjn")[0])
}

func TestLoadingElements(t *testing.T) {
	ubses, err := LoadAllSiteElements(jsonFile)
	assert.NoError(t, err)
	ubse := &StringBasedSiteElement{
		Name:                   "GitHub",
		UrlHome:                "https://www.github.com/",
		UrlUsername:            "https://www.github.com/%s",
		UnclaimedIfResponseHas: "",
		IsNSFW:                 false,
	}

	// github actions have an issue with this site for some reason. So we load it separately here
	assert.Equal(t,
		ubses[0],
		&StringBasedSiteElement{
			Name:                   "1337x",
			UrlHome:                "https://www.1337x.to/",
			UrlUsername:            "https://www.1337x.to/user/%s/",
			UnclaimedIfResponseHas: "\u003ctitle\u003eError something went wrong.\u003c/title\u003e",
			IsNSFW:                 false,
		},
		"loaded order differs from expected order",
	)

	siteTesters, err := LoadAllSiteUserTesters(jsonFile)
	assert.NoError(t, err)
	assert.Equal(t, ubse, siteTesters["GitHub"].dataElement.(*StringBasedSiteElement))

	sut := siteTesters["GitHub"]
	assert.True(t, sut.TestSourceHasAny("PyroHedgehog"))
	assert.False(t, sut.TestSourceHasAny("reallyasjnhnsujkdhnukdskjahjn"))
}
