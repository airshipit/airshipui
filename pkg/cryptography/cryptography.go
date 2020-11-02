/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cryptography

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"time"

	"opendev.org/airship/airshipui/pkg/log"
)

const (
	// keypair details
	keySize        = 4096 // 4k key
	privateKeyType = "RSA PRIVATE KEY"
	publicKeyType  = "CERTIFICATE"

	// certificate request details
	cn = "localhost"  // common name
	o  = "Airship UI" // organization
)

// GeneratePrivateKey will a pem encoded private key and an rsa private key object
func GeneratePrivateKey() ([]byte, *rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		log.Error("Problem generating private key", err)
		return nil, nil, err
	}

	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  privateKeyType,
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		log.Error("Problem generating private key pem", err)
		return nil, nil, err
	}

	return buf.Bytes(), privateKey, nil
}

// GeneratePublicKey will create a pem encoded cert
func GeneratePublicKey(privateKey *rsa.PrivateKey) ([]byte, error) {
	template := generateCSR()
	derCert, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  publicKeyType,
		Bytes: derCert,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// generateCSR creates the base information needed to create the certificate
func generateCSR() x509.Certificate {
	return x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{o},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
}

// TestCertValidity will check if the cert defined in the conf is not past its not after date
func TestCertValidity(pemFile string) error {
	r, err := ioutil.ReadFile(pemFile)
	if err != nil {
		log.Error(err)
		return err
	}

	block, _ := pem.Decode(r)
	_, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Error(err)
		return err
	}

	// calculate the validity of the cert
	// TODO: Add a cert check for time based validity here
	// fmt.Println(cert.NotAfter)
	return nil
}
