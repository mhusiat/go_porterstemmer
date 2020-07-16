// suffix-based stemmer as in Porter (1980)

package porterstemmer

import (
	"strings"
)

// string struct that defines most of suffix changes and removals
// (still, some more complex ones will be defined inside functions)
var suffixAggregate = struct {
	pluralsPast  [][2]string
	changeSuffix [][2]string
	shortRemove  [][2]string
	longRemove   [][2]string
}{
	pluralsPast: [][2]string{
		{"sses", "ss"},
		{"ies", "i"},
		{"ss", "ss"},
		{"s", ""},
	},
	changeSuffix: [][2]string{
		{"ational", "atione"},
		{"tional", "tion"},
		{"enci", "ence"},
		{"anci", "ance"},
		{"izer", "ize"},
		{"bli", "ble"},
		{"alli", "al"},
		{"entli", "ent"},
		{"eli", "e"},
		{"ousli", "ous"},
		{"ization", "ize"},
		{"ation", "ate"},
		{"ator", "ate"},
		{"alism", "al"},
		{"iveness", "ive"},
		{"fulness", "ful"},
		{"ousness", "ous"},
		{"aliti", "al"},
		{"iviti", "ive"},
		{"biliti", "ble"},
		{"logi", "log"},
	},
	shortRemove: [][2]string{
		{"icate", "ic"},
		{"ative", ""},
		{"alize", "al"},
		{"iciti", "ic"},
		{"ical", "ic"},
		{"ful", ""},
		{"ness", ""},
	},
	longRemove: [][2]string{
		{"al", ""},
		{"ance", ""},
		{"ence", ""},
		{"er", ""},
		{"ic", ""},
		{"able", ""},
		{"ible", ""},
		{"ant", ""},
		{"ement", ""},
		{"ment", ""},
		{"ent", ""},
		{"ou", ""},
		{"ism", ""},
		{"ate", ""},
		{"iti", ""},
		{"ous", ""},
		{"ive", ""},
		{"ize", ""},
	},
}

// takes whole word and location of letter due to 'y'
// (alternative is to supply letter and letter before it)
func vowel(word string, loc int) bool {

	vowels := []byte{'a', 'e', 'o', 'u', 'i'}
	for _, q1 := range vowels {
		if word[loc] == q1 {
			return true
		}
	}

	// 'y' can be consonant or vowel depending on
	// previous letter eg. lynx vs yacht; account for that
	if loc != 0 && word[loc] == 'y' {
		for _, q2 := range vowels {
			if word[loc-1] == q2 {
				return false
			}
		}
		return true
	}
	return false
}

// counts number of sylabs (or rather vowels-consonants blocks)
// that occur after each other; doesn't count first consonant
func sylabs(word string) int {
	num := 0
	previousSylab := false
	for pos := range word {
		response := vowel(word, pos)
		if response == true && previousSylab == false {
			num += 1
		}
		previousSylab = response
	}

	// last vowel is not counted in as a sylab; correction needed
	if previousSylab == true {
		num -= 1
	}

	return num

}

// simple function that checks if the word has ANY vowels
func hasVowels(word string) bool {
	for pos := range word {
		if vowel(word, pos) {
			return true
		}
	}
	return false
}

// checker for specific CVC (consonant-vowel-consonant) ending
// of the string/stem; has additional rules included
// eg. of such endings: row, bit, cut
func cvc(ending string) bool {
	if vowel(ending, 0) == false && vowel(ending, 2) == false &&
		vowel(ending, 1) == true {
		for _, l := range []byte{'y', 'w', 'x'} {
			if ending[2] == l {
				return false
			}
		}
		return true
	}
	return false
}

