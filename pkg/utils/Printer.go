package utils

import (
	"github.com/fatih/color"
)

type PrintColor struct {

}

func (p *PrintColor) Normal() *color.Color{
	return color.New(color.BgGreen)
}


func (p *PrintColor) Warings() *color.Color{
	return color.New(color.BgYellow)
}


func (p *PrintColor) Err() *color.Color{
	return color.New(color.BgRed)
}
