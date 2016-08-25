package heart

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/square/go-jose.v1"

	"github.com/stretchr/testify/assert"
)

func TestNewClientJWT(t *testing.T) {
	assert := assert.New(t)
	jwt := NewClientJWT("my_client_id", "http://server.com")
	assert.Equal(jwt.ISS, "my_client_id")
	assert.Equal(jwt.SUB, "my_client_id")
	assert.Equal(jwt.IAT+60, jwt.EXP)
	assert.Equal(len(jwt.JTI), 50)
}

func TestSignJWT(t *testing.T) {
	assert := assert.New(t)
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
	jws, err := jose.ParseSigned(signedBlob)
	assert.NoError(err)
	rsaPrivateKey := jwk.Key.(*rsa.PrivateKey)
	_, err = jws.Verify(&rsaPrivateKey.PublicKey)
	assert.NoError(err)
}
