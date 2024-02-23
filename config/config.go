package config

const (
	BrightnessMode = iota
	ColorMode
)

const BrightnessStepSize float32 = 25.5

var Colors = [...]string{
	"#ff0000", //red
	"#f79719", //orange
	"#ffff00", //yellow
	"#00ff00", //green
	"#0000ff", //blue
	"#00eeff", //light_blue
	"#a600ff", //purple
	"#ff00e6", //magenta
}
