package auth

import (
	"crypto/rsa"
	"os"
	"path/filepath"
	"strings"
	"testing"

	jwt "github.com/golang-jwt/jwt/v4"
)

func mustKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	root := moduleRoot(t)
	b, err := os.ReadFile(filepath.Join(root, "testdata", "dev_rsa_private.pem"))
	if err != nil {
		t.Fatal(err)
	}
	k, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	wd, _ := os.Getwd()
	dir := wd
	for i := 0; i < 15; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("root")
	return ""
}

func TestGenerateAndParseJWT(t *testing.T) {
	key := mustKey(t)
	tok, err := GenerateJWTToken(key, "user-1")
	if err != nil || tok.Token == "" {
		t.Fatal(err)
	}
	claims, ok, err := ParseAndValidateJWTtoken(key, tok.Token)
	if err != nil || !ok || claims.UserUUID != "user-1" {
		t.Fatalf("%v %v %+v", ok, err, claims)
	}
}

func TestParseJWT_Invalid(t *testing.T) {
	key := mustKey(t)
	_, _, err := ParseAndValidateJWTtoken(key, "not-a-jwt")
	if err == nil {
		t.Fatal()
	}
}

func TestParseJWT_Tampered(t *testing.T) {
	key := mustKey(t)
	tok, err := GenerateJWTToken(key, "user-1")
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(tok.Token, ".")
	if len(parts) < 3 {
		t.Fatal()
	}
	bad := parts[0] + "." + parts[1] + ".xxx"
	_, ok, err := ParseAndValidateJWTtoken(key, bad)
	if err == nil && ok {
		t.Fatal("expected failure")
	}
}

func TestParseJWT_WrongSigningMethod(t *testing.T) {
	key := mustKey(t)
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "x"})
	s, err := tkn.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = ParseAndValidateJWTtoken(key, s)
	if err == nil {
		t.Fatal("expected error")
	}
}
