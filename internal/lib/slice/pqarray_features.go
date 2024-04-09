package slice

import "github.com/lib/pq"

func ConvertToIntSlice(pqArr pq.Int32Array) []int {
	resultArray := make([]int, 0)
	for _, number := range pqArr {
		resultArray = append(resultArray, int(number))
	}

	return resultArray
}
