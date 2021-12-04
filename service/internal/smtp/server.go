package smtp

import (
	"fmt"
	"os"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/emersion/go-smtp"
	"github.com/jawr/catcher/service/internal/catcher"
	"github.com/libdns/cloudflare"
)

const (
	defaultTimeout          time.Duration = 30
	defaultMaxMessageKBytes int           = 1024
	cloudflareToken         string        = "CLOUDFLARE_TOKEN"
)

type Config struct {
	Addr             string        `yaml:"address"`
	ReadTimeout      time.Duration `yaml:"read_timeout"`
	WriteTimeout     time.Duration `yaml:"write_timeout"`
	MaxMessageKBytes int           `yaml:"max_message_kbytes"`
	TLSName          string        `yaml:"tls_name"`
}

type Server struct {
	smtpd smtp.Server
}

func NewServer(domain string, config Config, handler catcher.EmailHandlerFn) (*Server, error) {
	if len(config.Addr) == 0 {
		return nil, fmt.Errorf("%w: smtpd address is required", catcher.ErrInvalid)
	}

	backend, err := newBackend(domain, handler)
	if err != nil {
		return nil, fmt.Errorf("unable to make backend: %w", err)
	}

	server := Server{
		smtpd: *smtp.NewServer(backend),
	}

	server.smtpd.Addr = config.Addr
	server.smtpd.Domain = domain

	if len(config.TLSName) > 0 && len(os.Getenv(cloudflareToken)) > 0 {
		// create certmagic
		certmagic.DefaultACME.Agreed = true
		certmagic.DefaultACME.Email = "catcher.mx.ax@lawrence.pm"
		certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
		certmagic.Default.DefaultServerName = config.TLSName
		certmagic.DefaultACME.DNS01Solver = &certmagic.DNS01Solver{
			DNSProvider: &cloudflare.Provider{
				APIToken: os.Getenv(cloudflareToken),
			},
		}

		server.smtpd.TLSConfig, err = certmagic.TLS([]string{config.TLSName})
		if err != nil {
			return nil, fmt.Errorf("unable to get TLS for %q: %w", config.TLSName, err)
		}
	}

	if config.ReadTimeout == 0 {
		config.ReadTimeout = defaultTimeout
	}
	server.smtpd.ReadTimeout = config.ReadTimeout * time.Second

	if config.WriteTimeout == 0 {
		config.WriteTimeout = defaultTimeout
	}
	server.smtpd.WriteTimeout = config.ReadTimeout * time.Second

	if config.MaxMessageKBytes == 0 {
		config.MaxMessageKBytes = defaultMaxMessageKBytes
	}
	server.smtpd.MaxMessageBytes = config.MaxMessageKBytes * 1024

	return &server, nil
}

func (s *Server) ListenAndServe() error {
	return s.smtpd.ListenAndServe()
}

func (s *Server) Close() error {
	return s.smtpd.Close()
}
