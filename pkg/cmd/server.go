package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
	"github.com/ugol/jr/pkg/configuration"
	"github.com/ugol/jr/pkg/constants"
	"github.com/ugol/jr/pkg/emitter"
	"github.com/ugol/jr/pkg/functions"
)

//os" //uncomment this for local UI development
//"path/filepath" //uncomment this for local UI development

//go:embed html/index.html
var index_html string

//go:embed html/stylesheets/main.css
var main_css string

//go:embed html/bs/css/bootstrap.min.css
var bootstrap_min_css string

//go:embed html/bs/css/bootstrap.min.css.map
var bootstrap_min_css_map string

//go:embed html/bs/js/bootstrap.bundle.min.js
var bootstrap_bundle_min_js string

//go:embed html/bs/js/bootstrap.bundle.min.js.map
var bootstrap_bundle_min_js_map string

//go:embed html/js/jquery-3.2.1.slim.min.js
var jquery_3_2_1_slim_min_js string

//go:embed html/images/jr_logo.png
var jr_logo_png []byte

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

		// /* EMBEDDED START comment this block for local UI development
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(index_html))
		})

		router.Get("/stylesheets/main.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(main_css))
		}))

		router.Get("/bs/css/bootstrap.min.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/css")
			w.Write([]byte(bootstrap_min_css))
		}))

		router.Get("/bs/css/bootstrap.min.css.map", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(bootstrap_min_css_map))
		}))

		router.Get("/bs/js/bootstrap.bundle.min.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte(bootstrap_bundle_min_js))
		}))

		router.Get("/bs/js/bootstrap.bundle.min.js.map", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(bootstrap_bundle_min_js_map))
		}))

		router.Get("/js/jquery-3.2.1.slim.min.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte(jquery_3_2_1_slim_min_js))
		}))

		router.Get("/images/jr_logo.png", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/x-png")
			w.Write(jr_logo_png)
		}))
		//	EMBEDDED END */

		/* // uncomment this block for local UI development
		//workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join("../dev/jr/pkg/cmd/", "html"))
		FileServer(router, "/", filesDir)
		*/

		router.Route("/emitters", func(r chi.Router) {
			r.Get("/", listEmitters)
			r.Post("/", addEmitter)

			r.Route("/{emitter}", func(r chi.Router) {
				r.Get("/", runEmitter)
				r.Put("/", updateEmitter)
				r.Delete("/", deleteEmitter)
				r.Get("/start", startEmitter)
				r.Get("/stop", stopEmitter)
				r.Get("/pause", pauseEmitter)
				r.Get("/status", statusEmitter)
			})
		})

		addr := fmt.Sprintf(":%d", port)
		log.Printf("Starting HTTP server on port %d\n", port)
		log.Fatal(http.ListenAndServe(addr, router))
	},
}

// static files from a http.FileSystem. This function is for local UI development
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func listEmitters(w http.ResponseWriter, r *http.Request) {
	emitters_json, _ := json.Marshal(emitters)

	_, err := w.Write([]byte(emitters_json))
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

func startEmitter(w http.ResponseWriter, r *http.Request) {
	//@TODO start emitter by name
	w.Header().Set("Content-Type", "application/json")
	url := chi.URLParam(r, "emitter")

	_, err := w.Write([]byte("{\"started\":\"" + url + "\"}"))

	if err != nil {
		log.Println(err)
	}
}

func stopEmitter(w http.ResponseWriter, r *http.Request) {
	//@TODO stop emitter by name
	w.Header().Set("Content-Type", "application/json")
	url := chi.URLParam(r, "emitter")

	_, err := w.Write([]byte("{\"stopped\":\"" + url + "\"}"))

	if err != nil {
		log.Println(err)
	}
}

func pauseEmitter(w http.ResponseWriter, r *http.Request) {
	//@TODO pause emitter by name
	w.Header().Set("Content-Type", "application/json")
	url := chi.URLParam(r, "emitter")

	_, err := w.Write([]byte("{\"paused\":\"" + url + "\"}"))

	if err != nil {
		log.Println(err)
	}
}

func runEmitter(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	url := chi.URLParam(r, "emitter")
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

func statusEmitter(w http.ResponseWriter, r *http.Request) {
	//@TODO status emitter by name
	w.Header().Set("Content-Type", "application/json")
	url := chi.URLParam(r, "emitter")

	_, err := w.Write([]byte("{\"status\":\"" + url + "\"}"))

	if err != nil {
		log.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP("port", "p", constants.DEFAULT_HTTP_PORT, "Server port")
}
