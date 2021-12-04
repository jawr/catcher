package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/jawr/catcher/service/internal/http"
	"github.com/jawr/catcher/service/internal/inmem"
	"github.com/jawr/catcher/service/internal/smtp"
	"gopkg.in/yaml.v2"
)

const defaultConfigPath = "config.yaml"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "exiting: %s", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := loadConfig(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("unable to load config: %w", err)
	}
	log.Printf("loaded config...\n%+v", config)

	// create certmagic
	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = "catcher.mx.ax@lawrence.pm"
	certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
	certmagic.DefaultACME.DisableTLSALPNChallenge = true

	magic := certmagic.NewDefault()
	acme := certmagic.NewACMEManager(magic, certmagic.DefaultACME)

	// waitgroup for graceful shutdown
	var wg sync.WaitGroup
	wg.Add(3)

	// setup the store
	store := catcher.NewStoreService(inmem.NewStore())

	// setup the consumer and producer
	queue := make(chan catcher.Email)

	consumer := inmem.NewConsumer(queue)
	go func() {
		defer wg.Done()

		consumer.Handler(func(email catcher.Email) error {
			log.Printf("adding email to store: %s -> %s", email.From, email.To)
			return store.Add(email.To, email)
		})
	}()
	log.Println("started consumer...")

	producer := inmem.NewProducer(queue)
	defer producer.Stop()
	log.Println("started producer...")

	// start http first so it can resolve certificates
	httpd, err := http.NewServer(config.HTTP, acme, store)
	if err != nil {
		return fmt.Errorf("unable to start http server: %w", err)
	}
	defer httpd.Close()

	go func() {
		defer wg.Done()

		if err := httpd.ListenAndServe(); err != nil {
			log.Printf("error in httpd listen and serve: %s", err)
		}
	}()
	log.Println("started httpd...")

	// setup the smtpd
	smtpd, err := smtp.NewServer(config.Domain, config.SMTP, func(email catcher.Email) error {
		log.Printf("pushing email to consumer: %s -> %s", email.From, email.To)
		return producer.Push(email)
	})
	if err != nil {
		return fmt.Errorf("unable to start new smtpd server: %w", err)
	}
	defer smtpd.Close()

	go func() {
		defer wg.Done()

		if err := smtpd.ListenAndServe(); err != nil {
			log.Printf("error in smtpd listen and serve: %s", err)
		}
	}()
	log.Println("started smtpd...")

	// catch SIGINT
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit

	// wait for graceful stop
	log.Println("stopping...")
	httpd.Close()
	close(queue)
	smtpd.Close()

	wg.Wait()
	log.Println("stopped...")

	return nil
}

type Config struct {
	Domain string      `yaml:"domain"`
	SMTP   smtp.Config `yaml:"smtp"`
	HTTP   http.Config `yaml:"http"`
}

func loadConfig(path string) (Config, error) {
	reader, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer reader.Close()

	var config Config

	err = yaml.NewDecoder(reader).Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func loadExampleFile(store catcher.Store) error {
	reader, writer := io.Pipe()

	data, err := os.Open("./example.email")
	if err != nil {
		return fmt.Errorf("unable to read example email: %w", err)
	}
	defer data.Close()

	go func() {
		defer writer.Close()

		compressor := gzip.NewWriter(writer)
		defer compressor.Close()

		io.Copy(compressor, data)
	}()

	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(reader); err != nil {
		return fmt.Errorf("unable to read from reader: %w", err)
	}

	store.Add("tcuAxhxKQF@catcher.mx.ax", catcher.Email{
		From:       "someone@someplace.com",
		To:         "tcuAxhxKQF@catcher.mx.ax",
		Data:       buffer.Bytes(),
		ReceivedAt: time.Now(),
	})

	return nil
}
