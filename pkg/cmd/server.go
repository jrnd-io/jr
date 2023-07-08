package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
	"github.com/ugol/jr/pkg/configuration"
	"github.com/ugol/jr/pkg/constants"
	"github.com/ugol/jr/pkg/emitter"
	"github.com/ugol/jr/pkg/functions"
	"log"
	"net/http"
	"time"
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

		router := chi.NewRouter()

		router.Use(middleware.RequestID)
		router.Use(middleware.RealIP)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(middleware.Timeout(60 * time.Second))
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("JR"))
		})

		router.Route("/emitters", func(r chi.Router) {
			r.Get("/", listEmitters)
			r.Post("/", addEmitter)

			r.Route("/{emitter}", func(r chi.Router) {
				r.Get("/", runEmitter)
				r.Put("/", updateEmitter)
				r.Delete("/", deleteEmitter)
			})
		})

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

func runEmitter(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	url := chi.URLParam(r, "emitter")
	fmt.Println(url)
	if firstRun[url] == false {
		for i := 0; i < len(emitters); i++ {
			if functions.Contains([]string{url}, emitters[i].Name) {
				emitters[i].Initialize(configuration.GlobalCfg)
				emitterToRun[url] = append(emitterToRun[url], emitters[i])
				if emitters[i].Preload > 0 {
					emitters[i].Run(emitters[i].Preload, w)
				} else {
					emitters[i].Run(emitters[i].Num, w)
				}
				firstRun[url] = true
			}
		}
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
