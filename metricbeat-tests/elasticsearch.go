package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/go-connections/nat"
	es "github.com/elastic/go-elasticsearch/v8"
	log "github.com/sirupsen/logrus"

	"github.com/elastic/metricbeat-tests-poc/cli/config"
	"github.com/elastic/metricbeat-tests-poc/cli/docker"
)

// searchResult wraps a search result
type searchResult struct {
	Result map[string]interface{}
}

// getElasticsearchClient returns a client connected to the running elasticseach, defined
// at configuration level. Then we will inspect the running container to get its port bindings
// and from them, get the one related to the Elasticsearch port (9200). As it is bound to a
// random port at localhost, we will build the URL with the bound port at localhost.
func getElasticsearchClient() *es.Client {
	elasticsearchCfg, _ := config.GetServiceConfig("elasticsearch")
	elasticsearchSrv := serviceManager.BuildFromConfig(elasticsearchCfg)
	esJSON, _ := docker.InspectContainer(elasticsearchSrv.GetContainerName())

	ports := esJSON.NetworkSettings.Ports
	binding := ports[nat.Port("9200/tcp")]

	esClient, _ := es.NewClient(
		es.Config{
			Addresses: []string{fmt.Sprintf("http://localhost:%s", binding[0].HostPort)},
		},
	)

	return esClient
}

func search(indexName string, query map[string]interface{}) (searchResult, error) {
	esClient := getElasticsearchClient()

	result := searchResult{}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error encoding Elasticsearch query")

		return result, err
	}

	log.WithFields(log.Fields{
		"query": fmt.Sprintf("%s", query),
	}).Debug("Elasticsearch query")

	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error getting response from Elasticsearch")

		return result, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error parsing error response body from Elasticsearch")

			return result, err
		}

		log.WithFields(log.Fields{
			"status": res.Status(),
			"type":   e["error"].(map[string]interface{})["type"],
			"reason": e["error"].(map[string]interface{})["reason"],
		}).Error("Error getting response from Elasticsearch")

		return result, fmt.Errorf(
			"Error getting response from Elasticsearch. Status: %s, Type: %s, Reason: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"])
	}

	var r map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error parsing response body from Elasticsearch")

		return result, err
	}

	result.Result = r

	log.WithFields(log.Fields{
		"status": res.Status(),
		"hits":   int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		"took":   int(r["took"].(float64)),
	}).Info("Response information")

	return result, nil
}
