package vte

// #cgo pkg-config: vte-2.91-gtk4
// #include <vte/vte.h>
// #include "vte.go.h"
// #include <stdlib.h>
import "C"

import (
	"errors"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"unsafe"
)

var nilPtrErr = errors.New("cgo returned unexpected nil pointer")

type Terminal struct {
	gtk.Widget
}

func wrapWidget(obj *glib.Object) *gtk.Widget {
	return &gtk.Widget{
		InitiallyUnowned: glib.InitiallyUnowned{
			Object: obj,
		},
		Object: obj,
		Accessible: gtk.Accessible{
			Object: obj,
		},
		Buildable: gtk.Buildable{
			Object: obj,
		},
		ConstraintTarget: gtk.ConstraintTarget{
			Object: obj,
		},
	}
}

func wrapTerminal(obj *glib.Object) *Terminal {
	return &Terminal{*wrapWidget(obj)}
}

func (t *Terminal) native() *C.VteTerminal {
	p := unsafe.Pointer(t.Widget.Object.Native())
	return C.toVteTerminal(p)
}

func TerminalNew() (*Terminal, error) {
	t := C.vte_terminal_new()
	if t == nil {
		return nil, nilPtrErr
	}
	obj := glib.Take(unsafe.Pointer(t))
	return wrapTerminal(obj), nil
}

func (t *Terminal) Feed(text string) {
	str := C.CString(text)
	defer C.free(unsafe.Pointer(str))
	C.vte_terminal_feed(t.native(), str, C.gssize(len(text)))
}

func (t *Terminal) FeedChild(text string) {
	str := C.CString(text)
	defer C.free(unsafe.Pointer(str))
	C.vte_terminal_feed_child(t.native(), str, C.gssize(len(text)))
}

// GetColumnCount https://gnome.pages.gitlab.gnome.org/vte/gtk3/method.Terminal.get_column_count.html
func (t *Terminal) GetColumnCount() uint16 {
	return uint16(C.vte_terminal_get_column_count(t.native()))
}

// GetRowCount https://gnome.pages.gitlab.gnome.org/vte/gtk3/method.Terminal.get_row_count.html
func (t *Terminal) GetRowCount() uint16 {
	return uint16(C.vte_terminal_get_row_count(t.native()))
}

// SetSize https://gnome.pages.gitlab.gnome.org/vte/gtk3/method.Terminal.set_size.html
func (t *Terminal) SetSize(columns uint16, rows uint16) {
	C.vte_terminal_set_size(t.native(), C.long(columns), C.long(rows))
}

// SetSize https://gnome.pages.gitlab.gnome.org/vte/gtk3/method.Terminal.set_size.html
func (t *Terminal) SetFontScale(scale float64) {
	C.vte_terminal_set_font_scale(t.native(), C.double(scale))
}

// Format https://gnome.pages.gitlab.gnome.org/vte/gtk4/enum.Format.html
type Format int

const (
	FormatText Format = C.VTE_FORMAT_TEXT
	FormatHTML Format = C.VTE_FORMAT_HTML
)

// CopyClipboardFormat https://gnome.pages.gitlab.gnome.org/vte/gtk4/method.Terminal.copy_clipboard_format.html
func (t *Terminal) CopyClipboardFormat(format Format) {
	C.vte_terminal_copy_clipboard_format(t.native(), C.VteFormat(format))
}

// CopyPrimary https://gnome.pages.gitlab.gnome.org/vte/gtk4/method.Terminal.copy_primary.html
func (t *Terminal) CopyPrimary() {
	C.vte_terminal_copy_primary(t.native())
}

// PasteClipboard https://gnome.pages.gitlab.gnome.org/vte/gtk4/method.Terminal.paste_clipboard.html
func (t *Terminal) PasteClipboard() {
	C.vte_terminal_paste_clipboard(t.native())
}

// PastePrimary https://gnome.pages.gitlab.gnome.org/vte/gtk4/method.Terminal.paste_primary.html
func (t *Terminal) PastePrimary() {
	C.vte_terminal_paste_primary(t.native())
}

// SelectAll https://gnome.pages.gitlab.gnome.org/vte/gtk4/method.Terminal.select_all.html
func (t *Terminal) SelectAll() {
	C.vte_terminal_select_all(t.native())
}

// UnselectAll https://gnome.pages.gitlab.gnome.org/vte/gtk4/method.Terminal.unselect_all.html
func (t *Terminal) UnselectAll() {
	C.vte_terminal_unselect_all(t.native())
}

// GetHasSelection https://gnome.pages.gitlab.gnome.org/vte/gtk4/method.Terminal.get_has_selection.html
func (t *Terminal) GetHasSelection() bool {
	return C.vte_terminal_get_has_selection(t.native()) != 0
}

// ConnectSelectionChanged https://gnome.pages.gitlab.gnome.org/vte/gtk4/signal.Terminal.selection-changed.html
func (t *Terminal) ConnectSelectionChanged(f func()) glib.SignalHandle {
	return t.Connect("selection-changed", f)
}

// SetContextMenuModel https://gnome.pages.gitlab.gnome.org/vte/gtk4/method.Terminal.set_context_menu_model.html
// It requires vte 0.76 or later, and returns false (doing nothing) when built with an older vte.
func (t *Terminal) SetContextMenuModel(model gio.MenuModeller) bool {
	var m *C.GMenuModel
	if model != nil {
		m = (*C.GMenuModel)(unsafe.Pointer(glib.InternObject(model).Native()))
	}
	return C.set_context_menu_model(t.native(), m) != 0
}

func makeStrings(array []string) **C.char {
	cArray := C.make_strings(C.int(len(array) + 1))
	for i, e := range array {
		cstr := C.CString(e)
		C.set_string(cArray, C.int(i), (*C.char)(cstr))
	}
	C.set_string(cArray, C.int(len(array)), nil)
	return cArray
}

func destroyStrings(strings **C.char, count int) {
	C.destroy_strings(strings, C.int(count))
}

func (t *Terminal) SpawnAsyncSimple(workingDirectory string, argv, envv []string) error {
	wd := C.CString(workingDirectory)
	defer C.free(unsafe.Pointer(wd))

	argvStrings := makeStrings(argv)
	envvStrings := makeStrings(envv)
	defer destroyStrings(argvStrings, len(argv))
	defer destroyStrings(envvStrings, len(envv))

	timeout := -1 // in ms (wait indefinitely)

	C.vte_terminal_spawn_async(
		t.native(),     // terminal
		0,              // VTE_PTY_DEFAULT
		wd,             // working directory
		argvStrings,    // argv
		envvStrings,    // env
		0,              // G_SPAWN_DEFAULT
		nil,            // child_setup
		nil,            // child_setup_data
		nil,            // child_setup_data_destroy
		C.int(timeout), // timeout
		nil,            // cancellable
		nil,            // command_spawned
		nil,            // user_data
	)
	return nil
}
