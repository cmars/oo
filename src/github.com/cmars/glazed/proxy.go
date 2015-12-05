package glazed

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"gopkg.in/errgo.v1"
	"gopkg.in/macaroon-bakery.v1/bakery"
	"gopkg.in/macaroon-bakery.v1/bakery/checkers"
	"gopkg.in/macaroon-bakery.v1/httpbakery"
	"gopkg.in/macaroon.v1"
)

type Proxy struct {
	target    *url.URL
	condition string
	bakery    *bakery.Service

	mu       sync.Mutex
	passthru map[string]bool
}

func newIDBytes() ([]byte, error) {
	buf := make([]byte, 32)
	_, err := rand.Reader.Read(buf)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return buf, nil
}

func newID() (string, error) {
	buf, err := newIDBytes()
	if err != nil {
		return "", errgo.Mask(err)
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

func NewProxy(target string) (*Proxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	p := &Proxy{target: targetURL, passthru: map[string]bool{}}
	id, err := newID()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	p.bakery, err = bakery.NewService(bakery.NewServiceParams{
		Location: id,
	})
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return p, nil
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

type clientRoundTripper struct {
	*http.Client
}

func (rt clientRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	copyReq := *req
	copyReq.RequestURI = ""
	return rt.Client.Do(&copyReq)
}

func (p *Proxy) isAuthenticated(req *http.Request) bool {
	cookie, err := req.Cookie(p.bakery.Location())
	if err == http.ErrNoCookie {
		return false
	}
	msbuf, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Println("failed to base64-decode cookie: %v", err)
		return false
	}
	var ms macaroon.Slice
	err = ms.UnmarshalBinary(msbuf)
	if err != nil {
		log.Println("failed to unmarshal cookie: %v", err)
		return false
	}

	err = p.bakery.Check(ms, bakery.FirstPartyCheckerFunc(func(caveat string) error {
		return checkers.ErrCaveatNotRecognized
	}))
	if err != nil {
		log.Println("verification failed")
		return false
	}
	return true
}

func (p *Proxy) authenticate(req *http.Request) error {
	m, err := p.bakery.NewMacaroon("", nil, []checkers.Caveat{{
		Location: p.bakery.Location(), Condition: p.condition,
	}})
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	client := httpbakery.NewClient()
	ms, err := client.DischargeAll(m)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	err = p.bakery.Check(ms, bakery.FirstPartyCheckerFunc(func(caveat string) error {
		return checkers.ErrCaveatNotRecognized
	}))
	if err != nil {
		return p.handleError(req, err)
	}

	buf, err := ms.MarshalBinary()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	cookie := &http.Cookie{
		Name:  p.bakery.Location(),
		Value: base64.StdEncoding.EncodeToString(buf),
	}
	req.AddCookie(cookie)

	return nil
}
