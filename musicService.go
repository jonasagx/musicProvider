package main

import (
	. "backend"
)

func main() {
	Conf = LoadConfig()
	StartHTTPServer()
}
