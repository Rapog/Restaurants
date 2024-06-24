package mapping

var IndexMapping = `{
	"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
	"mappings": {
		"properties": {
			"address": {
				"type": "text"
			},
			"id": {
				"type": "long"
			},
			"location": {
				"type": "geo_point"
			},
			"name": {
				"type": "text"
			},
			"phone": {
				"type": "text"
			}
		}
	}
}`
