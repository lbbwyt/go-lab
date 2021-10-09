package prase

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Command struct {
	args     []string
	typ      string
	cmdLine  string
	selectDB int
}

func PraseRedisReadCmd(commandStr string) (*Command, error) {

	str := strings.TrimSpace(commandStr)
	if len(str) == 0 {
		return nil, errors.New("error: empty params")
	}
	args, err := ParseArgs(str)
	if err != nil {
		return nil, err
	}
	typ := strings.ToLower(args[0])
	//非读请求直接反复error
	if _, ok := RedisReadCmdMap[typ]; !ok {
		return nil, errors.New(fmt.Sprintf("错误: 不支持的命令 【 %s 】", commandStr))
	}
	c := &Command{
		args:    args,
		typ:     typ,
		cmdLine: str,
	}
	if typ == "select" {
		if len(args) != 2 {
			return nil, errors.New("error: CommandError")
		}
		n, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			return nil, errors.New("error: CommandError")
		}
		c.selectDB = int(n)
	}
	return c, nil
}

//redis 读命令
var RedisReadCmdMap = map[string]struct{}{

	//key 相关
	"esists":    {}, // 检查给定 key 是否存在:  EXISTS db
	"keys":      {}, //查找所有符合给定模式 pattern 的 key 。 不允许执行keys *
	"pttl":      {}, //以毫秒为单位返回 key 的剩余生存时间
	"randomkey": {}, //从当前数据库中随机返回(不删除)一个 key 。
	"ttl":       {}, //以秒为单位，返回给定 key 的剩余生存时间
	"type":      {}, //返回 key 所储存的值的类型。
	"scan":      {}, //迭代当前数据库中的数据库键

	//字符串相关
	"bitcount": {}, //计算给定字符串中，被设置为 1 的比特位的数量。
	"get":      {}, //返回 key 所关联的字符串值。
	"getbit":   {}, //对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
	"getrange": {}, //返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。
	"mget":     {}, //返回所有(一个或多个)给定 key 的值。
	"strlen":   {}, //返回 key 所储存的字符串值的长度。

	//hash 相关
	"hexists": {}, //查看哈希表 key 中，给定域 field 是否存在。
	"hget":    {}, //返回哈希表 key 中给定域 field 的值。
	"hgetall": {}, //返回哈希表 key 中，所有的域和值。
	"hkeys":   {}, //返回哈希表 key 中的所有域。
	"hlen":    {}, //返回哈希表 key 中域的数量。
	"hmget":   {}, //返回哈希表 key 中，一个或多个给定域的值。
	"hvals":   {}, //返回哈希表 key 中所有域的值。
	"hscan":   {},

	//list 相关
	"lindex": {}, //返回列表 key 中，下标为 index 的元素。 LINDEX mylist 0
	"llen":   {}, //返回列表 key 的长度。
	"lrange": {}, //返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。

	//set 集合相关
	"scard":       {}, //返回集合 key 的基数(集合中元素的数量)。SCARD tool
	"sdiff":       {}, //返回一个集合的全部成员，该集合是所有给定集合之间的差集。SDIFF peter's_movies joe's_movies
	"sinter":      {}, //返回一个集合的全部成员，该集合是所有给定集合的交集。SINTER group_1 group_2
	"sismember":   {}, //判断 member 元素是否集合 key 的成员。SISMEMBER joe's_movies "Fast Five"
	"smembers":    {}, //返回集合 key 中的s所有成员。
	"srandmember": {}, //如果命令执行时，只提供了 key 参数，那么返回集合中的一个随机元素。
	"sunion":      {}, //返回一个集合的全部成员，该集合是所有给定集合的并集。
	"sscan":       {}, //

	//有序集合相关元素
	"zcard":         {}, //返回有序集 key 的基数
	"zcount":        {}, //返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。  ZCOUNT salary 3000 5000
	"zrange":        {}, //返回有序集 key 中，指定区间内的成员。   ZRANGE salary 0 200000 WITHSCORES
	"zrangebyscore": {}, //返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。 有序集成员按 score 值递增(从小到大)次序排列。ZRANGEBYSCORE zset (1 5   返回所有符合条件 1 < score <= 5 的成员，
	"zrank":         {}, //返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。
	"zscore":        {}, //返回有序集 key 中，成员 member 的 score 值。
	"zscan":         {},
	//发布订阅
	//是否相关
	//脚本，
}

