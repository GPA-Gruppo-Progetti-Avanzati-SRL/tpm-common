package cfutil

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
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

var omocodiaMap = map[rune]string{
	'0': "0",
	'1': "1",
	'2': "2",
	'3': "3",
	'4': "4",
	'5': "5",
	'6': "6",
	'7': "7",
	'8': "8",
	'9': "9",
	'L': "0",
	'M': "1",
	'N': "2",
	'P': "3",
	'Q': "4",
	'R': "5",
	'S': "6",
	'T': "7",
	'U': "8",
	'V': "9",
}

/*
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
*/

var LetterCode2MonthMap = map[string]int{
	"A": 1,
	"B": 2,
	"C": 3,
	"D": 4,
	"E": 5,
	"H": 6,
	"L": 7,
	"M": 8,
	"P": 9,
	"R": 10,
	"S": 11,
	"T": 12,
}

var Month2LetterCodeMap = map[int]string{
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
	for _, ch := range s {
		if nc, ok := omocodiaMap[ch]; ok {
			resolved += nc
		} else {
			resolved += "*"
		}
	}

	return resolved
}

// CalculateCF birthDate in the format 20060102 MPRMLS62S21G337J
func CalculateCF(lastName, firstName string, sex string, birthDate string, foreign bool, birthTownHallCode string) (string, error) {
	cf := processLastName(lastName) + processFirstName(firstName) + birthDate[2:4]

	mnth, _ := strconv.Atoi(birthDate[4:6])
	cf += Month2LetterCodeMap[mnth]

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

	r, err := CalcolaCarattereControllo(cf)
	if err != nil {
		return "", err
	}

	// Check digit non computed
	cf += string(r)
	return cf, nil
}

func CheckCFAgainstCFInfo(cf, lastName, firstName string, maleFemale string, birthDate string, foreign bool, birthTownHallCode string) error {
	const semLogContext = "check-cf::check-cf"

	cf = strings.ToUpper(cf)
	err := CheckCF(cf)
	if err != nil {
		return err
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
	if cf[8:9] != Month2LetterCodeMap[mnth] {
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

type FiscalCodeInfo struct {
	FiscalCode   string
	SurnameCode  string
	NameCode     string
	BirthDate    time.Time
	Gender       string
	BelfioreCode string
	CheckCode    string
}

func decodeMonth(month string) (int, error) {
	value := LetterCode2MonthMap[month]
	if value == 0 {
		return 0, fmt.Errorf("Invalid month in FiscalCode: %d", value)
	}
	return value, nil
}

func decodeOmocodia(value string) (string, error) {
	output := ""
	for _, char := range value {
		decodedValue := omocodiaMap[char]
		if decodedValue == "" {
			return "", fmt.Errorf("Invalid char in FiscalCode: %c", char)
		}
		output += decodedValue
	}

	return output, nil
}

func ExtractInfo(fiscalCode string) (*FiscalCodeInfo, error) {
	fiscalCode = strings.ToUpper(fiscalCode)
	regex := regexp.MustCompile(`^(.{3})(.{3})(.{2})(.{1})(.{2})(.{1})(.{3})(.{1})$`)

	res := regex.FindAllStringSubmatch(fiscalCode, -1)

	if res == nil {
		return nil, fmt.Errorf("Invalid fiscal code: %s", fiscalCode)
	}

	surnameCode := res[0][1]
	nameCode := res[0][2]
	yearCode, err := decodeOmocodia(res[0][3])
	if err != nil {
		return nil, err
	}
	monthCode := res[0][4]
	dayCode, err := decodeOmocodia(res[0][5])
	if err != nil {
		return nil, err
	}
	belfioreCode, err := decodeOmocodia(res[0][7])
	if err != nil {
		return nil, err
	}
	belfioreCode = res[0][6] + belfioreCode
	checkCode := res[0][8]

	log.Debug().
		Str("surnameCode", surnameCode).
		Str("nameCode", nameCode).
		Str("yearCode", yearCode).
		Str("monthCode", monthCode).
		Str("dayCode", dayCode).
		Str("belfioreCode", belfioreCode).
		Str("checkCode", checkCode).
		Msg("FiscalCode parts")

	gender := "M"
	day, err := strconv.Atoi(dayCode)
	if err != nil {
		return nil, err
	}

	if day > 40 {
		gender = "F"
		day = day - 40
	}

	month, err := decodeMonth(monthCode)
	if err != nil {
		return nil, err
	}

	year, err := strconv.Atoi("19" + yearCode)
	if err != nil {
		return nil, err
	}

	birthDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	log.Debug().
		Str("gender", string(gender)).
		Int("day", day).
		Int("month", month).
		Int("year", year).
		Interface("birthDate", birthDate).
		Msg("FiscalCode data")

	fiscalCodeInfo := FiscalCodeInfo{
		FiscalCode:   fiscalCode,
		SurnameCode:  surnameCode,
		NameCode:     nameCode,
		BirthDate:    birthDate,
		Gender:       gender,
		BelfioreCode: belfioreCode,
		CheckCode:    checkCode,
	}

	log.Debug().
		Interface("fiscalCodeInfo", fiscalCodeInfo).
		Msg("FiscalCode data")

	return &fiscalCodeInfo, nil
}
