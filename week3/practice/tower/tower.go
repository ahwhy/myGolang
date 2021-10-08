package main

import "fmt"

func main() {
	tower("A", "B", "C", 6)
	// fmt.Println(num)
}

func tower(a, b, c string, layer int) {
	if layer <= 0 {
		return
	}

	if layer == 1 {
		fmt.Printf("%s - > %s\n", a, c)
		return
	}

	tower(a, c, b, layer-1)

	fmt.Printf("%s - > %s\n", a, c)

	tower(b, a, c, layer-1)
}
