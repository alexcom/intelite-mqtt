// Package compact contains implementation of Tasmota compact IR signal encoding.
// Written according to description here: https://tasmota.github.io/docs/IRSend-RAW-Encoding/
// So can be inaccurate or buggy.
package tasmota

import (
	"strconv"
	"strings"
)

const high = '+'
const low = '-'

func Compact(input []int) string {
	sb := strings.Builder{}
	// +3355 -> A
	direct := make(map[int]rune)
	sign := true
	letter := 'A'
	for _, n := range input {
		q := (n / 10) * 10
		r := n % 5
		if r >= 5 {
			q += 5
		}
		if l, ok := direct[n]; ok {
			sb.WriteRune(l)
		} else {
			if sign {
				sb.WriteRune(high)
			} else {
				sb.WriteRune(low)
			}
			sb.WriteString(strconv.Itoa(n))
			direct[n] = letter
			letter++
		}
		sign = !sign
	}
	return sb.String()
}

func Payload(in []int) string {
	// tasmota-ir command IRsend, 0-auto transmission frequency(default 38kHz), rest is compact form of IR code
	return "IRsend 0," + Compact(in)
}
