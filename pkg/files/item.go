package files

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/josa42/go-neovim/view"
)

var iconThemes = map[string][]rune{
	"nerdfont": {'', 'ﱮ', ''},
	"default":  {'▸', '▾', '•'},
}

// Interface Assertions
var _ view.TreeItem = (*FileItem)(nil)
var _ view.Openable = (*FileItem)(nil)
var _ view.Statusable = (*FileItem)(nil)

type FileItem struct {
	name        string
	path        string
	isDir       bool
	isOpen      bool
	children    []view.TreeItem
	matchIgnore *func(string) bool
	provider    *FileProvider
}

func NewFileItem(parentPath, name string, provider *FileProvider) *FileItem {
	path := filepath.Join(parentPath, name)
	item := &FileItem{
		name:     name,
		path:     path,
		isDir:    isDir(path),
		provider: provider,
	}

	return item
}

func (i *FileItem) Children() []view.TreeItem {
	names := childrenNames(i.path)
	children := []view.TreeItem{}

	for _, name := range names {

		found := false

		for _, c := range i.children {
			child, _ := c.(*FileItem)
			if child.name == name {
				children = append(children, child)
				found = true
				continue
			}
		}
		if !found {
			children = append(children, NewFileItem(i.path, name, i.provider))
		}
	}

	sort.Slice(children, func(i, j int) bool {
		a, _ := children[i].(*FileItem)
		b, _ := children[j].(*FileItem)

		if a.isDir == b.isDir {
			return a.name < b.name
		}

		return a.isDir
	})

	i.children = children

	filtered := []view.TreeItem{}

	for _, c := range children {
		i, ok := c.(*FileItem)
		if ok && !i.provider.isIgnored(i.path) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func (i *FileItem) String() string {
	icon := i.icon()
	if i.isDir {
		if i.isOpen {
			return fmt.Sprintf("%c %s/", icon, i.name)
		} else {
			return fmt.Sprintf("%c %s/", icon, i.name)
		}
	}
	return fmt.Sprintf("%c %s", icon, i.name)
}

// Openable Interface

func (i *FileItem) IsOpenable() bool {
	return i.isDir
}

func (i *FileItem) IsOpen() bool {
	return i.isOpen
}

func (i *FileItem) Open() {
	i.isOpen = true
}

func (i *FileItem) Close() {
	i.isOpen = false
}

func (i *FileItem) icon() rune {
	icons := iconThemes["default"]
	if i.provider.api.Global.Vars.Bool("nerdfont") {
		icons = iconThemes["nerdfont"]
	}
	if i.isDir {
		if !i.isOpen {
			return icons[0]
		} else {
			return icons[1]
		}
	}
	return icons[2]
}

// statusable interface

func (i *FileItem) Status() rune {
	switch i.provider.fileStatus.get(i.path, i.isDir) {
	case FileStatusChanged:
		return view.ItemStatusChanged
	case FileStatusUntracked:
		return view.ItemStatusAdded
	case FileStatusConflicted:
		return view.ItemStatusConflicted
	default:
		return ' '
	}
}
