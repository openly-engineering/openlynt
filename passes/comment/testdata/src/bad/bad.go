package main

import "fmt"

func main() {
	// having the magic W A N T keyword in comments for a linter that
	// lints comments is extremely weird. :-)
	// "" = regexp, `` = literal

	// TODO: bad // want "expected comment to contain /(.+)/"
	// second line
	fmt.Println("bad comments")
}

// TODO: this has http://, but it is one line. // want "expected at least 2 lines of comments around TODO, but have 1"
