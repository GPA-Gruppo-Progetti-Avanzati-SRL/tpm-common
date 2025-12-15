package cfutil

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	CONSONANTS = "bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ"
	VOWELS     = "aeiou"
)

func computeConsonanti(s string) string {
	var cns string
	for _, r := range s {
		if strings.IndexRune(CONSONANTS, r) != -1 {
			cns += string(r)
		}
	}

	return cns
}

func computeVocali(s string) string {
	var voc string
	for _, r := range s {
		if strings.IndexRune(VOWELS, r) != -1 {
			voc += string(r)
		}
	}

	return voc
}

var omocodieMap = map[string]string{
	"L": "0",
	"M": "1",
	"N": "2",
	"P": "3",
	"Q": "4",
	"R": "5",
	"S": "6",
	"T": "7",
	"U": "8",
	"V": "9",
}

var birthMonthMap = map[int]string{
	1:  "A",
	2:  "B",
	3:  "C",
	4:  "D",
	5:  "E",
	6:  "H",
	7:  "L",
	8:  "M",
	9:  "P",
	10: "R",
	11: "S",
	12: "T",
}

func processLastName(s string) string {
	var code string
	code = computeConsonanti(s)
	if len(code) >= 3 {
		code = code[:3]
	} else {
		vowels := computeVocali(s)
		code = code + vowels
		switch {
		case len(code) < 3:
			code = code + strings.Repeat("X", 3-len(code))
		case len(code) == 3:
		case len(code) > 3:
			code = code[:3]
		}
	}
	return strings.ToUpper(code)
}

func processFirstName(s string) string {
	var code string
	code = computeConsonanti(s)
	if len(code) > 3 {
		code = code[0:1] + code[2:4]
	} else {
		vowels := computeVocali(s)
		code = code + vowels
		switch {
		case len(code) < 3:
			code = code + strings.Repeat("X", 3-len(code))
		case len(code) == 3:
		case len(code) > 3:
			code = code[:3]
		}
	}

	return strings.ToUpper(code)
}

func resolveOmocodie(s string) string {
	var resolved string
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			resolved += s[i : i+1]
		} else {
			if nc, ok := omocodieMap[s[i:i+1]]; ok {
				resolved += nc
			} else {
				resolved += "*"
			}
		}
	}

	return resolved
}

// CalculateCF birthDate in the format 20060102 MPRMLS62S21G337J
func CalculateCF(lastName, firstName string, sex string, birthDate string, foreign bool, birthTownHallCode string) string {
	cf := processLastName(lastName) + processFirstName(firstName) + birthDate[2:4]

	mnth, _ := strconv.Atoi(birthDate[4:6])
	cf += birthMonthMap[mnth]

	day := birthDate[6:]
	switch strings.ToUpper(sex) {
	case "M":
		cf += day
	case "F":
		iday, _ := strconv.Atoi(day)
		cf += fmt.Sprintf("%02d", iday+40)
	}

	if foreign {
		cf = cf + strings.Repeat("_", 16-len(cf))
	} else {
		cf = cf + strings.ToUpper(birthTownHallCode)
	}

	// Check digit non computed
	cf += "_"
	return cf
}

func CheckCF(cf, lastName, firstName string, maleFemale string, birthDate string, foreign bool, birthTownHallCode string) error {
	const semLogContext = "check-cf::check-cf"

	cf = strings.ToUpper(cf)
	if len(cf) != 16 {
		return errors.New("length of cf is not 16")
	}

	if cf[:3] != processLastName(lastName) {
		return errors.New("last name is not correct")
	}

	if cf[3:6] != processFirstName(firstName) {
		return errors.New("first name is not correct")
	}

	//           111111
	// 0123456789012345
	// MPRMLS62S21G337J
	if resolveOmocodie(cf[6:8]) != birthDate[2:4] {
		return errors.New("birth year is not correct")
	}

	mnth, _ := strconv.Atoi(birthDate[4:6])
	if cf[8:9] != birthMonthMap[mnth] {
		return errors.New("birth month is not correct")
	}

	day := birthDate[6:]
	switch maleFemale {
	case "male":
	case "female":
		iday, _ := strconv.Atoi(day)
		day = fmt.Sprintf("%02d", iday+40)
	}

	if resolveOmocodie(cf[9:11]) != day {
		return errors.New("birth day is not correct")
	}

	if foreign {
		if cf[11:12] == "Z" {
			return nil
		}

		return errors.New("foreign place not correct")
	} else {
		codCat := cf[11:12] + resolveOmocodie(cf[12:15])
		if codCat != strings.ToUpper(birthTownHallCode) {
			return errors.New("birth townhall is not correct")
		}
	}

	return nil
}
