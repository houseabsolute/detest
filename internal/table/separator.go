package table

import (
	"errors"

	"github.com/houseabsolute/detest/internal/table/border"
	"github.com/houseabsolute/detest/internal/table/renderinfo"
	"github.com/houseabsolute/detest/internal/table/row"
)

type SeparatorType int

const (
	Start SeparatorType = iota
	AfterTitle
	AfterHeader
	InBody
	AfterBody
	End
)

type Separator struct {
	sepType SeparatorType
}

func (s Separator) Render(ri renderinfo.RI, before, after *row.Row) (string, error) {
	left, fill, right, err := s.base(ri)
	if err != nil {
		return "", err
	}

	rendered := make([]rune, ri.TotalWidth()+2)
	rendered[0] = left
	for i := 0; i < ri.TotalWidth(); i++ {
		rendered[i+1] = fill
	}
	rendered[ri.TotalWidth()+1] = right

	upDown := map[int]border.Border{}
	if before != nil {
		for _, i := range before.ColumnSeparatorPositions(ri) {
			if s.isHeavy() {
				upDown[i] = border.ColUpHeavy
			} else {
				upDown[i] = border.ColUp
			}
		}
	}
	if after != nil {
		for _, i := range after.ColumnSeparatorPositions(ri) {
			if s.isHeavy() {
				if _, e := upDown[i]; e {
					upDown[i] = border.ColBothHeavy
				} else {
					upDown[i] = border.ColDownHeavy
				}
			} else {
				if _, e := upDown[i]; e {
					upDown[i] = border.ColBoth
				} else {
					upDown[i] = border.ColDown
				}
			}
		}
	}

	for i, b := range upDown {
		rendered[i] = b.Rune()
	}

	return string(rendered), nil
}

func (s Separator) base(_ renderinfo.RI) (rune, rune, rune, error) {
	switch s.sepType {
	case Start:
		return border.TableTopLeft.Rune(), border.HorizontalHeavy.Rune(), border.TableTopRight.Rune(), nil
	case AfterTitle:
		return border.RowSepLeftHeavy.Rune(), border.HorizontalHeavy.Rune(), border.RowSepRightHeavy.Rune(), nil
	case AfterHeader:
		return border.RowSepLeft.Rune(), border.HorizontalLight.Rune(), border.RowSepRight.Rune(), nil
	case InBody:
		return border.RowSepLeft.Rune(), border.HorizontalDash.Rune(), border.RowSepRight.Rune(), nil
	case AfterBody:
		return border.RowSepLeftHeavy.Rune(), border.HorizontalHeavy.Rune(), border.RowSepRightHeavy.Rune(), nil
	case End:
		return border.TableBottomLeft.Rune(), border.HorizontalHeavy.Rune(), border.TableBottomRight.Rune(), nil
	}

	return 0, 0, 0, errors.New("unknown separator type")
}

func (s Separator) isHeavy() bool {
	switch s.sepType {
	case Start, AfterTitle, AfterBody, End:
		return true
	}

	return false
}
