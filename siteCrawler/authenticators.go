package siteCrawler

import (
	"strings"

	"github.com/carlmjohnson/requests"
)

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
