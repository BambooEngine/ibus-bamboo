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
	IBautoCommitWithMouseMovement
	_IBemojiDisabled //deprecated
	IBpreeditElimination
	_IBinputModeLookupTableEnabled //deprecated
	IBautoCapitalizeMacro
	_IBimQuickSwitchEnabled     //deprecated
	_IBrestoreKeyStrokesEnabled //deprecated
	IBmouseCapturing
	IBworkaroundForFBMessenger
	IBworkaroundForWPS
	IBstdFlags = IBspellCheckEnabled | IBspellCheckWithRules | IBautoNonVnRestore | IBddFreeStyle |
		IBmouseCapturing | IBautoCapitalizeMacro | IBnoUnderline | IBworkaroundForWPS
	IBUsStdFlags = 0
)
