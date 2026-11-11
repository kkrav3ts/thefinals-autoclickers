package keyboard

import "testing"

func TestKeyNamesCoverLetters(t *testing.T) {
	if KeyNames[0x41] != "A" || KeyNames[0x5A] != "Z" {
		t.Fatalf("letters = %q %q", KeyNames[0x41], KeyNames[0x5A])
	}
	if KeyNames[0xA3] != "Right Ctrl" || KeyNames[0x25] != "Left Arrow" {
		t.Fatal("extended keys missing from the bind list")
	}
}

func TestExtendedScan(t *testing.T) {
	if !extendedScan(0x25) || !extendedScan(0xA5) || !extendedScan(0x6F) {
		t.Fatal("arrows, right alt, and numpad divide are extended keys")
	}
	if extendedScan(0x41) || extendedScan(0x20) {
		t.Fatal("letters and space are not extended keys")
	}
}
