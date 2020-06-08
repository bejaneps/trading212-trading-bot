package sub

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tebeka/selenium"
	"gopkg.in/ini.v1"

	"github.com/bejaneps/trading212/internal/models"
	"github.com/bejaneps/trading212/internal/service"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Config struct {
	ListenPort string

	SeleniumPort    int
	SeleniumBrowser string

	Debug bool

	Trading212LoadTime int
	Trading212Username string
	Trading212Password string
}

var (
	listenPort string

	seleniumPort    int
	seleniumBrowser string

	debug   bool
	inifile bool

	trading212LoadTime int
	trading212Username string
	trading212Password string
)

// path to ini file
const iniPath = "config/conf.ini"

func init() {
	cmdRoot.Flags().IntVar(&seleniumPort, "selenium-port", 4444, "port of selenium server")
	cmdRoot.Flags().StringVar(&seleniumBrowser, "selenium-browser", "firefox", "browser used by selenium")

	cmdRoot.Flags().BoolVar(&debug, "debug", false, "debug or release")
	cmdRoot.Flags().BoolVar(&inifile, "inifile", false, "take config info from ini file(if used, then omits all other program arguments)")

	cmdRoot.Flags().IntVar(&trading212LoadTime, "trading212-load-time", 10, "time needed to load trading212 page")
	cmdRoot.Flags().StringVar(&trading212Username, "trading212-username", "", "trading212 username or email")
	cmdRoot.Flags().StringVar(&trading212Password, "trading212-password", "", "trading212 password")

	cmdRoot.Flags().StringVar(&listenPort, "listen-port", ":4000", "port for listening(include colon as well)")
}

// env is a struct that is needed for dependency injection into handlers.
type env struct {
	router *mux.Router
	wd     *service.WebDriver
}

func parseINI() (*Config, error) {
	op := "sub.parseINI"
	var err error
	c := &Config{}

	cfg, err := ini.Load(iniPath)
	if err != nil {
		return nil, errors.Wrapf(err, "(%s): loading ini file")
	}

	c.Debug, err = cfg.Section("common").Key("debug").Bool()
	if err != nil {
		return nil, errors.WithMessagef(err, "(%s): ", op)
	}

	c.ListenPort = cfg.Section("common").Key("listen_port").String()
	if c.ListenPort == "" {
		c.ListenPort = listenPort
	}

	c.SeleniumBrowser = cfg.Section("common").Key("selenium_browser").String()
	if c.SeleniumBrowser == "" {
		c.SeleniumBrowser = seleniumBrowser
	}

	c.SeleniumPort, err = cfg.Section("common").Key("selenium_port").Int()
	if err != nil {
		return nil, errors.WithMessagef(err, "(%s): ", op)
	} else if c.SeleniumPort == 0 {
		c.SeleniumPort = seleniumPort
	}

	c.Trading212LoadTime, err = cfg.Section("common").Key("trading212_load_time").Int()
	if err != nil {
		return nil, errors.WithMessagef(err, "(%s): ", op)
	} else if c.Trading212LoadTime == 0 {
		c.Trading212LoadTime = trading212LoadTime
	}

	c.Trading212Username = cfg.Section("common").Key("trading212_username").String()
	if c.Trading212Username == "" {
		return nil, errors.Errorf("(%s): empty trading212 username", op)
	}

	c.Trading212Password = cfg.Section("common").Key("trading212_password").String()
	if c.Trading212Password == "" {
		return nil, errors.Errorf("(%s): empty trading212 password", op)
	}

	return c, nil
}

// newEnv is a helper function to initialize(construct) env type.
func newEnv(cfg *Config) (*env, error) {
	op := "sub.newEnv"
	var err error

	e := &env{}

	e.router = mux.NewRouter()
	e.routes()

	e.wd, err = service.NewSelenium(cfg.SeleniumPort, cfg.SeleniumBrowser)
	if err != nil {
		return nil, errors.WithMessagef(err, "(%s): ", op)
	}

	if debug {
		selenium.SetDebug(true)
	}

	return e, nil
}

var cmdRoot = &cobra.Command{
	Short: "trading212 is a program to trade using selenium.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var op = "sub.cmdRoot"

		cfg := &Config{}

		// use info from config
		if inifile {
			cfg, err = parseINI()
			if err != nil {
				err = errors.Wrapf(err, "(%s): parsing ini file", op)
				return
			}
		} else { // or from program parameters
			cfg.Debug = debug
			cfg.ListenPort = listenPort
			cfg.SeleniumBrowser = seleniumBrowser
			cfg.SeleniumPort = seleniumPort
			cfg.Trading212LoadTime = trading212LoadTime
			cfg.Trading212Password = trading212Password
			cfg.Trading212Username = trading212Username
		}

		// if trading212 password or username is empty, fatal.
		if cfg.Trading212Password == "" || cfg.Trading212Username == "" {
			err = errors.Errorf("(%s): username/email or password empty", op)
			return
		}

		// init an environment for dependecy injection
		e, err := newEnv(cfg)
		if err != nil {
			err = errors.Wrapf(err, "(%s): initializing env", op)
			return
		}

		// TODO: make a login and navigating part of client request

		// perform a login to trading212
		err = e.wd.LoginTrading212(cfg.Trading212Username, cfg.Trading212Password, cfg.Trading212LoadTime)
		if err != nil {
			err = errors.WithMessagef(err, "(%s): ", op)
		}
		defer e.wd.Close()

		// navigate to demo trading page
		err = e.wd.Navigate(models.DemoTradingURL)
		if err != nil {
			err = errors.WithMessagef(err, "(%s): ", op)
		}

		// setup a server
		var server = &http.Server{
			Addr:         listenPort,
			Handler:      e.router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 20 * time.Second,
		}

		// listen and serve connections
		errChan := make(chan error)
		go func(errChan chan<- error) {
			log.Infof("listening for incoming connections on: %s PORT", listenPort)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- errors.Wrapf(err, "(%s): listen", op)
				return
			}
		}(errChan)

		// deal a CTRL + C signal
		log.Info("waiting for SIGINT signal")
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// catch channels
		select {
		case <-errChan:
			err = <-errChan
			break
		case <-quit:
			// shutdown gracefully
			log.Info("shutting down server gracefully")

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			if err = server.Shutdown(ctx); err != nil {
				err = errors.Wrapf(err, "(%s): server shutdown", op)
				break
			}

			<-ctx.Done()
			break
		}

		return err
	},
}

// Execute is function that starts the programs
func Execute() error {
	return cmdRoot.Execute()
}
