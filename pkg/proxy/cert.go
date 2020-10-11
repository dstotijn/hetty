package proxy

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// MaxSerialNumber is the upper boundary that is used to create unique serial
// numbers for the certificate. This can be any unsigned integer up to 20
// bytes (2^(8*20)-1).
var MaxSerialNumber = big.NewInt(0).SetBytes(bytes.Repeat([]byte{255}, 20))

// CertConfig is a set of configuration values that are used to build TLS configs
// capable of MITM
type CertConfig struct {
	ca     *x509.Certificate
	caPriv crypto.PrivateKey
	priv   *rsa.PrivateKey
	keyID  []byte
}

// NewCertConfig creates a MITM config using the CA certificate and
// private key to generate on-the-fly certificates.
func NewCertConfig(ca *x509.Certificate, caPrivKey crypto.PrivateKey) (*CertConfig, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	pub := priv.Public()

	// Subject Key Identifier support for end entity certificate.
	// https://www.ietf.org/rfc/rfc3280.txt (section 4.2.1.2)
	pkixPubKey, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	h := sha1.New()
	h.Write(pkixPubKey)
	keyID := h.Sum(nil)

	return &CertConfig{
		ca:     ca,
		caPriv: caPrivKey,
		priv:   priv,
		keyID:  keyID,
	}, nil
}

// LoadOrCreateCA loads an existing CA key pair from disk, or creates
// a new keypair and saves to disk if certificate or key files don't exist.
func LoadOrCreateCA(caKeyFile, caCertFile string) (*x509.Certificate, *rsa.PrivateKey, error) {
	tlsCA, err := tls.LoadX509KeyPair(caCertFile, caKeyFile)
	if err == nil {
		caCert, err := x509.ParseCertificate(tlsCA.Certificate[0])
		if err != nil {
			return nil, nil, fmt.Errorf("proxy: could not parse CA: %v", err)
		}
		caKey, ok := tlsCA.PrivateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, nil, errors.New("proxy: private key is not RSA")
		}
		return caCert, caKey, nil
	}
	if !os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("proxy: could not load CA key pair: %v", err)
	}

	// Create directories for files if they don't exist yet.
	keyDir, _ := filepath.Split(caKeyFile)
	if keyDir != "" {
		if _, err := os.Stat(keyDir); os.IsNotExist(err) {
			os.MkdirAll(keyDir, 0755)
		}
	}
	keyDir, _ = filepath.Split(caCertFile)
	if keyDir != "" {
		if _, err := os.Stat("keyDir"); os.IsNotExist(err) {
			os.MkdirAll(keyDir, 0755)
		}
	}

	// Create new CA keypair.
	caCert, caKey, err := NewCA("Hetty", "Hetty CA", time.Duration(365*24*time.Hour))
	if err != nil {
		return nil, nil, fmt.Errorf("proxy: could not generate new CA keypair: %v", err)
	}

	// Open CA certificate and key files for writing.
	certOut, err := os.Create(caCertFile)
	if err != nil {
		return nil, nil, fmt.Errorf("proxy: could not open cert file for writing: %v", err)
	}
	keyOut, err := os.OpenFile(caKeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, nil, fmt.Errorf("proxy: could not open key file for writing: %v", err)
	}

	// Write PEM blocks to CA certificate and key files.
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: caCert.Raw}); err != nil {
		return nil, nil, fmt.Errorf("proxy: could not write CA certificate to disk: %v", err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("proxy: could not convert private key to DER format: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return nil, nil, fmt.Errorf("proxy: could not write CA key to disk: %v", err)
	}

	return caCert, caKey, nil
}

// NewCA creates a new CA certificate and associated private key.
func NewCA(name, organization string, validity time.Duration) (*x509.Certificate, *rsa.PrivateKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	pub := priv.Public()

	// Subject Key Identifier support for end entity certificate.
	// https://www.ietf.org/rfc/rfc3280.txt (section 4.2.1.2)
	pkixpub, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, nil, err
	}
	h := sha1.New()
	h.Write(pkixpub)
	keyID := h.Sum(nil)

	// TODO: keep a map of used serial numbers to avoid potentially reusing a
	// serial multiple times.
	serial, err := rand.Int(rand.Reader, MaxSerialNumber)
	if err != nil {
		return nil, nil, err
	}

	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   name,
			Organization: []string{organization},
		},
		SubjectKeyId:          keyID,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(validity),
		DNSNames:              []string{name},
		IsCA:                  true,
	}

	raw, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, pub, priv)
	if err != nil {
		return nil, nil, err
	}

	// Parse certificate bytes so that we have a leaf certificate.
	x509c, err := x509.ParseCertificate(raw)
	if err != nil {
		return nil, nil, err
	}

	return x509c, priv, nil
}

// TLSConfig returns a *tls.Config that will generate certificates on-the-fly using
// the SNI extension in the TLS ClientHello.
func (c *CertConfig) TLSConfig() *tls.Config {
	return &tls.Config{
		GetCertificate: func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			if clientHello.ServerName == "" {
				return nil, errors.New("missing server name (SNI)")
			}
			return c.cert(clientHello.ServerName)
		},
		NextProtos: []string{"http/1.1"},
	}
}

func (c *CertConfig) cert(hostname string) (*tls.Certificate, error) {
	// Remove the port if it exists.
	host, _, err := net.SplitHostPort(hostname)
	if err == nil {
		hostname = host
	}

	serial, err := rand.Int(rand.Reader, MaxSerialNumber)
	if err != nil {
		return nil, err
	}

	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   hostname,
			Organization: []string{"Hetty"},
		},
		SubjectKeyId:          c.keyID,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
	}

	if ip := net.ParseIP(hostname); ip != nil {
		tmpl.IPAddresses = []net.IP{ip}
	} else {
		tmpl.DNSNames = []string{hostname}
	}

	raw, err := x509.CreateCertificate(rand.Reader, tmpl, c.ca, c.priv.Public(), c.caPriv)
	if err != nil {
		return nil, err
	}

	// Parse certificate bytes so that we have a leaf certificate.
	x509c, err := x509.ParseCertificate(raw)
	if err != nil {
		return nil, err
	}

	return &tls.Certificate{
		Certificate: [][]byte{raw, c.ca.Raw},
		PrivateKey:  c.priv,
		Leaf:        x509c,
	}, nil
}
