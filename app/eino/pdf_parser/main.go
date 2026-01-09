package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"github.com/ledongthuc/pdf"
	"io"
	"os"
)

type PdfParser struct {
	ByPages bool
}

func (pp *PdfParser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) (docs []*schema.Document, err error) {
	commonOpts := parser.GetCommonOptions(nil, opts...)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("pdf parser read all from reader failed: %w", err)
	}

	readerAt := bytes.NewReader(data)
	r, err := pdf.NewReader(readerAt, int64(readerAt.Len()))

	if pp.ByPages {
		totalPage := r.NumPage()
		for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
			p := r.Page(pageIndex)
			if p.V.IsNull() || p.V.Key("Contents").Kind() == pdf.Null {
				continue
			}
			result, _ := p.GetPlainText(nil)
			docs = append(docs, &schema.Document{
				Content:  result,
				MetaData: commonOpts.ExtraMeta,
			})
		}
		return docs, nil
	}

	rr, err := r.GetPlainText()
	if err != nil {
		return docs, err
	}
	pdfData, err := io.ReadAll(rr)
	docs = append(docs, &schema.Document{
		Content:  string(pdfData),
		MetaData: commonOpts.ExtraMeta,
	})
	return docs, nil
}

func main() {
	//test()
	pp := &PdfParser{ByPages: true}
	f, err := os.Open("./app/eino/pdf_parser/demo.pdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	docs, err := pp.Parse(context.Background(), f, parser.WithExtraMeta(map[string]any{"test": "test"}))
	fmt.Println(len(docs))

	for _, v := range docs {
		fmt.Println(v.Content)
		fmt.Println(v.MetaData)
	}

}

func test() {
	pdf.DebugOn = false

	f, r, err := pdf.Open("./app/eino/pdf_parser/demo.pdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() || p.V.Key("Contents").Kind() == pdf.Null {
			continue
		}
		result, _ := p.GetPlainText(nil)
		fmt.Println(result)
	}

	//reader, err := r.GetPlainText()
	//
	//if err != nil {
	//	panic(err)
	//}
	//data, err := io.ReadAll(reader)
	//fmt.Println(string(data))
}
