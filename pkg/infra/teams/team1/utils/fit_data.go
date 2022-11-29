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
	data := append(X, y)
	r.Train(regression.MakeDataPoints(data, 0)...)
	r.Run()

	// fmt.Printf("Regression formula:\n%v\n", r.Formula)
	return r.GetCoeffs()
}
