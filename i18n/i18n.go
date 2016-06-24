// Package i18n is used for basic translation of messages
package i18n

import "fmt"

var instance translator

type translator struct {
	bundle Bundle
}

type Bundle map[string][]string

func Load(bundle Bundle) error {
	instance = translator{bundle: bundle}
	return nil
}

type LangId uint8

// Internal function for getting translation
// It prints formatted string if args are provided
func translate(lang LangId, message string, args ...interface{}) string {
	if ms, ok := instance.bundle[message]; ok {
		// if empty default language return message also
		// this prevents empty messages
		if ms == nil || len(ms) == 0 || (lang == 0 && len(ms[int(lang)]) == 0) {
			if args != nil && len(args) > 0 {
				return fmt.Sprintf(message, args...)
			} else {
				return message
			}
		}
		if int(lang) > len(ms)-1 {
			if args != nil && len(args) > 0 {
				return fmt.Sprintf(ms[0], args...)
			} else {
				return ms[0]
			}
		}
		if args != nil && len(args) > 0 {
			return fmt.Sprintf(ms[int(lang)], args...)
		} else {
			return ms[int(lang)]
		}

	}
	return message
}

// Function L selects Language for the translation. It takes LangId and returns a function which
// takes a message
func L(langId LangId) func(string, ...interface{}) string {
	return func(message string, args ...interface{}) string {
		return translate(langId, message, args...)
	}
}

// Function T means 'Translate'. It takes message
func T(message string, args ...interface{}) string {
	return translate(0, message, args...)
}
