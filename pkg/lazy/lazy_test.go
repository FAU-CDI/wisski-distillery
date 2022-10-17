package lazy

import "fmt"

func ExampleLazy() {

	var lazy Lazy[int]

	// the first invocation to lazy will be called and set the value
	fmt.Println(lazy.Get(func() int { return 42 }))

	// the second invocation will not call init again, using the first value
	fmt.Println(lazy.Get(func() int { return 43 }))

	// Set can be used to set a specific value
	lazy.Set(0)
	fmt.Println(lazy.Get(func() int { panic("never called") }))

	// Output: 42
	// 42
	// 0
}

func ExampleLazy_nil() {
	var lazy Lazy[int]

	// passing nil as the initialization function causes the zero value to be set
	fmt.Println(lazy.Get(nil))

	// Output: 0
}
