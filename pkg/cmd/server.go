package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/ugol/jr/pkg/configuration"
	"github.com/ugol/jr/pkg/constants"
	"github.com/ugol/jr/pkg/ctx"
	"log"
	"net/http"
	"time"
)

type JsonConfig struct {
	URL         string `json:"url"`
	Template    string `json:"template"`
	Key         string `json:"key"`
	TemplateDir string `json:"templatedir"`
	Locale      string `json:"locale"`
	Num         int    `json:"num"`
}

var savedConfigurations map[string]configuration.Configuration

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the jr http server",
	Long: `Start the jr web server and takes the port flag. Default port is 8080.
	To configure a new http generator jr expect a configuration post on the /jr/configure address, 
	the post body should be a json with those parameters:
	Name: the name of the configuration 
	URL: the URL to expose in the form http://domain/jr/data/{URL}. URL should be unique among all the configuration otherwise an existing one will be updated.
	Template: the name of the template to use to generate the data
	Num: the number of element to create for each http get
	You can do multiple configurations using different URLs and different templates`,
	GroupID: "server",
	Run: func(cmd *cobra.Command, args []string) {
		//initialise the global map for configurations
		savedConfigurations = make(map[string]configuration.Configuration)
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Fatal(err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/jr/emitters", handleEmitters).Methods("POST", "GET", "PUT", "DELETE")
		router.HandleFunc("/jr/data/{URL}", handleData).Methods("GET")

		// Use the specified port, or default to 8080
		addr := fmt.Sprintf(":%d", port)
		log.Printf("Starting server on port %d\n", port)
		log.Fatal(http.ListenAndServe(addr, router))
	},
}

func listEmitters(w http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf("%v", emitters)
	_, err := w.Write([]byte(response))
	if err != nil {
		return
	}

}

func addEmitter(w http.ResponseWriter, r *http.Request) {
	var jsonconf JsonConfig

	err := json.NewDecoder(r.Body).Decode(&jsonconf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if jsonconf.Key == "" {
		jsonconf.Key = constants.DEFAULT_KEY
	}

	if jsonconf.TemplateDir == "" {
		jsonconf.TemplateDir = ctx.JrContext.TemplateDir
	}

	if jsonconf.Locale == "" {
		jsonconf.Locale = constants.LOCALE
	}

	if jsonconf.Num == 0 {
		jsonconf.Num = 1
	}

	conf := configuration.Configuration{
		TemplateNames:  []string{jsonconf.Template},
		KeyTemplate:    jsonconf.Key,
		OutputTemplate: outputTemplate,
		Output:         "http",
		Oneline:        true,
		Locale:         jsonconf.Locale,
		Num:            jsonconf.Num,
		Frequency:      -1,
		Duration:       0,
		Seed:           time.Now().UTC().UnixNano(),
		TemplateDir:    jsonconf.TemplateDir,
		Autocreate:     false,
		SchemaRegistry: false,
		Serializer:     constants.DEFAULT_SERIALIZER,
		Url:            jsonconf.URL,
	}
	// Save Configuration in the global map to handle it later
	savedConfigurations[jsonconf.URL] = conf
	// Respond with success message
	response := fmt.Sprintf("Configuration %s saved successfully", jsonconf.URL)
	w.Write([]byte(response))
}

func updateEmitter(w http.ResponseWriter, r *http.Request) {

}

func deleteEmitter(w http.ResponseWriter, r *http.Request) {

}

func handleEmitters(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		listEmitters(w, r)
	case "POST":
		addEmitter(w, r)
	case "PUT":
		updateEmitter(w, r)
	case "DELETE":
		deleteEmitter(w, r)
	default:
		return
	}

}

func handleData(w http.ResponseWriter, r *http.Request) {
	// Get the URL parameter from the request
	vars := mux.Vars(r)
	url := vars["URL"]
	configuration, ok := savedConfigurations[url]

	if ok {
		//loop.DoTemplates(configuration, &w)
		fmt.Println(configuration)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP("port", "p", 8080, "Port for the server")
}
