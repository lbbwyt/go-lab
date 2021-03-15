package cache_utils

//双向链表节点结构体
type Node struct {
	Key        string
	Value      interface{}
	Prev, Next *Node
}

//使用最近最少使用算法实现的缓存，
//该算法在数据库中使用时， put元素时通常不是放在头部， 而是放在中间靠后的位子，比如说5/8的位置， 因为put到头部的数据不一定是热点数据， 这样会把
// 真正热点数据挤出， put数据时放到5/8位置可以保证前5/8的区域永远是热点数据
type LRUCache struct {
	head, tail *Node
	Keys       map[string]*Node
	Cap        int
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		Keys: make(map[string]*Node),
		Cap:  capacity,
	}
}

func (this *LRUCache) Get(key string) interface{} {
	if node, ok := this.Keys[key]; ok {
		this.Remove(node)
		this.Add(node)
		return node.Value
	}
	return nil
}

func (this *LRUCache) Put(key string, value interface{}) {
	if node, ok := this.Keys[key]; ok {
		node.Value = value
		this.Remove(node)
		this.Add(node)
		return
	} else {
		node := &Node{
			Key:   key,
			Value: value,
		}
		this.Keys[key] = node
		this.Add(node)
	}

	if len(this.Keys) > this.Cap {
		delete(this.Keys, key)
		this.Remove(this.tail)
	}
}

//双向链表头部add一个节点, 注意 tail节点的next和head节点的prev 都为空
func (this *LRUCache) Add(node *Node) {
	node.Prev = nil
	node.Next = this.head

	if this.head != nil {
		this.head.Prev = node
	}
	this.head = node
	if this.tail == nil {
		this.tail = node
		this.tail.Next = nil
	}
}

func (this *LRUCache) Remove(node *Node) {
	//如果删除的是头节点
	if node == this.head {
		// 头指针后移
		this.head = node.Next
		//释放内存
		node.Next = nil
	}

	// 如果删除的是尾节点
	if node == this.tail {
		// 尾指针前移
		this.tail = node.Prev
		//释放内存
		node.Prev.Next = nil
		node.Next = nil
		return
	}
	//删除节点
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}
