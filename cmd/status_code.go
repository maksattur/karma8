package main

type statusCode uint8

const (
	StatusCodeOK statusCode = iota
	StatusCodeFailedLoadConfig
	StatusCodeFailedConnectPG
	StatusCodeFailedMigrations
	StatusCodeFailedStartServer
	StatusCodeFailedStopServer
)
