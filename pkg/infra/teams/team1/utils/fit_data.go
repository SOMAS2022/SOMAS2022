/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/
package utils

import (
	"fmt"

	"github.com/sajari/regression"
)

/*
Function fits least squares linear regression solution:
	argmin||Xw-y||_2^2
where X is the agent state and y is the reward for that state
Inputs:
	X (State - matrix with agent states in each rows)
	y (Reward - integer corresponding to reward given to that state)
Output:
	w (Optimal weights)
*/
func FitLinReg(X [][]float64, y []float64) []float64 {

	r := new(regression.Regression)

	// Concatenate observation array
	// var data [][]float64
	// for i, row := range X {
	// 	new_row := append(row, y[i])
	// 	data = append(data, new_row)
	// }

	fmt.Println(X)
	datapoints := make([]*regression.DataPoint, 0, len(a))
	for i, row := range X {
		datapoints = append(datapoints, regression.DataPoint(y[i], row))
	}
	r.Train(datapoints...)
	r.Run()

	fmt.Printf("Regression formula:\n%v\n", r.Formula)
	return r.GetCoeffs()
}
