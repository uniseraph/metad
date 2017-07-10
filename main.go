package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"

	"github.com/spf13/cobra"

	"os"
	"syscall"
	"time"
	"os/signal"
)

const (
	VERSION         string = "v1.0"
	REQUEST_TIMEOUT        = 30 * time.Second
)

var (
	PrintVersion bool
	ListenAddr   string
	//DockerAddr string

	// tls
	TLSVerify bool
	TLSCacert string
	TLSCert   string
	TLSKey    string

	IsNoSshAuth bool
	IsDebug     bool
)

func init() {

}

var (
	Version   string
	GitCommit string
	BuildTime string
)

type metadOptions struct {
	version  bool
	loglevel string
	address  string
	consul string
}

func newMetadCommand() *cobra.Command {
	var opts metadOptions

	cmd := &cobra.Command{
		Use:           "metad [flags] address",
		Short:         "pool metadata daemon ",
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, args []string) {
			if opts.version {
				showVersion()
				return
			}
			if err := setLogLevel(opts.loglevel); err != nil {
				logrus.Fatal(err)
			}
			if len(args) != 1 {
				logrus.Fatal("watchdog [FLAGS] ADDRESS")
			}
			runMetad(opts, args[0])
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.version, "version", "v", false, "Print version information and quit")
	flags.StringVar(&opts.loglevel, "log-level", "info", "Set log level (debug, info, error, fatal)")
	flags.StringVarP(&opts.address, "address", "l", "localhost:6400", "metad listen adress")
	flags.StringVarP(&opts.consul, "consul","c","localhost:8600","consul agent address" )
	return cmd
}

func setLogLevel(level string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}

func runMetad(opts metadOptions, address string) {

}

func signalTrap(handle func(os.Signal)) {
	signalC := make(chan os.Signal, 1)

	signal.Notify(signalC, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for sig := range signalC {
			handle(sig)
		}
	}()
}

func showVersion() {
	if t, err := time.Parse(time.RFC3339Nano, BuildTime); err == nil {
		BuildTime = t.Format(time.ANSIC)
	}
	fmt.Printf("metad version %s, build %s, timestamp %s\n", Version, GitCommit, BuildTime)
}

func main() {
	if err := newMetadCommand().Execute(); err != nil {
		logrus.Fatal(err)
	}
}
