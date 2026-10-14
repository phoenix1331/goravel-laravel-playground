package main

import (
	"goravel/bootstrap"
)

// main is the program's entry point. Unlike public/index.php, it runs once:
// the compiled binary boots the app and then serves every request from memory.
func main() {
	app := bootstrap.Boot()

	app.Start()
}
