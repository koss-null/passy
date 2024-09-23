package passgen

import (
	"time"

	"golang.org/x/exp/rand"
)

const (
	vowels             = "aeiouAEIOU"
	consonants         = "bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ"
	numbers            = "1234567890"
	separators         = `_-.`
	specialSymbols     = `!@#$%&?`
	verySpecialSymbols = `*^()+={[]}'";:/|\`
	russianLetters     = "абвгдежзийклмнопрстуфхцчъыьэюя"
	japaneeseLetters   = "あいえおかきくけこさしすせそたちつてとな漢字一二三四五六七八九十"
	arabLetters        = "أبجدهوزحطيكلمنسعفصقرشتثخذضظغ"
)

// GenReadablePass looks like *word**number**separator**word**number**separator**specialSymbol*
func GenReadablePass() string {
	rand.Seed(uint64(time.Now().UnixNano()))

	word := make([]rune, 0)

	// word
	word = append(word, generatePronounceableWord()...)
	// random number
	if rand.Intn(3) == 1 { // 25%
		word = append(word, randomNumber())
		if rand.Intn(3) == 1 { // 25%
			word = append(word, randomNumber())
			if rand.Intn(3) == 1 { // 25%
				word = append(word, randomNumber())
			}
		}
	}

	// separator
	word = append(word, randomSeparator())
	// word
	word = append(word, generatePronounceableWord()...)

	// random number
	word = append(word, randomNumber())
	word = append(word, randomNumber())
	// special symbol
	word = append(word, randomSpecialSymbol())
	return string(word)
}

func GenSafePass() string {
	const (
		minLength = 14
		maxLength = 18
	)
	rand.Seed(uint64(time.Now().UnixNano()))

	length := rand.Intn(maxLength-minLength) + minLength
	word := make([]rune, length)
	for i := range word {
		word[i] = randomSafeLetter()
	}
	return string(word)
}

func GenInsanePass() string {
	const (
		minLength = 27
		maxLength = 40
	)
	rand.Seed(uint64(time.Now().UnixNano()))

	length := rand.Intn(maxLength-minLength) + minLength
	word := make([]rune, length)
	for i := range word {
		word[i] = randomInsaneLetter()
	}
	return string(word)
}

func randomSafeLetter() rune {
	letterType := rand.Intn(100)
	switch {
	case letterType < 30:
		return randomVowel()
	case letterType < 60:
		return randomConsonant()
	case letterType < 67:
		return randomSeparator()
	case letterType < 82:
		return randomNumber()
	case letterType < 90:
		return randomSpecialSymbol()
	default:
		return randomVerySpecialSymbol()
	}
}

func randomInsaneLetter() rune {
	letterType := rand.Intn(100)
	switch {
	case letterType < 11:
		return randomVowel()
	case letterType < 22:
		return randomConsonant()
	case letterType < 33:
		return randomSeparator()
	case letterType < 44:
		return randomNumber()
	case letterType < 55:
		return randomSpecialSymbol()
	case letterType < 66:
		return randomVerySpecialSymbol()
	case letterType < 77:
		return randomRussianLetter()
	case letterType < 88:
		return randomJapaneseLetter()
	default:
		return randomArabicLetter()
	}
}

func randomSyllable() []rune {
	syllables := [][]rune{
		{randomConsonant(), randomVowel()},
		{randomVowel(), randomConsonant()},
		{randomVowel(), randomConsonant(), randomVowel()},
		{randomConsonant(), randomVowel(), randomConsonant()},
	}
	return syllables[rand.Intn(len(syllables))]
}

func generatePronounceableWord() []rune {
	const (
		minWordLen = 4
		maxWordLen = 8
	)

	word := make([]rune, 0, 20)
	length := rand.Intn(maxWordLen-minWordLen+1) + minWordLen
	for len(word) < length {
		word = append(word, randomSyllable()...)
	}
	return word[:maxWordLen]
}

func randomVowel() rune {
	return rune(vowels[rand.Intn(len(vowels))])
}

func randomConsonant() rune {
	return rune(consonants[rand.Intn(len(consonants))])
}

func randomSeparator() rune {
	return rune(separators[rand.Intn(len(separators))])
}

func randomNumber() rune {
	return rune(numbers[rand.Intn(len(numbers))])
}

func randomSpecialSymbol() rune {
	return rune(specialSymbols[rand.Intn(len(specialSymbols))])
}

func randomVerySpecialSymbol() rune {
	return rune(verySpecialSymbols[rand.Intn(len(verySpecialSymbols))])
}

func randomRussianLetter() rune {
	return rune(russianLetters[rand.Intn(len(russianLetters))])
}

func randomJapaneseLetter() rune {
	return rune(japaneeseLetters[rand.Intn(len(japaneeseLetters))])
}

func randomArabicLetter() rune {
	return rune(arabLetters[rand.Intn(len(arabLetters))])
}
