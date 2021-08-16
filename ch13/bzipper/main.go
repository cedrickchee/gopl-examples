// Bzipper reads input, bzip2-compresses it, and writes it out.
package main

import (
	"io"
	"log"
	"os"

	"gopl.io/ch13/bzip"
)

func main() {
	w := bzip.NewWriter(os.Stdout)
	if _, err := io.Copy(w, os.Stdin); err != nil {
		log.Fatalf("bzipper: %v\n", err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("bzipper: close: %v\n", err)
	}
}

/*
Run:
In the session below, we use `bzipper` to compress `/usr/share/dict/words`, the
system dictionary, from 938,848 bytes to 335,405 bytes--about a third of its
original size--then uncompress it with the system `bunzip2` command. The SHA256
hash is the same before and after, giving us confidence that the compressor is
working correctly.

$ go build gopl.io/ch13/bzipper
$ wc -c < /usr/share/dict/words
972398
$ sha256sum < /usr/share/dict/words
f6c94d35691b9c356f7e5072f94d23f127b168cf9b04f0f5b26e0cb1f6ef4414  -
$ ./bzipper < /usr/share/dict/words | wc -c
345883
$ ./bzipper < /usr/share/dict/words | bunzip2 | sha256sum
f6c94d35691b9c356f7e5072f94d23f127b168cf9b04f0f5b26e0cb1f6ef4414  -
*/
