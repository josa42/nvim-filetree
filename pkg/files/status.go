package files

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	FileStatusNormal = iota
	FileStatusIgnored
	FileStatusChanged
	FileStatusUntracked
	FileStatusConflicted
)

var (
	expStatusLine    = regexp.MustCompile(`^(..) (.*)$`)
	expTrailingSlash = regexp.MustCompile(`/*$`)
)

type status int

type statusMap map[string]status

func (s statusMap) get(path string, dir bool) status {

	path = strings.TrimRight(path, `/`)

	if dir {
		if s.dirContains(path, FileStatusConflicted) {
			return FileStatusConflicted
		}

		if s.dirContains(path, FileStatusChanged) {
			return FileStatusChanged
		}

		if s.dirContains(path, FileStatusUntracked) {
			return FileStatusUntracked
		}
	} else if fs, ok := s[path]; ok {
		return fs
	}
	return FileStatusNormal
}

func (s statusMap) dirContains(path string, fileStatus status) bool {
	path = expTrailingSlash.ReplaceAllString(path, "/")
	for p, fs := range s {
		if fs == fileStatus && strings.HasPrefix(p, path) {
			return true
		}
	}

	return false
}

func (s statusMap) hashChanges(s2 statusMap) bool {
	j1, _ := json.Marshal(s)
	j2, _ := json.Marshal(s2)

	return string(j1) != string(j2)
}

func updateStatus(dir string) statusMap {
	s := statusMap{}

	// TODO handle: not in repo; git bin does not exist
	cmdToplevel := git("rev-parse", "--show-toplevel")
	cmdToplevel.Dir = dir

	rootDirB, err := cmdToplevel.Output()
	if err != nil {
		log.Printf("git rev-parse - err: %v", err)
		return s
	}

	rootDir := strings.TrimSpace(string(rootDirB))

	cmdStatus := git("status", "--porcelain", "--ignored")
	cmdStatus.Dir = dir
	out, err := cmdStatus.Output()
	if err != nil {
		log.Printf("git status - err: %v", err)
		return s
	}

	for _, l := range strings.Split(string(out), "\n") {
		p := expStatusLine.FindStringSubmatch(l)
		if len(p) > 0 {
			path := strings.TrimRight(filepath.Join(rootDir, p[2]), `/`)
			s[path] = getStatus(p[1])
		}
	}

	return s
}

func isGitAvailable() bool {
	cmd := git("--version")
	err := cmd.Run()

	if err != nil {
		log.Printf("err: %v", err)
	}

	return err == nil
}

func isGitRepo(dir string) bool {

	gitPath := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitPath); err == nil {

		cmd := git("status")
		cmd.Dir = dir
		err := cmd.Run()

		return err == nil
	}
	return false
}

// ' ' = unmodified
// M   = modified
// A   = added
// D   = deleted
// R   = renamed
// C   = copied
// U   = updated but unmerged

var (
	//  [AMD]   = not updated
	// M[ MD]   = updated in index
	// A[ MD]   = added to index
	// R[ MD]   = renamed in index
	// C[ MD]   = copied in index
	// [MARC]   = index and work tree matches
	// [ MARC]M = work tree changed since index
	// [ D]R    = renamed in work tree
	// [ D]C    = copied in work tree
	expChanged = regexp.MustCompile(`^( [AMD]|M[ MD]|A[ MD]|R[ MD]|C[ MD]|[MARC] |[ MARC]M|[ D]R|[ D]C)$`)

	// [ MARC]D = deleted in work tree
	// D        =deleted from index
	expDeleted = regexp.MustCompile(`^([ MARC]D|D )$`)

	// DD = unmerged, both deleted
	// AU = unmerged, added by us
	// UD = unmerged, deleted by them
	// UA = unmerged, added by them
	// DU = unmerged, deleted by us
	// AA = unmerged, both added
	// UU = unmerged, both modified
	expConflicted = regexp.MustCompile(`^(DD|AU|UD|UA|DU|AA|UU)$`)

	// untracked
	expUntracked = regexp.MustCompile(`^\?\?$`)

	// ignored
	expIgnored = regexp.MustCompile(`^\!\!$`)
)

func getStatus(m string) status {

	if expConflicted.MatchString(m) {
		return FileStatusConflicted
	}

	if expUntracked.MatchString(m) {
		return FileStatusUntracked
	}

	if expIgnored.MatchString(m) {
		return FileStatusIgnored
	}

	if expChanged.MatchString(m) {
		return FileStatusChanged
	}

	if m != "  " {
		log.Printf("default: '%s'", m)
	}
	return FileStatusNormal
}

// TODO use go-git/go-git instead of spawning a process?
func git(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.Env = append(os.Environ(),
		"GIT_OPTIONAL_LOCKS=0", // ignored
	)
	return cmd
}
