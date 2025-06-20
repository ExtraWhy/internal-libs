package crwauth_test

import (
	"testing"

	crwauth "github.com/ExtraWhy/internal-libs/auth"
)

func Test_Auth(t *testing.T) {
	mtls_test := crwauth.MakeAuth(":8443")
	mtls_test.BeginAuth("mtls-test/server.crt", "mtls-test/server.key", "mtls-test/ca.crt")

}

//test with:
//cd mtls-test && curl --cert client.crt --key client.key --cacert ca.crt https://localhost:8443
