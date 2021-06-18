package cache_middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
	"net/http"
	"sync"
	"time"
)

type Options struct {
	CacheStore             CacheStore
	CacheDuration          time.Duration
	DisableSingleFlight    bool
	SingleFlightForgetTime time.Duration //DisableSingleFlight为false时有效
	Logger                 Logger
}

type Logger interface {
	Printf(format string, args ...interface{})
}

//return cacheKey needCache
type KeyGenerator func(c *gin.Context) (string, bool)

func Cache(keyGenerator KeyGenerator, options Options) gin.HandlerFunc {
	if options.CacheStore == nil {
		panic(ErrCacheMiss)
	}

	cacheHelper := newCacheHelper(options)
	return func(c *gin.Context) {
		cacheKey, needCache := keyGenerator(c)
		if !needCache {
			c.Next()
			return
		}
		//read cache
		{
			respCache := cacheHelper.getResponseCache()
			defer cacheHelper.putResponseCache(respCache)
			err := options.CacheStore.Get(cacheKey, &respCache)
			if err != nil {
				//直接返回缓存数据
				cacheHelper.respondWithCache(c, respCache)
			}
			if err != ErrCacheMiss {
				if options.Logger != nil {
					options.Logger.Printf("get cache: %v", err)
				}
			}
		}

		cacheWriter := &responseCacheWriter{}
		cacheWriter.reset(c.Writer)
		c.Writer = cacheWriter
		respCache := &responseCache{}
		if options.DisableSingleFlight {
			c.Next()
			//将后台logic返回结果写入respCache
			respCache.fill(cacheWriter)
		} else {
			handled := false
			// use singleflight to avoid Hotspot Invalid
			rawCacheWriter, _, _ := cacheHelper.sfGroup.Do(cacheKey, func() (interface{}, error) {
				if options.SingleFlightForgetTime > 0 {
					go func() {
						time.Sleep(options.SingleFlightForgetTime)
						cacheHelper.sfGroup.Forget(cacheKey)
					}()
				}

				c.Next()

				handled = true
				return cacheWriter, nil
			})
			cacheWriter = rawCacheWriter.(*responseCacheWriter)
			respCache.fill(cacheWriter)
			if !handled {
				cacheHelper.respondWithCache(c, respCache)
			}
		}

		//写缓存
		if err := options.CacheStore.Set(cacheKey, respCache, options.CacheDuration); err != nil {
			if options.Logger != nil {
				options.Logger.Printf("set cache error: %v", err)
			}
		}

	}
}

// CacheByURI a shortcut function for caching response with uri
func CacheByURI(options Options) gin.HandlerFunc {
	return Cache(
		func(c *gin.Context) (string, bool) {
			return c.Request.RequestURI, true
		},
		options,
	)
}

// CacheByPath a shortcut function for caching response with url path, discard the query params
func CacheByPath(options Options) gin.HandlerFunc {
	return Cache(
		func(c *gin.Context) (string, bool) {
			return c.Request.URL.Path, true
		},
		options,
	)
}

// responseCacheWriter
type responseCacheWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *responseCacheWriter) reset(writer gin.ResponseWriter) {
	w.body.Reset()
	w.ResponseWriter = writer
}

type cacheHelper struct {
	sfGroup           singleflight.Group
	responseCachePool *sync.Pool
	options           Options
}

func (h cacheHelper) getResponseCache() *responseCache {
	respCache := h.responseCachePool.Get().(*responseCache)
	respCache.reset()

	return respCache
}

func (h cacheHelper) putResponseCache(cache *responseCache) {
	h.responseCachePool.Put(cache)
}

func (h cacheHelper) respondWithCache(c *gin.Context, respCache *responseCache) {
	c.Writer.WriteHeader(respCache.Status)
	for k, v := range respCache.Header {
		for _, item := range v {
			c.Writer.Header().Set(k, item)
		}
	}
	if _, err := c.Writer.Write(respCache.Data); err != nil {
		if h.options.Logger != nil {
			h.options.Logger.Printf("write response error: %v", err)
		}
	}

	// abort handler chain and return directly
	c.Abort()
}

type responseCache struct {
	Status int
	Header http.Header
	Data   []byte
}

func (c responseCache) reset() {
	c.Data = c.Data[0:0]
	c.Header = make(http.Header)
}

func (c responseCache) fill(cacheWriter *responseCacheWriter) {
	c.Status = cacheWriter.Status()
	c.Data = cacheWriter.body.Bytes()
	c.Header = make(http.Header, len(cacheWriter.Header()))

	for key, value := range cacheWriter.Header() {
		c.Header[key] = value
	}
}

func newCacheHelper(options Options) *cacheHelper {
	return &cacheHelper{
		sfGroup:           singleflight.Group{},
		responseCachePool: newResponseCachePool(),
		options:           Options{},
	}
}

func newResponseCachePool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return &responseCache{
				Status: 0,
				Header: make(http.Header),
				Data:   nil,
			}
		},
	}
}
