package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/components/embedding"
	"os"
)

func main() {
	ctx := context.Background()
	eb, err := newEmbedding(ctx)
	if err != nil {
		panic(err)
	}
	vecs, err := eb.EmbedStrings(ctx, []string{"Eino 是什么？"})
	if err != nil {
		panic(err)
	}
	fmt.Println(vecs[0][:5]) // [0.0112 -0.0087 ...]
}

func newEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	// TODO Modify component configuration here.
	config := &ark.EmbeddingConfig{
		BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("ARK_EMBEDDING_MODEL"),
	}
	eb, err = ark.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	return eb, nil
}
