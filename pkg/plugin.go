package main

import (
	"log"
	"strings"

	"github.com/josa42/go-neovim"
	"github.com/josa42/go-neovim/view"
	"github.com/josa42/nvim-filetree/pkg/files"
	"github.com/josa42/nvim-filetree/pkg/opener"
)

func main() {
	defer neovim.SetupLogging()()
	neovim.Register(&TreePlugin{})
}

const (
	WindowWidth = 40
)

const (
	BufferVarIsTree        = "is_tree"
	BufferVarHideLightline = "lightline_hidden"
	GlobalVarTreeBufferID  = "tree_buffer_id"
	GlobalVarIsTreeOpen    = "tree_open"
	GlobalVarIsTreeOpening = "tree_opening"
)

type Tree interface {
	Render(*neovim.Api, neovim.Buffer)
	Action(*neovim.Api, int, string)
}

type TreePlugin struct {
	api      *neovim.Api
	treeView *view.TreeView
}

func (tp *TreePlugin) Register(api neovim.RegisterApi) {
	if a, ok := api.(*neovim.Api); ok {
		tp.api = a
	}

	api.Function("TreeOpen", tp.Open)
	api.Function("TreeClose", tp.Close)
	api.Function("TreeToggle", tp.Toggle)
	api.Function("TreeUnfocus", tp.Unfocus)
	api.Function("TreeFocus", tp.Focus)
	api.Function("TreeToggleFocus", tp.ToggleFocus)
}

func (tp *TreePlugin) Activate(api *neovim.Api) {
	tp.api = api

	tp.treeView = view.NewTreeView(files.NewFileProvider(api))

	api.Global.On(neovim.EventBufWinEnter, tp.onEnterSyncState)
	api.Global.On(neovim.EventWinEnter, tp.onEnterSyncState)
	api.Global.On(neovim.EventBufEnter, tp.onLeaveCloseLastTree)
	api.Global.On(neovim.EventWinLeave, tp.onLeaveUnfocusTree)
}

func (p *TreePlugin) Close() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Close() recover: %v\n", err)
		}
	}()
	p.api.Global.Vars.SetBool(GlobalVarIsTreeOpen, false)
	if b, found := p.getTreeBuffer(); found {
		b.Close()
	}
}

func (p *TreePlugin) Focus() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Focus() recover: %v\n", err)
		}
	}()

	if !p.treeBufferHasFocus() {
		if b, ok := p.getTreeBuffer(); ok {
			tab := p.api.CurrentTab()

			win, found := tab.FindWindow(func(win *neovim.Window) bool {
				wb := win.Buffer()
				return wb.ID() == b.ID()
			})

			if found {
				win.Focus()
			}
		}
	}
}

func (p *TreePlugin) Unfocus() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Unfocus() recover: %v\n", err)
		}
	}()
	if p.treeBufferHasFocus() {
		opener.FocusEditor(p.api)
	}
}

func (p *TreePlugin) Open() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Open() recover: %v\n", err)
		}
	}()
	if p.api.Global.Vars.Bool(GlobalVarIsTreeOpening) {
		return
	}

	p.api.Global.Vars.SetBool(GlobalVarIsTreeOpening, true)
	p.api.Global.Vars.SetBool(GlobalVarIsTreeOpen, true)

	buffer := p.getOrCreateBuffer()

	if !p.hasTreeBuffer() {
		p.attachTreeBuffer(buffer)
	}

	p.api.Global.Vars.SetBool(GlobalVarIsTreeOpening, false)
}

func (p *TreePlugin) onLeaveCloseLastTree() {
	if p.api.Global.Vars.Bool(GlobalVarIsTreeOpening) {
		return
	}

	if p.hasOnlyTreeBuffer() {
		tab := p.api.CurrentTab()
		tab.Close(true)
	}
}

// Sync open file tree across tabs
func (p *TreePlugin) onEnterSyncState() {
	if p.api.Global.Vars.Bool(GlobalVarIsTreeOpening) {
		return
	}

	if p.api.Global.Vars.Bool(GlobalVarIsTreeOpen) {
		focus := p.treeBufferHasFocus()
		p.Open()

		if !focus {
			p.Unfocus()
		}
	} else {
		p.Close()
	}
}

