package utils

func CalculatePermutation(prefix []string, strArr []string, permutation *List) {
	n := len(strArr)
	if n == 0 {
		wrapper := Element{Ids: prefix}
		*permutation = append(*permutation, wrapper)
	} else {
		for i := 0; i < n; i++ {
			newPrefix := append(prefix, strArr[i])
			// ** DEEP COPY NEEDED **
			temp := make([]string, len(strArr), 10)
			copy(temp, strArr)
			newStrArr := append(temp[0:i], temp[i+1:]...)
			CalculatePermutation(newPrefix, newStrArr, permutation)
		}
	}
}
