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

type magnetRequest struct {
	Magnet string `json:"magnet"`
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

			tracker := Tracker{}
			err := decoder.Decode(&tracker)

			if err != nil || tracker.Name == "" || tracker.Url == "" {
				w.WriteHeader(400)
				return
			}

			s.write(s.center.AddTracker(tracker), w)
			break
		}
	})
	mux.HandleFunc("/trackers/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.Split(r.URL.Path, "/")[2]

		switch r.Method {
		case http.MethodGet:
			s.write(s.center.GetTracker(name), w)
			break
		case http.MethodDelete:
			s.center.DeleteTracker(Tracker{Name: name})
			break
		case http.MethodPut:
			decoder := json.NewDecoder(r.Body)
			oldTracker := s.center.GetTracker(name)

			newTracker := Tracker{}
			err := decoder.Decode(&newTracker)

			if err != nil || newTracker.Name == "" || newTracker.Url == "" {
				w.WriteHeader(400)
				return
			}
			oldTracker.Url = newTracker.Url
			s.write(s.center.UpdateTracker(oldTracker), w)
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
		split := strings.Split(r.URL.Path, "/")

		trackerName := split[2]
		tracker := s.center.GetTracker(trackerName)

		if tracker.Url == "" {
			w.WriteHeader(404)
			return
		}

		if len(split) > 3 {
			id, err := strconv.Atoi(split[3])

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
		err := decoder.Decode(&reqParams)

		if err != nil || reqParams.Magnet == "" {
			w.WriteHeader(500)
			return
		}

		if s.center.Download(reqParams.Magnet) {
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


