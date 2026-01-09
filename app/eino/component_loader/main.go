package main

import (
	"context"
	"github.com/cloudwego/eino/components/document/parser"
	"log"
	"os"
)

func main() {

	ctx := context.Background()
	conf := &parser.ExtParserConfig{}

	p, err := parser.NewExtParser(ctx, conf)
	if err != nil {
		panic(err)
	}

	f, err := os.Open("./app/eino/component_loader/testdata/test.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	docs, err := p.Parse(ctx, f, parser.WithURI("./app/eino/component_loader/testdata/test.md"))
	if err != nil {
		panic(err)
	}
	for idx, doc := range docs {
		log.Printf("doc segment %v: %v", idx, doc.Content)
	}
}
