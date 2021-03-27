package limiter

//限流器（固定大小）
type ConnLimiter struct {
	MaxCurrentConn int
	Bucket         chan struct{}
}

func NewConnLimiter(maxCurrentConn int) *ConnLimiter {
	return &ConnLimiter{
		MaxCurrentConn: maxCurrentConn,
		Bucket:         make(chan struct{}, maxCurrentConn),
	}
}

func (c *ConnLimiter) GetConn() bool {
	if len(c.Bucket) >= c.MaxCurrentConn {
		return false
	}
	c.Bucket <- struct{}{}
	return true
}

func (c *ConnLimiter) ReleaseConn() {
	<-c.Bucket
}
