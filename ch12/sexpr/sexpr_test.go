package sexpr

import "testing"

// Test verifies that encoding and decoding a complex data value
// produces an equal result.
//
// The test does not make direct assertions about the encoded output
// because the output depends on map iteration order, which is
// nondeterministic.  The output of the t.Log statements can be
// inspected by running the test with the -v flag:
//
// 	$ go test -v gopl.io/ch12/sexpr
//
func Test(t *testing.T) {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	// Encode it
	data, err := Marshal(strangelove)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	t.Logf("Marshal() = %s\n", data)

	// Pretty-print it:
	data, err = MarshalIndent(strangelove)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("MarshalIndent() = \n%s\n", data)
}

/*
Output:

$ go test -v gopl.io/ch12/sexpr

=== RUN   Test
    sexpr_test.go:48: Marshal() = ((Title "Dr. Strangelove") (Subtitle "How I Learned to Stop Worrying and Love the Bomb") (Year 1964) (Actor (("Dr. Strangelove" "Peter Sellers") ("Grp. Capt. Lionel Mandrake" "Peter Sellers") ("Pres. Merkin Muffley" "Peter Sellers") ("Gen. Buck Turgidson" "George C. Scott") ("Brig. Gen. Jack D. Ripper" "Sterling Hayden") ("Maj. T.J. \"King\" Kong" "Slim Pickens"))) (Oscars ("Best Actor (Nomin.)" "Best Adapted Screenplay (Nomin.)" "Best Director (Nomin.)" "Best Picture (Nomin.)")) (Sequel nil))
    sexpr_test.go:55: MarshalIndent() =
        ((Title "Dr. Strangelove")
         (Subtitle "How I Learned to Stop Worrying and Love the Bomb") (Year 1964)
         (Actor
          (("Maj. T.J. \"King\" Kong" "Slim Pickens")
           ("Dr. Strangelove" "Peter Sellers")
           ("Grp. Capt. Lionel Mandrake" "Peter Sellers")
           ("Pres. Merkin Muffley" "Peter Sellers")
           ("Gen. Buck Turgidson" "George C. Scott")
           ("Brig. Gen. Jack D. Ripper" "Sterling Hayden")))
         (Oscars
          ("Best Actor (Nomin.)" "Best Adapted Screenplay (Nomin.)"
           "Best Director (Nomin.)" "Best Picture (Nomin.)")) (Sequel nil))
--- PASS: Test (0.00s)
PASS
ok      gopl.io/ch12/sexpr      0.002s
*/
