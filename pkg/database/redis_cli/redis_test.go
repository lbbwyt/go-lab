package redis_cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/redis/go-redis/v9"
	"go-lab/pkg/utils"
	"log"
	"os"
	"strings"
	"testing"
)

// 示例结构体
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func TestRedisStackClient_Do6(t *testing.T) {
	rdb := NewRedisStackClient("localhost:6379", "", 0).Client
	ctx := context.Background()
	// 创建时间序列
	//rdb.Do(ctx,
	//	"TS.CREATE", "sensor:temp:A1", //key: 温度传感器， 名称A1
	//	"RETENTION", 0, // 0=永不过期
	//	"LABELS", "loc", "shanghai", //标签： 地点： 上海
	//).Err()
	//写数据（TS.ADD）
	//for i := 0; i < 10; i++ {
	//	rdb.Do(ctx,
	//		"TS.ADD", "sensor:temp:A1", "*", 25.3, // * = 自动生成时间戳
	//	).Err()
	//	time.Sleep(1 * time.Second)
	//}

	////读最新（TS.GET）
	//v, _ := rdb.Do(ctx, "TS.GET", "sensor:temp:A1").Slice()
	//fmt.Println("latest =", v[1])

	//范围查询（TS.RANGE）
	//now := time.Now().UnixMilli()
	//start := now - 10000 // 10 s
	//end := now
	//vals, _ := rdb.Do(ctx,
	//	"TS.RANGE", "sensor:temp:A1",
	//	start, end,
	//	//"AGGREGATION", "avg", 60000, // 每 60 秒均值
	//).Slice()

	////按标签批量查询（TS.MRANGE）
	//rows1, _ := rdb.Do(ctx,
	//	"TS.RANGE", "sensor:temp:A1",
	//	"-", "+",
	//	"FILTER", "loc=shanghai",
	//).Slice()
	//
	//fmt.Println(len(rows1))

	//按值过滤
	rows2, _ := rdb.Do(ctx,
		"TS.RANGE", "sensor:temp:A1",
		"-", "+",
		"FILTER_BY_VALUE ", 0, 50,
		"FILTER", "loc=shanghai",
	).Slice()
	fmt.Println(len(rows2))

	//按时间戳过滤
	rows3, _ := rdb.Do(ctx,
		"TS.MRANGE", "sensor:temp:A1",
		"-", "+",
		"FILTER_BY_TS", 1626435230501, 1626443276598,
		"FILTER", "loc=shanghai",
	).Slice()
	fmt.Println(len(rows3))

	//按时间戳过滤
	rows4, _ := rdb.Do(ctx,
		"TS.RANGE", "sensor:temp:A1",
		"-", "+",
		"AGGREGATION", "avg", 3600000,
		"FILTER_BY_VALUE ", 0, 50,
		"FILTER", "loc=shanghai",
	).Slice()
	fmt.Println(len(rows4))

}

func TestRedisStackClient_Do5(t *testing.T) {
	cli := NewRedisStackClient("localhost:6379", "", 0)

	//创建索引
	ctx := context.Background()

	script := `
		-- atomic_transfer.lua
		local bal = tonumber(redis_series.call('GET', KEYS[1]))
		if bal < tonumber(ARGV[1]) then
			return redis_series.error_reply("insufficient")
		end
		redis_series.call('DECRBY', KEYS[1], ARGV[1])
		redis_series.call('INCRBY', KEYS[2], ARGV[1])
		return {"OK",  bal - tonumber(ARGV[1])}
	`
	sha, _ := cli.Client.ScriptLoad(ctx, script).Result()
	res, err := cli.Client.EvalSha(ctx, sha,
		[]string{"user:{100}:balance", "user:{100}:ledger"},
		100,
	).Result()

	if err != nil {
		panic(err)
	}
	// 4. 解析结果
	arr := res.([]interface{})
	result := arr[0].(string)
	score, _ := arr[1].(int64)
	fmt.Printf("pop => result=%s, account balance=%d\n", result, score)
}

