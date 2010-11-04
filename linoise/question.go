// Copyright 2010  The "go-linoise" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package linoise

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kless/go-term/term"
)


// Values by default
var (
	QuestionPrefix      = " + " // String placed before of questions
	QuestionErrPrefix   = "  "  // String placed before of error messages
	QuestionTrueString  = "y"   // String to represent 'true'
	QuestionFalseString = "n"   // String to represent 'false'

	QuestionFloatFmt  byte = 'g' // Format for float numbers
	QuestionFloatPrec int  = -1  // Precision for float numbers
)

// To pass strings in another languages.
var ExtraBoolString = make(map[string]bool)


// === Type
// ===

type Question struct {
	trueString, falseString string
}


// Gets a question type.
func NewQuestion() *Question {
	// === Check the strings that represent a boolean.
	_, err := atob(QuestionTrueString)
	if err != nil {
		panic(fmt.Sprintf("the string %q does not represent a boolean 'true'",
			QuestionTrueString))
	}

	_, err = atob(QuestionFalseString)
	if err != nil {
		panic(fmt.Sprintf("the string %q does not represent a boolean 'false'",
			QuestionFalseString))
	}

	return &Question{
		strings.ToLower(QuestionTrueString),
		strings.ToLower(QuestionFalseString),
	}
}

// Restores terminal settings.
func (q *Question) RestoreTerm() {
	term.RestoreTerm()
}
// ===


// Gets a line type ready to show questions.
func (q *Question) getLine(prompt, defaultAnswer string, hasDefault bool) *Line {
	prompt = QuestionPrefix + prompt

	// Add the value by default
	if hasDefault {
		prompt = fmt.Sprintf("%s [%s]", prompt, defaultAnswer)
	}

	// Add spaces
	if strings.HasSuffix(prompt, "?") {
		prompt += " "
	} else {
		prompt += ": "
	}

	return NewLinePrompt(prompt, nil) // No history.
}

// Base to read strings.
func (q *Question) _baseReadString(prompt, defaultAnswer string, hasDefault bool) string {
	line := q.getLine(prompt, defaultAnswer, hasDefault)

	for {
		answer, err := line.Read()
		if answer != "" || err == ErrCtrlD {
			return answer
		}

		if hasDefault {
			return defaultAnswer
		}
	}
	return ""
}

// Prints the question waiting until to press Return.
func (q *Question) ReadString(prompt string) string {
	return q._baseReadString(prompt, "", false)
}

// Prints the question waiting until to press Return. If input is nil then
// it returns the answer by default.
func (q *Question) ReadStringDefault(prompt, defaultAnswer string) string {
	return q._baseReadString(prompt, defaultAnswer, true)
}

// Prints the question until to get an integer number.
func (q *Question) _baseReadInt(prompt string, defaultAnswer int, hasDefault bool) int {
	line := q.getLine(prompt, strconv.Itoa(defaultAnswer), hasDefault)

	for {
		input, err := line.Read()
		if err == ErrCtrlD {
			break
		}

		if input == "" && hasDefault {
			return defaultAnswer
		}

		answer, err := strconv.Atoi(input)
		if err != nil {
			fmt.Fprintf(output, "%s%q: value has to be an integer\n",
				QuestionErrPrefix, input)
			continue
		} else {
			return answer
		}
	}
	return 0
}

// Prints the question until to get an integer number.
func (q *Question) ReadInt(prompt string) int {
	return q._baseReadInt(prompt, 0, false)
}

// Prints the question until to get an integer number. If input is nil then
// it returns the answer by default.
func (q *Question) ReadIntDefault(prompt string, defaultAnswer int) int {
	return q._baseReadInt(prompt, defaultAnswer, true)
}

// Prints the question until to get a float number.
func (q *Question) _baseReadFloat(prompt string, defaultAnswer float, hasDefault bool) float {
	line := q.getLine(
		prompt,
		strconv.Ftoa(defaultAnswer, QuestionFloatFmt, QuestionFloatPrec),
		hasDefault,
	)

	for {
		input, err := line.Read()
		if err == ErrCtrlD {
			break
		}

		if input == "" && hasDefault {
			return defaultAnswer
		}

		answer, err := strconv.Atof(input)
		if err != nil {
			fmt.Fprintf(output, "%s%q: value has to be a float\n",
				QuestionErrPrefix, input)
			continue
		} else {
			return answer
		}
	}
	return 0.0
}

// Prints the question until to get a float number.
func (q *Question) ReadFloat(prompt string) float {
	return q._baseReadFloat(prompt, 0.0, false)
}

// Prints the question until to get a float number. If input is nil then
// it returns the answer by default.
func (q *Question) ReadFloatDefault(prompt string, defaultAnswer float) float {
	return q._baseReadFloat(prompt, defaultAnswer, true)
}

// Prints the question until to get a string that represents a boolean.
func (q *Question) ReadBool(prompt string, defaultAnswer bool) bool {
	var options string

	if defaultAnswer {
		options = fmt.Sprintf("%s/%s", strings.ToUpper(q.trueString), q.falseString)
	} else {
		options = fmt.Sprintf("%s/%s", q.trueString, strings.ToUpper(q.falseString))
	}

	line := q.getLine(prompt, options, true)

	for {
		input, err := line.Read()
		if err == ErrCtrlD {
			break
		}

		if input == "" {
			return defaultAnswer
		}

		answer, err := atob(input)
		if err != nil {
			fmt.Fprintf(output, "%s%q: does not represent a boolean\n",
				QuestionErrPrefix, input)
			continue
		} else {
			return answer
		}
	}
	return false
}


// === Utility
// ===

// Returns the boolean value represented by the string.
// It accepts "y, Y, yes, YES, Yes, n, N, no, NO, No". And values in
// 'strconv.Atob', and 'ExtraBoolString'. Any other value returns an error.
func atob(str string) (value bool, err os.Error) {
	v, err := strconv.Atob(str)
	if err == nil {
		return v, nil
	}

	switch str {
	case "y", "Y", "yes", "YES", "Yes":
		return true, nil
	case "n", "N", "no", "NO", "No":
		return false, nil
	}

	// Check extra characters, if any.
	boolExtra, ok := ExtraBoolString[str]
	if ok {
		return boolExtra, nil
	}

	return false, os.NewError("wrong")
}

