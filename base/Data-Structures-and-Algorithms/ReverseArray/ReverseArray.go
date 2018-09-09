package ReverseArray

func ReverseArray(a []int) {
	i := 0
	u := len(a) - 1
	for i < u {
		a[i], a[u] = a[u], a[i]
		i, u = i+1, u-1
	}
}
