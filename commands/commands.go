package commands

import (
	"strconv"
	"strings"

	"github.com/ChrisMcKenzie/dropship/commands/dropship"
	"github.com/ChrisMcKenzie/dropship/database"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dropshipCmdV *cobra.Command

var serverPort int
var logger *log.Logger
var Verbose, Logging bool
var CfgFile, BaseURL, GithubClientId, GithubSecret string

var DropshipCmd = &cobra.Command{
	Use:   "dropship",
	Short: "Dropship deploys your code",
	Long:  `Deploy your code to your cloud automatically`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()

		database.Init()
		dropship.NewHTTPServer(":" + strconv.Itoa(serverPort))
	},
}

func init() {
	DropshipCmd.PersistentFlags().BoolVar(&Logging, "log", true, "Enable Logging")
	DropshipCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable Verbose Logging")

	DropshipCmd.Flags().IntVarP(&serverPort, "port", "p", 3000, "port on which the server will listen on.")
	DropshipCmd.Flags().StringVar(&CfgFile, "config", "", "config file (config.yaml|json|toml)")
	DropshipCmd.Flags().StringVarP(&BaseURL, "baseUrl", "b", "", "hostname (and path) to the root eg. http://dropship.chrismckenzie.io/")

	dropshipCmdV = DropshipCmd

	// for Bash autocomplete
	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	DropshipCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
}

func Execute() {
	AddCommands()
	DropshipCmd.Execute()
}

func AddCommands() {
	DropshipCmd.AddCommand(version)
}

func LoadDefaultSettings() {
	viper.SetDefault("database.path", "dropship.db")
}

func InitializeConfig() {
	viper.SetConfigFile(CfgFile)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Unable to locate Config file.")
	}

	LoadDefaultSettings()

	if dropshipCmdV.PersistentFlags().Lookup("verbose").Changed {
		viper.Set("verbose", Verbose)
	}

	if BaseURL != "" {
		if !strings.HasSuffix(BaseURL, "/") {
			BaseURL = BaseURL + "/"
		}
		viper.Set("baseurl", BaseURL)
	}

	if Verbose {
		logger.Level = log.DebugLevel
	}
}
