package main

import "github.com/1axup/catui/catui"

func main() {
	kitty := catui.Cat{
		Name:  "Schmidt",
		IsBot: false,
		Icon:  "cat1",
	}

	kitty.Say("Hello, there")
}
