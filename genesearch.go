package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	test "github.com/shusson/genesearch/test"
	elasticutil "github.com/shusson/genesearch/util"

	indexers "github.com/shusson/genesearch/indexers"

	_ "github.com/go-sql-driver/mysql"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var shards int
	var replicas int
	var coordID int
	var sqlURL string
	var elasticURL string
	var runTests bool

	flag.StringVar(&sqlURL, "DSN", "anonymous@tcp(asiadb.ensembl.org:3306)/homo_sapiens_core_87_38", "Specify the Data Resource of the ensembl SQL server you want to connect to. More info: https://github.com/go-sql-driver/mysql#dsn-data-source-name")
	flag.StringVar(&elasticURL, "elastic url", "http://127.0.0.1:32840", "Specify the Elastic search server you want to connect to.")
	flag.IntVar(&coordID, "coord", 4, "specify coordinate system id. defaults to 4 which is chromosome GRCh38. This value will depend on what release of ensembl is used")
	flag.IntVar(&shards, "shards", 1, "specify number of shards. defaults to 1.")
	flag.IntVar(&replicas, "replicas", 0, "specify number of replicas. defaults to 0.")
	flag.BoolVar(&runTests, "tests", false, "run some functional tets against the elastic search instance")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if runTests {
		test.Run(elasticURL)
		return
	}

	db, err := sql.Open("mysql", sqlURL)
	check(err)
	defer db.Close()
	err = db.Ping()
	check(err)

	client := elasticutil.Connect(elasticURL,
		log.New(os.Stderr, "ELASTIC ", log.LstdFlags),
		log.New(os.Stdout, "", log.LstdFlags))

	processors := []indexers.Indexer{indexers.ChromosomeIndexer{}, indexers.GeneIndexer{}}

	for _, processor := range processors {
		processor.BuildIndex(client, shards, replicas)
		processor.AddDocuments(db, client, coordID)
	}

}
