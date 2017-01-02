// Run through the years/months pulling death entries from Wikipedia and
// total them all up.
package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sadbox/mediawiki"
)

// *[[Name]], Age or *<!--blah--> [[Name]], Age  or <!--blah-->*[[Name]], Age
var re1 *regexp.Regexp = regexp.MustCompile(`^\s*(?:\<!--.+?--\>)*?\s*\*\s*(?:\<!--.+?--\>)*?\s*\[\[(.+?)\]\]\s*[,|\.]\s*(\d+)`)

// *[[Name]] or *<!--blah--> [[Name]] or <!--blah-->* [[Name]]
var re2 *regexp.Regexp = regexp.MustCompile(`^\s*(?:\<!--.+?--\>)*?\s*\*\s*(?:\<!--.+?--\>)*?\s*\[\[(.+?)\]\]\s*`)

// *SirName [[blah]], age, profession
var re5 *regexp.Regexp = regexp.MustCompile(`^\s*\*\s*(.+?)\s*\[\[(.+)\]\]\s*,\s*(\d+)`)

// ===Num=== Day of Month
var dayRE *regexp.Regexp = regexp.MustCompile(`^===\s*(\d+)\s*===$`)

// ===[[February 22|22]]===
var dayRE2 *regexp.Regexp = regexp.MustCompile(`^===\s*\[\[\w+\s+(\d+)\|\d+\]\]\s*===`)

// These are basically just various regexes that helped exclude parts of the pages
// so that we could more easily figure out what entries we were missing.
// Smashmouth programming, this could be more effecient
var garbageRE1 *regexp.Regexp = regexp.MustCompile(`^{{`)
var garbageRE2 *regexp.Regexp = regexp.MustCompile(`^<!--.+-->\s*$`)
var garbageRE3 *regexp.Regexp = regexp.MustCompile(`^==\s*References\s*==$`)
var garbageRE4 *regexp.Regexp = regexp.MustCompile(`^The following is a list of`)
var garbageRE5 *regexp.Regexp = regexp.MustCompile(`^Please place names under the date the person died`)
var garbageRE6 *regexp.Regexp = regexp.MustCompile(`^The intent of these pages is to report notable deaths`)
var garbageRE7 *regexp.Regexp = regexp.MustCompile(`^-->$`)
var garbageRE8 *regexp.Regexp = regexp.MustCompile(`^<!--$`)
var garbageRE9 *regexp.Regexp = regexp.MustCompile(`^=+$`)
var garbageRE10 *regexp.Regexp = regexp.MustCompile(`^Only those meeting the`)
var garbageRE11 *regexp.Regexp = regexp.MustCompile(`^NOTE: Causes of death such as`)
var garbageRE12 *regexp.Regexp = regexp.MustCompile(`^Alphabetical order please`)
var garbageRE13 *regexp.Regexp = regexp.MustCompile(`^\*\s*Name, age, country of citizenship and reason for notability,`)
var garbageRE14 *regexp.Regexp = regexp.MustCompile(`^Entries for each day are listed`)
var garbageRE15 *regexp.Regexp = regexp.MustCompile(`Please read before adding a name to this list`)
var garbageRE16 *regexp.Regexp = regexp.MustCompile(`^\*\s*\[\[\s*Deaths`)

// convertAge converts string based age to integer.
func convertAge(agestr string) int {
	age, err := strconv.Atoi(agestr)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return age
}

// Verifies that a wikipedia page for pageName does exist.
// NOTE: If there are multiple we only pull the first one due to false positives.
// It's worth noting that our regex isn't perfect, so always check false reports
// to be sure it's truely not found.
func validPage(pageName string) bool {

	// Why not share a client? I don't know.
	client, err := mediawiki.New("https://en.wikipedia.org/w/api.php", "My Wiki Death Bot")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	// Only want the first result
	page := strings.Split(pageName, "|")[0]
	result, err := client.Read(page)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if result.Pageid != 0 {
		return true
	}

	return false
}

