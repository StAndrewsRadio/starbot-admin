package utils

import (
	"math/rand"
	"strings"
	"unicode"
)

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

// Splits a string using FieldsFunc and a function that will only split n number of times
func FieldsN(s string, n int) []string {
	splits, lastResult := 0, false
	splitFunc := func(c rune) bool {
		if unicode.IsSpace(c) && splits < n {
			lastResult = true
		} else {
			if lastResult == true {
				splits++
			}

			lastResult = false
		}

		return lastResult
	}

	return strings.FieldsFunc(s, splitFunc)
}

// Checks if two string slices have equal contents
func StringSliceEquals(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

// Formats a list of discord users
func FormatUserList(userList []string) string {
	// check nil input
	if userList == nil || len(userList) == 0 {
		return ""
	}

	output := ""

	// create list
	for _, user := range userList {
		output += "<@" + user + ">, "
	}

	// remove trailing comma
	output = output[:len(output)-2]

	// get last index of ", "
	index := strings.LastIndex(output, ", ")
	if index != -1 {
		// replace it with an " and "
		output = output[:index] + " and " + output[index+2:]
	}

	return output
}

// Gets the index of a string in a slice
func StringSliceIndexOf(element string, slice []string) int {
	for k, v := range slice {
		if element == v {
			return k
		}
	}

	return -1
}
