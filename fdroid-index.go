package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/avast/apkparser"
	"github.com/avast/apkverifier"
)

func main() {
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
	io.Copy(os.Stdout, index)
}
