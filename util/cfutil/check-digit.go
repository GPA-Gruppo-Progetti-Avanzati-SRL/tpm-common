package cfutil

import (
	"fmt"
	"strings"
)

// Tabella dei valori per i caratteri in posizione dispari (1, 3, 5, 7, 9, 11, 13, 15)
var valoriDispari = map[rune]int{
	'0': 1, '1': 0, '2': 5, '3': 7, '4': 9, '5': 13, '6': 15, '7': 17, '8': 19, '9': 21,
	'A': 1, 'B': 0, 'C': 5, 'D': 7, 'E': 9, 'F': 13, 'G': 15, 'H': 17, 'I': 19, 'J': 21,
	'K': 2, 'L': 4, 'M': 18, 'N': 20, 'O': 11, 'P': 3, 'Q': 6, 'R': 8, 'S': 12, 'T': 14,
	'U': 16, 'V': 10, 'W': 22, 'X': 25, 'Y': 24, 'Z': 23,
}

// Tabella dei valori per i caratteri in posizione pari (2, 4, 6, 8, 10, 12, 14)
var valoriPari = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4, 'F': 5, 'G': 6, 'H': 7, 'I': 8, 'J': 9,
	'K': 10, 'L': 11, 'M': 12, 'N': 13, 'O': 14, 'P': 15, 'Q': 16, 'R': 17, 'S': 18, 'T': 19,
	'U': 20, 'V': 21, 'W': 22, 'X': 23, 'Y': 24, 'Z': 25,
}

// Tabella di conversione del resto in carattere di controllo
var caratteriControllo = []rune{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

// CalcolaCarattereControllo calcola il carattere di controllo per un codice fiscale
// Input: i primi 15 caratteri del codice fiscale
// Output: il carattere di controllo (16° carattere)
func CalcolaCarattereControllo(codiceFiscale15 string) (rune, error) {

	// Converte in maiuscolo
	cf := strings.ToUpper(codiceFiscale15)

	// Verifica che la lunghezza sia esattamente 15
	if len(cf) != 15 {
		return 0, fmt.Errorf("il codice fiscale deve essere di 15 caratteri (ricevuti %d)", len(cf))
	}

	somma := 0

	// Itera sui 15 caratteri
	for i, char := range cf {
		if i%2 == 0 { // Posizione dispari (0-indexed, quindi 0, 2, 4...)
			if val, ok := valoriDispari[char]; ok {
				somma += val
			} else {
				return 0, fmt.Errorf("carattere non valido alla posizione %d: %c", i+1, char)
			}
		} else { // Posizione pari (1, 3, 5...)
			if val, ok := valoriPari[char]; ok {
				somma += val
			} else {
				return 0, fmt.Errorf("carattere non valido alla posizione %d: %c", i+1, char)
			}
		}
	}

	// Calcola il resto della divisione per 26
	resto := somma % 26

	// Restituisce il carattere di controllo corrispondente
	return caratteriControllo[resto], nil
}

// VerificaCodiceFiscale verifica se un codice fiscale completo (16 caratteri) è valido come check digit
func CheckCF(codiceFiscale string) error {
	cf := strings.ToUpper(codiceFiscale)

	if len(cf) != 16 {
		return fmt.Errorf("il codice fiscale deve essere di 16 caratteri")
	}

	// Calcola il carattere di controllo atteso
	carattereCalcolato, err := CalcolaCarattereControllo(cf[:15])
	if err != nil {
		return err
	}

	// Confronta con il carattere di controllo presente
	caratterePresente := rune(cf[15])
	if carattereCalcolato != caratterePresente {
		err = fmt.Errorf("unmatched check digit: wanted %q, found %q", carattereCalcolato, caratterePresente)
	}

	return err
}
