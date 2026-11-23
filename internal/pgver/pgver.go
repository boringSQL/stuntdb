package pgver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var pattern = regexp.MustCompile(`(\d+)\.(\d+)`)

type Version struct {
	Major int
	Minor int
}

func Parse(text string) (Version, error) {
	match := pattern.FindStringSubmatch(text)
	if match == nil {
		return Version{}, fmt.Errorf("no version in %q", strings.TrimSpace(text))
	}
	major, err := strconv.Atoi(match[1])
	if err != nil {
		return Version{}, err
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return Version{}, err
	}
	return Version{Major: major, Minor: minor}, nil
}

func (v Version) String() string {
	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor)
}

func (v Version) Less(other Version) bool {
	if v.Major != other.Major {
		return v.Major < other.Major
	}
	return v.Minor < other.Minor
}
