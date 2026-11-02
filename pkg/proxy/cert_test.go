package proxy

import (
	"crypto/tls"
	"testing"
	"time"
)

func TestCertConfigTLSConfigUsesFallbackHostWhenSNIIsMissing(t *testing.T) {
	caCert, caKey, err := NewCA("Hetty", "Hetty CA", time.Hour)
	if err != nil {
		t.Fatalf("NewCA() error = %v", err)
	}

	certConfig, err := NewCertConfig(caCert, caKey)
	if err != nil {
		t.Fatalf("NewCertConfig() error = %v", err)
	}

	cert, err := certConfig.TLSConfig("example.com:443").GetCertificate(&tls.ClientHelloInfo{})
	if err != nil {
		t.Fatalf("GetCertificate() error = %v", err)
	}
	if cert == nil {
		t.Fatal("GetCertificate() returned nil certificate")
	}
	if len(cert.Certificate) == 0 {
		t.Fatal("GetCertificate() returned empty certificate chain")
	}
}

func TestCertConfigTLSConfigRequiresHostWhenSNIIsMissing(t *testing.T) {
	caCert, caKey, err := NewCA("Hetty", "Hetty CA", time.Hour)
	if err != nil {
		t.Fatalf("NewCA() error = %v", err)
	}

	certConfig, err := NewCertConfig(caCert, caKey)
	if err != nil {
		t.Fatalf("NewCertConfig() error = %v", err)
	}

	_, err = certConfig.TLSConfig().GetCertificate(&tls.ClientHelloInfo{})
	if err == nil {
		t.Fatal("GetCertificate() expected an error when SNI and fallback host are both missing")
	}
}