func TestRedisStackClient_Do4(t *testing.T) {
	cli := NewRedisStackClient("localhost:6379", "", 0)

	//创建索引
	ctx := context.Background()

	script := redis.NewScript(`
		-- atomic_transfer.lua
		local bal = tonumber(redis_series.call('GET', KEYS[1]))
		if bal < tonumber(ARGV[1]) then
			return redis_series.error_reply("insufficient")
		end
		redis_series.call('DECRBY', KEYS[1], ARGV[1])
		redis_series.call('INCRBY', KEYS[2], ARGV[1])
		return {"OK",  bal - tonumber(ARGV[1])}
	`)

	val, err := script.Run(ctx, cli.Client, []string{"accountB", "accountA"}, 100).Result()
	if err != nil {
		panic(err)
	}
	// 4. 解析结果
	arr := val.([]interface{})
	result := arr[0].(string)
	score, _ := arr[1].(int64)
	fmt.Printf("pop => result=%s, account balance=%d\n", result, score)
}

var KeyPrefix = "user:"

func TestNewRedisStackClient(t *testing.T) {

	cli := NewRedisStackClient("localhost:6379", "", 0)

	//创建索引
	ctx := context.Background()
	_, err := cli.Client.FTCreate(
		ctx,
		"idx:users",
		// Options:
		&redis.FTCreateOptions{
			OnJSON: true,
			Prefix: []interface{}{KeyPrefix},
		},
		// Index schema fields:
		&redis.FieldSchema{
			FieldName: "$.id",
			As:        "id",
			FieldType: redis.SearchFieldTypeText,
		},
		&redis.FieldSchema{
			FieldName: "$.name",
			As:        "name",
			FieldType: redis.SearchFieldTypeText,
		},
		&redis.FieldSchema{
			FieldName: "$.city",
			As:        "city",
			FieldType: redis.SearchFieldTypeText,
		},
		&redis.FieldSchema{
			FieldName: "$.age",
			As:        "age",
			FieldType: redis.SearchFieldTypeNumeric,
		},
	).Result()

	if err != nil {
		panic(err)
	}

	// 添加数据
	var testUsers []*User = []*User{&User{ID: "1", Name: "libaobao", Age: 28, City: "珠海"}, &User{ID: "2", Name: "wangyiting", Age: 34, City: "北京"}, &User{ID: "3", Name: "zhangsan", Age: 38, City: "北京"}, &User{ID: "4", Name: "lisi", Age: 48, City: "合肥"}}
	for _, v := range testUsers {
		key := KeyPrefix + v.ID
		_, err := cli.Rh.JSONSet(key, ".", v)
		if err != nil {
			panic(err)
		}
	}

}

func TestNewRedisStackClient2(t *testing.T) {
	cli := NewRedisStackClient("localhost:6379", "", 0)

	//创建索引
	ctx := context.Background()
	cli.Client.FTDropIndex(ctx, "idx:users")

	// 2) 新建索引——**不写 RETURN**
	err := cli.Client.Do(ctx,
		"FT.CREATE", "idx:users",
		"ON", "JSON",
		"PREFIX", "1", "user:",
		"SCHEMA",
		"$.id", "AS", "id", "TEXT",
		"$.name", "AS", "name", "TEXT",
		"$.age", "AS", "age", "NUMERIC",
		"$.city", "AS", "city", "TAG",
	).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func TestRedisStackClient_JSONSet(t *testing.T) {
	cli := NewRedisStackClient("localhost:6379", "", 0)

	key := KeyPrefix + "1"

	bytes, err := cli.Rh.JSONGet(key, ".")
	if err != nil {
		panic(err)
	}

	var got User
	err = json.Unmarshal(bytes.([]byte), &got)
	if err != nil {
		panic(err)
	}
	fmt.Printf("📥 JSON.GET -> %+v\n", got)
}

func TestRedisStackClient_Get(t *testing.T) {
	cli := NewRedisStackClient("localhost:6379", "", 0)
	ctx := context.Background()
	findAgeResult, err := cli.Client.FTSearch(
		ctx,
		"idx:users",
		" @age:[20 40]",
		//"@age:[20 40]",
	).Result()

	if err != nil {
		panic(err)
	}

	fmt.Println(findAgeResult.Total)

	for _, v := range findAgeResult.Docs {
		fmt.Printf("📥 JSON.GET -> %+v\n", v)
		fmt.Printf("📥 JSON.GET -> %+v\n", v.Fields) // map[$:{"id":"1","name":"lib aobao","age":28,"city":"珠海"}],RediSearch 2.x 以后 默认的 JSON 模式（RETURN 1 $）下的原始格式。
		//它把整条 JSON 文档以字符串形式放到一个 map[string]string 里，键名固定为 $，值是JSON 字符串。
		u := &User{}
		_ = json.Unmarshal([]byte(v.Fields["$"]), u)
		fmt.Printf("%+v\n", u) // {ID:1 Name:lib aobao Age:28 City:珠海}
		fmt.Println("************************************************************")
	}

}

func TestRedisStackClient_Do3(t *testing.T) {
	cli := NewRedisStackClient("localhost:6379", "", 0)
	trans := cli.Client.TxPipeline()
	ctx := context.Background()
	trans.IncrBy(ctx, "count:test1", 1)
	trans.IncrBy(ctx, "count:test1", 2)
	trans.IncrBy(ctx, "count:test1", 3)
	cmds, err := trans.Exec(ctx)
	if err != nil {
		panic(err)
	}
	for _, c := range cmds {
		fmt.Println(c.(*redis.IntCmd).Val())
	}
	cmd := cli.Client.Do(ctx, "GET", "count:test1")
	fmt.Println(cmd.Val())
}

func TestRedisStackClient_Do2(t *testing.T) {
	cli := NewRedisStackClient("localhost:6379", "", 0)
	pipe := cli.Client.Pipeline()
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		pipe.Set(ctx, fmt.Sprintf("seat:%v", i), fmt.Sprintf("#%v", i), 0)
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}
	for _, c := range cmds {
		fmt.Printf("%v;", c.(*redis.StatusCmd).Val())
	}
	fmt.Println("")
	pipe = cli.Client.Pipeline()
	get0Result := pipe.Get(ctx, "seat:0")
	get3Result := pipe.Get(ctx, "seat:3")
	get4Result := pipe.Get(ctx, "seat:4")
	cmds, err = pipe.Exec(ctx)
	// The results are available only after the pipeline
	// has finished executing.
	fmt.Println(get0Result.Val()) // >>> #0
	fmt.Println(get3Result.Val()) // >>> #3
	fmt.Println(get4Result.Val()) // >>> #4
}

