package app

import "strings"

func Ok() map[string]bool {
	return map[string]bool{"ok": true}
}

func OkWithField(field interface{}) map[string]interface{} {
	return map[string]interface{}{"ok": field}
}

func Err(msg string) map[string]string {
	return map[string]string{"error": msg}
}

// StringContains returns true if s contains any of substrings
func StringContains(s string, substrings ...string) bool {
	if len(substrings) > 0 {
		for _, v := range substrings {
			if strings.Contains(s, v) {
				return true
			}
		}
	}
	return false
}
