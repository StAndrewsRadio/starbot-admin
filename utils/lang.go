package utils

import "math/rand"

var (
	alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// Returns the substring of a string, cut from the start index to the end index.
func Substring(s string, start, end int) string {
	var startIndex, currentIndex int

	// iterate through every index in the string
	for rangeIndex := range s {
		// mark down the start or end indexes if reached
		if currentIndex == start {
			startIndex = rangeIndex
		} else if currentIndex == end {
			return s[startIndex:rangeIndex]
		}

		currentIndex++
	}

	return s[startIndex:]
}

// Checks if a slice of strings contains a given string
func StringSliceContains(slice []string, string string) bool {
	for _, element := range slice {
		if element == string {
			return true
		}
	}

	return false
}

// Returns a random string of letters of a given length.
func RandomString(length int) string {
	result := make([]rune, length)

	// iterate through every letter we need
	for i := range result {
		result[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(result)
}
