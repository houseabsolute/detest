package style

type Style struct {
	IncludeANSI bool
}

var Default = Style{
	IncludeANSI: true,
}

func New(includeANSI bool) Style {
	return Style{includeANSI}
}
