// Laravel equivalent: public/index.php - see docs/chapters/01-installation.md
package main

import (
	"goravel/bootstrap"
)

// #region main
// main is the program's entry point. Unlike public/index.php, it runs once:
// the compiled binary boots the app and then serves every request from memory.
func main() {
	app := bootstrap.Boot()

	app.Start()
}

// #endregion main
