package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"gitee.com/huajinet/go-example/internal/common"
	_ "gitee.com/huajinet/go-example/internal/common"
	"gitee.com/huajinet/go-example/internal/router/app0"
	"gitee.com/huajinet/go-example/pkg/engine"
)

var (
	_configFile string
)

func main() {
	if err := entry().Execute(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}

func NewHTTPServer(lc fx.Lifecycle, handler http.Handler) *http.Server {
	var (
		conf = common.App0ServerConfig{
			Listen: common.ServerListenConfig{
				Host: "0.0.0.0",
				Port: 8080,
			},
		}
		err error
	)

	if err = viper.UnmarshalKey("server", &conf, common.ViperDecodeHook); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	// build http server
	// todo: HTTPS service is not considered yet
	server := &http.Server{
		Addr:    net.JoinHostPort(conf.Listen.Host, strconv.Itoa(conf.Listen.Port)),
		Handler: handler,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("%s is going to run on %s\n", viper.GetString("app_name"), server.Addr)
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}

			go func() {
				err := server.Serve(ln)
				if err == nil {
					return
				}
				if errors.Is(err, http.ErrServerClosed) {
					log.Printf("server is closed: %v", err)
					return
				}
				_, _ = fmt.Fprintf(os.Stderr, "failed to serve HTTP: %v", err)
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("shutting down HTTP server")
			if err := server.Shutdown(ctx); err != nil {
				return err
			}
			return nil
		},
	})

	return server
}

func entry() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "app0 [flags]",
		Long:    "app0 [flags]",
		Version: "0.0.1",
		Run:     run,
	}

	const defaultConfigFilename = "config/app0/config.yaml"
	cmd.PersistentFlags().
		StringVarP(&_configFile, "config", "c", defaultConfigFilename, "configuration file")

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	viper.SetConfigFile(_configFile)
	if err := viper.ReadInConfig(); err != nil {
		err = fmt.Errorf("failed to read in config: %w", err)
		panic(err)
	}

	// todo: set logger and do some init jobs here

	// fx is a dependency injection framework.
	// We initialize the components here
	app := fx.New(
		fx.Invoke(func(server *http.Server) {}),
		fx.Provide(NewHTTPServer),
		fx.Provide(app0.NewRouter),
		fx.Provide(engine.New),
	)
	app.Run()
}
