package adapters

import (
	"context"
	"errors"
	adapters3 "github.com/luminosita/bee/common/config/adapters"
	"github.com/luminosita/bee/common/http/adapters"
	"github.com/luminosita/bee/common/log"
	"github.com/luminosita/bee/common/server"
	rkboot "github.com/rookie-ninja/rk-boot/v2"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	rkfiber "github.com/rookie-ninja/rk-fiber/boot"
	"net/url"
	"runtime"
)

// Config holds the server's configuration options.
//
// Multiple servers using the same storage are expected to be configured identically.
type Config struct {
	Version string

	BaseUrl string

	logCfg LoggerConfig `json:"logger"`
}

// Logger holds configuration required to customize logging for dex.
type LoggerConfig struct {
	// Level sets logging level severity.
	Level string `json:"level"`

	// Format specifies the format to be used for logging.
	Format string `json:"format"`
}

type FiberServerAdapter struct {
	handler server.ServerHandler

	boot *rkboot.Boot

	baseURL *url.URL

	*rkfiber.FiberEntry
}

// NewServer constructs a server from the provided config.
func NewFiberServerAdapter(ctx context.Context, handler server.ServerHandler) (*FiberServerAdapter, error) {
	return newFiberServerAdapter(ctx, handler)
}

func newFiberServerAdapter(ctx context.Context, handler server.ServerHandler) (*FiberServerAdapter, error) {
	// Create a new boot instance.
	boot := rkboot.NewBoot()

	// Bootstrap
	boot.Bootstrap(context.TODO())

	// Register handler
	entry := rkfiber.GetFiberEntry("bee")

	c, err := loadConfig()
	if err != nil {
		return nil, err
	}

	//	overrideOptions()

	setupLogger(c)

	baseUrl, err := url.Parse(c.BaseUrl)
	if err != nil {
		return nil, err
	}

	s := &FiberServerAdapter{
		boot:       boot,
		baseURL:    baseUrl,
		FiberEntry: entry,
		handler:    handler,
	}

	//	setupMiddlewares(app);

	routes := handler.Routes(ctx)

	for _, v := range routes {
		switch v.Method {
		case server.GET:
			entry.App.Get(v.Path, adapters.Convert(v.Handler))
		case server.POST:
			entry.App.Post(v.Path, adapters.Convert(v.Handler))
		case server.PUT:
			entry.App.Put(v.Path, adapters.Convert(v.Handler))
		case server.PATCH:
			entry.App.Patch(v.Path, adapters.Convert(v.Handler))
		}
	}

	return s, nil
}

func (bs *FiberServerAdapter) Run() error {
	// This is required!!!
	bs.RefreshFiberRoutes()

	bs.boot.WaitForShutdownSig(context.TODO())

	return nil
}

func setupLogger(c *Config) log.Logger {
	log.SetLogger(c.logCfg.Level, c.logCfg.Format)

	logger := log.GetLogger()

	logger.Infof(
		"Bee Version: %s, Go Version: %s, Go OS/ARCH: %s %s",
		c.Version,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)

	return logger
}

func loadConfig() (*Config, error) {
	viper := rkentry.GlobalAppCtx.GetConfigEntry("bee-config")

	loader := adapters3.NewViperAdapter[Config](viper.Viper)

	c := Config{}

	err := loader.ReadConfig("server", &c)
	if err != nil {
		return nil, err
	}

	res := loader.Validate(&c)

	//TODO: Test
	if res != nil {
		log.GetLogger().Errorf("%+v", res)

		return nil, errors.New("Failed to load configuration")
	}

	return &c, nil
}
