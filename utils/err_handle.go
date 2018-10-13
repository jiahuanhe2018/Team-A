package utils

import "log"

func PanicErr(err error)  {
	if err != nil {
		log.Panic(err)
	}
}
