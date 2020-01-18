package main

var x_y_z = "abc"

type x struct {
	Bad_Var string
	OKVar string
}

func main() {
	x_y_z = "efg"

	// these are ignored and should be handled by a struct-type rule
	y := &x{}
	y.Bad_Var = "no"
	y.OKVar = "lol"

	x_y_z, abc, d_e_f = "efg", "xyz", "abc"
}
