package config

import (
	"fmt"
	"os"
)

// LoadCertificatePEM returns PEM bytes from inline config or from a file path.
func LoadCertificatePEM(app App) ([]byte, error) {
	if app.CertificateKeyPath != "" {
		b, err := os.ReadFile(app.CertificateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("read certificate_key_path: %w", err)
		}
		return b, nil
	}
	if app.CertificateKey != "" {
		return []byte(app.CertificateKey), nil
	}
	return nil, fmt.Errorf("set app.certificate_key or app.certificate_key_path")
}
