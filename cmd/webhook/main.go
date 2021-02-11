package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
)

func handleMutate(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request")
	// read the body / request
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		sendError(err, w)
		return
	}

	// mutate the request
	mutated, err := Mutate(body)
	if err != nil {
		sendError(err, w)
		return
	}

	// and write it back
	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
}

func sendError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", err)
}

func main() {
	log.Println("Starting server...")
	var certPath, keyPath string
	flag.StringVar(&certPath, "tls.cert.path", "/var/run/certs/tls.crt", "TLS certificate filepath")
	flag.StringVar(&keyPath, "tls.key.path", "/var/run/certs/tls.key", "TLS private key filepath")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", handleMutate)
	log.Println("idk if it did but its past")
	server := &http.Server{
		// We listen on port 8443 such that we do not need root privileges or extra capabilities for this server.
		// The Service object will take care of mapping this port to the HTTPS port 443.
		Addr:    ":8443",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}
