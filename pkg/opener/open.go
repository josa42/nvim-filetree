package opener

import (
	"time"

	"github.com/josa42/go-neovim"
)

func Activate(api *neovim.Api, path string) {
	FocusEditor(api)

	if win, found := findWindow(api, path); found {
		win.Focus()

	} else if b := api.CurrentBuffer(); b.IsEmpty() {
		Open(api, path)

	} else {
		OpenTab(api, path)
	}
}

func Open(api *neovim.Api, path string) {
	FocusEditor(api)
	api.Executef("silent edit %s", path)
}

func OpenTab(api *neovim.Api, path string) {
	api.Executef("silent tabe %s", path)
	time.Sleep(200 * time.Millisecond)
	FocusEditor(api)
}

func OpenVerticalSplit(api *neovim.Api, path string) {
	FocusEditor(api)
	api.Executef("silent vsplit %s", path)
}

func OpenHoricontalSplit(api *neovim.Api, path string) {
	FocusEditor(api)
	api.Executef("silent split %s", path)
}

func FocusEditor(api *neovim.Api) {
	defer func() {
		if err := recover(); err != nil {
			// n.Command(fmt.Sprintf("echom '%v'", err))
		}
	}()

	tab := api.CurrentTab()

	for _, win := range tab.Windows() {
		if !win.Buffer().Vars.Bool("is_tree") {
			win.Focus()
			return
		}
	}
}

func findWindow(api *neovim.Api, path string) (*neovim.Window, bool) {

	return api.FindWindow(func(win *neovim.Window) bool {
		return win.Buffer().Path() == path
	})
}

// TODO Usage window?
// if bnum !=# -1 && getbufvar(bnum, '&buftype') ==# ''
//                     \ && !getwinvar(i, '&previewwindow')
//                     \ && (!getbufvar(bnum, '&modified') || &hidden)
//             return i
//         endif
