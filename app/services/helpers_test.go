package appServices

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
	jwt "github.com/golang-jwt/jwt/v4"
)

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

func rsaConfig(t *testing.T) config.Config {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(moduleRoot(t), "testdata", "dev_rsa_private.pem"))
	if err != nil {
		t.Fatal(err)
	}
	k, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		t.Fatal(err)
	}
	return config.Config{App: config.App{PrivateKey: k}}
}
