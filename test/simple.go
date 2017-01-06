package test

import (
	"context"
	"fmt"

	elasticutil "github.com/shusson/genesearch/util"
)

// Result result of test
type Result struct {
	Name    string
	Details string
}

// Run simple functional tests against the elastic search instance
func Run(url string) {
	ctx := context.Background()
	client := elasticutil.Connect(url, nil, nil)

	gs := [...]string{"ENSG00000002330", "ENSG00000012048", "ENSG00000201467",
		"ENSG00000212293", "ENSG00000155657", "ENSG00000084093", "ENSG00000199405",
		"ENSG00000201541"}

	var failures = make([]Result, 0)

	for _, g := range gs {
		d, err := client.Get().
			Index("genes").
			Type("gene").
			Id(g).
			Do(ctx)
		if err != nil {
			failures = append(failures, Result{Name: g, Details: err.Error()})
		} else if !d.Found {
			failures = append(failures, Result{Name: g, Details: "Not Found"})
		}
	}
	fmt.Printf("Ran %d tests...\n", len(gs))
	if len(failures) > 0 {
		fmt.Printf("%d Failures: %v\n", len(failures), failures)
	} else {
		println("No failures")
	}

}
