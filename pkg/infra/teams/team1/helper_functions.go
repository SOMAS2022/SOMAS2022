package team1

func argmax(array []float64) int {

	maxIndex := 0

	for index, element := range array {
		if element > array[maxIndex] {
			maxIndex = index
		}
	}

	return maxIndex
}
