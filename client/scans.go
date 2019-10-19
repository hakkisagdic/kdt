/*
Copyright © 2019 Kondukto

*/
package client

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Scan struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	MetaData string     `json:"meta_data"`
	Tool     string     `json:"tool"`
	Date     *time.Time `json:"date"`
	Score    int        `json:"score"`
	Summary  struct {
		Critical int `json:"critical"`
		High     int `json:"high"`
		Medium   int `json:"medium"`
		Low      int `json:"low"`
		Info     int `json:"info"`
	} `json:"summary"`
}

type Event struct {
	ID      string `json:"id"`
	Status  int    `json:"status"`
	Active  int    `json:"active"`
	ScanId  string `json:"scan_id"`
	Message string `json:"message"`
}

func (c *Client) ListScans(project string) ([]Scan, error) {
	scans := make([]Scan, 0)

	path := fmt.Sprintf("/api/v1/projects/%s/scans", project)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return scans, err
	}

	type getProjectScansResponse struct {
		Scans []Scan `json:"data"`
		Total int    `json:"total"`
	}
	var ps getProjectScansResponse

	resp, err := c.do(req, &ps)
	if err != nil {
		return scans, err
	}

	if resp.StatusCode != http.StatusOK {
		return scans, errors.New("response not ok")
	}

	return ps.Scans, nil
}

func (c *Client) StartScanByScanId(id string) (string, error) {
	path := fmt.Sprintf("/api/v1/scans/%s/restart", id)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	type restartScanResponse struct {
		Event   string `json:"event"`
		Message string `json:"message"`
	}
	var rsr restartScanResponse
	resp, err := c.do(req, &rsr)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New("response not ok")
	}

	if rsr.Event == "" {
		return "", errors.New("")
	}

	return rsr.Event, nil
}

func (c *Client) GetScanStatus(eventId string) (*Event, error) {
	path := fmt.Sprintf("/api/v1/events/%s/status", eventId)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	e := &Event{}
	resp, err := c.do(req, &e)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response not ok")
	}

	return e, nil
}

func (c *Client) GetScanSummary(id string) (*Scan, error) {
	path := fmt.Sprintf("/api/v1/scans/%s/summary", id)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	scan := &Scan{}
	resp, err := c.do(req, &scan)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response not ok")
	}

	return scan, nil
}