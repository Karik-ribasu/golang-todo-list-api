package config

import (
	"os"
	"path/filepath"
	"testing"

	viper "github.com/spf13/viper"
)

func TestReadConfig_MissingFile(t *testing.T) {
	dir := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	_, err = ReadConfig()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReadConfig_OK(t *testing.T) {
	root := findModuleRoot(t)
	pemPath := filepath.Join(root, "testdata", "dev_rsa_private.pem")
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	content := "[db]\nuser = \"a\"\npasswd = \"b\"\naddr = \"c\"\nport = \"3306\"\nname = \"d\"\n\n[app]\ncertificate_key_path = \"" + pemPath + "\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	cfg, err := ReadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Db.User != "a" || cfg.App.CertificateKeyPath != pemPath {
		t.Fatalf("%+v", cfg)
	}
}

func TestReadConfig_UnmarshalError(t *testing.T) {
	v := viper.New()
	v.Set("db", 12345)
	_, err := readMergedConfig(v)
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
}

func findModuleRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := wd
	for i := 0; i < 15; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("go.mod not found from %s", wd)
	return ""
}
