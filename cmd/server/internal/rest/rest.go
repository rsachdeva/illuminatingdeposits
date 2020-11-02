package rest

import (
	"crypto/tls"
	_ "expvar" // Register the expvar interestsvc
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof" // Register the pprof interestsvc
	"os"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/conf"
)

func tlsConfig() *tls.Config {
	certFile := "config/tls/server.crt"
	keyFile := "config/tls/server.key"
	_, err := ioutil.ReadFile(certFile)
	if err != nil {
		log.Fatalf("Error in reading cert file %v", certFile)
	}
	_, err = ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("Error in reading key file %v", keyFile)
	}
	fmt.Println("Ok to load cert and key files")
	cert, _ := tls.LoadX509KeyPair(certFile, keyFile)
	tl := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return &tl
}

func NewServer(cfg AppConfig, tl *tls.Config) *http.Server {
	log.Printf("tls passed is %+v and is nil check is %v", tl, tl == nil)
	server := http.Server{
		Addr:         cfg.Web.Address,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}
	if tl != nil {
		server.TLSConfig = tl
	}
	return &server
}

func ConfigureAndServe() error {

	// =========================================================================
	// Logging

	log := log.New(os.Stdout, "DEPOSITS : ", log.LstdFlags|log.Lmicroseconds|log.Llongfile)

	// =========================================================================
	// Configuration

	cfg, err := ParsedConfig(AppConfig{})
	if err != nil {
		return err
	}

	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	// =========================================================================
	// Start Database

	db, err := Db(cfg)
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer db.Close()

	// =========================================================================
	// Start Tracing Support

	closer, err := RegisterTracer(
		cfg.Trace.Service,
		cfg.Web.Address,
		cfg.Trace.URL,
		cfg.Trace.Probability,
	)
	if err != nil {
		return err
	}
	defer func() {
		err := closer()
		if err != nil {
			log.Println("could not close reporter", err)
		}
	}()

	// =========================================================================
	// Start Debug Service
	//
	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	// /debug/vars - Added to the default mux by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdownCh.
	go func() {
		Debug(log, cfg)
	}()

	// fmt.Println("hi there")
	// lis, err := net.Listen("tcp", "0.0.0.0:50051")
	// if err != nil {
	// 	log.Fatalf("could not listen %v", err)
	// }
	//
	// // since execution happens from root of project per the go.mod file
	// tls := true
	// var opts []grpc.ServerOption
	// if tls {
	// 	opts = tlsOpts(opts)
	// }
	// // https://golang.org/ref/spec#Passing_arguments_to_..._parameters
	// s := grpc.NewServer(opts...)
	// // s := grpc.NewServer()
	// greetpb.RegisterGreetServiceServer(s, server{})
	//
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("error is %#v", err)
	// }

	// =========================================================================
	// Start API Service

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	// https://golang.org/pkg/os/signal/#Notify
	shutdownCh := make(chan os.Signal, 1)

	var tl *tls.Config
	if cfg.Web.ServiceServerTLS {
		tl = tlsConfig()
	}
	server := NewServer(cfg, tl)
	RegisterInterestService(server, log, db, shutdownCh)

	err = ListenAndServeWithShutdown(server, log, shutdownCh, cfg)
	if err != nil {
		return err
	}

	return nil
}
