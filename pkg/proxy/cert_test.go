package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"net"
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

func TestCertConfigCertIncludesIPSANForHostIP(t *testing.T) {
	caCert, caKey, err := NewCA("Hetty", "Hetty CA", time.Hour)
	if err != nil {
		t.Fatalf("NewCA() error = %v", err)
	}

	certConfig, err := NewCertConfig(caCert, caKey)
	if err != nil {
		t.Fatalf("NewCertConfig() error = %v", err)
	}

	cert, err := certConfig.TLSConfig("158.191.147.84:443").GetCertificate(&tls.ClientHelloInfo{})
	if err != nil {
		t.Fatalf("GetCertificate() error = %v", err)
	}

	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatalf("ParseCertificate() error = %v", err)
	}
	if len(leaf.IPAddresses) == 0 {
		t.Fatal("expected IP SANs for IP hosts")
	}
	if !leaf.IPAddresses[0].Equal(net.ParseIP("158.191.147.84")) {
		t.Fatalf("unexpected IP SAN %v", leaf.IPAddresses[0])
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
