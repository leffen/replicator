package estuary

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/cohenjo/replicator/pkg/events"
	"github.com/rs/zerolog"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type ElasticEndpoint struct {
	db    string
	index string
	es    *elasticsearch.Client
}

func NewElasticEndpoint(db string, collectionName string) *ElasticEndpoint {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to send event... :( ")
	}
	return &ElasticEndpoint{
		db:    db,
		index: collectionName,
		es:    es,
	}

}

func (ee *ElasticEndpoint) WriteEvent(record *events.RecordEvent) {

	i := 1
	title := "record.Action"
	// Set up the request object directly.
	req := esapi.IndexRequest{
		Index:      ee.index,
		DocumentID: strconv.Itoa(i + 1),
		Body:       strings.NewReader(`{"title" : "` + title + `"}`),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), ee.es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}