# Solution

## Introduction
This code is based in several existing providers, basically gplus, yahoo and amazon where we have borrowed code and concepts.

## Oauth2 Server
To make it run you only have to:

``` bash
$ cd server
$ go run main.go
```

Any key-password combination must works.

## Demo Client

To compile and start the authentication-client server please run the following command.

``` bash
$ cd client
$ go run main.go
```

## Aps provider
It consists in the following methods and a pair of aux methods created for code organitazion and cleanless
purposes.

* In session.go
    * func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error)
* In aps.go
    * func New(clientKey, secret, callbackURL string, scopes ...string) *Provider
    * func (p *Provider) FetchUser(session goth.Session) (goth.User, error)
    * func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error)
    * func (p *Provider) RefreshTokenAvailable() bool
    * func (p *Provider) BeginAuth(state string) (goth.Session, error)
    * func (p *Provider) SetPrompt(prompt ...string)
* Aux Methods (both in aps.go)
	* func userFromReader(reader io.Reader, user *goth.User) error
	* func newConfig(provider *Provider, scopes []string) *oauth2.Config
* Other modifications
	* Some needed imports in aps.go and session.go
	* The path of aps import in main.go in order to match my tree.
	
## Tests
Some test are provided with the code just based in the same providers examples as above.
To execute the test just run the following command.

```bash
$ cd client
$ cd goth
$ cd aps
$ go test -v
```

##NOTES
* Godep have not been used in this repository, so it is possible you need to get some package.
* It have been tested both in Win10 using go 1.6.3 and also in Debian 8 using same go version. Both
64 bits OS.
