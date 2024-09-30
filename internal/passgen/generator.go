package passgen

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/pkg/errors"
)

var (
	vowels             = []rune("aeiouAEIOU")
	consonants         = []rune("bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ")
	numbers            = []rune("1234567890")
	separators         = []rune(`_-.`)
	specialSymbols     = []rune(`!@#$%&?`)
	verySpecialSymbols = []rune(`*^()+={[]}'";:/|\~<>`)
	wierdSignsPack1    = []rune("ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖØÙÚÛÜÝÞß")
	wierdSignsPart2    = []rune("¡¢£¤¥§©¦¨«¬®¯°µ¶·¸»¿")
	wierdSignsPack3    = []rune("²³¹ºª¼½¾×±")
)

type Generator struct {
	randomInts []int
}

// New returns Generator with 128 random generated values.
func New() (*Generator, error) {
	const randomBatchLen = 128

	rnd := make([]byte, randomBatchLen*4) // 4 byte for 1 int32
	_, err := rand.Read(rnd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read random sequence")
	}

	g := Generator{randomInts: make([]int, 0, randomBatchLen)}
	for i := 0; i < len(rnd); i += 4 {
		// Read 4 bytes into num
		num := binary.BigEndian.Uint32(rnd[i : i+4])
		g.randomInts = append(g.randomInts, int(num))
	}
	return &g, nil
}

func (g *Generator) RandIntn(n int) int {
	r := g.randomInts[0]
	if len(g.randomInts) == 1 {
		// skipping error as this case is working only when we trying to get
		// >64 random values from the Generator, which is should not ever happen.
		newGen, _ := New()
		g.randomInts = newGen.randomInts
	}

	g.randomInts = g.randomInts[1:]
	return r % n
}

// GenReadablePass looks like *word**number**separator**word**number**separator**specialSymbol*
func (g *Generator) GenReadablePass() string {
	const minWordLen = 12
	word := make([]rune, 0)

	// word
	word = append(word, g.generatePronounceableWord()...)
	// separator
	word = append(word, g.randomSeparator())
	// word
	word = append(word, g.generatePronounceableWord()...)

	num := func() []rune {
		num := []rune{g.randomNumber()}
		if g.RandIntn(2) == 1 { // 50%
			num = append(num, g.randomNumber())
			if g.RandIntn(2) == 1 { // 25%
				num = append(num, g.randomNumber())
			}
		}
		return num
	}()
	randomPlace := g.RandIntn(len(word))
	word = append(word[:randomPlace], append(num, word[randomPlace:]...)...)

	randomPlace = g.RandIntn(len(word))
	word = append(word[:randomPlace], append([]rune{g.randomSpecialSymbol()}, word[randomPlace:]...)...)
	for len(word) < minWordLen {
		randomPlace = g.RandIntn(len(word))
		word = append(word[:randomPlace], append([]rune{g.randomSpecialSymbol()}, word[randomPlace:]...)...)
	}

	return string(word)
}

func (g *Generator) GenSafePass() string {
	const (
		minLength = 18
		maxLength = 25
	)

	word := make([]rune, 0)
	// word
	word = append(word, g.generatePronounceableWord()...)
	// separator
	word = append(word, g.randomSeparator())
	// word
	word = append(word, g.generatePronounceableWord()...)
	// separator
	word = append(word, g.randomSeparator())
	// word
	word = append(word, g.generatePronounceableWord()...)

	num := func() []rune {
		num := []rune{g.randomNumber()}
		if g.RandIntn(2) == 1 { // 50%
			num = append(num, g.randomNumber())
			if g.RandIntn(2) == 1 { // 25%
				num = append(num, g.randomNumber())
			}
		}
		return num
	}()
	randomPlace := g.RandIntn(len(word))
	word = append(word[:randomPlace], append(num, word[randomPlace:]...)...)
	randomPlace = g.RandIntn(len(word))
	word = append(word[:randomPlace], append([]rune{g.randomVerySpecialSymbol()}, word[randomPlace:]...)...)

	length := g.RandIntn(maxLength-minLength) + minLength
	for i := len(word); i < length; i++ {
		randomPlace = g.RandIntn(len(word))
		word = append(word[:randomPlace], append([]rune{g.randomSafeLetter()}, word[randomPlace:]...)...)
	}
	return string(word)
}

func (g *Generator) GenInsanePass() string {
	const (
		minLength = 27
		maxLength = 40
	)

	length := g.RandIntn(maxLength-minLength) + minLength
	word := make([]rune, length)
	for i := range word {
		word[i] = g.randomInsaneLetter()
	}
	return string(word)
}

func (g *Generator) randomSafeLetter() rune {
	letterType := g.RandIntn(100)
	switch {
	case letterType < 50:
		return g.randomSpecialSymbol()
	case letterType < 66:
		return g.randomVerySpecialSymbol()
	case letterType < 88:
		return g.randomNumber()
	default:
		return g.randomSeparator()
	}
}

func (g *Generator) randomInsaneLetter() rune {
	letterType := g.RandIntn(100)
	switch {
	case letterType < 11:
		return g.randomVowel()
	case letterType < 22:
		return g.randomConsonant()
	case letterType < 33:
		return g.randomSeparator()
	case letterType < 44:
		return g.randomNumber()
	case letterType < 55:
		return g.randomSpecialSymbol()
	case letterType < 66:
		return g.randomVerySpecialSymbol()
	case letterType < 77:
		return g.randomWierdPack1()
	case letterType < 88:
		return g.randomWierdPack2()
	default:
		return g.randomWierdPack3()
	}
}

func (g *Generator) randomSyllable() []rune {
	syllables := [][]rune{
		{g.randomConsonant(), g.randomVowel()},
		{g.randomVowel(), g.randomConsonant()},
		{g.randomVowel(), g.randomConsonant(), g.randomVowel()},
		{g.randomConsonant(), g.randomVowel(), g.randomConsonant()},
	}
	return syllables[g.RandIntn(len(syllables))]
}

func (g *Generator) generatePronounceableWord() []rune {
	const (
		minWordLen = 4
		maxWordLen = 8
	)

	word := make([]rune, 0, 20)
	length := g.RandIntn(maxWordLen-minWordLen) + minWordLen
	for len(word) < length {
		word = append(word, g.randomSyllable()...)
	}
	return word[:length]
}

func (g *Generator) randomVowel() rune {
	return rune(vowels[g.RandIntn(len(vowels))])
}

func (g *Generator) randomConsonant() rune {
	return rune(consonants[g.RandIntn(len(consonants))])
}

func (g *Generator) randomSeparator() rune {
	return rune(separators[g.RandIntn(len(separators))])
}

func (g *Generator) randomNumber() rune {
	return rune(numbers[g.RandIntn(len(numbers))])
}

func (g *Generator) randomSpecialSymbol() rune {
	return rune(specialSymbols[g.RandIntn(len(specialSymbols))])
}

func (g *Generator) randomVerySpecialSymbol() rune {
	return rune(verySpecialSymbols[g.RandIntn(len(verySpecialSymbols))])
}

func (g *Generator) randomWierdPack1() rune {
	return wierdSignsPack1[g.RandIntn(len(wierdSignsPack1))]
}

func (g *Generator) randomWierdPack2() rune {
	return wierdSignsPart2[g.RandIntn(len(wierdSignsPart2))]
}

func (g *Generator) randomWierdPack3() rune {
	return wierdSignsPack3[g.RandIntn(len(wierdSignsPack3))]
}
