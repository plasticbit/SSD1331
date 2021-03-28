package OLED

func (m DisplayMode) String() string {
	switch m {
	case Nomal:
		return "Normal"
	case EntireON:
		return "EntireON"
	case EntireOFF:
		return "EntireOFF"
	case Inverse:
		return "Inverse"

	default:
		return ""
	}
}

func (s DisplayOnOff) String() string {
	switch s {
	case DisplayON:
		return "DisplayON"
	case DisplayOnInDim:
		return "DisplayOnInDim"
	case DisplayOff:
		return "DisplayOff"

	default:
		return ""
	}
}

func (s ScrollStep) String() string {
	switch s {
	case Frames6:
		return "Frames6"
	case Frames10:
		return "Frames10"
	case Frames100:
		return "Frames100"
	case Frames200:
		return "Frames200"

	default:
		return ""
	}
}
