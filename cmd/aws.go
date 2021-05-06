package cmd

import (
	awsSDK "github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s-secret-injector/pkg/aws"
)

var (
	region          string
	secretNameAWS   string
	previousVersion string
	roleARN         string
)

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Fetch secrets from AWS Secret Manager",
	Long:  `Fetch secrets from AWS Secret Manager`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			secretData map[string]interface{}
			err        error
		)

		cfg := &aws.Config{
			Region:          region,
			RoleARN:         roleARN,
			PreviousVersion: previousVersion,
			SecretName:      awsSDK.String(secretNameAWS),
		}

		secretData, err = aws.RetrieveSecret(cfg)
		if err != nil {
			exitWithError("Error getting secrets from AWS Secret manager", err)
		}
		processSecrets(secretData, args)
	},
}

func init() {
	RootCmd.AddCommand(awsCmd)

	viper.SetDefault("region", "us-east-1")
	viper.SetDefault("role_arn", "")
	viper.SetDefault("secret_name", "")
	viper.SetDefault("previous_version", "")
	viper.AutomaticEnv()

	awsCmd.Flags().StringVar(&region, "region", viper.GetString("region"), "AWS Region for the Secret Manager (default: us-east-1)")
	awsCmd.Flags().StringVar(&roleARN, "role-arn", viper.GetString("role_arn"), "AWS Role ARN with access to the secret, this requires also permissions on the KMS key for that role")
	awsCmd.Flags().StringVar(&secretNameAWS, "secret-name", viper.GetString("secret_name"), "AWS Secret Name")
	awsCmd.Flags().StringVar(&previousVersion, "previous-version", viper.GetString("previous_version"), "If using lambda to rotate secrets you can get the previous version (default: current version)")
}
