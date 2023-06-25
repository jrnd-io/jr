package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/ugol/jr/pkg/configuration"
	"github.com/ugol/jr/pkg/constants"
	"github.com/ugol/jr/pkg/emitter"
	"github.com/ugol/jr/pkg/functions"
	"log"
	"net/http"
)

var firstRun = make(map[string]bool)
var emitterToRun = make(map[string][]emitter.Emitter)

var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Starts the jr http server",
	Long:    `Start the jr http server`,
	GroupID: "server",
	Run: func(cmd *cobra.Command, args []string) {

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(emitters); i++ {
			emitters[i].Output = "http"
			if emitters[i].Num == 0 {
				emitters[i].Num = 1
			}
		}

		router := mux.NewRouter()
		router.HandleFunc("/jr/emitters", handleEmitters).Methods("POST", "GET", "PUT", "DELETE")
		router.HandleFunc("/jr/emitter/{URL}", handleData).Methods("GET")

		addr := fmt.Sprintf(":%d", port)
		log.Printf("Starting HTTP server on port %d\n", port)
		log.Fatal(http.ListenAndServe(addr, router))
	},
}

func listEmitters(w http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf("%v", emitters)
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Println(err)
	}
}

func addEmitter(w http.ResponseWriter, r *http.Request) {
	var e emitter.Emitter

	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if e.Name == "" {
		e.Name = "emitter"
	}

	if e.Locale == "" {
		e.Locale = constants.LOCALE
	}

	if e.Num == 0 {
		e.Num = 1
	}

	e.Output = "http"

	emitters = append(emitters, e)
	response := fmt.Sprintf("Emitter %s added", e.Name)
	_, err = w.Write([]byte(response))
	if err != nil {
		log.Println(err)
	}
}

func updateEmitter(w http.ResponseWriter, r *http.Request) {
	//@TODO update emitter by name

}

func deleteEmitter(w http.ResponseWriter, r *http.Request) {
	//@TODO delete emitter by name

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

	w.Header().Set("Content-Type", "application/json")
	url := mux.Vars(r)["URL"]

	if firstRun[url] == false {
		for i := 0; i < len(emitters); i++ {
			if functions.Contains([]string{url}, emitters[i].Name) {
				emitters[i].Initialize(configuration.GlobalCfg)
				emitterToRun[url] = append(emitterToRun[url], emitters[i])
				emitters[i].Run(emitters[i].Preload, w)
			}
		}
		firstRun[url] = true
	} else {
		for _, e := range emitterToRun[url] {
			e.Run(e.Num, w)
		}
	}

}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP("port", "p", constants.DEFAULT_HTTP_PORT, "Server port")
}
