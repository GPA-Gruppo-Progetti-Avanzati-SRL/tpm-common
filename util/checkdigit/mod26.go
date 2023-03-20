package checkdigit

/*
 * Algoritmo utilizzato all'interno dei codici a bruciatura. si tratta di una variazione del calcolo del check-digit per codice fiscale.... piuttosto semplificato to say the least
 * Check Digit: 1 Carattere alfabetico “Codice Controllo” calcolato assegnando un peso
 * - i caratteri numerici 'pari'  (incluso lo '0') hanno un peso pari a 0
 * - i caratteri numerici 'dispari' hanno un peso pari a 1
 * - i caratteri alfabetici hanno un peso pari a 2
 * Viene calcolata la somma dei pesi e il risultato è fornito dal modulo 26 della somma + 'A' -1
 */

func ComputeMod26CheckDigit(s string) string {
	v := 0
	for i, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			v += 2
		} else if c >= '0' && c <= '9' {
			if (i+1)%2 == 1 {
				v += 1
			}
		}
	}

	d := modulo2Character[v%26]
	return string(d)
}
