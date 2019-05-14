package main

type Config struct {
	port int
	wssPort int
	numCoreThread int
	baseCycle int
	logLevel int
	isTest int
	cluster string
	local string
	backendUri string
	conf string

	// wss sign certificate & private key
	serverCrt string
	serverKey string
}

var(

)