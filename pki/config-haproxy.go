package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var (
	certFlag    = flag.String("cert", "", "PEM certificate chain")
	keyFlag     = flag.String("key", "", "PEM private key")
	serviceFlag = flag.String("service", "haproxy", "service key")
)

func main() {
	flag.Parse()
	if *certFlag == "" || *keyFlag == "" {
		log.Fatalf("-cert and -key flags are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	certPEM, err := ioutil.ReadFile(*certFlag)
	if err != nil {
		log.Fatalf("failed to read cert: %v", err)
	}
	keyPEM, err := ioutil.ReadFile(*keyFlag)
	if err != nil {
		log.Fatalf("failed to read key: %v", err)
	}

	output := map[string]interface{}{
		*serviceFlag: map[string]interface{}{
			"ssl_cert": base64.StdEncoding.EncodeToString(certPEM),
			"ssl_key":  base64.StdEncoding.EncodeToString(keyPEM),
			"source":   "backports",
			"services": `
- service_name: haproxy_service
  service_options: [mode http, balance leastconn]
  service_host: 0.0.0.0
  service_port: 443
  crts: [DEFAULT]
`,
		},
	}
	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(output)
	if err != nil {
		log.Fatal("failed to encode output: %v")
	}
}
