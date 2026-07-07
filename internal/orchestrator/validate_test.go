package orchestrator

import "testing"

func TestValidateWorkers(t *testing.T) {
	cases := []struct {
		name    string
		in      int
		want    int
		wantErr bool
	}{
		{"borne basse valide", 1, 1, false},
		{"valeur nominale", 3, 3, false},
		{"borne haute valide", 100, 100, false},
		{"zéro invalide", 0, 3, true},
		{"négatif invalide", -5, 3, true},
		{"trop grand invalide", 101, 3, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ValidateWorkers(c.in)
			if got != c.want {
				t.Errorf("valeur = %d, want %d", got, c.want)
			}
			if (err != nil) != c.wantErr {
				t.Errorf("err = %v, wantErr = %v", err, c.wantErr)
			}
		})
	}
}
