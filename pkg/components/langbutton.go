package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type LangButton struct {
	widget.Button
}

func NewLangButton(label string, icon fyne.Resource, tapped func()) *LangButton {
	ret := &LangButton{}
	ret.ExtendBaseWidget(ret)
	ret.Text = label
	ret.Icon = icon
	ret.OnTapped = tapped
	return ret
}

func (l *LangButton) MinSize() fyne.Size {
	return fyne.NewSize(35, 35)
}
