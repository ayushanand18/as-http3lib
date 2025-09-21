package utils

import "log"

func PrintStartBanner() {
	log.Println("\033[1;34m")
	log.Println("")
	log.Println(" .o88b. d8888b.  .d8b.  d88888D db    db db   db d888888b d888888b d8888b. ")
	log.Println("d8P  Y8 88  `8D d8' `8b YP  d8' `8b  d8' 88   88 `~~88~~' `~~88~~' 88  `8D ")
	log.Println("8P      88oobY' 88ooo88    d8'   `8bd8'  88ooo88    88       88    88oodD' ")
	log.Println("8b      88`8b   88~~~88   d8'      88    88~~~88    88       88    88~~~   ")
	log.Println("Y8b  d8 88 `88. 88   88  d8' db    88    88   88    88       88    88      ")
	log.Println(" `Y88P' 88   YD YP   YP d88888P    YP    YP   YP    YP       YP    88      ")
	log.Println("                                                                         ")
	log.Println("\033[0m\033[1;32mcrazyhttp - Faster HTTP/3-native alternative to FastAPI in Go\033[0m")
	log.Println("\033[1;32m(Crazy-ingly fast performance)\033[0m")
}
