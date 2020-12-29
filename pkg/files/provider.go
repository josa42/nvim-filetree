package files

import (
	"log"
	"path/filepath"
	"time"

	"github.com/josa42/go-gitignore"
	"github.com/josa42/go-neovim"
	"github.com/josa42/go-neovim/view"
	"github.com/josa42/nvim-filetree/pkg/actions"
	"github.com/josa42/nvim-filetree/pkg/opener"
)

// Interface Assertions
var _ view.TreeProvider = (*FileProvider)(nil)
var _ neovim.Updatable = (*FileProvider)(nil)
var _ view.ActionableTree = (*FileProvider)(nil)
var _ view.Changable = (*FileProvider)(nil)

type FileProvider struct {
	api           *neovim.Api
	root          *FileItem
	visibleItems  []*FileItem
	gitignore     gitignore.Gitignore
	changeTrigger *func()
	fileStatus    statusMap
}

func NewFileProvider(api *neovim.Api) *FileProvider {
	root := &FileItem{}

	p := &FileProvider{
		api:  api,
		root: root,
	}

	root.provider = p

	return p
}

func (p *FileProvider) FileType() string {
	return "tree"
}

func (p *FileProvider) Root() view.TreeItem {
	return p.root
}

// Update Interface

func (p *FileProvider) Update() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("FileProvider:Update() recover: %v\n", err)
		}
	}()

	p.updateRootPath()

	// TODO refactor gitignore handling
	p.gitignore, _ = gitignore.NewGitignoreFromFile(filepath.Join(p.root.path, ".gitignore"))
}

func (p *FileProvider) updateRootPath() bool {
	path := p.api.Cwd()
	if p.root.path != path {
		p.root.path = path

		return true
	}
	return false
}

func (p *FileProvider) isIgnored(path string) bool {
	pr, _ := filepath.Rel(p.root.path, path)
	return p.gitignore.Match(pr)
}

func (p *FileProvider) status(path string) rune {
	return ' '
}

// Actionable Interface

func (p *FileProvider) Actions() []view.TreeAction {

	handler := func(action string) func(i view.TreeItem) {
		return func(i view.TreeItem) {
			if f, ok := i.(*FileItem); ok {
				p.handleAction(f, action)
			}
		}
	}

	return []view.TreeAction{
		{Keys: "<CR>", Handler: handler(actions.Activate)},
		{Keys: "<2-LeftMouse>", Handler: handler(actions.ActivateFile)},
		{Keys: "<LeftRelease>", Handler: handler(actions.ToggleDir)},
		{Keys: "o", Handler: handler(actions.Activate)},
		{Keys: "e", Handler: handler(actions.Open)},
		{Keys: "t", Handler: handler(actions.OpenTab)},
		{Keys: "v", Handler: handler(actions.OpenVerticalSplit)},
		{Keys: "s", Handler: handler(actions.OpenHorizontalSplit)},
		{Keys: "<ESC>", Handler: handler(actions.Unfocus)},
		{Keys: "h", Handler: handler(actions.Help)},
	}
}

func (p *FileProvider) handleAction(i *FileItem, action string) {
	switch action {
	case actions.Activate:
		if i.isDir {
			i.isOpen = !i.isOpen
		} else {
			opener.Activate(p.api, i.path)
		}

	case actions.ToggleDir:
		if i.isDir {
			i.isOpen = !i.isOpen
		}

	case actions.ActivateFile:
		if !i.isDir {
			opener.Activate(p.api, i.path)
		}

	case actions.Open:
		opener.Open(p.api, i.path)

	case actions.OpenTab:
		opener.OpenTab(p.api, i.path)

	case actions.OpenHorizontalSplit:
		opener.OpenHoricontalSplit(p.api, i.path)

	case actions.OpenVerticalSplit:
		opener.OpenVerticalSplit(p.api, i.path)

	case actions.Unfocus:
		opener.FocusEditor(p.api)

	case actions.Help:
		p.api.Out.Print("?: Help - (o)pen - (e)dit - (t)ab - (s)plit - (v)ertical split - ESC unfocus")
	}
}

func (p *FileProvider) Listen(changed func()) {
	if p.changeTrigger == nil {
		p.runChangeListener()
	}
	p.changeTrigger = &changed
}

func (p *FileProvider) Unlisten() {
	p.changeTrigger = nil
}

func (p *FileProvider) runChangeListener() {

	go func() {
		gitAvailable := isGitAvailable()

		nextGitRun := time.Now()

		// TODO Fix this! Listen tochangechanges instead?
		for {
			time.Sleep(1 * time.Second)
			if p.changeTrigger == nil {
				break
			}

			pc := p.updateRootPath()

			sc := false
			if gitAvailable {
				sc, nextGitRun = p.updateFileStatus(nextGitRun)
			}

			if pc || sc {
				t := *p.changeTrigger
				t()
			}
		}
	}()
}

func (p *FileProvider) updateFileStatus(nextRun time.Time) (bool, time.Time) {

	if nextRun.After(time.Now()) {
		return false, nextRun
	}

	if !isGitRepo(p.root.path) {
		return false, time.Now().Add(30 * time.Second)
	}

	fs := updateStatus(p.root.path)

	if p.fileStatus.hashChanges(fs) {
		p.fileStatus = fs
		return true, time.Now().Add(5 * time.Second)
	}

	return false, time.Now().Add(5 * time.Second)
}
