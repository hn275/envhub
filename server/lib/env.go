package lib

import (
	"log"
	"os"
)

func Getenv(k string) (v string) {
	v = os.Getenv(k)
	if v == "" {
		log.Fatalf("env [%s] not set", k)
	}
	return v
}
