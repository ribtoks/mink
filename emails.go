package main

import (
	"regexp"
	"strconv"
	"strings"
)

var reg = regexp.MustCompile(`([a-zA-Z0-9._-]+@([a-zA-Z0-9_-]+\.)+[a-zA-Z0-9_-]+)`)

// Check if email looks valid
func isValidEmail(email string) bool {
	split := strings.Split(email, ".")
	if len(split) < 2 {
		return false
	}

	ending := split[len(split)-1]

	if len(ending) < 2 {
		return false
	}

	if _, err := strconv.Atoi(ending); err == nil {
		return false
	}

	return true
}

func parseEmails(body []byte) []string {
	res := reg.FindAll(body, -1)
	scrapedEmails := make([]string, 0, 10)

	for _, r := range res {
		email := string(r)
		if !isValidEmail(email) {
			continue
		}

		var found bool
		for _, existingEmail := range scrapedEmails {
			if existingEmail == email {
				found = true
				break
			}
		}

		if found {
			continue
		}

		scrapedEmails = append(scrapedEmails, email)
	}
	return scrapedEmails
}
