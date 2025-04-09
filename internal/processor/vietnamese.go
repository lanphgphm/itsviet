package processor 

import (
	"strings"
)

type InputMethod string 

const (
	Telex 	InputMethod = "telex"
)

type VietnameseProcessor struct {
	buffer 			[]rune
	method 			InputMethod
	maxBufferSize	int
}

func NewVietnameseProcessor(method InputMethod) *VietnameseProcessor {
	return &VietnameseProcessor {
		buffer: 		make([]rune, 0, 10), 
		method: 		method, 
		maxBufferSize: 	10,
	}
}

func (vp *VietnameseProcessor) ProcessKey(key rune) ([]rune, bool) {
	vp.buffer = append(vp.buffer, key)

	if len(vp.buffer) > vp.maxBufferSize {
		vp.buffer = vp.buffer[1:]
	}

	// process current buffer string
	bufferStr := string(vp.buffer)

	var result string 
	var processed bool 

	if (vp.method == Telex) {
		result, processed = vp.processTelexInput(bufferStr)
	} // else if chain of more input type here 

	if (processed) {
		vp.buffer = vp.buffer[:0] // clear buffer
		return []rune(result), true
	}

	return []rune{key}, false
}

func (vp *VietnameseProcessor) processTelexInput(input string) (string, bool) {
	
	switch {
	case strings.HasSuffix(input, "aa"):
		return "â", true
	case strings.HasSuffix(input, "aw"):
		return "ă", true
	case strings.HasSuffix(input, "ee"):
		return "ê", true
	case strings.HasSuffix(input, "oo"):
		return "ô", true
	case strings.HasSuffix(input, "ow"):
		return "ơ", true
	case strings.HasSuffix(input, "dd"):
		return "đ", true
	}

	switch {
	case strings.HasSuffix(input, "as"):
		return "á", true
	case strings.HasSuffix(input, "af"):
		return "à", true
	case strings.HasSuffix(input, "ar"):
		return "ả", true
	case strings.HasSuffix(input, "ax"):
		return "ã", true
	case strings.HasSuffix(input, "aj"):
		return "ạ", true
	}

	return input, false
}