func (p *TreePlugin) onLeaveUnfocusTree() {
	b := p.api.CurrentBuffer()

	if b.Vars.Bool(BufferVarIsTree) {
		tab := p.api.CurrentTab()
		window, _ := tab.FindWindow(func(window *neovim.Window) bool {
			return !window.Buffer().Vars.Bool(BufferVarIsTree)
		})

		window.Focus()
	}
}

func (p *TreePlugin) Toggle() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Toggle() recover: %v\n", err)
		}
	}()
	if p.api.Global.Vars.Bool(GlobalVarIsTreeOpen) {
		p.Close()
	} else {
		p.Open()
		p.Focus()
	}
}

func (p *TreePlugin) ToggleFocus() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("ToggleFocus() recover: %v\n", err)
		}
	}()
	if p.treeBufferHasFocus() {
		p.Unfocus()
	} else if p.hasTreeBuffer() {
		p.Focus()
	} else {
		p.Open()
	}
}

func (p *TreePlugin) getOrCreateBuffer() *neovim.Buffer {
	if b, ok := p.getTreeBuffer(); ok {
		return b
	}
	return p.createTreeBuffer()
}

func (p *TreePlugin) getTreeBuffer() (*neovim.Buffer, bool) {
	return p.api.FindBuffer(func(buffer *neovim.Buffer) bool {
		return buffer.Vars.Bool(BufferVarIsTree)
	})
}

func (p *TreePlugin) createTreeBuffer() *neovim.Buffer {
	p.api.Global.Vars.SetBool(GlobalVarIsTreeOpening, true)
	defer p.api.Global.Vars.SetBool(GlobalVarIsTreeOpening, false)

	buffer := p.api.CreateSplitBuffer(WindowWidth, neovim.SplitTopLeft, neovim.SplitVertical)
	buffer.Vars.SetBool(BufferVarIsTree, true)
	buffer.Vars.SetBool(BufferVarHideLightline, true)
	buffer.Options.SetFileType("tree")
	buffer.SetTitle("פּ")

	p.api.Executef("setlocal %s", strings.Join([]string{
		"cursorline",
		"foldcolumn=0",
		"nonumber",
		"foldmethod=manual",
		"nocursorcolumn",
		"nofoldenable",
		"nolist",
		"norelativenumber",
		"nospell",
		"nowrap",
		"signcolumn=no",
		"colorcolumn=",
	}, " "))
	p.api.Execute("iabclear <buffer>")

	p.api.Renderer.Attach(buffer, p.treeView)

	p.api.Global.Vars.SetInt(GlobalVarTreeBufferID, buffer.ID())
	p.api.Global.Vars.SetBool(GlobalVarIsTreeOpening, false)

	return buffer
}

func (p *TreePlugin) attachTreeBuffer(b *neovim.Buffer) {
	p.api.Global.Vars.SetInt(GlobalVarTreeBufferID, b.ID())

	p.api.Executef("topleft vertical %d new | buffer %d", WindowWidth, b.ID())

	// window
	win := p.api.CurrentWindow()
	wo := win.Options
	wo.SetFixWidth(true)
}

func (p *TreePlugin) treeBufferHasFocus() bool {
	if idx := p.api.Global.Vars.Int(GlobalVarTreeBufferID); idx > 0 {
		buffer := p.api.CurrentBuffer()
		return idx == buffer.ID()
	}
	return false
}

func (p *TreePlugin) hasTreeBuffer() bool {
	tab := p.api.CurrentTab()
	idx := p.api.Global.Vars.Int(GlobalVarTreeBufferID)

	return idx > 0 && tab.HasBufferID(idx)
}

func (p *TreePlugin) hasOnlyTreeBuffer() bool {
	tab := p.api.CurrentTab()
	if bID := p.api.Global.Vars.Int(GlobalVarTreeBufferID); bID != 0 {
		return tab.HasBufferID(bID) && len(tab.Windows()) == 1
	}
	return false
}
