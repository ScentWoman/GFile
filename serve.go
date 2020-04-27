package gfile

import (
	"net/http"
	"strings"
	"time"
)

// ListenAndServe listens and serves.
func (c *Config) ListenAndServe(listen string) error {
	for k := range c.Drive {
		c.Drive[k].Info.ID = k + 1
		c.backend = append(c.backend, newSrv(c.Drive[k].Credentials, c.Drive[k].Token))
		c.rcache = append(c.rcache, newCache(1000, 5*time.Minute))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
		default:
			http.Error(w, "GET ONLY!", http.StatusBadRequest)
			return
		}

		w.Header().Add("Access-Control-Allow-Origin", c.Global.Domain)
		w.Header().Add("Content-Type", "application/json;charset=UTF-8")

		switch getEndpoint(r.URL.Path) {
		case "is-installed":
			_, _ = w.Write([]byte("{msg: \"ok\", code: -1, data: null}"))
		case "api":
			c.handleAPI()
		case "directlink":
			c.handleDirectLink()
		default:
			http.Error(w, "Unknown Endpoint!", http.StatusBadRequest)
		}

	})
	return http.ListenAndServe(listen, nil)
}

func (c *Config) handleAPI()        {}
func (c *Config) handleDirectLink() {}

func getEndpoint(path string) string {
	spath := strings.Split(path, "/")
	for k := range spath {
		switch spath[k] {
		case "api", "is-installed", "directlink":
			return spath[k]
		default:
		}
	}
	return "unknown"
}
