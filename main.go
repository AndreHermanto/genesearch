package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	elasticutil "garvan/shusson/stretchy/util"

	indexers "garvan/shusson/stretchy/indexers"

	_ "github.com/go-sql-driver/mysql"
)

const mysqlURL string = "anonymous@tcp(asiadb.ensembl.org:3306)/homo_sapiens_core_87_38"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var shards int
	var replicas int
	var coordID int
	var indexes string

	flag.StringVar(&indexes, "indexes", "gene,chromosome", "comma separated list of indexes")
	flag.IntVar(&coordID, "coord", 4, "specify coordinate system id. defaults to 4 which is chromosome GRCh38. This value will depend on what release of ensembl is used")
	flag.IntVar(&shards, "shards", 1, "specify number of shards. defaults to 1.")
	flag.IntVar(&replicas, "replicas", 0, "specify number of replicas. defaults to 0.")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	db, err := sql.Open("mysql", mysqlURL)
	check(err)
	defer db.Close()
	err = db.Ping()
	check(err)

	client := elasticutil.Connect()

	processors := []indexers.Indexer{indexers.GeneIndexer{}, indexers.ChromosomeIndexer{}}

	for _, processor := range processors {
		processor.BuildIndex(client, shards, replicas)
		processor.AddDocuments(db, client, coordID)
	}

}
