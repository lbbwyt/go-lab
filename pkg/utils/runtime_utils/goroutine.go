package runtime_utils

import (
	"fmt"
	"github.com/transaction-wg/seata-golang/pkg/util/log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var debugIgnoreStdout = false

// GoWithRecover wraps a `go func()` with recover()

//runtime_utils.GoWithRecover(func() {
//	session.Close()
//}, nil)
func GoWithRecover(handler func(), recoverHandler func(r interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// TODO: log
				if !debugIgnoreStdout {
					fmt.Fprintf(os.Stderr, "%s goroutine panic: %v\n%s\n",
						time.Now(), r, string(debug.Stack()))
				}
				if recoverHandler != nil {
					go func() {
						defer func() {
							if p := recover(); p != nil {
								if !debugIgnoreStdout {
									fmt.Fprintf(os.Stderr, "recover goroutine panic:%v\n%s\n", p, string(debug.Stack()))
								}
							}
						}()
						recoverHandler(r)
					}()
				}
			}
		}()
		handler()
	}()
}

func ExecWithRecover(handler func(), recoverHandler func(r interface{})) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "%s goroutine panic: %v\n%s\n",
				time.Now(), r, string(debug.Stack()))
			if recoverHandler != nil {
				go func() {
					defer func() {
						if p := recover(); p != nil {
							fmt.Fprintf(os.Stderr, "recover goroutine panic:%v\n%s\n", p, string(debug.Stack()))
						}
					}()
					recoverHandler(r)
				}()
			}
		}
	}()
	handler()
}

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		log.Warn(fmt.Sprintf("cannot get goroutine id: %v", err))
		return 0
	}
	return id
}
