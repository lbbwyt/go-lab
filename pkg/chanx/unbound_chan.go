package chanx

type T interface{}

type UnboundChan struct {
	In     chan<- T //chan for write
	Out    <-chan T //chan for read
	buffer *RingBuffer
}

// Len returns len of In plus len of Out plus len of buffer.
func (c UnboundChan) Len() int {
	return len(c.In) + c.buffer.Len() + len(c.Out)
}

// BufLen returns len of the buffer.
func (c UnboundChan) BufLen() int {
	return c.buffer.Len()
}

func NewUnboundedChan(initCapacity int) UnboundChan {
	return NewUnboundChanSize(initCapacity, initCapacity, initCapacity)
}

func NewUnboundChanSize(initInCapacity, initOutCapacity, initBufCapacity int) UnboundChan {
	in := make(chan T, initInCapacity)
	out := make(chan T, initOutCapacity)
	ch := UnboundChan{In: in, Out: out, buffer: NewRingBuffer(initBufCapacity)}

	go process(in, out, ch)

	return ch
}

func process(in chan T, out chan T, ch UnboundChan) {
	defer close(out)
loop:
	for {
		val, ok := <-in
		if !ok { // in is closed
			break loop
		}

		//out is not full
		select {
		case out <- val:
			continue
		default:

		}

		//out is full
		ch.buffer.Write(val)
		for !ch.buffer.IsEmpty() {
			select {
			case val, ok := <-in:
				if !ok {
					break loop
				}
				ch.buffer.Write(val)
			case out <- ch.buffer.Peek():
				ch.buffer.Pop()
				if ch.buffer.IsEmpty() {
					ch.buffer.Reset()
				}
			}

		}
	}

	//防止in 关闭之后， out关闭之前， buffer存在数据未写入out
	for !ch.buffer.IsEmpty() {
		out <- ch.buffer.Pop()
	}
	ch.buffer.Reset()
}
