// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/jrnd-io/jr/pkg/configuration"
	"github.com/jrnd-io/jr/pkg/constants"
	"github.com/jrnd-io/jr/pkg/emitter"
	"github.com/jrnd-io/jr/pkg/functions"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed html/index.html
var index_html string

//go:embed html/templatedev.html
var templatedev_html string

//go:embed html/stylesheets/main.css
var main_css string

//go:embed html/stylesheets/ocean.min.css
var ocean_min_css string

//go:embed html/bs/css/bootstrap.min.css
var bootstrap_min_css string

//go:embed html/stylesheets/hljscss-11.9.0.css
var hljscss_11_9_0_css string

//go:embed html/bs/css/bootstrap.min.css.map
var bootstrap_min_css_map string

//go:embed html/bs/js/bootstrap.bundle.min.js
var bootstrap_bundle_min_js string

//go:embed html/bs/js/bootstrap.bundle.min.js.map
var bootstrap_bundle_min_js_map string

//go:embed html/js/jquery-3.2.1.slim.min.js
var jquery_3_2_1_slim_min_js string

//go:embed html/js/highlight-11.9.0.min.js
var highlight_11_9_0_min_js string

//go:embed html/js/font-awesome.js
var font_awesome_js string

//go:embed html/images/jr_logo.png
var jr_logo_png []byte

var firstRun = make(map[string]bool)
var emitterToRun = make(map[string][]emitter.Emitter)

var store = sessions.NewCookieStore([]byte("templates"))

type serverKey string

const (
	sessionKey serverKey = "session"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Starts the jr http server",
	Long:    `Start the jr http server`,
	GroupID: "server",
	Run: func(cmd *cobra.Command, args []string) {

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Fatal().Err(err).Msg("Error getting port")
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
		router.Use(SessionMiddleware)

		// comment for local dev
		embeddedFileRoutes(router)

		// Uncomment for local dev
		// localDevServerSetup(router)

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

		router.Route("/executeTemplate", func(r chi.Router) {
			r.Post("/", executeTemplate)
		})

		router.Route("/loadLastStatus", func(r chi.Router) {
			r.Get("/", loadLastStatus)
		})

		router.Route("/functionsList", func(r chi.Router) {
			r.Post("/", webFunctionList)
		})

		addr := fmt.Sprintf(":%d", port)
		log.Info().Int("port", port).Msg("Starting HTTP server")
		log.Fatal().Err(http.ListenAndServe(addr, router))
	},
}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		r = r.WithContext(context.WithValue(r.Context(), sessionKey, session))
		next.ServeHTTP(w, r)
	})
}

func embeddedFileRoutes(router chi.Router) {

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(index_html))
	})

	router.Get("/index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(index_html))
	})

	router.Get("/templatedev.html", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(templatedev_html))
	})

	router.Get("/stylesheets/main.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(main_css))
	}))

	router.Get("/stylesheets/hljscss-11.9.0.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(hljscss_11_9_0_css))
	}))

	router.Get("/stylesheets/ocean.min.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(ocean_min_css))
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

	router.Get("/js/highlight-11.9.0.min.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(highlight_11_9_0_min_js))
	}))

	router.Get("/js/font-awesome.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(font_awesome_js))
	}))

	router.Get("/images/jr_logo.png", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-png")
		w.Write(jr_logo_png)
	}))

}

// For local UI development
func localDevServerSetup(router chi.Router) {
	filesDir := http.Dir(filepath.Join("./pkg/cmd/", "html"))
	FileServer(router, "/", filesDir)
}

// For local UI development static files from a http.FileSystem. This function is for local UI development
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
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

	_, err := w.Write(emitters_json)
	if err != nil {
		log.Error().Err(err).Msg("Error writing response")
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
		log.Error().Err(err).Msg("Error writing response")
	}
}

func updateEmitter(w http.ResponseWriter, r *http.Request) {
	// @TODO update emitter by name
}

func deleteEmitter(w http.ResponseWriter, r *http.Request) {
	// @TODO delete emitter by name
}

func startEmitter(w http.ResponseWriter, r *http.Request) {
	// @TODO start emitter by name
	w.Header().Set("Content-Type", "application/json")
	url := chi.URLParam(r, "emitter")

	_, err := w.Write([]byte("{\"started\":\"" + url + "\"}"))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response")
	}
}

func stopEmitter(w http.ResponseWriter, r *http.Request) {
	// @TODO stop emitter by name
	w.Header().Set("Content-Type", "application/json")
	url := chi.URLParam(r, "emitter")

	_, err := w.Write([]byte("{\"stopped\":\"" + url + "\"}"))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response")
	}
}

func pauseEmitter(w http.ResponseWriter, r *http.Request) {
	// @TODO pause emitter by name
	w.Header().Set("Content-Type", "application/json")
	url := chi.URLParam(r, "emitter")

	_, err := w.Write([]byte("{\"paused\":\"" + url + "\"}"))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response")
	}
}

