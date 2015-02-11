package import_processor

import (
	"go/token"
	. "goprep"
)

// Funciones para trabajo con tokens.

// Funcion que deja pasar un token.
func passToken(tok Token, out chan string, sync chan interface{}) {
	out <- tok.Str
	sync <- nil
}

// Funcion que elimina un token.
func eliminateToken(out chan string, sync chan interface{}) {
	out <- ""
	sync <- nil
}

// Funcion que trata la declaracion de imports.
func Imports_declare(tok Token, in chan Token, out chan string, sync chan interface{}) {
	var enc bool = false
	for tok.Token == token.IMPORT {
		passToken(tok, out, sync)
		tok = <-in
		switch tok.Token {
		case token.LPAREN:
			passToken(tok, out, sync)
			tok = <-in
			for tok.Token != token.RPAREN {
				switch tok.Token {
				case token.PERIOD:
					passToken(tok, out, sync)
					tok = <-in
					if tok.Str == "\"runtime\"" {
						enc = true
						passToken(tok, out, sync)
						tok = <-in
						passToken(tok, out, sync)
						tok = <-in
					} else {
						passToken(tok, out, sync)
						tok = <-in
						passToken(tok, out, sync)
						tok = <-in
					}
				case token.STRING:
					if tok.Str == "\"runtime\"" {
						enc = true
						passToken(tok, out, sync)
						tok = <-in
						passToken(tok, out, sync)
						tok = <-in
					} else {
						passToken(tok, out, sync)
						tok = <-in
						passToken(tok, out, sync)
						tok = <-in
					}
				case token.COMMENT: // Ignora comentarios
					passToken(tok, out, sync)
					tok = <-in
				}
			}
		case token.PERIOD:
			passToken(tok, out, sync)
			tok = <-in
			if tok.Str == "\"runtime\"" {
				enc = true
				passToken(tok, out, sync)
				tok = <-in
				passToken(tok, out, sync)
				tok = <-in
			} else {
				passToken(tok, out, sync)
				tok = <-in
				passToken(tok, out, sync)
				tok = <-in
			}
		case token.STRING:
			if tok.Str == "\"runtime\"" {
				enc = true
				passToken(tok, out, sync)
				tok = <-in
				passToken(tok, out, sync)
				tok = <-in
			} else {
				passToken(tok, out, sync)
				tok = <-in
				passToken(tok, out, sync)
				tok = <-in
			}
		case token.COMMENT: // Ignora comentarios
			passToken(tok, out, sync)
			tok = <-in
		}
	}
	if !enc { // Incluye paquete "runtime"
		passToken(tok, out, sync)
		tok = <-in
		out <- tok.Str + "\n" + "import \"runtime\"\n"
		sync <- nil
	} else {
		passToken(tok, out, sync)
	}
}
