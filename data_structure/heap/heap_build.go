package heap

/*
 - AdjustTraingle 如果只是修改slice里的元素，不需要传slice的指针；
 - 如果要往slice里append或让slice指向新的子切片，则需要传slice指针
*/
func AdjustTraingle(arr []int, parent int) {
	left := 2*parent + 1
	if left >= len(arr) {
		return
	}
	minIndex := parent
	minValue := arr[minIndex]
	if arr[left] < minValue {
		minValue = arr[left]
		minIndex = left
	}

	right := 2*parent + 2
	if right < len(arr) {
		if arr[right] < minValue {
			minValue = arr[right]
			minIndex = right
		}
	}

	if minIndex != parent {
		arr[minIndex], arr[parent] = arr[parent], arr[minIndex]
		// 递归 每当有元素调整下来时，要对以它为父节点的三角形区域进行调整
		AdjustTraingle(arr, minIndex) 
	}
}

func ReverseAdjust(arr []int) {
	n := len(arr)
	if n <= 1 {
		return
	}
	
	lastIndex := (n + 1) / 2 * 2
	//逆序检查每一个三角形区域   left(i)  = 2i + 1   &&   right(i)  = 2i + 2
	for i := lastIndex; i > 0; i -= 2 {
		right := i
		parent := (right - 1) / 2
		AdjustTraingle(arr, parent)
	}
}
