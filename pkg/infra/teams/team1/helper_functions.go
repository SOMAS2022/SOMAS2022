package team1

import "math"

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

// Ensures a float is between -1 and 1
func boundFloat(inputNumber float64) float64 {
	if inputNumber > 1.0 {
		return 1.0
	} else if inputNumber < -1.0 {
		return -1.0
	} else {
		return inputNumber
	}
}

// Ensures array values are between -1 and 1
func boundArray(inputArray [4]float64) [4]float64 {
	return [4]float64{
		boundFloat(inputArray[0]),
		boundFloat(inputArray[1]),
		boundFloat(inputArray[2]),
		boundFloat(inputArray[3]),
	}
}

// Add two arrays
func addArrays(A [4]float64, B [4]float64) [4]float64 {
	return [4]float64{
		A[0] + B[0],
		A[1] + B[1],
		A[2] + B[2],
		A[3] + B[3],
	}
}

func decayNumber(inputNumber float64) float64 {
	if inputNumber < 0 {
		return 0.70 * inputNumber
	} else {
		return 0.90 * inputNumber
	}
}

func decayArray(inputArray [4]float64) [4]float64 {
	return [4]float64{
		decayNumber(inputArray[0]),
		decayNumber(inputArray[1]),
		decayNumber(inputArray[2]),
		decayNumber(inputArray[3]),
	}
}

func softmax(inputArray [3]float64) [3]float64 {
	expValues := [3]float64{
		math.Exp(inputArray[0]),
		math.Exp(inputArray[1]),
		math.Exp(inputArray[2]),
	}

	// Sum exponential array
	sum := 0.0
	for i := 0; i < 3; i++ {
		sum += expValues[i]
	}

	// Divide each element in input array by sum
	for i := 0; i < 3; i++ {
		expValues[i] /= sum
	}

	return expValues
}

func makeIncremental(inputArray [3]float64) [3]float64 {
	var outputArray [3]float64

	outputArray[0] = inputArray[0]

	for i := 1; i < 3; i++ {
		outputArray[i] = outputArray[i-1] + inputArray[i]
	}

	return outputArray
}

// Make lowest value -1, highest 1 and everything else interpolation between
func normalise(array [3]float64) [3]float64 {
	max := max(array[:])
	min := min(array[:])

	var normArray [3]float64

	for index, value := range array {
		if value == max {
			normArray[index] = 1.0
		} else if value == min {
			normArray[index] = -1.0
		} else {
			// Interpolate between -1 and 1
			normArray[index] = (value - (max+min)/2) / (max - min)
		}
	}

	return normArray
}
