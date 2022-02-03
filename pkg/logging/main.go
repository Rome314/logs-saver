package logging

import (
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

type Entry struct {
	*log.Entry
}

func (e *Entry) WithMethod(method string) *Entry {
	return &Entry{e.Entry.WithField("method", method)}
}
func (e *Entry) WithPlace(place string) *Entry {
	return &Entry{e.Entry.WithField("place", place)}
}

func GetLogger(module, submodule string) *Entry {
	return &Entry{Entry: logger.WithFields(log.Fields{
		"module":    module,
		"submodule": submodule,
	})}
}

func init() {

	logger = log.New()
	logger.SetFormatter(&log.TextFormatter{
		// SortingFunc: func(strings []string) {
		// 	for index, key := range strings {
		// 		switch key {
		// 		case "module":
		// 			strings[3], strings[index] = strings[index], strings[3]
		// 		case "submodule":
		// 			strings[4], strings[index] = strings[index], strings[4]
		// 		case "method":
		// 			strings[5], strings[index] = strings[index], strings[5]
		// 		case "place":
		// 			strings[6], strings[index] = strings[index], strings[6]
		// 		}
		// 	}
		// },
	})

}
