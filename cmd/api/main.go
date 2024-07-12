package main

const VERSION = "1.0.0"

func main() {
	app := application{}

	app.newConfig()

	app.initLoggers()

	if err := app.initDb(); err != nil {
		app.errLog.Fatalf("database connection failed: %s", err)
	}
	defer app.db.Close()

	app.initModels()

	if err := app.initServer(); err != nil {
		app.errLog.Fatalf("server connection failed: %s", err)
	}
}
