package main

const VERSION = "1.0.0"

func main() {
	app := application{}

	app.NewConfig()

	app.NewLoggers()

	if err := app.initDb(); err != nil {
		app.errLog.Fatalf("database connection failed: %s", err)
	}

	if err := app.initServer(); err != nil {
		app.errLog.Fatalf("server connection failed: %s", err)
	}
}
