package led

type Color struct {
	Red   uint8 `mapstructure:"red"`
	Green uint8 `mapstructure:"green"`
	Blue  uint8 `mapstructure:"blue"`
}
