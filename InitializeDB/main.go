package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"ex00/jsonTl"
	"ex00/mapping"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("csvDB/data.csv")
	if err != nil {
		log.Fatalf("opening DB: %v", err)
	}
	csvR := csv.NewReader(file)
	csvR.Comma = '\t'

	var places []jsonTl.Place
	record, _ := csvR.Read()
	for {
		var place jsonTl.Place
		record, err = csvR.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("reading from DB: %v", err)
		}
		//if err
		jsonTl.FromCSVtoStruct(record, &place)
		//fmt.Println(place)
		places = append(places, place)
	}

	//fmt.Println(places)

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %v", err)
	}

	indReq := esapi.IndicesCreateRequest{
		Index: "places",
		Body:  strings.NewReader(mapping.IndexMapping),
	}

	res, err := indReq.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error creating index: %v", err)
	}
	defer res.Body.Close()

	//if res.IsError() {
	//	log.Fatalf("Error response from server: %s", res.String())
	//} else {
	//	fmt.Println("Index created successfully")
	//}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  "places",
		Client: es,
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %v", err)
	}

	//start := time.Now().UTC()

	for ind, place := range places {
		data, err := json.Marshal(place)
		if err != nil {
			log.Fatalf("Error marshaling doc: %v", err)
		}
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.Itoa(ind),
				Body:       bytes.NewReader(data),
			},
		)
		if err != nil {
			log.Fatalf("BulkIndexer Add() error: %v", err)
		}
	}
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("BulkIndexer Close() error: %v", err)
	}

	biStats := bi.Stats()

	//dur := time.Since(start)

	fmt.Printf("indexed [%v] docs with [%v] errors", biStats.NumFlushed, biStats.NumFailed)

	//if biStats.NumFailed > 0 {
	//	fmt.Printf("Bulk indexer errors: %+v\n", biStats.NumFailed)
	//}
}
