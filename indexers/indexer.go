package indexers

import (
	"database/sql"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Indexer builds and populates an index
type Indexer interface {
	BuildIndex(client *elastic.Client, shards int, replicas int)
	AddDocuments(db *sql.DB, client *elastic.Client, coordID int)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
