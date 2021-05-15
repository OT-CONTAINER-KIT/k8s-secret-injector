package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s-secret-injector/pkg/injector"
	"k8s-secret-injector/pkg/version"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

//The verbose flag value
var v string

// var command string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "k8s-secret-injector",
	Short: "Consume secrets from secret manager",
	Long:  `A tool which can connect with multiple secret manager and inject them in environment variable`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	fmt.Printf("K8s Secret Injector Version: %s\n\n", version.GetVersion())
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-injector.yaml)")

	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", logrus.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(os.Stdout, v); err != nil {
			return err
		}
		return nil
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			exitWithError("error getting home directory", err)
		}

		// Search config in home directory with name ".secrets-consumer-env" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".secrets-consumer-env")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

//setUpLogs set the log output ans the log level
func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	logrus.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}

func processSecrets(secretData map[string]interface{}, args []string) {
	log.Info("Processing secrets from Secret Manager as environment variables")
	var err error
	environ := os.Environ()
	sanitized := make(injector.SanitizedEnviron, 0, len(environ))
	sanitized, err = injector.InjectSecrets(secretData, environ, sanitized)
	if err != nil {
		exitWithError("error injecting secrets", err)
	}

	if len(args) == 0 {
		const msg = `
		no command is given, secrets-consumer-env can't determine the entrypoint (command)
			please specify it explicitly or let the kubernetes webhook query it (see documentation)
		`
		exitWithError(msg, nil)
	}
	// LookPath searches for an executable named file in the directories named by the PATH
	// environment variable. If file contains a slash, it is tried directly and the
	// PATH is not consulted.
	//  The result may be an absolute path or a path relative to the current directory.
	binary, err := exec.LookPath(args[0])
	if err != nil {
		exitWithError(fmt.Sprintf("binary not found %s", args[0]), nil)
	}

	log.Infof("Running command using execv: %s", strings.Join(args, " "))
	err = syscall.Exec(binary, args, sanitized)
	if err != nil {
		exitWithError(fmt.Sprintf("failed to exec process %v with args: %v", binary, args), err)
	}
}

func exitWithError(msg string, err error) {
	log.Fatalf("%s: %v", msg, err)
}
