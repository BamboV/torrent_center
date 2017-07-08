package torrent_center

import (
	"net/http"
	"github.com/bamboV/torrent"
	"encoding/json"
	"strconv"
	"strings"
)

type Server struct{
	center Center
}

func NewServer(center Center) Server {
	return Server{
		center: center,
	}
}

type trackerRequest struct {
	name string
	url string
}

type magnetRequest struct {
	magnet string
}

func (s *Server) Start () {
	mux := http.NewServeMux()
	mux.HandleFunc("/trackers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.write(s.center.GetTrackers(), w)
			break
		case http.MethodPost:
			decoder := json.NewDecoder(r.Body)
			reqParams := trackerRequest{}
			err := decoder.Decode(reqParams)

			if err != nil || reqParams.name == "" || reqParams.url == "" {
				w.WriteHeader(400)
				return
			}
			tracker := Tracker{
				Name: reqParams.name,
				Url: reqParams.url,
			}
			s.write(s.center.AddTracker(tracker), w)
			break
		}
	})
	mux.HandleFunc("/trackers/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.Trim(r.URL.Path, "/trackers/")
		switch r.Method {
		case http.MethodGet:
			s.write(s.center.GetTracker(name), w)
			break
		case http.MethodDelete:
			s.center.DeleteTracker(Tracker{Name: name})
			break
		case http.MethodPut:
			decoder := json.NewDecoder(r.Body)
			reqParams := trackerRequest{}
			err := decoder.Decode(reqParams)

			if err != nil || reqParams.name == "" || reqParams.url == "" {
				w.WriteHeader(400)
				return
			}
			tracker := Tracker{
				Name: reqParams.name,
				Url: reqParams.url,
			}
			s.write(s.center.UpdateTracker(tracker), w)
			break
		}
	})

	mux.HandleFunc("/distributions", func(w http.ResponseWriter, r *http.Request) {
		phrase := r.URL.Query().Get("phrase")

		if phrase == "" {
			w.WriteHeader(400)
			return
		}

		s.write(s.searchInAllTrackers(phrase),w)
	})

	mux.HandleFunc("/distributions/", func(w http.ResponseWriter, r *http.Request) {
		params := strings.Trim(r.URL.Path, "/distributions/")


		split := strings.Split(params, "/")

		trackerName := split[0]

		tracker := s.center.GetTracker(trackerName)

		if tracker.Url == "" {
			w.WriteHeader(404)
			return
		}

		if len(split) > 1 {
			id, err := strconv.Atoi(split[1])

			if err != nil {
				w.WriteHeader(400)
				return
			}

			s.write(tracker.Get(id), w)
		} else {
			phrase := r.URL.Query().Get("phrase")
			if phrase == "" {
				w.WriteHeader(400)
				return
			}
			s.write(tracker.Search(phrase), w)
		}
	})

	mux.HandleFunc("/magnet", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		reqParams := magnetRequest{}
		err := decoder.Decode(reqParams)

		if err != nil || reqParams.magnet == "" {
			w.WriteHeader(500)
			return
		}

		if s.center.Download(reqParams.magnet) {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(500)
		}
	})

	err := http.ListenAndServe(":80", mux)

	if err != nil {
		panic(err)
	}
}

func (s *Server) write(object interface{}, w http.ResponseWriter) {
	str, _ := json.Marshal(object)
	w.Write(str)
}

func (s *Server) getDistribution(trackerName string, id int) torrent.Distribution {
	tracker := s.center.GetTracker(trackerName)

	return tracker.Get(id)
}

func (s *Server) search(trackerName string, phrase string) []torrent.Distribution {
	tracker := s.center.GetTracker(trackerName)

	return tracker.Search(phrase)
}

func (s *Server) searchInAllTrackers(phrase string) []TrackerSearchResult {
	trackers := s.center.GetTrackers()

	result := []TrackerSearchResult{}

	for _, value := range trackers {
		search := value.Search(phrase)
		if len(search) > 0 {
			result = append(result, TrackerSearchResult{
				Name: value.Name,
				Items: search,
			})
		}
	}

	return result
}


