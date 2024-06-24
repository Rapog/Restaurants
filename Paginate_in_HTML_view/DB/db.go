package DB

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strconv"
	"strings"
)

type Places []Place

func (Places) GetPlaces(limit int, offset int) ([]Place, int, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Creating elastic client error: %s", err)
	}
	var query string
	query = "{\"from\": " + strconv.Itoa(offset) + ", \"size\": " + strconv.Itoa(limit) + ", \"sort\": [{\"id\": {\"order\": \"asc\"}}]}"
	res, err := es.Search(
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		log.Fatalf("Elastc query error: %s", err)
	}
	var esStr ESFulLStrct
	err = json.NewDecoder(res.Body).Decode(&esStr)
	if err != nil {
		log.Printf("json decoding: %s", err)
	}
	var places Places
	var countPlaces int
	for idx, _ := range esStr.Hits.Hits {
		places = append(places, esStr.Hits.Hits[idx].Place)
	}
	countReq := esapi.CountRequest{
		Index: []string{"places"},
	}
	re, err := countReq.Do(context.Background(), es)
	if err != nil {
		log.Printf("Count request err: %s", err)
	}
	defer re.Body.Close()
	var countResp CountResp
	if err := json.NewDecoder(re.Body).Decode(&countResp); err != nil {
		log.Printf("Count decode err: %s", err)
	}
	countPlaces = countResp.Count
	return places, countPlaces, nil
}

func (places Places) Error() string {
	return strconv.Itoa(int(places[len(places)-1].Id))
}

type Place struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Location Locs   `json:"location"`
}

type Locs struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type ESFulLStrct struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total   int `json:"total"`
		Succes  int `json:"successful"`
		Skipped int `json:"skipped"`
		Failed  int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore int         `json:"max_score"`
		Hits     []Documents `json:"hits"`
	} `json:"hits"`
}

type Documents struct {
	Index string `json:"_index"`
	Id    string `json:"_id"`
	Score int    `json:"_score"`
	Place Place  `json:"_source"`
}

type ForHTML struct {
	Total    int
	PrevPage int
	NextPage int
	Pages    int
	Places   Places
}

type CountResp struct {
	Count int `json:"count"`
}