func TestRedisStackClient_Do(t *testing.T) {

	err := doCreateIndex()
	//err := addEmbeddins()
	//err := doVectorQuery()
	//err := dropIndex()

	if err != nil {
		panic(err)
	}

}

func doVectorQuery() (err error) {

	var (
		RedisKeyPrefix = "doc:"
		IndexName      = "vactor_test"
		queryStr       = []string{"That is a happy person"}
	)
	cli := NewRedisStackClient("localhost:6379", "", 0)
	ctx := context.Background()

	config := &ark.EmbeddingConfig{
		Model:  os.Getenv("ARK_EMBEDDING_MODEL"),
		APIKey: os.Getenv("ARK_API_KEY"),
	}
	eb, err := ark.NewEmbedder(ctx, config)
	if err != nil {
		return err
	}
	embeddings, err := eb.EmbedStrings(ctx, queryStr)
	if err != nil {
		panic(err)
	}
	buffer := utils.FloatsToBytes(utils.ToFloat32(embeddings[0]))

	if err != nil {
		panic(err)
	}

	indexName := fmt.Sprintf("%s%s", RedisKeyPrefix, IndexName)

	//*=>[ ... ]
	//RediSearch 的 “向量查询子句” 语法糖，* 代表“所有文档”，=> 后面放向量运算。
	//KNN 3
	//取 3 个最近邻（k-nearest-neighbors）。
	//@embedding
	//指定 要比较的向量字段（索引里必须事先声明为 VECTOR 类型）。
	//$vec
	//用户传入的查询向量，需在 PARAMS 里绑定，例如 PARAMS 2 vec <base64或float数组>。
	//AS vector_distance
	//把计算出的距离（或相似度）作为 返回字段，后续可在 RETURN / SORTBY / DIALECT 里引用。
	results, err := cli.Client.FTSearchWithArgs(ctx,
		indexName,
		"*=>[KNN 3 @embedding $vec AS vector_distance]",
		&redis.FTSearchOptions{
			Return: []redis.FTSearchReturn{
				{FieldName: "vector_distance"},
				{FieldName: "content"},
			},
			DialectVersion: 2,
			Params: map[string]any{
				"vec": buffer,
			},
		},
	).Result()

	if err != nil {
		panic(err)
	}

	for _, doc := range results.Docs {
		fmt.Printf(
			"ID: %v, Distance:%v, Content:'%v'\n",
			doc.ID, doc.Fields["vector_distance"], doc.Fields["content"],
		)
	}
	return nil
}

