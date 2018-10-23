package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/google/subcommands"
)

// extra script to run before running entry point
// FIXME: It only supports container image with bash in it.
const script = `
mount -t proc proc /proc && \
exec %s
`

type RunCmd struct {
	root string
}

func (*RunCmd) Name() string     { return "run" }
func (*RunCmd) Synopsis() string { return "run command in container" }
func (*RunCmd) Usage() string {
	return `run <cmd>:
  Run command in container.
`
}

func (c *RunCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.root, "root", ".", "root folder")
}

func (c *RunCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := f.Args()

	if len(args) == 0 {
		return subcommands.ExitSuccess
	}

	command := strings.Join(args, " ")
	cmd := exec.Command("bash", "-c", fmt.Sprintf(script, command))
	cmd.Dir = "/"
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
		Chroot:     c.root,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
