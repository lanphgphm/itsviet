package main

import (
	"fmt"
	"unicode"
)

type VProcessor struct {
	buffer		string
	// lastWord	string
	enabled		bool
}

func NewVProcessor() *VProcessor {
	return &VProcessor{
		buffer:		"", 
		// lastWord:	"", 
		enabled: 	true,
	}
}

func (vp *VProcessor) Toggle() {
	vp.enabled = !vp.enabled
	vp.buffer = ""
	// vp.lastWord = ""
}

func (vp *VProcessor) Process(keyStr string, keyCode byte, state uint16) (bool, string) {
	/*
		Reads the keypress sequence, and decide if it should be intercepted. 
		Intercept rules: 
		- Non-shift Modifiers (Ctrl, Alt, Super): Skip (forward as-is)
		- Shift+Char: Treat like regular character (add to buffer) 
		- Backspace: Delete word from buffer (entire word, even punctation)
		- Other things (Del, Ins, PrtScr, PgUp/Down, etc.): Skip (forward as-is)
		- Alphabetical character: Intercept (this is the target)
	*/

	// if detect non-letter char, its time to inject the finished word 
	// clear buffer for next word
	if !unicode.IsLetter(rune(keyStr[0])) {
		remainBuffer := vp.buffer + keyStr
		fmt.Printf("Non-letter key detected: '%s' - clearing buffer\n", keyStr)
		vp.buffer = ""
		return true, remainBuffer
	}

	// intercept only alphabetical characters 
	vp.buffer += keyStr
	fmt.Printf("Buffer: '%s'\n", vp.buffer)
	needsTransform, transformed := vp.applyTelex()

	if needsTransform {
		fmt.Printf("Transform triggered! Result: '%s'\n", transformed)
		vp.buffer = transformed
		return false, ""
	}

	return false, ""
}

func (vp *VProcessor) applyTelex() (bool, string) {
	if len(vp.buffer) < 2 {
		return false, ""
	}

	n := len(vp.buffer)
	lastChar := rune(vp.buffer[n-1])
	pb := vp.buffer[:n-1] // pb = processedBuffer 

	targetIndex := -1 
	for i := len(pb)-1; i>=0; i-- {
	// for i := 0; i < len(pb); i++ {
		r := rune(pb[i])
		if isVietTarget(r) {
			targetIndex = i
			break
		}
	}

	if targetIndex > -1 {
		targetChar := rune(pb[targetIndex])
		key := VietChar{
			Target:		targetChar, 
			Modifier:	lastChar, 
		}

		if replacement, exists := telexMod[key]; exists {
			transformed := pb[:targetIndex] + replacement + pb[targetIndex+1:]
			return true, transformed
		}
	}
	
	return false, ""
}

func isVietTarget(r rune) bool {
	targets := []rune{
        'a', 'e', 'i', 'o', 'u', 'd',
        'A', 'E', 'I', 'O', 'U', 'D',
        'ă', 'â', 'ê', 'ô', 'ơ', 'ư',
        'Ă', 'Â', 'Ê', 'Ô', 'Ơ', 'Ư',
    }

	for _, t := range targets {
		if r == t {
			return true
		}
	}

	return false
}
