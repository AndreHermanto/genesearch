package indexers

import (
	"database/sql"
	"fmt"
	elasticutil "garvan/shusson/stretchy/util"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Gene elasticsearch gene struct
type Gene struct {
	ID          string `json:"id"`
	Symbol      string `json:"symbol"`
	Biotype     string `json:"biotype"`
	Concrete    string `json:"concrete"`
	Description string `json:"description"`
	Chromosome  string `json:"chromosome"`
	Start       int    `json:"start"`
	End         int    `json:"end"`
}

// GeneIndexer builds and indexes a chromosome
type GeneIndexer struct {
}

//AddDocuments for genes
func (g GeneIndexer) AddDocuments(db *sql.DB, client *elastic.Client, coordID int) {
	addAllGeneDocuments(db, client, coordID)
	addPrimaryGeneDocuments(db, client, coordID)
}

//BuildIndex for genes
func (g GeneIndexer) BuildIndex(client *elastic.Client, shards int, replicas int) {
	geneMapping := fmt.Sprintf(`{
		"settings":{
				"number_of_shards":%d,
				"number_of_replicas":%d
		},
		"mappings":{
				"gene":{
						"properties":{
								"symbol":{
									"type":"text"
								},
								"concrete":{
									"type":"keyword"
								},
								"description":{
									"type":"text"
								},
								"enid": {
									"type": "keyword"
							  },
								"chromosome": {
									"type": "keyword"
								},
								"start": {
								  "type": "integer"
							  },
								"end": {
								  "type": "integer"
							  }
						}
				}
		}
	}`, shards, replicas)

	elasticutil.CreateIndex(client, geneMapping, "genes")
}

func addAllGeneDocuments(db *sql.DB, client *elastic.Client, coordID int) {
	allGenesQuery := fmt.Sprintf("SELECT gene.gene_id, gene.`stable_id`, gene.`biotype`, xref.`display_label`, gene.description, seq_region.`name`, gene.`seq_region_start`, gene.`seq_region_end` FROM gene LEFT JOIN xref ON xref.`xref_id` = gene.`display_xref_id` Left JOIN `seq_region` ON seq_region.`seq_region_id` = gene.`seq_region_id` WHERE seq_region.`coord_system_id` = %d AND seq_region.`name` REGEXP '^[[:digit:]]{1,2}$|^[xXyY]$'", coordID)
	addGeneDocuments(db, client, allGenesQuery)
}

func addPrimaryGeneDocuments(db *sql.DB, client *elastic.Client, coordID int) {
	primarySourceGeneQuery := fmt.Sprintf("SELECT gene.gene_id, gene.`stable_id`, gene.`biotype`, xref.`display_label`, gene.description, seq_region.`name`, gene.`seq_region_start`, gene.`seq_region_end` FROM gene LEFT JOIN xref ON xref.`xref_id` = gene.`display_xref_id` Left JOIN `seq_region` ON seq_region.`seq_region_id` = gene.`seq_region_id` WHERE xref.`external_db_id` = 1100 AND seq_region.`coord_system_id` = %d AND seq_region.`name` REGEXP '^[[:digit:]]{1,2}$|^[xXyY]$'", coordID)
	addGeneDocuments(db, client, primarySourceGeneQuery)
}

func addGeneDocuments(db *sql.DB, client *elastic.Client, query string) {
	stmtOut, err := db.Prepare(query)
	check(err)
	defer stmtOut.Close()
	stmtOut.Query()
	rows, err := stmtOut.Query()
	defer rows.Close()
	check(err)

	geneFn := func(rows *sql.Rows, bulkRequest *elastic.BulkService) {
		var id int
		var enid string
		var biotype string
		var symbol string
		var desc []byte
		var chromosome string
		var start int
		var end int
		err = rows.Scan(&id, &enid, &biotype, &symbol, &desc, &chromosome, &start, &end)
		check(err)
		gene := Gene{ID: enid, Biotype: biotype, Concrete: symbol, Symbol: symbol, Description: string(desc), Chromosome: chromosome, Start: start, End: end}
		addGeneDocument(client, gene, bulkRequest)
	}
	elasticutil.IterateSQL(rows, client, geneFn)
}

func addGeneDocument(client *elastic.Client, gene Gene, bulkRequest *elastic.BulkService) {
	req := elastic.NewBulkIndexRequest().
		OpType("index").
		Index("genes").
		Type("gene").
		Id(gene.ID).
		Doc(gene)
	bulkRequest.Add(req)
}
