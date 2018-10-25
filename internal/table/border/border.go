package border

type Border rune

const (
	TableTopLeft     Border = '┍'
	TableTopRight    Border = '┑'
	HorizontalDash   Border = '╌'
	HorizontalLight  Border = '─'
	HorizontalHeavy  Border = '━'
	Vertical         Border = '│'
	ColBoth          Border = '┼'
	ColDown          Border = '┬'
	ColUp            Border = '┴'
	ColBothHeavy     Border = '┿'
	ColDownHeavy     Border = '┯'
	ColUpHeavy       Border = '┷'
	RowSepLeft       Border = '├'
	RowSepRight      Border = '┤'
	RowSepLeftHeavy  Border = '┝'
	RowSepRightHeavy Border = '┥'
	TableBottomLeft  Border = '┕'
	TableBottomRight Border = '┙'
)

func (b Border) Rune() rune {
	return rune(b)
}

func (b Border) String() string {
	return string(b)
}

// ┍━━━━━━━━━━━━━━━━━━━┑
// │       Title       │
// ┝━━━━━━━━━━━━━━━━━━━┥
