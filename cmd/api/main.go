package main

import "log"

func main() {
	config := &config{
		address: ":3000",
	}
	app := &application{
		config: *config,
	}

	mux := app.mount()

	log.Fatal(app.run(&mux))
}
