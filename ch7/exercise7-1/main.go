// Exercise 7.1: Using the ideas from `ByteCounter`, implement counters for
// words and for lines. You will find `bufio.ScanWords` useful.
package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // convert int to ByteCounter
	return len(p), nil
}

type LineCounter int

func (l *LineCounter) Write(p []byte) (int, error) {
	// for _, b := range p {
	// 	if b == '\n' {
	// 		l.lines++
	// 	}
	// }

	// A better way
	b := bytes.NewBuffer(p)
	s := bufio.NewScanner(b)
	for s.Scan() {
		*l++
	}
	return len(p), s.Err()
}

func (l *LineCounter) N() int {
	return int(*l)
}

func (l *LineCounter) String() string {
	return fmt.Sprintf("%d line(s)", *l)
}

func (l *LineCounter) Reset() {
	*l = 0
}

type WordCounter int

func (w *WordCounter) Write(p []byte) (int, error) {
	n := 0
	b := bytes.NewBuffer(p)
	s := bufio.NewScanner(b)
	s.Split(bufio.ScanWords) // break the input into words.
	for s.Scan() {
		n++
	}
	*w = WordCounter(n)

	return len(p), s.Err()
}

func (w *WordCounter) N() int {
	return int(*w)
}

func (w *WordCounter) String() string {
	return fmt.Sprintf("%d word(s)", *w)
}

func (w *WordCounter) Reset() {
	*w = 0
}

func main() {
	// txt := `hello world. how are you
	// test one.
	// test two.
	// lorem ipsum.`
	txt := "hello world. how are you\ntest one.\ntest two.\nlorem ipsum."

	// Line count
	var l LineCounter
	n, err := l.Write([]byte(txt))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("char count =", n)
	fmt.Printf("line count = %d\n", l.N()) // "4"

	l.Reset()
	fmt.Fprint(&l, txt)
	fmt.Printf("line count = %d\n", l.N()) // "4"
	fmt.Printf("%v\n", &l)                 // "4 line(s)"

	// Count words
	var w WordCounter
	n, err = w.Write([]byte(txt))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("char count =", n)
	fmt.Printf("word count = %d\n", w.N()) // "11"

	w.Reset()
	fmt.Fprint(&w, txt)
	fmt.Printf("word count = %d\n", w.N()) // "11"
	fmt.Printf("%v\n", &w)                 // "11 word(s)"
}
