package prase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

//解析redis 命令并返回执行结果：
type RedisExecutor struct {
	Addr     string
	UserName string
	Password string
	Cmd      *Command
	Client   *redis.Client
}

func NewRedisExecutor(addr, userName, password string) *RedisExecutor {
	c := &RedisExecutor{
		Addr:     addr,
		UserName: userName,
		Password: password,
		Cmd:      nil,
		Client:   nil,
	}
	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       0,
	})
	c.Client = client
	return c
}

// 执行redis命令
func (e *RedisExecutor) Do(cmd *Command) (*ExecuteResult, error) {
	e.Cmd = cmd
	vargs := make([]interface{}, 0, len(cmd.args))
	for _, arg := range cmd.args {
		vargs = append(vargs, arg)
	}
	res, err := e.Client.Do(context.Background(), vargs...).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
		return e.formatResult("(nil)", cmd.cmdLine), nil
	} else {
		return e.formatResult(res, cmd.cmdLine), nil
	}
}

// 执行redis命令
func (e *RedisExecutor) DoRaw(cmd *Command) (interface{}, error) {
	e.Cmd = cmd
	vargs := make([]interface{}, 0, len(cmd.args))
	for _, arg := range cmd.args {
		vargs = append(vargs, arg)
	}
	res, err := e.Client.Do(context.Background(), vargs...).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
		return e.format("(nil)\n", cmd.cmdLine), nil
	} else {
		return e.format(res, cmd.cmdLine), nil
	}
}

//执行resdis 命令， 并返回结果， 错误时， 返回错误信息字符串
func (e *RedisExecutor) Execute(cdmStr string) string {

	log.Info(fmt.Sprintf("命令:%s", cdmStr))

	cmd, err := PraseRedisReadCmd(cdmStr)
	if err != nil {
		return err.Error()
	}
	res, err := e.DoRaw(cmd)
	if err != nil {
		return err.Error()
	}
	return printRawReply(0, res)
}

func jsonFormat(data interface{}) string {
	s, _ := json.Marshal(data)

	return string(s)
}

func printRawReply(level int, reply interface{}) string {

	res := ""

	switch reply := reply.(type) {
	case int64:
		res = fmt.Sprintf("%d", reply)
	case string:
		res = fmt.Sprintf("%s", reply)
	case []byte:
		res = fmt.Sprintf("%s", reply)
	case nil:
		// do nothing
	case redis.Error:
		res = fmt.Sprintf("%s\n", reply.Error())
	case []interface{}:
		for i, v := range reply {
			if i != 0 {
				res = fmt.Sprintf("%s", strings.Repeat(" ", level*4))
			}

			printRawReply(level+1, v)
			if i != len(reply)-1 {
				res = fmt.Sprintf("\n")
			}
		}
	default:
		res = fmt.Sprintf("Unknown reply type: %+v", reply)
	}

	return res
}

func (e *RedisExecutor) format(res interface{}, cmdLine string) string {
	var output string
	if reflect.TypeOf(res).Kind() == reflect.Slice {
		buff := bytes.NewBuffer(make([]byte, 512))
		s := reflect.ValueOf(res)
		for i := 0; i < s.Len(); i++ {
			if i > 1000 {
				buff.WriteString(fmt.Sprintf("%d) ...\n", i+1))
				break
			}
			ele := s.Index(i)
			buff.WriteString(fmt.Sprintf("%d) %v\n", i+1, ele.Interface()))
		}
		output = buff.String()
	} else {
		output = fmt.Sprintf("%v", res)
	}
	return fmt.Sprintf("%s\n", output)
}

func (e *RedisExecutor) formatResult(res interface{}, cmdLine string) *ExecuteResult {
	var output string
	if reflect.TypeOf(res).Kind() == reflect.Slice {
		buff := bytes.NewBuffer(make([]byte, 512))
		s := reflect.ValueOf(res)
		for i := 0; i < s.Len(); i++ {
			if i > 1000 {
				buff.WriteString(fmt.Sprintf("%d) ...\n", i+1))
				break
			}
			ele := s.Index(i)
			buff.WriteString(fmt.Sprintf("%d) %v\n", i+1, ele.Interface()))
		}
		output = buff.String()
	} else {
		output = fmt.Sprintf("%v", res)
	}
	return &ExecuteResult{
		Data: output,
		Cmd:  cmdLine,
	}
}

type ExecuteResult struct {
	Data  string       `json:"data"`
	Error *ResultError `json:"error"`
	Cmd   string       `json:"cmd"`
}

type ResultError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
