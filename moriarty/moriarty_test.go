package moriarty

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const jsonFile = "resources/data.json"

func TestSiteCheckWithStatus(t *testing.T) {
	ubse := &StringBasedSiteElement{
		Name:                   "academia.edu",
		UrlHome:                "https://academia.edu",
		UrlUsername:            "https://independent.academia.edu/%s",
		UnclaimedIfResponseHas: "", //uses the status code, so doesn't matter.
		IsNSFW:                 false,
	}
	sut := NewSiteUserTester(ubse)
	assert.True(t, sut.TestSiteHas("Name")[0])
	assert.False(t, sut.TestSiteHas("reallyreallyrarename")[0])
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
	assert.True(t, sut.TestSiteHas("helloWorld")[0])
	assert.False(t, sut.TestSiteHas("reallyasjnhnsujkdhnukdskjahjn")[0])
}

func TestLoadingElements(t *testing.T) {
	ubses, err := LoadAllSiteElements(jsonFile)
	assert.NoError(t, err)
	ubse := &StringBasedSiteElement{
		Name:                   "1337x",
		UrlHome:                "https://www.1337x.to/",
		UrlUsername:            "https://www.1337x.to/user/%s/",
		UnclaimedIfResponseHas: "\u003ctitle\u003eError something went wrong.\u003c/title\u003e",
		IsNSFW:                 false,
	}
	assert.Equal(t, ubse, ubses[0])
	siteTesters, err := LoadAllSiteUserTesters(jsonFile)
	assert.NoError(t, err)
	assert.Equal(t, ubse, siteTesters["1337x"].dataElement.(*StringBasedSiteElement))

	sut := siteTesters["1337x"]
	assert.True(t, sut.TestSiteHasAny("helloWorld"))
	assert.False(t, sut.TestSiteHasAny("reallyasjnhnsujkdhnukdskjahjn"))

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
		sut.TestSiteHas("Name")
	}
}
func TestMoriartyMinimally(t *testing.T) {
	user := NewUserRecordings("", nil, nil, nil)
	user.AddKnownUsername("helloWorld")

	arty, err := NewMoriarty(jsonFile)
	if err != nil {
		t.Fatal(err)
	}
	arty.AssignNewUser(user)
	knownSites, likelySites, possibleSites := arty.TrackUserAcrossSites()

	assert.Contains(t, knownSites, "Apple Discussions")
	assert.Equal(t, likelySites, []string{})
	assert.Equal(t, possibleSites, []string{})
}

// test Moriarty's whole functionality!
func TestMoriarty(t *testing.T) {
	user := NewUserRecordings("", nil, nil, nil)
	user.AddKnownUsername("helloWorld")

	arty, err := NewMoriarty(jsonFile)
	if err != nil {
		t.Fatal(err)
	}
	arty.AssignNewUser(user)
	knownSites, likelySites, possibleSites := arty.TrackUserAcrossSites()
	assert.Contains(t, knownSites, "Apple Discussions")
	assert.Equal(t, likelySites, []string{})
	assert.Equal(t, possibleSites, []string{})

	// now test it again but with the names moved to likely
	user.LikelyUsernames = user.KnownUsernames
	user.KnownUsernames = []string{}
	knownSites, likelySites, possibleSites = arty.TrackUserAcrossSites()
	assert.Contains(t, likelySites, "Apple Discussions")
	assert.Equal(t, knownSites, []string{})
	assert.Equal(t, possibleSites, []string{})

	// and now with possible
	user.PossibleUsernames = user.LikelyUsernames
	user.LikelyUsernames = []string{}
	knownSites, likelySites, possibleSites = arty.TrackUserAcrossSites()
	assert.Contains(t, possibleSites, "Apple Discussions")
	assert.Equal(t, knownSites, []string{})
	assert.Equal(t, likelySites, []string{})
	// AddSiteElement(jsonFile)
}

func TestMoriartyChannels(t *testing.T) {
	user := NewUserRecordings("", nil, nil, nil)
	user.AddKnownUsername("helloWorld")

	arty, err := NewMoriarty(jsonFile)
	if err != nil {
		t.Fatal(err)
	}
	arty.AssignNewUser(user)

	shouldBeFull, shouldBeEmpty1, shouldBeEmpty2, done := arty.GetUserResultsFromSites()
	<-done
	full := make([]string, 0, len(shouldBeFull))
	for len(shouldBeFull) != 0 {
		full = append(full, <-shouldBeFull)
	}
	empty1 := make([]string, 0, len(shouldBeEmpty1))
	for len(shouldBeEmpty1) != 0 {
		empty1 = append(full, <-shouldBeEmpty1)
	}
	empty2 := make([]string, 0, len(shouldBeEmpty2))
	for len(shouldBeEmpty2) != 0 {
		empty2 = append(empty2, <-shouldBeEmpty2)
	}
	assert.Contains(t, full, "Apple Discussions")
	assert.Equal(t, empty1, []string{})
	assert.Equal(t, empty2, []string{})

	// now test it again but with the names moved to likely
	user.LikelyUsernames = user.KnownUsernames
	user.KnownUsernames = []string{}
	shouldBeEmpty1, shouldBeFull, shouldBeEmpty2, done = arty.GetUserResultsFromSites()
	<-done
	full = make([]string, 0, len(shouldBeFull))
	for len(shouldBeFull) != 0 {
		full = append(full, <-shouldBeFull)
	}
	empty1 = make([]string, 0, len(shouldBeEmpty1))
	for len(shouldBeEmpty1) != 0 {
		empty1 = append(full, <-shouldBeEmpty1)
	}
	empty2 = make([]string, 0, len(shouldBeEmpty2))
	for len(shouldBeEmpty2) != 0 {
		empty2 = append(empty2, <-shouldBeEmpty2)
	}
	assert.Contains(t, full, "Apple Discussions")
	assert.Equal(t, empty1, []string{})
	assert.Equal(t, empty2, []string{})

	// and now with possible
	user.PossibleUsernames = user.LikelyUsernames
	user.LikelyUsernames = []string{}
	shouldBeEmpty1, shouldBeEmpty2, shouldBeFull, done = arty.GetUserResultsFromSites()
	<-done
	full = make([]string, 0, len(shouldBeFull))
	for len(shouldBeFull) != 0 {
		full = append(full, <-shouldBeFull)
	}
	empty1 = make([]string, 0, len(shouldBeEmpty1))
	for len(shouldBeEmpty1) != 0 {
		empty1 = append(full, <-shouldBeEmpty1)
	}
	empty2 = make([]string, 0, len(shouldBeEmpty2))
	for len(shouldBeEmpty2) != 0 {
		empty2 = append(empty2, <-shouldBeEmpty2)
	}
	assert.Contains(t, full, "Apple Discussions")
	assert.Equal(t, empty1, []string{})
	assert.Equal(t, empty2, []string{})

}

func TestMoriartyUsability(t *testing.T) {
	arty, err := NewMoriarty(jsonFile)
	if err != nil {
		t.Fatal(err)
	}
	arty.AssignNewUser(NewUserRecordings(
		"jonah",
		[]string{"Jonah{*}Wilmsmeyer", "wilm8026", "Pyro{*}Hedgehog"},
		[]string{"jpwilmsmeyer@gmail.com", ""},
		nil,
	))
	known, likely, possible := arty.TrackUserAcrossSites()
	fmt.Println("Known")
	fmt.Println(known)
	fmt.Println("Likely")
	fmt.Println(likely)
	fmt.Println("Possible")
	fmt.Println(possible)
	fmt.Println("done")

}
