package heart

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v1"
)

func TestNewOpenIDProvider(t *testing.T) {
	assert := assert.New(t)
	g := gin.New()
	g.StaticFile(".well-known/openid-configuration", "fixtures/op_config.json")
	server := httptest.NewServer(g)
	defer server.Close()
	op, err := NewOpenIDProvider(server.URL)
	assert.NoError(err)
	assert.Equal("http://localhost:8080/openid-connect-server-webapp/jwk", op.Config.JWKSURI)
}

func TestFetchKey(t *testing.T) {
	assert := assert.New(t)
	g := gin.New()
	g.StaticFile("jwks", "fixtures/jw_key_set.json")
	server := httptest.NewServer(g)
	defer server.Close()
	config := OPConfig{JWKSURI: server.URL + "/jwks"}
	op := OpenIDProvider{Config: config}
	err := op.FetchKey()
	assert.NoError(err)
	assert.Equal("1471357142", op.Keys.Keys[0].KeyID)
}

func TestUnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
	tokenResp, err := os.Open("fixtures/oid_token_response.json")
	assert.NoError(err)
	defer tokenResp.Close()
	tokenBytes, err := ioutil.ReadAll(tokenResp)
	assert.NoError(err)
	token := OpenIDTokenResponse{}
	now := time.Now()
	json.Unmarshal(tokenBytes, &token)
	assert.NotEmpty(token.Expiration)
	assert.InDelta(now.Unix()+600, token.Expiration.Unix(), 3.0)
}

func TestUserInfo(t *testing.T) {
	assert := assert.New(t)
	accessToken := "random_characters"
	g := gin.New()
	g.Use(func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		expectedAuthHeader := fmt.Sprintf("Bearer %s", accessToken)
		assert.Equal(expectedAuthHeader, authHeader)
	})
	g.StaticFile("user_info", "fixtures/user_info.json")

	server := httptest.NewServer(g)
	defer server.Close()
	config := OPConfig{UserInfoEndpoint: server.URL + "/user_info"}
	op := OpenIDProvider{Config: config}
	ui, err := op.UserInfo(accessToken)
	assert.NoError(err)
	assert.Equal("foo@bar", ui.SUB)
}

func TestExchange(t *testing.T) {
	assert := assert.New(t)
	code := "random_characters"
	g := gin.New()
	g.POST("/token", func(c *gin.Context) {
		providedCode := c.PostForm("code")
		assert.Equal(code, providedCode)
		c.File("fixtures/oid_token_response.json")
	})

	server := httptest.NewServer(g)
	defer server.Close()
	config := OPConfig{TokenEndpoint: server.URL + "/token"}
	op := OpenIDProvider{Config: config}
	resp, err := op.Exchange(code, setupClient())
	assert.NoError(err)
	assert.Equal("random_characters", resp.AccessToken)
}

func TestFailedExchange(t *testing.T) {
	assert := assert.New(t)
	code := "random_characters"
	g := gin.New()
	g.POST("/token", func(c *gin.Context) {
		providedCode := c.PostForm("code")
		assert.Equal(code, providedCode)
		c.String(http.StatusInternalServerError, "I'm sad")
	})

	server := httptest.NewServer(g)
	defer server.Close()
	config := OPConfig{TokenEndpoint: server.URL + "/token"}
	op := OpenIDProvider{Config: config}
	_, err := op.Exchange(code, setupClient())
	assert.Error(err)
	assert.True(errors.IsUnauthorized(err))
}

func TestVerify(t *testing.T) {
	assert := assert.New(t)
	jwt := NewClientJWT("my_client_id", "http://server.com")
	jwkJSON, err := os.Open("fixtures/verify_private_key.json")
	assert.NoError(err)
	defer jwkJSON.Close()
	jwkBytes, err := ioutil.ReadAll(jwkJSON)
	assert.NoError(err)
	jwk := jose.JsonWebKey{}
	json.Unmarshal(jwkBytes, &jwk)
	signedBlob, err := SignJWT(jwt, jwk)
	assert.NoError(err)
	assert.NotEmpty(signedBlob)
	jwksJSON, err := os.Open("fixtures/verify_pk_set.json")
	assert.NoError(err)
	defer jwksJSON.Close()
	jwksBytes, err := ioutil.ReadAll(jwksJSON)
	assert.NoError(err)
	jwks := jose.JsonWebKeySet{}
	err = json.Unmarshal(jwksBytes, &jwks)
	assert.NoError(err)
	op := &OpenIDProvider{Keys: jwks}
	valid, err := op.Validate(signedBlob)
	assert.NoError(err)
	assert.True(valid)
}

func setupClient() Client {
	jwkJSON, _ := os.Open("fixtures/jwk.json")
	jwkBytes, _ := ioutil.ReadAll(jwkJSON)
	jwk := jose.JsonWebKey{}
	json.Unmarshal(jwkBytes, &jwk)
	return Client{
		ISS:         "simple",
		AUD:         "http://localhost:8080/openid-connect-server-webapp/",
		RedirectURI: "http://localhost:3333/redirect",
		PrivateKey:  jwk,
	}
}
