package main

import "fmt"

func main() {
	// FIXME: this todo is bad
	fmt.Println("bad comment file")

}

// FIXME(@chrsm): implement `x` such that `x` doesn't not do something, because not
// doing anything is not as useful as not doing nothing.
// see https://github.com/openly-engineering/openlynt/issues/1
func x() {
	// FIXME: bad comment
}

// TODO(@chrsm): todo without link!! oh no

// and another TODO without a link!! oh no!O!o1o1o1

// https://github.com/openly-engineering/openlynt/issues/1 badly formatted TODO(@chrsm)

/* and a multiline one TODO
 * yep
 */
