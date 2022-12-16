package math_test

import (
	"testing"

	"infra/game/math"
)

func TestCalculateMonsterHealth(t *testing.T) {
	t.Parallel()

	type args struct {
		N  uint
		ST uint
		L  uint
		CL uint
	}
	tests := []struct {
		name    string
		args    args
		wantmin uint
		wantmax uint
	}{
		{
			name: "Case 1",
			args: args{N: 100, ST: 2000, L: 60, CL: 1},
			// min for delta = 0.8 and max for delta = 1.2
			wantmin: 1423,
			wantmax: 2134,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := math.CalculateMonsterHealth(tt.args.N, tt.args.ST, tt.args.L, tt.args.CL); !(tt.wantmin <= got && got <= tt.wantmax) {
				t.Errorf("CalculateMonsterHealth() = %v, wanted between %v and %v", got, tt.wantmin, tt.wantmax)
			}
		})
	}
}

func TestCalculateMonsterDamage(t *testing.T) {
	t.Parallel()

	type args struct {
		N  uint
		HP uint
		ST uint
		TH float32
		L  uint
		CL uint
	}
	tests := []struct {
		name    string
		args    args
		wantmin uint
		wantmax uint
	}{
		{
			name: "Case 1",
			args: args{N: 100, HP: 1000, ST: 2000, TH: 0.1, L: 60, CL: 1},
			// min for delta = 0.8 and max for delta = 1.2
			wantmin: 1922,
			wantmax: 2848,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := math.CalculateMonsterDamage(tt.args.N, tt.args.HP, tt.args.ST, tt.args.TH, tt.args.L, tt.args.CL); !(tt.wantmin <= got && got <= tt.wantmax) {
				t.Errorf("CalculateMonsterDamage() = %v,  wanted between %v and %v", got, tt.wantmin, tt.wantmax)
			}
		})
	}
}
