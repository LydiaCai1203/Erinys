package main

import "erinys/erinys"

func main() {
	engine := erinys.NewHTTPEngine("/cache")
	engine.Run(":9001")
}
