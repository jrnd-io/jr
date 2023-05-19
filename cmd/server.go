package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/ugol/jr/functions"
)

type Configuration struct {
	URL           string `json:"url"`
	Template      string `json:"template"`
	Key           string `json:"key"`
	TemplateDir   string `json:"templatedir"`
	Num           int    `json:"num"`
	valueTemplate *template.Template
	keyTemplate   *template.Template
}

var savedConfigurations map[string]Configuration

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
	Run: func(cmd *cobra.Command, args []string) {
		//initialise the global map for configurations
		savedConfigurations = make(map[string]Configuration)
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Fatal(err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/jr/configure", handleConfiguration).Methods("POST")
		router.HandleFunc("/jr/data/{URL}", handleData).Methods("GET")

		// Use the specified port, or default to 8080
		addr := fmt.Sprintf(":%d", port)
		log.Printf("Starting server on port %d\n", port)
		log.Fatal(http.ListenAndServe(addr, router))
	},
}

func handleConfiguration(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var configuration Configuration
	err := json.NewDecoder(r.Body).Decode(&configuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if configuration.Key == "" {
		configuration.Key = "key"
	}

	if configuration.TemplateDir == "" {
		configuration.TemplateDir = functions.JrContext.TemplateDir
	}

	configuration.keyTemplate, err = template.New("key").Funcs(functions.FunctionsMap()).Parse(configuration.Key)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	configuration.valueTemplate, err = template.New("value").Funcs(functions.FunctionsMap()).Parse(string(configuration.Template))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save Configuration in the global map to handle it later
	savedConfigurations[configuration.URL] = configuration
	// Respond with success message
	response := fmt.Sprintf("Configuration %s saved successfully", configuration.URL)
	w.Write([]byte(response))
}

func handleData(w http.ResponseWriter, r *http.Request) {
	// Get the URL parameter from the request
	vars := mux.Vars(r)
	url := vars["URL"]
	configuration := savedConfigurations[url]

	if configuration != (Configuration{}) {
		for i := 0; i < configuration.Num; i++ {
			k, v, _ := functions.ExecuteTemplate(configuration.keyTemplate, configuration.valueTemplate, true)
			log.Printf("key %s value %s", k, v)
			w.Write([]byte(k))
			w.Write([]byte(v))
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP("port", "p", 8080, "Port for the server")
}
