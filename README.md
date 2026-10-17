# gotk4-vte

This is a binding for libvte (the GTK4 build, `vte-2.91-gtk4`) and gotk4

The binding acts as an extension to https://github.com/diamondburned/gotk4

This is based on napsy/go-vte, it aims to be the most complete binding for libvte.

Work in progress!

## Requirements

- cgo
- the GTK4 build of libvte with its development files: `libvte-2.91-gtk4-dev` on Debian/Ubuntu, `vte4` on Arch/Manjaro
- `Terminal.SetContextMenuModel` requires vte 0.76 or later. With an older vte it compiles, but does nothing and
  returns false.

## Supported API

| Area          | Methods                                                                       |
|---------------|-------------------------------------------------------------------------------|
| Create, spawn | `TerminalNew`, `SpawnAsyncSimple`                                             |
| Input/output  | `Feed`, `FeedChild`                                                           |
| Size          | `GetColumnCount`, `GetRowCount`, `SetSize`, `SetFontScale`                    |
| Clipboard     | `CopyClipboardFormat`, `CopyPrimary`, `PasteClipboard`, `PastePrimary`        |
| Selection     | `SelectAll`, `UnselectAll`, `GetHasSelection`, `ConnectSelectionChanged`      |
| Context menu  | `SetContextMenuModel` (vte 0.76+)                                             |

`Terminal` embeds `gtk.Widget`, so all widget methods (event controllers, action groups etc.) are available too.

## Environment variables

There is currently no official TERM value that denotes truecolor support, because terminfo/termcap does not
have a standard entry for 24-bit color. You should set the following environment variables to take advantage
of all VTE features:

```bash
export TERM=xterm-256color
export COLORTERM=truecolor
```


## Example


```go
term, err := vte.TerminalNew()
if err != nil {
	log.Fatal("Unable to create terminal:", err)
}
wd, _ := os.Getwd()
term.SpawnAsyncSimple(wd, []string{"/bin/bash"}, []string{"TERM=xterm"})
```

### Copy and paste

Selecting text with the mouse and pasting the primary selection with the middle mouse button are handled by vte.
Clipboard shortcuts and the context menu are up to the application, for example:

```go
actions := gio.NewSimpleActionGroup()
copyAction := gio.NewSimpleAction("copy", nil)
copyAction.ConnectActivate(func(_ *glib.Variant) { term.CopyClipboardFormat(vte.FormatText) })
actions.AddAction(copyAction)
pasteAction := gio.NewSimpleAction("paste", nil)
pasteAction.ConnectActivate(func(_ *glib.Variant) { term.PasteClipboard() })
actions.AddAction(pasteAction)
term.InsertActionGroup("term", actions)

// Ctrl+Shift+C and Ctrl+Shift+V. The callbacks return true, so the key press is never sent to the child process
// (vte would send Ctrl+Shift+C as Ctrl+C).
shortcuts := gtk.NewShortcutController()
shortcuts.SetPropagationPhase(gtk.PhaseCapture)
for trigger, action := range map[string]string{"<Control><Shift>c": "copy", "<Control><Shift>v": "paste"} {
	shortcuts.AddShortcut(gtk.NewShortcut(
		gtk.NewShortcutTriggerParseString(trigger),
		gtk.NewCallbackAction(func(_ gtk.Widgetter, _ *glib.Variant) bool {
			actions.ActivateAction(action, nil)
			return true
		}),
	))
}
term.AddController(shortcuts)

// Right click menu (vte 0.76+)
menu := gio.NewMenu()
menu.Append("Copy", "term.copy")
menu.Append("Paste", "term.paste")
term.SetContextMenuModel(menu)
```

Use `ConnectSelectionChanged` and `GetHasSelection` to enable the copy action only when there is a selection.
