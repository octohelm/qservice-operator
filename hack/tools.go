package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var reMod = regexp.MustCompile(`^module ([^\n]+)\n`)

func main() {
	cwd, _ := os.Getwd()

	data, err := ioutil.ReadFile(filepath.Join(cwd, "go.mod"))
	if err != nil {
		panic(err)
	}

	pkg := string(reMod.FindAllSubmatch(data, -1)[0][1])
	pkgAPIBase := fmt.Sprintf("%s/pkg/apis", pkg)
	pkgClient := fmt.Sprintf("%s/pkg/client", pkg)
	pkgClientClientset := fmt.Sprintf("%s/clientset", pkgClient)
	pkgClientListers := fmt.Sprintf("%s/listers", pkgClient)
	pkgClientInformers := fmt.Sprintf("%s/informers", pkgClient)
	clientsetName := "versioned"

	args := os.Args[1:]
	gens := args[0]

	apiGroups := strings.Split(os.Getenv("API_GROUPS"), ",")

	pkgAPIs := pkgPathsFix(pkgAPIBase, apiGroups...)

	if strings.Contains(gens, "deepcopy") {
		run("go", "install", "k8s.io/code-generator/cmd/deepcopy-gen")

		run(
			"deepcopy-gen",
			"--go-header-file", path.Join(cwd, "./hack/boilerplate.go.txt"),
			"--input-dirs", pkgPathsFix(pkg, args[1:]...),
			"--output-file-base", "zz_generated.deepcopy",
			"--bounding-dirs", pkgAPIBase,
		)
	}

	if strings.Contains(gens, "client") {
		run("go", "install", "k8s.io/code-generator/cmd/client-gen")
		run(
			"client-gen",
			"--go-header-file", path.Join(cwd, "./hack/boilerplate.go.txt"),
			"--input-base", `""`,
			"--input", pkgAPIs,
			"--clientset-name", clientsetName,
			"--output-package", pkgClientClientset,
		)
	}

	if strings.Contains(gens, "lister") {
		run("go", "install", "k8s.io/code-generator/cmd/lister-gen")
		run(
			"lister-gen",
			"--go-header-file", path.Join(cwd, "./hack/boilerplate.go.txt"),
			"--input-dirs", pkgAPIs,
			"--output-package", pkgClientListers,
		)
	}

	if strings.Contains(gens, "informer") {
		run("go", "install", "k8s.io/code-generator/cmd/informer-gen")
		run(
			"informer-gen",
			"--go-header-file", path.Join(cwd, "./hack/boilerplate.go.txt"),
			"--input-dirs", pkgAPIs,
			"--versioned-clientset-package", fmt.Sprintf("%s/%s", pkgClientClientset, clientsetName),
			"--listers-package", pkgClientListers,
			"--output-package", pkgClientInformers,
		)
	}
}

func pkgPathsFix(base string, subPkgs ...string) string {
	buf := bytes.NewBuffer(nil)
	for i := range subPkgs {
		if i != 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(filepath.Join(base, subPkgs[i]))
	}
	return buf.String()
}

func run(args ...string) {
	sh := "sh"
	if runtime.GOOS == "windows" {
		sh = "bash"
	}
	execute(exec.Command(sh, "-c", strings.Join(args, " ")))
}

func execute(cmd *exec.Cmd) {
	cwd, _ := os.Getwd()

	fmt.Fprintf(os.Stdout, "%s %s\n", path.Join(cwd, cmd.Dir), strings.Join(cmd.Args, " "))

	{
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		go scanAndStdout(bufio.NewScanner(stdoutPipe))
	}
	{
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		go scanAndStderr(bufio.NewScanner(stderrPipe))
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func scanAndStdout(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Fprintln(os.Stdout, scanner.Text())
	}
}

func scanAndStderr(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Fprintln(os.Stderr, scanner.Text())
	}
}
