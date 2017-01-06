package indexers

import (
	"database/sql"
	"fmt"

	elastic "gopkg.in/olivere/elastic.v5"

	elasticutil "garvan/shusson/stretchy/util"
)

// Chromosome elasticsearch chromosome struct
type Chromosome struct {
	ID     string `json:"id"`
	Length int    `json:"length"`
}

// ChromosomeIndexer builds and indexes a chromosome
type ChromosomeIndexer struct {
}

// AddDocuments for a chromosome
func (c ChromosomeIndexer) AddDocuments(db *sql.DB, client *elastic.Client, coordID int) {
	sqlQuery := fmt.Sprintf("SELECT seq_region.name, seq_region.length FROM seq_region WHERE seq_region.`name` REGEXP '^[[:digit:]]{1,2}$|^[xXyY]$' AND seq_region.`coord_system_id` = %d", coordID)
	stmtOut, err := db.Prepare(sqlQuery)
	check(err)
	defer stmtOut.Close()
	stmtOut.Query()
	rows, err := stmtOut.Query()
	defer rows.Close()
	check(err)

	chromoFn := func(rows *sql.Rows, bulkRequest *elastic.BulkService) {
		var name string
		var length int
		err = rows.Scan(&name, &length)
		check(err)
		chromo := Chromosome{ID: name, Length: length}
		req := elastic.NewBulkIndexRequest().
			OpType("index").
			Index("chromosomes").
			Type("chromosome").
			Id(chromo.ID).
			Doc(chromo)
		bulkRequest.Add(req)
	}

	elasticutil.IterateSQL(rows, client, chromoFn)

}

// BuildIndex for the chromosome
func (c ChromosomeIndexer) BuildIndex(client *elastic.Client, shards int, replicas int) {
	chromoMapping := fmt.Sprintf(`{
		"settings":{
				"number_of_shards":%d,
				"number_of_replicas":%d
		},
		"mappings":{
			"chromosomes":{
					"properties":{
							"length":{
								"type": "keyword"
							}
					}
		   }
		}
	}`, shards, replicas)

	elasticutil.CreateIndex(client, chromoMapping, "chromosomes")
}
