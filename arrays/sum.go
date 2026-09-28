package arrays

func Sum(array []int) (ret int) {
	ret = 0
	for i := range array {
		ret += array[i]
	}
	return
}

func SumAll(arrayList ...[]int) []int {
	var res []int
	for _, array := range arrayList {
		res = append(res, Sum(array))
	}
	return res
}

func SumAllTails(arrayList ...[]int) []int {
	var res []int
	for _, array := range arrayList {
		if len(array) > 0 {
			tails := array[1:]
			res = append(res, Sum(tails))
		} else {
			res = append(res, 0)
		}
	}
	return res
}
