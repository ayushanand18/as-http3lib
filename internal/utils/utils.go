package utils

import "log"

func PrintStartBanner() {
	log.Println("\033[1;34m")
	log.Println(" █████╗ ███████╗██╗  ██╗██████╗ ")
	log.Println("██╔══██╗██╔════╝██║  ██║╚════██╗")
	log.Println("███████║███████╗███████║ █████╔╝")
	log.Println("██╔══██║╚════██║██╔══██║ ╚═══██╗")
	log.Println("██║  ██║███████║██║  ██║██████╔╝")
	log.Println("╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═════╝ ")
	log.Println("\033[0m\033[1;32mas-htp3lib - Faster HTTP/3-native alternative to FastAPI in Go\033[0m")
	log.Println("\033[1;32m(Asynchronous HTTP/3 Lib)\033[0m")
}
