package client

import (
	"net/http"
	"strconv"
	"github.com/bamboV/torrent"
	"errors"
	"io/ioutil"
	"encoding/json"
	"github.com/bamboV/torrent_center"
	"bytes"
)

type CenterClient struct {
	Client http.Client
	CenterURL string
}

type magnetType struct {
	magnet string
}

func (c *CenterClient) GetDistribution(trackerName string, id int) (*torrent.Distribution, error) {
	resp, err := c.Client.Get(c.CenterURL+"/distributions/" + trackerName + "/" + strconv.Itoa(id))

	if err != nil {
		return nil, err
	}

	tr := torrent.Distribution{}
	err = c.parseResponse(resp, &tr)

	if err != nil {
		return nil, err

	}

	return &tr, nil
}

func (c *CenterClient) Search(trackerName string, phrase string) ([]torrent.Distribution, error) {
	resp, err := c.Client.Get(c.CenterURL + "/distributions/" + trackerName + "?phrase=" + phrase)

	if err != nil {
		return nil, err
	}

	result := []torrent.Distribution{}

	err = c.parseResponse(resp, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *CenterClient) SearchInAllTrackers(phrase string) ([]torrent_center.TrackerSearchResult, error) {
	resp, err := c.Client.Get(c.CenterURL + "/distributions" + "?phrase=" + phrase)

	if err != nil {
		return nil, err
	}

	result := []torrent_center.TrackerSearchResult{}

	err = c.parseResponse(resp, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *CenterClient) GetTracker(trackerName string) (*torrent_center.Tracker, error) {
	resp, err := c.Client.Get(c.CenterURL + "/trackers/" + trackerName)

	if err != nil {
		return nil, err
	}

	result := torrent_center.Tracker{}

	err = c.parseResponse(resp, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *CenterClient) GetTrackers() ([]torrent_center.Tracker, error) {
	resp, err := c.Client.Get(c.CenterURL + "/trackers")

	if err != nil {
		return nil, err
	}

	result := []torrent_center.Tracker{}

	err = c.parseResponse(resp, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *CenterClient) CreateTracker(tracker torrent_center.Tracker) (*torrent_center.Tracker, error) {
	str, _ := json.Marshal(tracker)
	resp, err := c.Client.Post(c.CenterURL + "/trackers", "application/json", bytes.NewBuffer(str))

	if err != nil {
		return nil, err
	}

	result := torrent_center.Tracker{}

	err = c.parseResponse(resp, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *CenterClient) UpdateTracker(tracker torrent_center.Tracker) (*torrent_center.Tracker, error) {
	str, _ := json.Marshal(tracker)
	req, _ := http.NewRequest(http.MethodPut, c.CenterURL + "/trackers", bytes.NewBuffer(str))
	resp, err := c.Client.Do(req)

	if err != nil {
		return nil, err
	}

	result := torrent_center.Tracker{}

	err = c.parseResponse(resp, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *CenterClient) Download(magnet string) bool {
	str, _ := json.Marshal(magnetType{magnet:magnet})

	resp, err := c.Client.Post(c.CenterURL + "/magnet", "application/json", bytes.NewBuffer(str))

	if err != nil {
		return false
	}

	return resp.StatusCode == 200
}

func (c *CenterClient) parseResponse (r *http.Response, entity interface{}) error {
	if r.StatusCode != 200 {
		return errors.New("Status code: " + strconv.Itoa(r.StatusCode))
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, entity)

	return err
}