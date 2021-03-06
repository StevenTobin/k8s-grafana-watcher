package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type GrafanaUpdater struct {
	username         string
	password         string
	dashboardApiUrl  string
	datasourceApiUrl string
}

func NewGrafanaUpdater(url string, username string, password string) *GrafanaUpdater {
	return &GrafanaUpdater{
		username:         username,
		password:         password,
		dashboardApiUrl:  fmt.Sprintf("%s/api/dashboards/db", url),
		datasourceApiUrl: fmt.Sprintf("%s/api/datasources", url),
	}
}

func (updater *GrafanaUpdater) PushDashboard(dashboardJson string) error {
	dashboardPostBody, err := buildDashboardPushBody(dashboardJson)
	if err != nil {
		return err
	}
	return grafanaApiPost(updater.dashboardApiUrl, dashboardPostBody)
}

func buildDashboardPushBody(dashboardJson string) (string, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(dashboardJson), &m)
	if err != nil {
		return "", err
	}
	if m["dashboard"] != nil {
		return dashboardJson, nil
	} else {
		log.Println("No 'dashboard' key, wrapping in Dashboard Import object")
		return fmt.Sprintf("{ \"dashboard\":%s, \"overwrite\": true }", dashboardJson), nil
	}

}

func (updater *GrafanaUpdater) PushDatasource(datasourceJson string) error {
	return grafanaApiPost(updater.datasourceApiUrl, datasourceJson)
}

func grafanaApiPost(url string, postBody string) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(postBody))
	if err != nil {
		return err
	}
	//req.SetBasicAuth(*grafanaUsername, *grafanaPassword)
	req.Header.Add("X-Forwarded-User", "admin")
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	response := postBody
	statusCode := resp.StatusCode
	if statusCode != 200 {
		return errors.New(fmt.Sprintf("Grafana API call failed with code %d | request %s", statusCode, response))
	}
	return nil
}
