package heap

// 实现大顶堆
type Heap struct {
	m   []int
	len int // 堆中元素数量
}

func NewIntHeap(items []int) *Heap {
	return buildHeap(items)
}

/*
 - 建堆，就是在原切片上操作，形成堆结构
 - 只要按照顺序，把切片下标为n/2到1的节点依次堆化，最后就会把整个切片堆化
*/
func buildHeap(m []int) *Heap {
	n := len(m) - 1
	for i := n / 2; i > 0; i-- {
		heapf(m, n, i)
	}

	return &Heap{m, n}
}

/*
 - 对下标为i的节点进行堆化， n表示堆的最后一个节点下标
 - 2i,2i+1
*/
func heapf(m []int, n, i int) {
	for {
		maxPos := i

		if 2*i <= n && m[2*i] > m[i] { // 确定父节点 parent(i) = floor((i - 1)/2)
			maxPos = 2 * i
		}

		if 2*i+1 <= n && m[2*i+1] > m[maxPos] { // 确定子树左右 2i + 1 或 2i + 2
			maxPos = 2*i + 1
		}

		if maxPos == i { // 如果i节点位置正确，则退出
			break
		}

		m[i], m[maxPos] = m[maxPos], m[i]
		i = maxPos
	}
}

func (h *Heap) Push(data int) {
	h.len++
	h.m = append(h.m, data) // 向切片尾部插入数据(推断出父节点下标为i/2)

	i := h.len
	for i/2 > 0 && h.m[i/2] < h.m[i] { // 自下而上的堆化
		h.m[i/2], h.m[i] = h.m[i], h.m[i/2]
		i = i / 2
	}
}

// 弹出堆顶元素，为防止出现数组空洞，需要把最后一个元素放入堆顶，然后从上到下堆化
func (h *Heap) Pop() int {
	if h.len < 1 {
		return -1
	}

	// 把最后一个元素给堆顶
	out := h.m[1]
	h.m[1] = h.m[h.len]
	h.m = h.m[:h.len]
	h.len--
	// 对堆顶节点进行堆化即可
	heapf(h.m, h.len, 1)
	// for i := h.len / 2; i > 0; i-- {
	// 	heapf(h.m, h.len, i)
	// }

	return out
}

func (h *Heap) Items() []int {
	return h.m
}
