// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Commencing countdown.")
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-tick
	}
	launch()
}

func launch() {
	fmt.Println("Lift off!")
}

/*
Run:
$ go run gopl.io/ch8/countdown1

Output:
Commencing countdown.
10
9
8
7
6
5
4
3
2
1
Lift off!
*/
