package sherlock

import "strings"

const WildcardPlaceholder string = "{*}"

var WildcardReplacements []string = []string{"_", "-", " ", ".", ""}

// The recordings of a single user.
type UserRecordings struct {
	KnownUsernames    []string
	KnownEmails       []string
	KnownSites        []string
	LikelyUsernames   []string
	LikelyEmails      []string
	LikelySites       []string
	PossibleUsernames []string
	PossibleEmails    []string
	PossibleSites     []string

	CheckingNSFW bool
}

// generate a user based on known items.
func NewUserRecordings(usernames, emails, sites []string) *UserRecordings {
	if emails == nil {
		emails = []string{}
	}
	if sites == nil {
		sites = []string{}
	}
	// for usernames, we check if there is a wildcard in there. If there is, we replace it with all equivalents.
	ur := &UserRecordings{
		KnownUsernames:    []string{},
		KnownEmails:       emails,
		KnownSites:        sites,
		LikelyUsernames:   []string{},
		LikelyEmails:      []string{},
		LikelySites:       []string{},
		PossibleUsernames: []string{},
		PossibleEmails:    []string{},
		PossibleSites:     []string{},
	}
	ur.AddKnownUsername(usernames...)
	return ur
}

// adds usernames, as well as smart generating associated usernames
func (ur *UserRecordings) AddKnownUsername(names ...string) {
	for _, name := range names {
		if strings.Contains(name, WildcardPlaceholder) {
			for _, replacement := range WildcardReplacements {
				likelyName := strings.Replace(name, WildcardPlaceholder, replacement, 1)
				ur.LikelyUsernames = append(ur.LikelyUsernames, likelyName)
				// we add the likely name to usernames to allow for multiple wildcards
				names = append(names, likelyName)
			}
		} else {
			ur.KnownUsernames = append(ur.KnownUsernames, name)
		}
	}
}
