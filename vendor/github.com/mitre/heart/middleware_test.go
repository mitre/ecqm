package heart

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gopkg.in/square/go-jose.v1"
)

type MiddlewareSuite struct {
	suite.Suite
	PrivateKey jose.JsonWebKey
}

func (suite *MiddlewareSuite) SetupTest() {
	jwkJSON, err := os.Open("fixtures/jwk.json")
	suite.NoError(err)
	defer jwkJSON.Close()
	jwkBytes, err := ioutil.ReadAll(jwkJSON)
	suite.NoError(err)
	jwk := jose.JsonWebKey{}
	json.Unmarshal(jwkBytes, &jwk)
	suite.PrivateKey = jwk
}

func (suite *MiddlewareSuite) TestNoAuthHeader() {
	handler := OAuthIntrospectionHandler("", "", "", suite.PrivateKey)
	g := gin.New()
	g.GET("/", handler)
	server := httptest.NewServer(g)
	defer server.Close()
	resp, err := http.Get(server.URL)
	suite.NoError(err)
	defer resp.Body.Close()
	suite.Equal(http.StatusForbidden, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	suite.NoError(err)
	suite.Equal("No Authorization header provided", string(bodyBytes))
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}
