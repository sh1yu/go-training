package common

import "log"

func MustNilErr(err error, errPrefix string) {
	if err != nil {
		log.Fatalf("%v:%v", errPrefix, err)
	}
}
