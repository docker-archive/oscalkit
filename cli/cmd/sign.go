package cmd

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"gopkg.in/square/go-jose.v2"

	"github.com/urfave/cli"
)

var privKey string
var alg string

// Sign ...
var Sign = cli.Command{
	Name:      "sign",
	Usage:     "sign OSCAL JSON artifacts",
	ArgsUsage: "[files...]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "key, k",
			Usage:       "private key file for signing. Must be in PEM or DER formats. Supports RSA/EC keys and X.509 certificats with embedded RSA/EC keys",
			Destination: &privKey,
		},
		cli.StringFlag{
			Name:        "alg, a",
			Usage:       "algorithm for signing. Supports RSASSA-PKCS#1v1.5, RSASSA-PSS, HMAC, ECDSA and Ed25519",
			Destination: &alg,
		},
	},
	Before: func(c *cli.Context) error {
		if privKey == "" {
			return cli.NewExitError("oscalkit sign is missing the --key flag", 1)
		}

		if alg == "" {
			return cli.NewExitError("oscalkit sign is missing the --alg flag", 1)
		}

		if c.NArg() < 1 {
			return cli.NewExitError("oscalkit sign requires at least one argument", 1)
		}

		return nil
	},
	Action: func(c *cli.Context) error {
		for _, srcFile := range c.Args() {
			privKeyFile, err := ioutil.ReadFile(privKey)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Error loading private key file %s: %s", privKey, err), 1)
			}

			srcFileData, err := ioutil.ReadFile(srcFile)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Error reading source file %s: %s", srcFile, err), 1)
			}

			input := privKeyFile
			block, _ := pem.Decode(privKeyFile)
			if block != nil {
				input = block.Bytes
			}

			var priv interface{}
			var msg string
			priv, err0 := x509.ParsePKCS1PrivateKey(input)
			if err0 == nil {
				msg, err = sign(priv, srcFileData)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Signing error: %s", err), 1)
				}
			}

			priv, err1 := x509.ParsePKCS8PrivateKey(input)
			if err1 == nil {
				msg, err = sign(priv, srcFileData)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Signing error: %s", err), 1)
				}
			}

			priv, err2 := x509.ParseECPrivateKey(input)
			if err2 == nil {
				msg, err = sign(priv, srcFileData)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Signing error: %s", err), 1)
				}
			}

			if msg != "" {
				splitPath := strings.Split(path.Base(srcFile), ".")
				filePath := fmt.Sprintf("%s-SIGNED.%s", splitPath[0], splitPath[1])
				if err := ioutil.WriteFile(filePath, []byte(msg), 0644); err != nil {
					return cli.NewExitError(fmt.Sprintf("Error writing signed file: %s", err), 1)
				}

				continue
			}

			return cli.NewExitError(fmt.Sprintf("Error parsing private key: %s, %s, %s", err0, err1, err2), 1)
		}

		return nil
	},
}

func sign(key interface{}, payload []byte) (string, error) {
	sigAlg := jose.SignatureAlgorithm(alg)
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: sigAlg, Key: key}, nil)
	if err != nil {
		return "", err
	}

	obj, err := signer.Sign(payload)
	if err != nil {
		return "", err
	}

	return obj.FullSerialize(), nil
}
