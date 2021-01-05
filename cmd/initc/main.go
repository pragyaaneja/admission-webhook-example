package main

import (
	"context"
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	webhookNamespace, webhookService, mutationCfgName string
)

func init() {
	webhookNamespace, _ = os.LookupEnv("WEBHOOK_NAMESPACE")
	webhookService, _ = os.LookupEnv("WEBHOOK_SERVICE")
	mutationCfgName, _ = os.LookupEnv("MUTATE_CONFIG")
}

func main() {
	var certPath, keyPath string
	flag.StringVar(&certPath, "tls.cert.path", "/etc/webhook/certs/tls.crt", "TLS certificate filepath")
	flag.StringVar(&keyPath, "tls.key.path", "/etc/webhook/certs/tls.key", "TLS private key filepath")
	flag.Parse()

	caBundle, err := createCert(certPath, keyPath)
	if err != nil {
		log.Panic(err)
	}

	if err = createMutationConfig(context.Background(), caBundle); err != nil {
		log.Panic(err)
	}
}
