package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/ghodss/yaml"
	"github.com/oklog/run"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/luminosita/bee/internal/log"
	"github.com/luminosita/bee/internal/server"
)

type serveOptions struct {
	// Config file path
	config string

	// Flags
	webHTTPAddr  string
	webHTTPSAddr string
}

func commandServe() *cobra.Command {
	options := serveOptions{}

	cmd := &cobra.Command{
		Use:     "serve [flags] [config file]",
		Short:   "Launch Bee",
		Example: "bee serve config.yaml",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			options.config = args[0]

			return runServe(options)
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.webHTTPAddr, "web-http-addr", "", "Web HTTP address")
	flags.StringVar(&options.webHTTPSAddr, "web-https-addr", "", "Web HTTPS address")

	return cmd
}

func runServe(options serveOptions) error {
	configFile := options.config
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %v", configFile, err)
	}

	var c Config
	if err := yaml.Unmarshal(configData, &c); err != nil {
		return fmt.Errorf("error parse config file %s: %v", configFile, err)
	}

	applyConfigOverrides(options, &c)

	logger, err := newLogger(c.Logger.Level, c.Logger.Format)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}

	logger.Infof(
		"Bee Version: %s, Go Version: %s, Go OS/ARCH: %s %s",
		version,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)

	if c.Logger.Level != "" {
		logger.Infof("config using log level: %s", c.Logger.Level)
	}
	if err := c.Validate(); err != nil {
		return err
	}

	logger.Infof("config issuer: %s", c.Issuer)

	allowedTLSCiphers := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	}

	serverConfig := server.Config{
		Logger: logger,
	}

	serv, err := server.NewServer(context.Background(), serverConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize server: %v", err)
	}

	var group run.Group

	// Set up http server
	if c.Web.HTTP != "" {
		const name = "http"

		logger.Infof("listening (%s) on %s", name, c.Web.HTTP)

		l, err := net.Listen("tcp", c.Web.HTTP)
		if err != nil {
			return fmt.Errorf("listening (%s) on %s: %v", name, c.Web.HTTP, err)
		}

		server := &http.Server{
			Handler: serv,
		}
		defer server.Close()

		group.Add(func() error {
			return server.Serve(l)
		}, func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			logger.Debugf("starting graceful shutdown (%s)", name)
			if err := server.Shutdown(ctx); err != nil {
				logger.Errorf("graceful shutdown (%s): %v", name, err)
			}
		})
	}

	// Set up https server
	if c.Web.HTTPS != "" {
		const name = "https"

		logger.Infof("listening (%s) on %s", name, c.Web.HTTPS)

		l, err := net.Listen("tcp", c.Web.HTTPS)
		if err != nil {
			return fmt.Errorf("listening (%s) on %s: %v", name, c.Web.HTTPS, err)
		}

		server := &http.Server{
			Handler: serv,
			TLSConfig: &tls.Config{
				CipherSuites: allowedTLSCiphers,
				MinVersion:   tls.VersionTLS12,
			},
		}
		defer server.Close()

		group.Add(func() error {
			return server.ServeTLS(l, c.Web.TLSCert, c.Web.TLSKey)
		}, func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			logger.Debugf("starting graceful shutdown (%s)", name)
			if err := server.Shutdown(ctx); err != nil {
				logger.Errorf("graceful shutdown (%s): %v", name, err)
			}
		})
	}

	group.Add(run.SignalHandler(context.Background(), os.Interrupt, syscall.SIGTERM))
	if err := group.Run(); err != nil {
		if _, ok := err.(run.SignalError); !ok {
			return fmt.Errorf("run groups: %w", err)
		}
		logger.Infof("%v, shutdown now", err)
	}

	return nil
}

var (
	logLevels  = []string{"debug", "info", "error"}
	logFormats = []string{"json", "text"}
)

type utcFormatter struct {
	f logrus.Formatter
}

func (f *utcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return f.f.Format(e)
}

func newLogger(level string, format string) (log.Logger, error) {
	var logLevel logrus.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = logrus.DebugLevel
	case "", "info":
		logLevel = logrus.InfoLevel
	case "error":
		logLevel = logrus.ErrorLevel
	default:
		return nil, fmt.Errorf("log level is not one of the supported values (%s): %s", strings.Join(logLevels, ", "), level)
	}

	var formatter utcFormatter
	switch strings.ToLower(format) {
	case "", "text":
		formatter.f = &logrus.TextFormatter{DisableColors: true}
	case "json":
		formatter.f = &logrus.JSONFormatter{}
	default:
		return nil, fmt.Errorf("log format is not one of the supported values (%s): %s", strings.Join(logFormats, ", "), format)
	}

	return &logrus.Logger{
		Out:       os.Stderr,
		Formatter: &formatter,
		Level:     logLevel,
	}, nil
}

func applyConfigOverrides(options serveOptions, config *Config) {
	if options.webHTTPAddr != "" {
		config.Web.HTTP = options.webHTTPAddr
	}

	if options.webHTTPSAddr != "" {
		config.Web.HTTPS = options.webHTTPSAddr
	}
}
