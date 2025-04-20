package main

import (
	"fmt"
	// "strings"
	"unicode"

	"github.com/BurntSushi/xgb/xproto"
)

type VProcessor struct {
	buffer		string
	lastWord	string
	enabled		bool
}

func NewVProcessor() *VProcessor {
	return &VProcessor{
		buffer:		"", 
		lastWord:	"", 
		enabled: 	true,
	}
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

	if !vp.enabled {
		return false, ""
	}

	// check for modifier, skip all but allow shift 
	// if shift is pressed --> false 
	// if ctrl+shift is pressed --> true
	if state&^xproto.ModMaskShift != 0 {
		return false, ""
	}

	// BACKSPACE := 22 
	// if keyCode == BACKSPACE 
	if keyCode == 22 {
		n := len(vp.buffer)
		if n > 0 {
			vp.buffer = vp.buffer[:n-1]
			fmt.Printf("Buffer after backspace: '%s'\n", vp.buffer)
		}
		return false, ""
	}

	// placeholder to skip special keys, function keys 
	if keyCode < 10 || len(keyStr) != 1 {
		return false, ""
	}

	// dont intercept anything thats not a letter 
	if !unicode.IsLetter(rune(keyStr[0])) {
		fmt.Printf("Non-letter key detected: '%s' - clearing buffer\n", keyStr)
		vp.buffer = ""
		return false, ""
	}

	// intercept only alphabetical characters 
	vp.buffer += keyStr
	fmt.Printf("Buffer: '%s'\n", vp.buffer)
	needsTransform, transformed := vp.applyTelex()

	if needsTransform {
		fmt.Printf("Transform triggered! Result: '%s'\n", transformed)
		vp.buffer = transformed
		fmt.Printf("Buffer set to remaining: '%s'\n", vp.buffer)
		return true, transformed
	}

	// return true, keyStr
	return false, ""
}

func (vp *VProcessor) applyTelex() (bool, string) {
	if len(vp.buffer) < 2 {
		return false, ""
	}

	n := len(vp.buffer)
	lastChar := vp.buffer[n-1]

	switch lastChar {
	case 'a', 'A':
		processedBuffer := vp.buffer[:n-1]
		for i := len(processedBuffer)-1; i >= 0; i-- {
			if processedBuffer[i] == 'a' || processedBuffer[i] == 'A' {
				replacement := "â"
				if processedBuffer[i] == 'A' {
					replacement = "Â"
				}
				transformed := processedBuffer[:i] + replacement + processedBuffer[i+1:]
				return true, transformed
			}
		}
	}

	return false, ""
}

func (vp *VProcessor) Toggle() {
	vp.enabled = !vp.enabled
	vp.buffer = ""
	vp.lastWord = ""
}