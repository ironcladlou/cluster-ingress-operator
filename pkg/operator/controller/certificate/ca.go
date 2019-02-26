package certificate

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	// caCertSecretName is the name of the secret that holds the CA certificate
	// that the operator will use to create default certificates for
	// clusteringresses.
	caCertSecretName = "router-ca"
)

// CASecretName returns the namespaced name for the router CA secret.
func CASecretName(operatorNamespace string) types.NamespacedName {
	return types.NamespacedName{
		Namespace: operatorNamespace,
		Name:      caCertSecretName,
	}
}

func (r *reconciler) ensureRouterCASecret() (*corev1.Secret, error) {
	current, err := r.currentRouterCASecret()
	if err != nil {
		return nil, err
	}
	if current != nil {
		return current, nil
	}
	desired, err := desiredRouterCASecret(r.operatorNamespace)
	if err != nil {
		return nil, err
	}
	if err := r.createRouterCASecret(desired); err != nil {
		return nil, fmt.Errorf("failed to create CA secret: %v", err)
	}
	r.recorder.Event(desired, "Normal", "CreatedWildcardCACert", "Created a default wildcard CA certificate")
	return desired, nil
}

// currentRouterCASecret returns the current router CA secret.
func (r *reconciler) currentRouterCASecret() (*corev1.Secret, error) {
	name := CASecretName(r.operatorNamespace)
	secret := &corev1.Secret{}
	if err := r.client.Get(context.TODO(), name, secret); err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return secret, nil
}

// generateRouterCA generates and returns a CA certificate and key.
func generateRouterCA() ([]byte, []byte, error) {
	signerName := fmt.Sprintf("%s@%d", "cluster-ingress-operator", time.Now().Unix())

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key: %v", err)
	}

	root := &x509.Certificate{
		Subject: pkix.Name{CommonName: signerName},

		SignatureAlgorithm: x509.SHA256WithRSA,

		NotBefore:    time.Now().Add(-1 * time.Second),
		NotAfter:     time.Now().Add(2 * 365 * 24 * time.Hour),
		SerialNumber: big.NewInt(1),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,

		IsCA: true,

		// Don't allow the CA to be used to make another CA.
		MaxPathLen:     0,
		MaxPathLenZero: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, root, root, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %v", err)
	}

	certs, err := x509.ParseCertificates(derBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	if len(certs) != 1 {
		return nil, nil, fmt.Errorf("expected a single certificate")
	}

	certBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certs[0].Raw,
	})

	keyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return certBytes, keyBytes, nil
}

// desiredRouterCASecret returns the desired router CA secret.
func desiredRouterCASecret(namespace string) (*corev1.Secret, error) {
	certBytes, keyBytes, err := generateRouterCA()
	if err != nil {
		return nil, fmt.Errorf("failed to generate certificate: %v", err)
	}

	name := CASecretName(namespace)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
		Data: map[string][]byte{
			"tls.crt": certBytes,
			"tls.key": keyBytes,
		},
		Type: corev1.SecretTypeTLS,
	}
	return secret, nil
}

// createRouterCASecret creates the router CA secret.
func (r *reconciler) createRouterCASecret(secret *corev1.Secret) error {
	if err := r.client.Create(context.TODO(), secret); err != nil {
		return err
	}
	log.Info("created secret", "namespace", secret.Namespace, "name", secret.Name)
	return nil
}
