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

func argmin(array []float64) int {
	minIndex := 0

	for index, element := range array {
		if element < array[minIndex] {
			minIndex = index
		}
	}

	return minIndex
}

func max(array []float64) float64 {
	if len(array) == 0 {
		return 0
	}

	max := array[0]

	for _, element := range array {
		if element > max {
			max = element
		}
	}

	return max
}

func min(array []float64) float64 {
	if len(array) == 0 {
		return 0
	}

	min := array[0]

	for _, element := range array {
		if element < min {
			min = element
		}
	}

	return min
}