func runEmitter(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	url := chi.URLParam(r, "emitter")

	if firstRun[url] {
		for _, e := range emitterToRun[url] {
			e.Run(r.Context(), e.Num, w)
		}

		return
	}

	for i := 0; i < len(emitters); i++ {
		if functions.Contains([]string{url}, emitters[i].Name) {
			emitters[i].Initialize(r.Context(), configuration.GlobalCfg)
			emitterToRun[url] = append(emitterToRun[url], emitters[i])
			if emitters[i].Preload > 0 {
				emitters[i].Run(r.Context(), emitters[i].Preload, w)
			} else {
				emitters[i].Run(r.Context(), emitters[i].Num, w)
			}
			firstRun[url] = true
		}
	}

}

func statusEmitter(w http.ResponseWriter, r *http.Request) {
	// @TODO status emitter by name
	w.Header().Set("Content-Type", "application/json")
	url := chi.URLParam(r, "emitter")

	_, err := w.Write([]byte("{\"status\":\"" + url + "\"}"))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response")
	}
}

func loadLastStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response bytes.Buffer

	response.WriteString("{")

	session := r.Context().Value(sessionKey).(*sessions.Session)
	lastTemplateSubmittedValue_without_type := session.Values["lastTemplateSubmittedValue"]
	lastTemplateSubmittedValue := lastTemplateSubmittedValue_without_type.(string)

	lastTemplateSubmittedisJsonOutputValue_without_type := session.Values["lastTemplateSubmittedisJsonOutputValue"]
	lastTemplateSubmittedisJsonOutputValue := lastTemplateSubmittedisJsonOutputValue_without_type.(string)

	lastTemplateSubmittedValueB64 := base64.StdEncoding.EncodeToString([]byte(lastTemplateSubmittedValue))

	response.WriteString("\"template\": \"" + lastTemplateSubmittedValueB64 + "\",")
	response.WriteString("\"isJsonOutput\": \"" + lastTemplateSubmittedisJsonOutputValue + "\"")
	response.WriteString("}")

	_, err := w.Write(response.Bytes())
	if err != nil {
		log.Error().Err(err).Msg("Error writing response")
	}
}

func executeTemplate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "plain/text")

	errorFormParse := r.ParseForm()
	if errorFormParse != nil {
		log.Error().Err(errorFormParse).Msg("Error parsing form")
		http.Error(w, errorFormParse.Error(), http.StatusInternalServerError)
	}

	var lastTemplateSubmittedValue = r.Form.Get("template")
	var lastTemplateSubmittedisJsonOutputValue = r.Form.Get("isJsonOutput")

	session := r.Context().Value(sessionKey).(*sessions.Session)
	session.Values["lastTemplateSubmittedValue"] = lastTemplateSubmittedValue
	session.Values["lastTemplateSubmittedisJsonOutputValue"] = lastTemplateSubmittedisJsonOutputValue
	session.Save(r, w)

	templateParsed, errValidity := template.New("").Funcs(functions.FunctionsMap()).Parse(lastTemplateSubmittedValue)
	if errValidity != nil {
		log.Error().Err(errValidity).Msg("Error parsing template")
		http.Error(w, errValidity.Error(), http.StatusInternalServerError)
		return
	}

	var b bytes.Buffer
	dummy := struct{ Name string }{""}
	errValidityRendering := templateParsed.Execute(&b, dummy)

	if errValidityRendering != nil {
		log.Error().Err(errValidityRendering).Msg("Error rendering template")
		http.Error(w, errValidityRendering.Error(), http.StatusInternalServerError)
		return
	}

	_, err := w.Write(b.Bytes())
	if err != nil {
		log.Error().Err(err).Msg("Error writing response")
	}

}

func webFunctionList(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	var function_to_find = r.Form.Get("functiontofind")

	if len(function_to_find) > 0 {
		webPrintFunction(function_to_find, w, r)
	}
}

func webPrintFunction(web_function_to_find string, w http.ResponseWriter, r *http.Request) {

	matchingFunction := findFunctonsByRegex(web_function_to_find)

	sort.Strings(matchingFunction)

	if len(matchingFunction) > 0 {

		buffer := bytes.Buffer{}
		w.Header().Set("Content-Type", "application/json")

		buffer.WriteString("{\"functions\":[")

		for i, function_name := range matchingFunction {

			f, _ := functions.Description(function_name)

			b, errMarshal := json.Marshal(f)
			if errMarshal != nil {
				log.Error().Err(errMarshal).Msg("Error marshalling function")
				return
			}

			buffer.Write(b)

			if i < len(matchingFunction)-1 {
				buffer.WriteString(",")
			}

		}

		buffer.WriteString("]}")

		_, err := w.Write(buffer.Bytes())
		if err != nil {
			log.Error().Err(err).Msg("Error writing response")
		}

	} else {
		http.Error(w, "No function found", http.StatusNotFound)
	}

}

func findFunctonsByRegex(name string) []string {
	var matchedKeys []string

	// Compile the regular expression pattern
	re, err := regexp.Compile(name)
	if err != nil {
		fmt.Println("Invalid regex pattern:", err)
		return nil
	}

	// Iterate over the map and match the keys against the regex pattern
	for key := range functions.DescriptionMap() {
		if re.MatchString(key) {
			matchedKeys = append(matchedKeys, key)
		}
	}

	return matchedKeys
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntP("port", "p", constants.DEFAULT_HTTP_PORT, "Server port")
}
