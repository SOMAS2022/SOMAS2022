/*******************************************************
* Copyright (C) 2022 Team 1 @ SOMAS2022
*
* This file is part of SOMAS2022.
*
* This file or its contents can not be copied and/or used
* without the express permission of Team 1, SOMAS2022
*******************************************************/
package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinReg(t *testing.T) {
	// Test w=[0,1,1] (first term is bias)
	X := [][]float64{
		{1, 3},
		{2, 4},
		{5, 8},
	}
	y := []float64{4, 6, 13}

	w := FitLinReg(X, y)
	assert.InDelta(t, w[0], 0, 0.001)
	assert.InDelta(t, w[1], 1, 0.001)
	assert.InDelta(t, w[2], 1, 0.001)
}
