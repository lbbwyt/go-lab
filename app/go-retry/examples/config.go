package examples

import "time"

type Config struct {
	attempts int
	delay    time.Duration
	maxDelay time.Duration
	backoff  bool
}

func newDefaultConfig() *Config {
	return &Config{
		attempts: 3,
		delay:    time.Second,
		maxDelay: time.Second * 10,
	}
}

type Option func(*Config)

// 尝试次数
func WithAttempts(attempts int) Option {
	return func(c *Config) {
		c.attempts = attempts
	}
}

// Dealy 间隔
func WithDelay(delay time.Duration) Option {
	return func(c *Config) {
		c.delay = delay
	}
}

func WithBackoff(maxDelay time.Duration) Option {
	return func(c *Config) {
		c.maxDelay = maxDelay
		c.backoff = true
	}
}
