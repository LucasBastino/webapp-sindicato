package validation

import "testing"


func TestValidateDni(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "valid empty", in: "", want: ""},
		{name: "valid short dni", in: "1234567", want: ""},
		{name: "valid 8 digits dni", in: "12345678", want: ""},
		{name: "valid only numbers", in: "30111222", want: ""},
		{name: "reject letters", in: "1234567a", want: "Caracter no válido."},
		{name: "reject dot", in: "12.345.678", want: "Caracter no válido."},
		{name: "reject space", in: "1234 567", want: "Caracter no válido."},
		{name: "reject more than 8 digits", in: "123456789", want: "No puede contener mas de 8 caracteres."},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ValidateDni(c.in)
			if got != c.want {
				t.Fatalf("ValidateDni(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestValidateEmail (t *testing.T) {
	cases := []struct{
		name string
		in	 string
		want string
	}{
		{name: "valid empty", in: "", want: ""},
		{name: "valid email", in: "prueba@gmail.com", want: ""},
		{name: "valid email with numbers", in: "pru45eba@gmail.com", want: ""},
		{name: "valid email with dot", in: "pru.eba@gmail.com", want: ""},
		{name: "valid email with characters", in: "pru$e-ba@gmail.com", want: ""},
		{name: "valid double dot in domain", in: "prueba@g.mail.com", want: ""},
		{name: "reject without username", in: "@gmail.com", want:"Formato no válido."},
		{name: "reject without '@'", in: "pruebagmail.com", want: "Formato no válido."},
		{name: "reject without domain", in: "prueba@", want: "Formato no válido."},
		{name: "reject without extension", in: "prueba@gmail", want: "Formato no válido."},
		{name: "reject invalid extension", in: "prueba@gmail.c", want: "Formato no válido."},
		{name: "reject consecutive double dot", in: "prueba@gmail..com", want: "Formato no válido."},
		{name: "reject dot first domain character", in: "prueba@.gmail.com", want: "Formato no válido."},
		{name: "reject dot last domain character", in: "prueba@gmail.com.", want: "Formato no válido."},
		{name: "reject invalid character", in: "prueba@g$mail.com", want: "Caracter no válido."},
	}

	for _, c := range cases{
		t.Run(c.name, func(t *testing.T){
			got := ValidateEmail(c.in)
			if got != c.want{
				t.Fatalf("ValidateEmail(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}