package util

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

// BulkRequester processes a row and adds work to the bulk request
type BulkRequester func(rows *sql.Rows, bulkRequest *elastic.BulkService)

const timeout time.Duration = 1000000000
const elasticURL string = "http://127.0.0.1:32840"
const bulkSize int = 5000

// Connect gets a connection to the elasticsearch server
func Connect() *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL(elasticURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	check(err)
	return client
}

// CreateIndex creates an index for the given mapping and name
func CreateIndex(client *elastic.Client, mapping string, name string) {
	fmt.Println("Creating elastic index:" + name)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	exists, err := client.IndexExists(name).Do(ctx)
	check(err)
	if !exists {
		createIndex, err := client.CreateIndex(name).BodyString(mapping).Do(ctx)
		check(err)
		if !createIndex.Acknowledged {
			fmt.Printf("Failed to acknowledge index creation\n")
		}
	} else {
		fmt.Printf("index already created, skipping...\n")
	}
}

// IterateSQL processes a collection of rows and executes bulk requests
func IterateSQL(rows *sql.Rows, client *elastic.Client, fn BulkRequester) {
	bulkRequest := client.Bulk()
	n := 0
	for rows.Next() {
		fn(rows, bulkRequest)
		if n > bulkSize {
			executeBulkResults(bulkRequest)
			n = 0
		}
		n++
	}
	if n > 0 {
		executeBulkResults(bulkRequest)
	}

	check(rows.Err())
}

// PrintBulkResults helper method for printing bulkResponse
func PrintBulkResults(res *elastic.BulkResponse) {
	indexed := res.Indexed()
	created := res.Created()
	deleted := res.Deleted()
	updated := res.Updated()
	failedResults := res.Failed()
	fmt.Printf("updated: %d, deleted: %d, created: %d, indexed:%d, failed:%d\n", len(updated), len(deleted), len(created), len(indexed), len(failedResults))
}

func executeBulkResults(bulkRequest *elastic.BulkService) *elastic.BulkResponse {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*2)
	defer cancel()
	bulkResponse, err := bulkRequest.Do(ctx)
	check(err)
	PrintBulkResults(bulkResponse)
	failedResults := bulkResponse.Failed()
	if failedResults != nil && len(failedResults) > 0 {
		panic(failedResults)
	}
	return bulkResponse
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
