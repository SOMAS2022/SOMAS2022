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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinReg(t *testing.T) {
	// assert := assert.New(t)
	X := [][]float64{
		{1, 1},
		{2, 2},
	}
	y := []float64{1, 2}

	w := FitLinReg(X, y)
	assert.InDelta(t, w[0], 0.5, 0.001)
	assert.InDelta(t, w[1], 0.5, 0.001)
}
