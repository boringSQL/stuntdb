package stuntdb

import "fmt"

type Skipped struct {
	Subject string
	Reason  string
}

func (s Skipped) String() string {
	return fmt.Sprintf("%s: %s", s.Subject, s.Reason)
}
