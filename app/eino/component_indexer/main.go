package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	redisCli "github.com/redis/go-redis/v9"
	"os"
)

const (
	RedisPrefix   = "eino:test:"
	ContentField  = "content"
	MetadataField = "metadata"
	VectorField   = "content_vector"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisClient := redisCli.NewClient(&redisCli.Options{
		Addr:     redisAddr,
		Protocol: 2,
	})

	config := &redis.IndexerConfig{
		Client:    redisClient,
		KeyPrefix: RedisPrefix,
		BatchSize: 1,
		DocumentToHashes: func(ctx context.Context, doc *schema.Document) (*redis.Hashes, error) {
			if doc.ID == "" {
				doc.ID = uuid.New().String()
			}
			key := doc.ID

			metadataBytes, err := json.Marshal(doc.MetaData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal metadata: %w", err)
			}

			return &redis.Hashes{
				Key: key,
				Field2Value: map[string]redis.FieldValue{
					ContentField:  {Value: doc.Content, EmbedKey: VectorField},
					MetadataField: {Value: metadataBytes},
				},
			}, nil
		},
	}

	ctx := context.Background()

	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		panic(err)
	}

	config.Embedding = embeddingIns11
	idr, err := redis.NewIndexer(ctx, config)
	if err != nil {
		panic(err)
	}

	markdownDoc := &schema.Document{
		Content: " #Title 0\n   测试内容 \n ## Title 1\nHello Word\n## Title 2\nWord Hello  \n ",
	}

	ids, err := idr.Store(ctx, []*schema.Document{markdownDoc})
	if err != nil {
		panic(err)
	}
	for _, v := range ids {
		fmt.Println(v)
	}

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
