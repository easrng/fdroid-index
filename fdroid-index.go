package main

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"encoding/json"

	"github.com/avast/apkparser"
	"github.com/avast/apkverifier"
	"github.com/itchyny/gojq"
	"github.com/pelletier/go-toml/v2"
)

//go:embed transform.jq
var transform string

func main() {
	query, _ := gojq.Parse(transform)

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <fdroid repo url including fingerprint>", os.Args[0])
		os.Exit(1)
	}

	repo := os.Args[1]
	repoURL, err := url.Parse(repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse repo url: %s\n", err.Error())
		os.Exit(1)
	}
	repoBase := repoURL.Scheme + ":/" + "/" + repoURL.Host + repoURL.Path
	repoFingerprint := repoURL.Query()["fingerprint"][0]

	resp, err := http.Get(repoBase + "/index-v1.jar")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get repo index: %s\n", err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read repo index: %s\n", err.Error())
		os.Exit(1)
	}
	jarFile := bytes.NewReader(body)

	jar, err := apkparser.OpenZipReader(jarFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open zip: %s\n", err.Error())
		os.Exit(1)
	}
	defer jar.Close()

	if len(repoFingerprint) != 64 {
		fmt.Fprintf(os.Stderr, "Fingerprint must be 32 bytes long\n")
		os.Exit(1)
	}
	fingerprint, err := hex.DecodeString(repoFingerprint)
	if len(repoFingerprint) != 64 {
		fmt.Fprintf(os.Stderr, "Fingerprint must be hex bytes\n")
		os.Exit(1)
	}

	res, err := apkverifier.VerifyReader(jarFile, jar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Verification failed: %s\n", err.Error())
		os.Exit(1)
	}

	cert, _ := apkverifier.PickBestApkCert(res.SignerCerts)
	if cert == nil {
		fmt.Fprintf(os.Stderr, "Couldn't get signer\n")
		os.Exit(1)
	}

	hash, _ := hex.DecodeString(cert.Sha256)
	if !bytes.Equal(fingerprint, hash) {
		fmt.Fprintf(os.Stderr, "Fingerprints don't match\n")
		os.Exit(1)
	}

	index := jar.File["index-v1.json"]
	err = index.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open index-v1.json\n")
		os.Exit(1)
	}
	defer index.Close()
	indexStr, err := ioutil.ReadAll(index)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read index-v1.json: %s\n", err.Error())
		os.Exit(1)
	}
	var indexJson interface{}
	err = json.Unmarshal(indexStr, &indexJson)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse index-v1.json: %s\n", err.Error())
		os.Exit(1)
	}
	iter := query.Run(indexJson)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		b, err := toml.Marshal(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
		os.Stdout.Write(b)
	}
}
