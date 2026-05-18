package config

import (
	"path/filepath"
	"testing"
)

func TestLoadCertificatePEM_Inline(t *testing.T) {
	b, err := LoadCertificatePEM(App{CertificateKey: "PEM"})
	if err != nil || string(b) != "PEM" {
		t.Fatal()
	}
}

func TestLoadCertificatePEM_File(t *testing.T) {
	root := findModuleRoot(t)
	p := filepath.Join(root, "testdata", "dev_rsa_private.pem")
	b, err := LoadCertificatePEM(App{CertificateKeyPath: p})
	if err != nil || len(b) < 50 {
		t.Fatal(err)
	}
}

func TestLoadCertificatePEM_MissingPath(t *testing.T) {
	_, err := LoadCertificatePEM(App{CertificateKeyPath: "/nonexistent/pem"})
	if err == nil {
		t.Fatal()
	}
}

func TestLoadCertificatePEM_NoSource(t *testing.T) {
	_, err := LoadCertificatePEM(App{})
	if err == nil {
		t.Fatal()
	}
}
