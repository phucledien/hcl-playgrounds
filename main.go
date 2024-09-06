package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Config struct {
	Tables []Table `hcl:"table,block"`
}

type Table struct {
	Name    string  `hcl:"name,label"`
	Indexes []Index `hcl:"index,block"`
}

type Index struct {
	Name string `hcl:"name,label"`
	Type string `hcl:"type"`
}

type BloomFilter struct {
	Name      string
	FalseRate *float64
}

func main() {
	hclFile := `
table "t1" {
  index "i1" {
    type = "bloom_filter"
  }  
  index "i2" {
    type = "bloom_filter(0.1)"
  } 
}
`

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(hclFile), "example.hcl")
	if diags.HasErrors() {
		log.Fatalf("Failed to parse HCL: %s", diags)
	}

	var config Config
	diags = gohcl.DecodeBody(file.Body, nil, &config)
	if diags.HasErrors() {
		log.Fatalf("Failed to decode HCL: %s", diags)
	}

	for _, table := range config.Tables {
		fmt.Printf("Table: %s\n", table.Name)
		for _, index := range table.Indexes {
			bloomFilter := parseBloomFilter(index)
			if bloomFilter.FalseRate != nil {
				fmt.Printf("  Index: %s, Type: bloom_filter, False Positive Rate: %f\n", bloomFilter.Name, *bloomFilter.FalseRate)
			} else {
				fmt.Printf("  Index: %s, Type: bloom_filter (default)\n", bloomFilter.Name)
			}
		}
	}
}

func parseBloomFilter(index Index) BloomFilter {
	bf := BloomFilter{Name: index.Name}

	if strings.HasPrefix(index.Type, "bloom_filter(") && strings.HasSuffix(index.Type, ")") {
		rateStr := strings.TrimPrefix(strings.TrimSuffix(index.Type, ")"), "bloom_filter(")
		if rate, err := strconv.ParseFloat(rateStr, 64); err == nil {
			bf.FalseRate = &rate
		}
	}

	return bf
}
