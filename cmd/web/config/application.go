package config

import "log"

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

type FlagsConfig struct {
	Port string
}
