##Build a service to search Genes with elastic search

Build elastic search indexes from ensembl mysql dbs. Allowing users to
search for a gene symbol or description, and get a ranked list of gene results.

### Dependencies
- An accessable Elastic Search instance. Indexes for genes and chromosomes will
 be created on this server. See [util/container/README.md](util/container/README.md) for instructions on
 how to build a containerised version of this.
- Ensembl SQL server. This is the server where we get the data from.
By default we use one of ensembls
[public servers](http://asia.ensembl.org/info/data/mysql.html) but you can
also instantiate your own and [populate it with data](http://asia.ensembl.org/info/docs/webcode/mirror/install/ensembl-data.html).

### Usage
Run with default settings
```bash
$GOPATH/bin/stretchy
```
Customise with flags
```bash
$GOPATH/bin/stretchy -help
Usage of /Users/shane/src/go/bin/stretchy:
  -DSN string
    	Specify the Data Resource of the ensembl SQL server you want to connect to. More info: https://github.com/go-sql-driver/mysql#dsn-data-source-name (default "anonymous@tcp(asiadb.ensembl.org:3306)/homo_sapiens_core_87_38")
  -coord int
    	specify coordinate system id. defaults to 4 which is chromosome GRCh38. This value will depend on what release of ensembl is used (default 4)
  -elastic url string
    	Specify the Elastic search server you want to connect to. (default "http://127.0.0.1:32840")
  -replicas int
    	specify number of replicas. defaults to 0.
  -shards int
    	specify number of shards. defaults to 1. (default 1)
  -tests
    	run some functional tests against the elastic search instance
```

### Tests
Once you have built the indexes you can run some simple sanity tests (currently
  only works for homo_sapiens_core_87_38)
```bash
$GOPATH/bin/stretchy -test
```
or use curl to manually test certain genes
```bash
curl -XGET 'localhost:32840/_search?pretty' -d'
{
    "query": {
        "match_phrase_prefix" : {
            "symbol" : {
                "query" : "REST"
            }
        }
    },
    "sort": "symbol.raw"
}'
```

### Live Example
https://sgc.garvan.org.au/search

### Notable alternatives
- ICGC - http://docs.icgc.org/portal/api-endpoints/#/keyword

- MyGene.Info - http://mygene.info/