// this is a post-processing rule after first step in the
// algorithm that cleans some suffix endings; necessary in few cases
func postRule(word string) string {
	if sylabs(word) == 1 && len(word) >= 3 && cvc(word[len(word)-3:len(word)]) {
		return word + "e"
	}
	for _, i := range []string{"at", "bl", "iz"} {
		if strings.HasSuffix(word, i) {
			return word + "e"
		}
	}
	if len(word) >= 2 && word[len(word)-1] == word[len(word)-2] &&
		vowel(word, len(word)-1) == false {
		for _, i := range []byte{'l', 'z', 's'} {
			if word[len(word)-1] == i {
				return word
			}
		}
		return word[:len(word)-1]
	}
	return word
}

// an iterator that checks for suffices and replaces them if the rule
// does match. uses the global struct from the top
func stripSuffix(word string, nrSylab int, suffices [][2]string) string {
	for _, pair := range suffices {
		if strings.HasSuffix(word, pair[0]) {
			if sylabs(word[:len(word)-len(pair[0])]) >= nrSylab {
				return word[:len(word)-len(pair[0])] + pair[1]
			}
			return word
		}
	}
	return word
}

// removes plurals and past form and deals with initial suffices
func pluralsPastStem(word string) string {

	word = stripSuffix(word, 0, suffixAggregate.pluralsPast)
	switch {

	case strings.HasSuffix(word, "eed"):
		if sylabs(word[:len(word)-3]) > 0 {
			return word[:len(word)-1]
		}
		return word

	case strings.HasSuffix(word, "ed"):
		if hasVowels(word[:len(word)-2]) {
			return postRule(word[:len(word)-2])
		}
		return word

	case strings.HasSuffix(word, "ing"):
		if hasVowels(word[:len(word)-3]) {
			return postRule(word[:len(word)-3])
		}
		return word
	default:
		return word
	}
}

// simple 'y' to 'i' suffix change function
func yToI(word string) string {
	if strings.HasSuffix(word, "y") && hasVowels(word[:len(word)-1]) {
		return word[:len(word)-1] + "i"
	}
	return word

}

// applies another set of rules based on the iterator and suffix list
func changeSuffixStem(word string) string {
	stem := stripSuffix(word, 1, suffixAggregate.changeSuffix)
	return stem
}

// applies another set of rules based on the iterator and suffix list
func shortRemoveStem(word string) string {
	stem := stripSuffix(word, 1, suffixAggregate.shortRemove)
	return stem

}

// applies another set of rules based on the iterator and suffix list
func longRemoveStem(word string) string {
	if strings.HasSuffix(word, "ion") {
		if sylabs(word[:len(word)-3]) > 1 && (word[len(word)-4] == 's' ||
			word[len(word)-4] == 't') {
			return word[:len(word)-3]
		}
		return word
	}
	stem := stripSuffix(word, 2, suffixAggregate.longRemove)
	return stem

}

// post-stemming cleanup rules that remove some 'e' from endings
func additionalCleanup(word string) string {
	if strings.HasSuffix(word, "e") {
		if (sylabs(word[:len(word)-1]) > 1) || (len(word) >= 4 &&
			sylabs(word[:len(word)-1]) == 1 &&
			cvc(word[len(word)-4:len(word)-1]) == false) ||
			(sylabs(word[:len(word)-1]) == 1 && len(word) < 4) {
			word = word[:len(word)-1]
		}
	}
	if sylabs(word) > 1 && strings.HasSuffix(word, "ll") {
		return word[:len(word)-1]
	}

	return word

}

// wrapper for all the rules applied in turns
func Stem(word string) string {
	var stem string
	// for words shorter than two letters no stemming needed
	if len(word) <= 2 {
		return word
	}
	// deal with plurals and past tense
	stem = pluralsPastStem(word)
	// change y to i at the end
	stem = yToI(stem)
	// replaces longer suffices by appropriate forms
	stem = changeSuffixStem(stem)
	// if stem sylabs >= 1, remove/change suffix
	stem = shortRemoveStem(stem)
	// if stem sylabs >= 2, remove suffix
	stem = longRemoveStem(stem)
	// additional cleaning after stemming and return stem
	return additionalCleanup(stem)
}
