package main

import "github.com/1axup/catui/catui"

func main() {
	kitty := catui.Cat{ // needed for printing!!
		Name:  "Schmidt",
		IsBot: false,
		Icon:  "cat1",
	}

	botti := catui.Cat{
		Name:  "Botti",
		IsBot: true,
		Icon:  "bot1",
	}

	kitty.Say("Hello, there")
	botti.Say("Hello, you cat")
}
