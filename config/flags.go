package config

const (
	PreeditIM = iota + 1
	SurroundingTextIM
	BackspaceForwardingIM
	ShiftLeftForwardingIM
	ForwardAsCommitIM
	XTestFakeKeyEventIM
	UsIM
)

var ImLookupTable = map[int]string{
	PreeditIM:             "Cấu hình mặc định (Pre-edit)",
	SurroundingTextIM:     "Sửa lỗi gạch chân (Surrounding Text)",
	BackspaceForwardingIM: "Sửa lỗi gạch chân (ForwardKeyEvent I)",
	ShiftLeftForwardingIM: "Sửa lỗi gạch chân (ForwardKeyEvent II)",
	ForwardAsCommitIM:     "Sửa lỗi gạch chân (Forward as commit)",
	XTestFakeKeyEventIM:   "Sửa lỗi gạch chân (XTestFakeKeyEvent)",
	UsIM:                  "Thêm vào danh sách loại trừ",
}

var ImShortNames = map[int]string{
	PreeditIM:             "preedit",
	SurroundingTextIM:     "surrounding-text",
	BackspaceForwardingIM: "forward-key-1",
	ShiftLeftForwardingIM: "forward-key-2",
	ForwardAsCommitIM:     "forward-as-commit",
	XTestFakeKeyEventIM:   "xtest",
	UsIM:                  "exclude",
}

var ImBackspaceList = []int{
	SurroundingTextIM,
	BackspaceForwardingIM,
	ShiftLeftForwardingIM,
	ForwardAsCommitIM,
	XTestFakeKeyEventIM,
}

const (
	IBautoCommitWithVnNotMatch uint = 1 << iota
	IBmacroEnabled
	_IBautoCommitWithVnFullMatch //deprecated
	_IBautoCommitWithVnWordBreak //deprecated
	IBspellCheckEnabled
	IBautoNonVnRestore
	IBddFreeStyle
	IBnoUnderline
	IBspellCheckWithRules
	IBspellCheckWithDicts
	IBautoCommitWithDelay
	_IBautoCommitWithMouseMovement //deprecated
	_IBemojiDisabled               //deprecated
	IBpreeditElimination
	_IBinputModeLookupTableEnabled //deprecated
	IBautoCapitalizeMacro
	_IBimQuickSwitchEnabled     //deprecated
	_IBrestoreKeyStrokesEnabled //deprecated
	_IBmouseCapturing           //deprecated
	IBworkaroundForFBMessenger
	IBworkaroundForWPS
	IBstdFlags = IBspellCheckEnabled | IBspellCheckWithRules | IBautoNonVnRestore | IBddFreeStyle |
		IBautoCapitalizeMacro | IBnoUnderline | IBworkaroundForWPS
	IBUsStdFlags = 0
)
