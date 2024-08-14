package main

import (
	"fmt"

	"github.com/armiariyan/assessment-tsel/cmd"

	"github.com/labstack/gommon/color"
)

const banner = `
████████╗   ███████╗███████╗██╗
╚══██╔══╝   ██╔════╝██╔════╝██║
   ██║█████╗███████╗█████╗  ██║
   ██║╚════╝╚════██║██╔══╝  ██║
   ██║      ███████║███████╗███████╗
   ╚═╝      ╚══════╝╚══════╝╚══════╝

armiariyan.

`

func main() {
	fmt.Print(color.Yellow(banner))
	cmd.Run()
}
