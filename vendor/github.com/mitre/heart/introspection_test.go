package heart

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v1"
)

func mockIntrospectionEndpoint(assert *assert.Assertions, jwk jose.JsonWebKey) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("POST", r.Method)
		clientID := r.FormValue("client_id")
		assert.Equal("my_client_id", clientID)
		token := r.FormValue("token")
		assert.Equal("foo", token)
		signedBlob := r.FormValue("client_assertion")
		assert.NotEmpty(signedBlob)
		jws, err := jose.ParseSigned(signedBlob)
		assert.NoError(err)
		rsaPrivateKey := jwk.Key.(*rsa.PrivateKey)
		_, err = jws.Verify(&rsaPrivateKey.PublicKey)
		assert.NoError(err)
		ir := &IntrospectionResponse{Active: true, Scope: "foo bar", EXP: time.Now().Unix(), SUB: "steve", ClientID: "heart-watch"}
		encoder := json.NewEncoder(w)
		encoder.Encode(ir)
	})
}

func introspectionSetUp(assert *assert.Assertions) (string, *httptest.Server) {
	jwt := NewClientJWT("my_client_id", "http://server.com")
	jwkJSON, err := os.Open("fixtures/jwk.json")
	assert.NoError(err)
	defer jwkJSON.Close()
	jwkBytes, err := ioutil.ReadAll(jwkJSON)
	assert.NoError(err)
	jwk := jose.JsonWebKey{}
	json.Unmarshal(jwkBytes, &jwk)
	signedBlob, err := SignJWT(jwt, jwk)
	assert.NoError(err)
	assert.NotEmpty(signedBlob)
	server := httptest.NewServer(mockIntrospectionEndpoint(assert, jwk))
	return signedBlob, server
}

func TestIntrospectToken(t *testing.T) {
	assert := assert.New(t)
	signedBlob, server := introspectionSetUp(assert)
	defer server.Close()
	ir, err := IntrospectToken(server.URL, "foo", "my_client_id", signedBlob)
	assert.NoError(err)
	assert.True(ir.Active)
	assert.Equal("foo bar", ir.Scope)
}
