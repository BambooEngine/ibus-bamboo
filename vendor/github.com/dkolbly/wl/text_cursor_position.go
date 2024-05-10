package wl

const (
	_TEXT_CURSOR_POSITION_NOTIFY = 0
)

type TextCursorPosition struct {
	BaseProxy
}

func NewTextCursorPosition(conn *Context) *TextCursorPosition {
	ret := new(TextCursorPosition)
	conn.Register(ret)
	return ret
}

func (p *TextCursorPosition) Notify(surface *Surface, x float32, y float32) error {
	return p.Context().SendRequest(p, _TEXT_CURSOR_POSITION_NOTIFY, surface, x, y)
}
