package tlsClient

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net/http"
)

const _RsaPrivateKey = "PRIVATE KEY"

func decodeCertFile(certFile string) (pemBlock []byte, err error) {
	bytes, err := ioutil.ReadFile(certFile)
	if err != nil {
		return
	}

	var (
		v         *pem.Block
		pemBlocks []*pem.Block
	)

	for {
		v, bytes = pem.Decode(bytes)
		if v == nil {
			break
		}

		if v.Type != _RsaPrivateKey {
			pemBlocks = append(pemBlocks, v)
		}
	}

	return pem.EncodeToMemory(pemBlocks[0]), nil
}

func decodeKeyFile(keyFile string, passPhrase *string) (pkey []byte, err error) {
	bytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return
	}
	var v *pem.Block

	for {
		v, bytes = pem.Decode(bytes)
		if v == nil {
			break
		}
		if v.Type == _RsaPrivateKey {
			if x509.IsEncryptedPEMBlock(v) {
				pkey, _ = x509.DecryptPEMBlock(v, []byte(*passPhrase))
				pkey = pem.EncodeToMemory(&pem.Block{
					Type:  v.Type,
					Bytes: pkey,
				})
			} else {
				pkey = pem.EncodeToMemory(v)
			}
		}
	}

	return
}

// New (certFile, keyFile string, passPhrase *string) (client http.Client, err error)
// Configure http client with TLS transport
func New(certFile, keyFile string, passPhrase *string) (client *http.Client, err error) {
	pemBlock, err := decodeCertFile(certFile)
	if err != nil {
		return
	}

	pkey, err := decodeKeyFile(keyFile, passPhrase)
	if err != nil {
		return
	}

	cert, _ := tls.X509KeyPair(pemBlock, pkey)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// RootCAs:      caCertPool,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client = &http.Client{Transport: transport}

	// conn, err := tls.Dial("tcp", "acqapi-test.tinkoff.ru:80", tlsConfig)
	// if err != nil {
	// 	log.Fatalf("client: dial: %s", err)
	// }
	// defer conn.Close()
	// log.Println("client: connected to: ", conn.RemoteAddr())
	// state := conn.ConnectionState()
	// for _, v := range state.PeerCertificates {
	// 	fmt.Println("Client: Server public key is:")
	// 	fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
	// }
	// log.Println("client: handshake: ", state.HandshakeComplete)
	// log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

	return
}
