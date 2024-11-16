package ui

import (
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const APP_ID = "ibus-bamboo.mode-options"

var modeOptions = []string{
	"1. Pre-edit (có gạch chân)",
	"2. Surrounding Text (không gạch chân)",
	"3. ForwardKeyEvent I (không gạch chân)",
	"4. ForwardKeyEvent II (không gạch chân)",
	"5. Forward as Commit (không gạch chân)",
	"6. XTestFakeKeyEvent (không gạch chân)",
}

func saveShortCut() error {
	return nil
}

func renderShortcut(window *gtk.ApplicationWindow) *gtk.Box {
	// Button
	buttonClose := gtk.NewButtonWithLabel("Đóng")
	buttonClose.ConnectClicked(func() {
		window.Close()
	})
	buttonReset := gtk.NewButtonWithLabel("Đặt lại")
	buttonSave := gtk.NewButtonWithLabel("Lưu")

	// Box
	buttonBox := gtk.NewBox(gtk.OrientationHorizontal, 10)
	buttonBox.SetHAlign(gtk.AlignEnd)
	buttonBox.SetHExpand(true)
	buttonBox.Append(buttonClose)
	buttonBox.Append(buttonReset)
	buttonBox.Append(buttonSave)

	// MainBox
	mainBox := gtk.NewBox(gtk.OrientationVertical, 10)
	mainBox.SetMarginTop(10)
	mainBox.SetMarginBottom(10)
	mainBox.SetMarginStart(10)
	mainBox.SetMarginEnd(10)
	mainBox.Append(buttonBox)

	return mainBox
}

func saveInputTextView() error {
	return nil
}

func renderInputTextView() *gtk.Box {
	// TextView
	textView := gtk.NewTextView()
	textView.SetWrapMode(gtk.WrapNone)
	textView.SetVExpand(true)
	textView.SetHExpand(true)
	textView.SetTopMargin(10)
	textView.SetBottomMargin(10)
	textView.SetLeftMargin(10)
	textView.SetRightMargin(10)

	// CSS for textView
	textView.AddCSSClass("bordered-textview")
	css := `
	.bordered-textview text {
		border: 1px solid #999999;
		border-radius: 5px;
		padding: 5px;
	}
	`
	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromString(css)
	gtk.StyleContextAddProviderForDisplay(textView.Display(), cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	// ScrollView for TextView input
	scrollView := gtk.NewScrolledWindow()
	scrollView.SetVExpand(true)
	scrollView.SetMinContentHeight(200)
	scrollView.SetChild(textView)

	// Button
	buttonSave := gtk.NewButtonWithLabel("Lưu")
	buttonBox := gtk.NewBox(gtk.OrientationHorizontal, 0)
	buttonBox.SetMarginBottom(10)
	buttonBox.SetMarginEnd(10)
	buttonBox.SetHAlign(gtk.AlignEnd)
	buttonBox.SetHExpand(true)
	buttonBox.Append(buttonSave)

	// MainBox
	mainBox := gtk.NewBox(gtk.OrientationVertical, 10)
	mainBox.SetMarginTop(10)
	mainBox.SetMarginBottom(10)
	mainBox.SetMarginStart(10)
	mainBox.SetMarginEnd(10)
	mainBox.Append(scrollView)
	mainBox.Append(buttonBox)

	return mainBox
}

func renderOther(window *gtk.ApplicationWindow) *gtk.Box {
	// Mode Choose
	labelDefaultMode := gtk.NewLabel("Chế độ gõ mặc định")
	dropdownMode := gtk.NewDropDownFromStrings(modeOptions)
	modeBox := gtk.NewBox(gtk.OrientationHorizontal, 10)
	modeBox.Append(labelDefaultMode)
	modeBox.Append(dropdownMode)

	// Checkbox
	checkboxFixFB := gtk.NewCheckButtonWithLabel("Sửa lỗi lặp chữ trong FB")
	checkboxFixWPS := gtk.NewCheckButtonWithLabel("Sửa lỗi không hiện chữ trong WPS")

	// Button
	buttonClose := gtk.NewButtonWithLabel("Đóng")
	buttonClose.ConnectClicked(func() {
		window.Close()
	})
	buttonBox := gtk.NewBox(gtk.OrientationHorizontal, 10)
	buttonBox.Append(buttonClose)
	buttonBox.SetHAlign(gtk.AlignEnd)
	buttonBox.SetHExpand(true)

	// MainBox
	mainBox := gtk.NewBox(gtk.OrientationVertical, 10)
	mainBox.SetMarginTop(10)
	mainBox.SetMarginBottom(10)
	mainBox.SetMarginStart(10)
	mainBox.SetMarginEnd(10)
	mainBox.Append(modeBox)
	mainBox.Append(checkboxFixFB)
	mainBox.Append(checkboxFixWPS)
	mainBox.Append(buttonBox)

	return mainBox
}

func OpenGUI(engName string) {
	app := gtk.NewApplication(APP_ID, gio.ApplicationDefaultFlags)
	app.ConnectActivate(func() {
		// Main app
		window := gtk.NewApplicationWindow(app)
		window.SetTitle("ibus-bamboo shortcut options")
		window.SetDecorated(true)
		window.SetDefaultSize(600, 300)

		// Tabs Notebook
		notebook := gtk.NewNotebook()
		notebook.AppendPage(renderShortcut(window), gtk.NewLabel("Phím tắt"))
		notebook.AppendPage(renderInputTextView(), gtk.NewLabel("Gõ tắt"))
		notebook.AppendPage(renderInputTextView(), gtk.NewLabel("Tự định nghĩa kiểu gõ"))
		notebook.AppendPage(renderOther(window), gtk.NewLabel("Khác"))
		window.SetChild(notebook)

		window.ConnectCloseRequest(func() bool {
			window.Destroy()
			app.Quit()
			return true
		})

		window.SetVisible(true)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
