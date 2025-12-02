package forum

// Entry point for the program
func Main(args []string) {
	// TODO pass args and populate them maybe, otherwise default
	// maybe initialize something first if needed...
	// for example the database?!?!
	// maybe that could be a flag...
	// e.g. --init

	var ip, port string
	// in case the flag is not passed, the auto intialization could be triggered
	if len(args) != 3 {
		ip = "127.0.0.1"
		port = "8080"
	} else {
		ip = args[1]
		port = args[2]
	}

	// TODO Initialization of needed stuff

	startServer(ip, port) // Unused values...
}
