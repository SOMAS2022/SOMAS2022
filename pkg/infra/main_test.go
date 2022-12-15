package main

import "testing"

func TestFoo(t *testing.T) {
	t.Parallel()
	main()
}

func TestMain(m *testing.M) {
	m.Run()
}
