package utils

import "log"

func Assert(condition bool) {
	if !condition {
		log.Panicln("Assertion failed")
	}
}

func AssertMsg(condition bool, msg string) {
	if !condition {
		log.Panicln("Assertion failed: " + msg)
	}
}

func PanicIfErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
