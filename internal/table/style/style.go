package style

type Style struct {
	IncludeANSI bool
}

var Default Style = Style{
	IncludeANSI: true,
}
