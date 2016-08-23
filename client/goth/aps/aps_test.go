// aps_test.go
package aps_test

import (
	"fmt"
	"testing"

	"github.com/markbates/goth"
	"github.com/silviofdez/oauth/client/goth/aps"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := provider()

	a.Equal(p.ClientKey, "bawdy-reindeers-14-56dd2bcc2ba94")
	a.Equal(p.Secret, "6454acedc7024fdfa743c5407da7ad44")
	a.Equal(p.CallbackURL, "http://localhost:3000/auth/aps/callback")
}

func Test_Implements_Provider(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	a.Implements((*goth.Provider)(nil), provider())
}

func Test_BeginAuth(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := provider()
	session, err := p.BeginAuth("test_state")
	s := session.(*aps.Session)
	a.NoError(err)
	a.Contains(s.AuthURL, "http://localhost:9096/authorize")
	a.Contains(s.AuthURL, fmt.Sprintf("client_id=%s", "bawdy-reindeers-14-56dd2bcc2ba94"))
	a.Contains(s.AuthURL, "state=test_state")
}

func Test_BeginAuthWithPrompt(t *testing.T) {
	// This exists because there was a panic caused by the oauth2 package when
	// the AuthCodeOption passed was nil. This test uses it, Test_BeginAuth does
	// not, to ensure both cases are covered.
	t.Parallel()
	a := assert.New(t)

	p := provider()
	p.SetPrompt("test", "prompts")
	session, err := p.BeginAuth("test_state")
	s := session.(*aps.Session)
	a.NoError(err)
	a.Contains(s.AuthURL, "http://localhost:9096/authorize")
	a.Contains(s.AuthURL, fmt.Sprintf("client_id=%s", "bawdy-reindeers-14-56dd2bcc2ba94"))
	a.Contains(s.AuthURL, "state=test_state")
	a.Contains(s.AuthURL, "prompt=test+prompts")
}

func Test_SessionFromJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	p := provider()
	session, err := p.UnmarshalSession(`{"AuthURL":"http://localhost:9096/authorize","AccessToken":"1234567890"}`)
	a.NoError(err)

	s := session.(*aps.Session)
	a.Equal(s.AuthURL, "http://localhost:9096/authorize")
	a.Equal(s.AccessToken, "1234567890")
}

func provider() *aps.Provider {
	return aps.New("bawdy-reindeers-14-56dd2bcc2ba94", "6454acedc7024fdfa743c5407da7ad44", "http://localhost:3000/auth/aps/callback")
}