var HelpCommands = [][]string{
	{"APPEND", "key value", "KV"},
	{"AUTH", "password", "Server"},
	{"BITCOUNT", "key [start] [end]", "KV"},
	{"BITOP", "operation destkey key [key ...]", "KV"},
	{"BITPOS", "key bit [start] [end]", "KV"},
	{"BLPOP", "key [key ...] timeout", "List"},
	{"BRPOP", "key [key ...] timeout", "List"},
	{"CONFIG GET", "parameter", "Server"},
	{"CONFIG REWRITE", "-", "Server"},
	{"DECR", "key", "KV"},
	{"DECRBY", "key decrement", "KV"},
	{"DEL", "key [key ...]", "KV"},
	{"DUMP", "key", "KV"},
	{"ECHO", "message", "Server"},
	{"EVAL", "script numkeys key [key ...] arg [arg ...]", "Script"},
	{"EVALSHA", "sha1 numkeys key [key ...] arg [arg ...]", "Script"},
	{"EXISTS", "key", "KV"},
	{"EXPIRE", "key seconds", "KV"},
	{"EXPIREAT", "key timestamp", "KV"},
	{"FLUSHALL", "-", "Server"},
	{"FLUSHDB", "-", "Server"},
	{"FULLSYNC", "[NEW]", "Replication"},
	{"GET", "key", "KV"},
	{"GETBIT", "key offset", "KV"},
	{"GETRANGE", "key start end", "KV"},
	{"GETSET", " key value", "KV"},
	{"HCLEAR", "key", "Hash"},
	{"HDEL", "key field [field ...]", "Hash"},
	{"HDUMP", "key", "Hash"},
	{"HEXISTS", "key field", "Hash"},
	{"HEXPIRE", "key seconds", "Hash"},
	{"HEXPIREAT", "key timestamp", "Hash"},
	{"HGET", "key field", "Hash"},
	{"HGETALL", "key", "Hash"},
	{"HINCRBY", "key field increment", "Hash"},
	{"HKEYEXISTS", "key", "Hash"},
	{"HKEYS", "key", "Hash"},
	{"HLEN", "key", "Hash"},
	{"HMCLEAR", "key [key ...]", "Hash"},
	{"HMGET", "key field [field ...]", "Hash"},
	{"HMSET", "key field value [field value ...]", "Hash"},
	{"HPERSIST", "key", "Hash"},
	{"HSET", "key field value", "Hash"},
	{"HTTL", "key", "Hash"},
	{"HVALS", "key", "Hash"},
	{"INCR", "key", "KV"},
	{"INCRBY", "key increment", "KV"},
	{"INFO", "[section]", "Server"},
	{"KEYS", "pattern", "KV"},
	{"LCLEAR", "key", "List"},
	{"LDUMP", "key", "List"},
	{"LEXPIRE", "key seconds", "List"},
	{"LEXPIREAT", "key timestamp", "List"},
	{"LINDEX", "key index", "List"},
	{"LKEYEXISTS", "key", "List"},
	{"LLEN", "key", "List"},
	{"LMCLEAR", "key [key ...]", "List"},
	{"LPERSIST", "key", "List"},
	{"LPOP", "key", "List"},
	{"LPUSH", "key value [value ...]", "List"},
	{"LRANGE", "key start stop", "List"},
	{"LTTL", "key", "List"},
	{"MGET", "key [key ...]", "KV"},
	{"MSET", "key value [key value ...]", "KV"},
	{"PERSIST", "key", "KV"},
	{"PING", "-", "Server"},
	{"RESTORE", "key ttl value", "Server"},
	{"ROLE", "-", "Server"},
	{"RPOP", "key", "List"},
	{"RPUSH", "key value [value ...]", "List"},
	{"SADD", "key member [member ...]", "Set"},
	{"SCARD", "key", "Set"},
	{"SCLEAR", "key", "Set"},
	{"SCRIPT EXISTS", "script [script ...]", "Script"},
	{"SCRIPT FLUSH", "-", "Script"},
	{"SCRIPT LOAD", "script", "Script"},
	{"SDIFF", "key [key ...]", "Set"},
	{"SDIFFSTORE", "destination key [key ...]", "Set"},
	{"SDUMP", "key", "Set"},
	{"SELECT", "index", "Server"},
	{"SET", "key value", "KV"},
	{"SETBIT", "key offset value", "KV"},
	{"SETEX", "key seconds value", "KV"},
	{"SETNX", "key value", "KV"},
	{"SETRANGE", "key offset value", "KV"},
	{"SEXPIRE", "key seconds", "Set"},
	{"SEXPIREAT", "key timestamp", "Set"},
	{"SINTER", "key [key ...]", "Set"},
	{"SINTERSTORE", "destination key [key ...]", "Set"},
	{"SISMEMBER", "key member", "Set"},
	{"SKEYEXISTS", "key", "Set"},
	{"SLAVEOF", "host port [RESTART] [READONLY]", "Replication"},
	{"SMCLEAR", "key [key ...]", "Set"},
	{"SMEMBERS", "key", "Set"},
	{"SPERSIST", "key", "Set"},
	{"SREM", "key member [member ...]", "Set"},
	{"STRLEN", "key", "KV"},
	{"STTL", "key", "Set"},
	{"SUNION", "key [key ...]", "Set"},
	{"SUNIONSTORE", "destination key [key ...]", "Set"},
	{"SYNC", "logid", "Replication"},
	{"TIME", "-", "Server"},
	{"TTL", "key", "KV"},
	{"XHSCAN", "key cursor [MATCH match] [COUNT count] [ASC|DESC]", "Hash"},
	{"XLSORT", "key [BY pattern] [LIMIT offset count] [GET pattern [GET pattern ...]] [ASC|DESC] [ALPHA] [STORE destination]", "List"},
	{"XSCAN", "type cursor [MATCH match] [COUNT count] [ASC|DESC]", "Server"},
	{"XSSCAN", "key cursor [MATCH match] [COUNT count] [ASC|DESC]", "Set"},
	{"XSSORT", "key [BY pattern] [LIMIT offset count] [GET pattern [GET pattern ...]] [ASC|DESC] [ALPHA] [STORE destination]", "Set"},
	{"XZSCAN", "key cursor [MATCH match] [COUNT count] [ASC|DESC]", "ZSet"},
	{"XZSORT", "key [BY pattern] [LIMIT offset count] [GET pattern [GET pattern ...]] [ASC|DESC] [ALPHA] [STORE destination]", "ZSet"},
	{"ZADD", "key score member [score member ...]", "ZSet"},
	{"ZCARD", "key", "ZSet"},
	{"ZCLEAR", "key", "ZSet"},
	{"ZCOUNT", "key min max", "ZSet"},
	{"ZDUMP", "key", "ZSet"},
	{"ZEXPIRE", "key seconds", "ZSet"},
	{"ZEXPIREAT", "key timestamp", "ZSet"},
	{"ZINCRBY", "key increment member", "ZSet"},
	{"ZINTERSTORE", "destkey numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE SUM|MIN|MAX]", "ZSet"},
	{"ZKEYEXISTS", "key", "ZSet"},
	{"ZLEXCOUNT", "key min max", "ZSet"},
	{"ZMCLEAR", "key [key ...]", "ZSet"},
	{"ZPERSIST", "key", "ZSet"},
	{"ZRANGE", "key start stop [WITHSCORES]", "ZSet"},
	{"ZRANGEBYLEX", "key min max [LIMIT offset count]", "ZSet"},
	{"ZRANGEBYSCORE", "key min max [WITHSCORES] [LIMIT offset count]", "ZSet"},
	{"ZRANK", "key member", "ZSet"},
	{"ZREM", "key member [member ...]", "ZSet"},
	{"ZREMRANGBYLEX", "key min max", "ZSet"},
	{"ZREMRANGEBYRANK", "key start stop", "ZSet"},
	{"ZREMRANGEBYSCORE", "key min max", "ZSet"},
	{"ZREVRANGE", "key start stop [WITHSCORES]", "ZSet"},
	{"ZREVRANGEBYSCORE", "key max min  [WITHSCORES][LIMIT offset count]", "ZSet"},
	{"ZREVRANK", "key member", "ZSet"},
	{"ZSCORE", "key member", "ZSet"},
	{"ZTTL", "key", "ZSet"},
	{"ZUNIONSTORE", "destkey numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE SUM|MIN|MAX]", "ZSet"},
}
