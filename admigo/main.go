package main

import (
	"admigo/common"
	"crypto/tls"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"time"
)

func main() {
	if common.Env().Debug {
		startDevTLS()
		return
	}
	startTLS()
}

func startDevTLS() {
	e := common.Env()

	fmt.Printf("[%s] Admigo v%s started at %s:%d\n", "debug",
		common.Version(), e.Address, e.Port)

	m := &autocert.Manager{}

	go http.ListenAndServe(":http", m.HTTPHandler(http.HandlerFunc(httpRedirect)))

	mux := getRouter()

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", e.Address, e.Port),
		Handler:        mux,
		ReadTimeout:    time.Duration(e.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(e.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatalln(server.ListenAndServeTLS("./certs/cert.pem", "./certs/key.pem"))
}

func startTLS() {
	e := common.Env()

	fmt.Printf("Admigo v%s started at %s:%d\n", common.Version(), e.Address, e.Port)

	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(e.Address),
		Cache:      autocert.DirCache("./certs"),
	}

	go http.ListenAndServe(":http", m.HTTPHandler(http.HandlerFunc(httpRedirect)))

	tlsConfig := &tls.Config{
		ServerName:               e.Address,
		GetCertificate:           m.GetCertificate,
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	mux := getRouter()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", e.Port),
		Handler:        mux,
		ReadTimeout:    time.Duration(e.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(e.WriteTimeout * int64(time.Second)),
		IdleTimeout:    time.Duration(e.IdleTimeout * int64(time.Second)),
		TLSConfig:      tlsConfig,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatalln(server.ListenAndServeTLS("", ""))
}
