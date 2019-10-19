package gitsrv

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type srv struct {
	rootDir string
	gitBinPath string
}

func GitServer(rootDir string, prefix string) http.Handler {
	g := &srv{
		rootDir: rootDir,
		gitBinPath: "/usr/bin/git",
	}

	mux := httprouter.New()

	// TODO: Add Auth
	mux.HandlerFunc(http.MethodGet, "/git/:repo/info/refs", infoRefsHandler(g))
	mux.HandlerFunc(http.MethodPost, "/git/:repo/git-receive-pack", recievePackHandler(g))

	return mux
}

func parseServiceType(input string) string {
	if !strings.HasPrefix(input, "git-") {
		return ""
	}

	return strings.Replace(input, "git-", "", 1)
}

func infoRefsHandler(g *srv) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("gsv 1")
		params := httprouter.ParamsFromContext(r.Context())

		repo := params.ByName("repo")

		serviceName := parseServiceType(r.URL.Query().Get("service"))
		log.Println("gsv 2")
		repoPath, err := g.repoPath(repo)
		log.Println("gsv 3", repoPath)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return
		}

		refs, err := g.infoRefs(repoPath, serviceName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", serviceName))
		w.WriteHeader(http.StatusOK)
		w.Write(refs)

	}
}

func recievePackHandler(g *srv) http.HandlerFunc {
	return func(w http.ResponseWriter, r * http.Request) {
		params := httprouter.ParamsFromContext(r.Context())

		repo := params.ByName("repo")
		serviceName := "receive-pack"
		w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-result", serviceName))
		g.serviceRpc(repo, serviceName, r.Body, w)
	}
}

func (g *srv) repoPath(name string) (string, error){
	repoPath := path.Join(g.rootDir, name)
	stat, err := os.Stat(repoPath)
	if  err != nil {
		return repoPath, err
	}

	if !stat.IsDir() {
		return repoPath, fmt.Errorf("not a valid repo")
	}

	return repoPath, nil
}

func (g *srv) infoRefs(repoPath string, serviceName string ) ([]byte, error) {
	args := []string{serviceName, "--stateless-rpc", "--advertise-refs", "."}
	refs, err := g.gitCommand(repoPath, args...)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	formatPacket(&buf, "# service=git-%s\n", serviceName)
	buf.Write(packetFlush())
	buf.Write(refs)

	return buf.Bytes(), nil
}

func (g *srv) serviceRpc(repo string, serviceName string, input io.ReadCloser, out io.Writer) error {
	defer input.Close()

	args := []string{serviceName, "--stateless-rpc", "."}
	cmd := exec.Command(g.gitBinPath, args...)
	cmd.Dir = path.Join(g.rootDir, repo)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		return err
	}

	io.Copy(stdin, input)
	stdin.Close()

	io.Copy(out, stdout)

	err = cmd.Wait()

	return err
}

func (g *srv) gitCommand(dir string, args ...string) ([]byte, error) {
	command := exec.Command(g.gitBinPath, args...)
	command.Dir = dir

	return command.Output()
}

func formatPacket(w io.Writer, format string, a ...interface{}) (int, error) {

	str := fmt.Sprintf(format, a...)

	s := strconv.FormatInt(int64(len(str)+4), 16)

	if len(s)%4 != 0 {
		s = strings.Repeat("0", 4-len(s)%4) + s
	}

	return w.Write([]byte(s + str))
}

func packetFlush() []byte {
	return []byte("0000")
}

