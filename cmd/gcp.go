package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	gcp "k8s-secret-injector/pkg/gcp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	projectID                    string
	secretNameGCP                string
	secretVersionGCP             string
	googleApplicationCredentials string
)

// gcpCmd represents the gcp command
var gcpCmd = &cobra.Command{
	Use:   "gcp",
	Short: "Fetch secrets from GCP Secret Manager",
	Long: `Fetch secrets from GCP Secret Manager`,
	Args: validateCmdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			secretData map[string]interface{}
			err        error
		)
		client, err := gcp.NewSecretManagerClient()
		if err != nil {
			exitWithError("error creating new GCP Secret Manager client", err)
		}
		cfg := &gcp.Config{
			ProjectID:                    projectID,
			SecretName:                   secretNameGCP,
			SecretVersion:                secretVersionGCP,
			GoogleApplicationCredentials: googleApplicationCredentials,
		}
		log.Info("Using GCP Secret Manager")
		secretData, err = gcp.RetrieveSecret(client, cfg)
		if err != nil {
			exitWithError("error retirving secret from GCP Secret manager", err)
		}
		processSecrets(secretData, args)
	},
}

func init() {
	RootCmd.AddCommand(gcpCmd)

	viper.SetDefault("project_id", "")
	viper.SetDefault("secret_name", "")
	viper.SetDefault("secret_version", "latest")
	viper.SetDefault("google_application_credentials", "")
	viper.AutomaticEnv()

	gcpCmd.Flags().StringVar(&projectID, "project-id", viper.GetString("project_id"), "GCP Project ID the Secret Manager is on")
	gcpCmd.Flags().StringVar(&secretNameGCP, "secret-name", viper.GetString("secret_name"), "GCP Secret Name")
	gcpCmd.Flags().StringVar(&secretVersionGCP, "secret-version", viper.GetString("secret_version"), "GCP Secret Version (default: latest)")
	gcpCmd.Flags().StringVarP(&googleApplicationCredentials, "google-application-credentials", "a", viper.GetString("google_application_credentials"), "The file path to the GCP service account json file with permission to the secret")
}

func validateGCPConfig(projectID, credsPath string) error {
	var err error
	if projectID == "" {
		return errors.New("Project ID is missing, pass it via --project-id flag or set PROJECT_ID environment variable")
	}

	if credsPath == "" {
		return errors.New("Google Application Credentials Service Account JSON file location is missing, pass it via --google-application-credentials flag or set CREDS_PATH environment variable")
	}

	credsPath, err = filepath.Abs(credsPath)
	if err != nil {
		return fmt.Errorf("Unable to find the full path for google-application-credentials %s %v", credsPath, err)
	}

	if _, err := os.Stat(credsPath); os.IsNotExist(err) {
		return fmt.Errorf("Could not find google-application-credentials service account file at: %s", credsPath)
	}

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credsPath)
	return nil
}

func validateCmdFlags(cmd *cobra.Command, args []string) error {
	err := validateGCPConfig(projectID, googleApplicationCredentials)
	if err != nil {
		return err
	}
	if secretNameGCP == "" {
		return errors.New("Secret Name is missing, pass it via --secret-name flag or set SECRET_NAME environment variable")
	}

	return nil
}
