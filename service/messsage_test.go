package service

import "testing"

func TestSendEmail(t *testing.T) {
	s := Service{}
	err := s.sendEmail(nil)
	if err != nil {
		panic(err)
	}
}
