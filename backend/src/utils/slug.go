package utils

import (
	"regexp"
	"strings"
)

// GenerateSlug converts a string to a URL-friendly slug
// Example: "Food & Drinks" -> "food-drinks"
func GenerateSlug(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")

	// Remove leading and trailing hyphens
	s = strings.Trim(s, "-")

	// If empty after processing, return a default
	if s == "" {
		s = "category"
	}

	return s
}

// GenerateUniqueSlug generates a unique slug by appending a number if needed
// This should be called with a function that checks if the slug exists
func GenerateUniqueSlug(baseSlug string, exists func(string) bool) string {
	slug := baseSlug
	counter := 1

	for exists(slug) {
		slug = baseSlug + "-" + string(rune('0'+counter))
		counter++
		if counter > 9 {
			slug = baseSlug + "-" + string(rune('0'+counter/10)) + string(rune('0'+counter%10))
		}
	}

	return slug
}

// IsValidSlug checks if a string is a valid slug format
func IsValidSlug(s string) bool {
	if s == "" {
		return false
	}

	// Check if it contains only lowercase letters, numbers, and hyphens
	reg := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	return reg.MatchString(s)
}
