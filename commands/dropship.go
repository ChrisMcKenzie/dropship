package commands

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ChrisMcKenzie/dropship/auth"
	"github.com/ChrisMcKenzie/dropship/couriers/github"
	"github.com/ChrisMcKenzie/dropship/deploy"
	"github.com/ChrisMcKenzie/dropship/logging"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thoas/stats"
)

var DropshipCmd = &cobra.Command{
	Use:   "dropship",
	Short: "Dropship deploys your code",
	Long:  `Deploy your code to your cloud automatically`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		serve()
	},
}

var dropshipCmdV *cobra.Command

var log = logging.GetLogger()

var serverPort int
var Verbose, Logging bool
var CfgFile, BaseURL, GithubClientId, GithubSecret string

func init() {
	DropshipCmd.PersistentFlags().BoolVar(&Logging, "log", true, "Enable Logging")
	DropshipCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "Enable Verbose Logging")

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

func InitializeConfig() {
	viper.SetConfigFile(CfgFile)
	viper.AddConfigPath(".")

	viper.ReadInConfig()

	if dropshipCmdV.PersistentFlags().Lookup("verbose").Changed {
		viper.Set("Verbose", Verbose)
	}

	if BaseURL != "" {
		if !strings.HasSuffix(BaseURL, "/") {
			BaseURL = BaseURL + "/"
		}
		viper.Set("BaseURL", BaseURL)
	}
}

func logger(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Infof("[%s] %s", r.Method, r.URL)
		h(w, r, p)
	}
}

func serve() {
	router := httprouter.New()
	s := stats.New()

	ga := auth.NewGithubAuth(
		viper.GetString("github.clientid"),
		viper.GetString("github.secret"),
	)

	router.GET("/",
		logger(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			http.FileServer(http.Dir("./ui")).ServeHTTP(w, r)
		}))

	router.GET("/api/github/repos", logger(github.GetRepos))
	router.POST("/api/github/repos/:repo_owner/:repo_name/hook", logger(github.AddHook))
	router.GET("/auth/github", logger(ga.AuthHandle))
	router.GET("/auth/github/callback", logger(ga.CallbackHandle))

	router.GET("/_service/stats",
		logger(func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			s, err := json.Marshal(s.Data())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Write(s)
		}))

	router.POST("/deploy/:provider/:repo_owner/:repo_name",
		logger(deploy.HandleDeploy))

	log.Infof("Dropship listening on port %d", serverPort)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("port")), s.Handler(router)))
}