// parsePage will take the downloaded death page and pull out any entries.
func parsePage(page *mediawiki.Page) int {
	total := 0
	scanner := bufio.NewScanner(strings.NewReader(page.Revisions[0].Body))
	dayOfMonth := 0
	for scanner.Scan() {
		matchFound := false
		pageIsValid := false
		str := scanner.Text()
		name := "Unknown Name"
		age := 0
		// Needs to be first since it conflicts with our entry regexes
		if garbageRE16.MatchString(str) {
			// noop
		} else if re1.MatchString(str) {
			matches := re1.FindStringSubmatch(str)
			matchFound = true
			if validPage(matches[1]) {
				pageIsValid = true
			}
			name = matches[1]
			age = convertAge(matches[2])
		} else if re5.MatchString(str) {
			// Sirname Version
			matches := re5.FindStringSubmatch(str)
			matchFound = true
			if validPage(matches[2]) {
				pageIsValid = true
			}
			name = fmt.Sprintf("%s %s", matches[1], matches[2])
			age = convertAge(matches[3])
		} else if re2.MatchString(str) {
			matches := re2.FindStringSubmatch(str)
			matchFound = true
			if validPage(matches[1]) {
				pageIsValid = true
			}
			name = matches[1]
		} else if dayRE.MatchString(str) {
			var err error
			dayOfMonth, err = strconv.Atoi(dayRE.FindStringSubmatch(str)[1])
			if err != nil {
				fmt.Println(err.Error())
			}
		} else if dayRE2.MatchString(str) {
			var err error
			dayOfMonth, err = strconv.Atoi(dayRE2.FindStringSubmatch(str)[1])
			if err != nil {
				fmt.Println(err.Error())
			}
		} else if str == "" {
			// noop
		} else if garbageRE1.MatchString(str) {
			// noop
		} else if garbageRE2.MatchString(str) {
			// noop
		} else if garbageRE3.MatchString(str) {
			// noop
		} else if garbageRE4.MatchString(str) {
			// noop
		} else if garbageRE5.MatchString(str) {
			// noop
		} else if garbageRE6.MatchString(str) {
			// noop
		} else if garbageRE7.MatchString(str) {
			// noop
		} else if garbageRE8.MatchString(str) {
			// noop
		} else if garbageRE9.MatchString(str) {
			// noop
		} else if garbageRE10.MatchString(str) {
			// noop
		} else if garbageRE11.MatchString(str) {
			// noop
		} else if garbageRE12.MatchString(str) {
			// noop
		} else if garbageRE13.MatchString(str) {
			// noop
		} else if garbageRE14.MatchString(str) {
			// noop
		} else if garbageRE15.MatchString(str) {
			// noop
		} else {
			//fmt.Println("----> ", scanner.Text())
		}
		if matchFound && pageIsValid {
			total++
			fmt.Printf(">>> %d - %s - %d\n", dayOfMonth, name, age)
		} else if matchFound && !pageIsValid {
			fmt.Printf("=== %d - %s - %d\n", dayOfMonth, name, age)
		}
	}
	return total
}

func main() {
	var totalDead int

	client, err := mediawiki.New("https://en.wikipedia.org/w/api.php", "My Wiki Death Bot")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Smashmouth programming, this could be more effecient.
	years := []string{"2004", "2005", "2006", "2007", "2008", "2009", "2010", "2011", "2012", "2013", "2014", "2015", "2016"}
	months := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

	for _, year := range years {
		yearDead := 0
		for _, month := range months {
			// Format of page name by practice, verify to be sure.
			pageToRead := fmt.Sprintf("Deaths_in_%s_%s", month, year)

			page, err := client.Read(pageToRead)
			if err != nil {
				fmt.Println(err)
				return
			}
			monthDead := parsePage(page)
			yearDead += monthDead
			totalDead += monthDead

			fmt.Printf("\t\tTotal Dead for month %s: %d\n", month, monthDead)
		}
		fmt.Printf("Total Dead for year %s: %d\n", year, yearDead)
	}
	fmt.Printf("Total Dead: %d\n", totalDead)
}
