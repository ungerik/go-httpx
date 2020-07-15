package httperr

import "fmt"

func ExampleDontLog() {
	fmt.Println(ShouldLog(BadRequest))
	fmt.Println(ShouldLog(DontLog(BadRequest)))

	// Output:
	// true
	// false
}