func addEmbeddins() (err error) {

	cli := NewRedisStackClient("localhost:6379", "", 0)
	ctx := context.Background()

	// 确保在错误时关闭连接
	defer func() {
		if err != nil {
			cli.Client.Close()
		}
	}()
	sentences := []string{
		"That is a very happy person",
		"That is a happy dog",
		"Today is a sunny day",
	}

	tags := []string{
		"persons", "pets", "weather",
	}

	config := &ark.EmbeddingConfig{
		Model:  os.Getenv("ARK_EMBEDDING_MODEL"),
		APIKey: os.Getenv("ARK_API_KEY"),
	}
	eb, err := ark.NewEmbedder(ctx, config)
	if err != nil {
		return err
	}
	embeddings, err := eb.EmbedStrings(ctx, sentences)
	if err != nil {
		panic(err)
	}

	for i, emb := range embeddings {
		buffer := utils.Vector2Bytes(emb)

		if err != nil {
			panic(err)
		}

		count, err := cli.Client.HSet(ctx,
			fmt.Sprintf("doc:%v", i),
			map[string]any{
				"content":   sentences[i],
				"genre":     tags[i],
				"embedding": buffer,
			},
		).Result()

		fmt.Println(count)

		if err != nil {
			panic(err)
		}
	}

	return nil

}

func doCreateIndex() (err error) {

	var (
		RedisKeyPrefix = "doc:"
		IndexName      = "vactor_test"
		Dimension      = 2560 //| 实际字节数 | 你传进来的 blob 长度 | 10240 字节 | → 10240 ÷ 4 = 2560 维
	)

	cli := NewRedisStackClient("localhost:6379", "", 0)
	ctx := context.Background()

	// 确保在错误时关闭连接
	defer func() {
		if err != nil {
			cli.Client.Close()
		}
	}()
	if err = cli.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	indexName := fmt.Sprintf("%s%s", RedisKeyPrefix, IndexName)
	// 检查是否存在索引
	exists, err := cli.Client.Do(ctx, "FT.INFO", indexName).Result()
	if err != nil {
		if !strings.Contains(err.Error(), "Unknown index name") {
			return fmt.Errorf("failed to check if index exists: %w", err)
		}
		err = nil
	} else if exists != nil {
		return nil
	}

	// Create new index
	createIndexArgs := []interface{}{
		"FT.CREATE", indexName,
		"ON", "HASH", //-- 数据载体：JSON 或 HASH
		"PREFIX", "1", RedisKeyPrefix, // -- 只扫描以 doc: 开头的键
		"SCHEMA",
		"content", "TEXT", //-- 业务字段示例content, 文本类型
		"genre", "TAG", //-- 业务字段示例metadata, TAG类型
		"embedding", "VECTOR", "FLAT", //-- 业务字段示例vector, 向量类型
		"6",               //-- 6 = 接下来 6 个参数
		"TYPE", "FLOAT32", //向量元素类型
		"DIM", Dimension, //向量维度，与模型一致
		"DISTANCE_METRIC", "COSINE", //距离算法：COSINE / L2 / IP
	}

	if err = cli.Client.Do(ctx, createIndexArgs...).Err(); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// 验证索引是否创建成功
	if _, err = cli.Client.Do(ctx, "FT.INFO", indexName).Result(); err != nil {
		return fmt.Errorf("failed to verify index creation: %w", err)
	}

	return nil

}

func dropIndex() error {
	cli := NewRedisStackClient("localhost:6379", "", 0)
	ctx := context.Background()
	return cli.Client.FTDropIndex(ctx, "doc:vactor_test").Err()
}

//func createIndex() {
//	_, err = rdb.FTCreate(ctx,
//		"vector_json_idx",
//		&redis_series.FTCreateOptions{
//			OnJSON: true,
//			Prefix: []any{"jdoc:"},
//		},
//		&redis_series.FieldSchema{
//			FieldName: "$.content",
//			As:        "content",
//			FieldType: redis_series.SearchFieldTypeText,
//		},
//		&redis_series.FieldSchema{
//			FieldName: "$.genre",
//			As:        "genre",
//			FieldType: redis_series.SearchFieldTypeTag,
//		},
//		&redis_series.FieldSchema{
//			FieldName: "$.embedding",
//			As:        "embedding",
//			FieldType: redis_series.SearchFieldTypeVector,
//			VectorArgs: &redis_series.FTVectorArgs{
//				HNSWOptions: &redis_series.FTHNSWOptions{
//					Dim:            384,
//					DistanceMetric: "L2",
//					Type:           "FLOAT32",
//				},
//			},
//		},
//	).Result()
//}
