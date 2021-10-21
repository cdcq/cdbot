package helpers

import "log"

func AddLog(where string, when string, what error) {
	log.Println("Error at " + where + " when " + when + ":")
	log.Println(what)
}
