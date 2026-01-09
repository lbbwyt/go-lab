package scheduler

import "fmt"

type DAG struct {
	adj   map[string][]string // 邻接表
	indeg map[string]int      //入度
}

func NewDAG() *DAG { return &DAG{adj: make(map[string][]string), indeg: make(map[string]int)} }

func (d *DAG) AddEdge(from, to string) {

	// 自动补全缺失顶点
	d.AddVertex(from)
	d.AddVertex(to)

	d.adj[from] = append(d.adj[from], to)
	d.indeg[to]++
}

// 添加顶点（如果已存在则忽略）
func (g *DAG) AddVertex(v string) {
	if _, ok := g.adj[v]; !ok {
		g.adj[v] = nil
	}
	if _, ok := g.indeg[v]; !ok {
		g.indeg[v] = 0
	}
}

// 拓扑排序（Kahn/BFS 法）,O(V+E) 时间复杂度，O(V) 空间复杂度。
func (d *DAG) Topological() ([]string, error) {
	q := []string{}
	//统计每个顶点的入度。把入度 0 的顶点入队。
	for id, d := range d.indeg {
		if d == 0 {
			q = append(q, id)
		}
	}
	order := []string{}

	//每弹出一个顶点，将其邻接顶点的入度减 1；若减到 0 则入队。
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		order = append(order, cur)
		for _, nxt := range d.adj[cur] {
			d.indeg[nxt]--
			if d.indeg[nxt] == 0 {
				q = append(q, nxt)
			}
		}
	}
	//若最终弹出的顶点数 不等于 总顶点数，说明有环。
	if len(order) != len(d.indeg) {
		return nil, fmt.Errorf("cycle detected")
	}
	return order, nil
}
