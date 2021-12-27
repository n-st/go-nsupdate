package nsupdate // import "golang.voidptr.de/nsupdate"

import (
	"io/ioutil"
	"regexp"
	"strings"
)

type TSigOptions struct {
	Name      string
	Algorithm string
	Secret    string
}

var keyfileKeyBlock =
// ["key", "<key name>", "<key attributes>"]
regexp.MustCompile(`\b(key)\s+"([^"]+)"\s+{([^}]*)};`)

var keyfileAlgorithmLine =
// ["algorithm", "<algorithm name (non-FQDN form)>"]
regexp.MustCompile(`\b(algorithm)\s+"?([^";]+)"?;`)

var keyfileSecretLine =
// ["secret", "<base64-encoded secret>"]
regexp.MustCompile(`\b(secret)\s+"?([^";]+)"?;`)

func ReadKeyFile(filename string) (result TSigOptions, err error) {
	data, readerr := ioutil.ReadFile(filename)
	if readerr != nil {
		err = readerr
		return
	}
	fileContent := string(data)

	keyBlock := keyfileKeyBlock.FindStringSubmatch(fileContent)
	if len(keyBlock) > 0 {
		result.Name = keyBlock[2]
		keyAttributes := keyBlock[3]

		// The file format specification doesn't say what happens if multiple
		// "algorithm" or "secret" attributes are specified. We'll just take
		// the first one in that case.

		algorithmLine := keyfileAlgorithmLine.FindStringSubmatch(keyAttributes)
		if len(algorithmLine) > 0 {
			result.Algorithm = algorithmLine[2]
		}

		secretLine := keyfileSecretLine.FindStringSubmatch(keyAttributes)
		if len(secretLine) > 0 {
			result.Secret = secretLine[2]
		}
	}

	if !strings.HasSuffix(result.Name, ".") {
		result.Name = result.Name + "."
	}

	if !strings.HasSuffix(result.Algorithm, ".") {
		result.Algorithm = result.Algorithm + "."
	}

	return
}
