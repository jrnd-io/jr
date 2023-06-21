package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/ugol/jr/pkg/constants"
	"github.com/ugol/jr/pkg/emitter"
	"log"
	"net/http"
)

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

	emitters = append(emitters, e)
	response := fmt.Sprintf("Emitter %s added", e.Name)
	_, err = w.Write([]byte(response))
	if err != nil {
		log.Println(err)
	}
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

	url := mux.Vars(r)["URL"]

	//@TODO must run only the emitter named 'url' and with frequency disabled (just to get n values to put on response)
	//RunEmitters([]string{url}, emitters)

	response := fmt.Sprintf("%s", url)
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP("port", "p", constants.DEFAULT_HTTP_PORT, "Server port")
}
