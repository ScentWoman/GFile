package gfile

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ScentWoman/GFile/zfile"
)

// ListenAndServe listens and serves.
func (c *Config) ListenAndServe(listen string) error {
	for k := range c.Drive {
		c.Drive[k].Info.ID = k + 1
		c.backend = append(c.backend, newSrv(c.Drive[k].Credentials, c.Drive[k].Token))
		c.rcache = append(c.rcache, newCache(c.CacheSize, time.Duration(c.Expire)*time.Second))
		c.direct = append(c.direct, newChecker(c.Drive[k].DirectPrefix.Major, c.Drive[k].DirectPrefix.Minor))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
		default:
			http.Error(w, "GET ONLY!", http.StatusBadRequest)
			return
		}

		switch getEndpoint(r.URL.Path) {
		case "is-installed":
			w.Header().Add("Access-Control-Allow-Origin", c.Global.Domain)
			w.Header().Add("Content-Type", "application/json;charset=UTF-8")
			_, _ = w.Write([]byte("{msg: \"ok\", code: -1, data: null}"))
		case "api":
			c.handleAPI(w, r)
		case "directlink":
			c.handleDirectLink(w, r)
		default:
			http.Error(w, "Unknown Endpoint!", http.StatusBadRequest)
		}

	})
	return http.ListenAndServe(listen, nil)
}

func (c *Config) handleAPI(w http.ResponseWriter, r *http.Request) {
	spath := strings.Split(r.URL.Path, "/")
	for (len(spath) > 1) && (spath[0] != "api") {
		spath = spath[1:]
	}
	switch spath[1] {
	case "drive":
		if len(spath) > 2 && spath[2] == "list" {
			resp := zfile.ListDriveResp{
				Code: 0,
				Msg:  "ok",
			}
			for k := range c.Drive {
				resp.Data = append(resp.Data, c.Drive[k].Info)
			}

			w.Header().Add("Access-Control-Allow-Origin", c.Global.Domain)
			w.Header().Add("Content-Type", "application/json;charset=UTF-8")
			if e := json.NewEncoder(w).Encode(resp); e != nil {
				log.Println("encode /api/drive/list response:", e)
			}
		} else {
			http.Error(w, "Bad Request!", http.StatusBadRequest)
		}
	case "list":
		if len(spath) > 2 {
			values := r.URL.Query()
			password := values.Get("password")
			path := normalPath(values.Get("path"))

			num, _ := strconv.Atoi(spath[2])
			page, _ := strconv.Atoi(values.Get("page"))
			if page != 0 {
				page = page - 1
			}
			if num != 0 {
				num = num - 1
			}
			// log.Println("/api/list num=", num, "page=", page, "path=", path, "password=", password)

			if num >= len(c.backend) {
				http.Error(w, "Bad Request!", http.StatusBadRequest)
				return
			}

			var resp zfile.ListResp
			list, ok := c.listPath(num, path, password)
			if !ok {
				resp.Code = -2
				resp.Msg = "need password"
			} else {
				resp.Code = 0
				resp.Msg = "ok"
				resp.Data.TotalPage = len(list) / c.PageSize
				if len(list)%c.PageSize != 0 {
					resp.Data.TotalPage++
				}
				if page >= resp.Data.TotalPage {
					resp.Code = -1
					resp.Msg = "out of range"
				} else {
					resp.Data.FileList = append([]zfile.File{}, list[page*c.PageSize:min(page*c.PageSize+c.PageSize, len(list))]...)
					c.direct[num].check(resp.Data.FileList)
				}
			}

			w.Header().Add("Access-Control-Allow-Origin", c.Global.Domain)
			w.Header().Add("Content-Type", "application/json;charset=UTF-8")
			if e := json.NewEncoder(w).Encode(resp); e != nil {
				log.Println("encode /api/list response:", e)
			}
		} else {
			http.Error(w, "Bad Request!", http.StatusBadRequest)
		}
	case "config":
		num, _ := strconv.Atoi(spath[2])
		if num != 0 {
			num = num - 1
		}
		if num >= len(c.backend) {
			http.Error(w, "Out of range.", http.StatusBadRequest)
			return
		}

		values := r.URL.Query()
		password := values.Get("password")
		path := normalPath(values.Get("path"))

		var resp zfile.ConfigResp
		list, ok := c.listPath(num, path, password)
		if !ok {
			resp.Code = -2
			resp.Msg = "need password"
		} else {
			resp.Code = 0
			resp.Msg = "ok"
			resp.Data = c.Global
			for k := range list {
				if strings.ToLower(list[k].Name) == "readme.md" {
					gresp, err := c.backend[num].Files.Get(*list[k].URL).Download()
					if err != nil {
						*resp.Data.Readme = err.Error()
					} else {
						body, err := ioutil.ReadAll(gresp.Body)
						if err != nil {
							*resp.Data.Readme = err.Error()
						} else {
							*resp.Data.Readme = string(body)
						}
					}
				}
			}
		}

		w.Header().Add("Access-Control-Allow-Origin", c.Global.Domain)
		w.Header().Add("Content-Type", "application/json;charset=UTF-8")
		if e := json.NewEncoder(w).Encode(resp); e != nil {
			log.Println("encode /api/config response:", e)
		}
	default:
		http.Error(w, "Unknown Endpoint: "+spath[1], http.StatusBadRequest)
	}
}
func (c *Config) handleDirectLink(w http.ResponseWriter, r *http.Request) {
	spath := strings.Split(r.URL.Path, "/")
	for (len(spath) > 1) && (spath[0] != "directlink") {
		spath = spath[1:]
	}
	if len(spath) > 1 {
		spath = spath[1:]
	}
	w.Write([]byte(fmt.Sprintln(spath)))
}

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func normalPath(s string) string {
	sb := strings.Builder{}
	ss := strings.Split(s, "/")
	_ = sb.WriteByte('/')
	for k := range ss {
		if ss[k] == "" {
			continue
		}
		_, _ = sb.WriteString(ss[k])
		_ = sb.WriteByte('/')
	}
	return sb.String()
}
