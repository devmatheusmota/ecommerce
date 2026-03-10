package handlers

import "os"

var Version = "dev"

func VersionString() string {
	if version := os.Getenv("VERSION"); version != "" {
		return version
	}
	return Version
}
