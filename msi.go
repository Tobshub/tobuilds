package tobuilds

import (
	"fmt"
	"os/exec"
)

type MSI struct {
	ctx      *Ctx
	platform Platform
	name     string
}

func newMSI(ctx *Ctx, name string) *MSI {
	fmt.Println("INFO: registered new MSI", name)
	return &MSI{ctx, PlatformWindows, name}
}

type MsiExecOptions struct {
	I, // install
	A, // administrative install
	X, // uninstall
	Quiet, // quiet mode
	Passive, // unattended mode
	Qn, // no ui
	Qb, // basic ui
	Qr, // reduced ui
	Qf, // full ui
	H, // help
    NoRestart, PromptRestart, ForceRestart bool
}

func (o MsiExecOptions) parse() []string {
    var opts []string
    if o.Quiet { opts = append(opts, "/quiet") }
    if o.Passive { opts = append(opts, "/passive") }
    if o.Qn { opts = append(opts, "/qn") }
    if o.Qb { opts = append(opts, "/qb") }
    if o.Qr { opts = append(opts, "/qr") }
    if o.Qf { opts = append(opts, "/qf") }
    if o.H { opts = append(opts, "/h") }
    if o.NoRestart { opts = append(opts, "/norestart") }
    if o.PromptRestart { opts = append(opts, "/promptrestart") }
    if o.ForceRestart { opts = append(opts, "/forcerestart") }
    if o.I { opts = append(opts, "/i") }
    if o.A { opts = append(opts, "/a") }
    if o.X { opts = append(opts, "/x") }
    return opts
}

func (m *MSI) Exec(opts MsiExecOptions) error {
	if !m.platform.isCurrent() {
		fmt.Printf("INFO: skipped run (%s) for different platform\n", m.name)
		return nil
	}
	f, err := m.ctx.GetFile(m.name)
	if err != nil {
		return err
	}

    args := opts.parse()
    args = append(args, f.Name())
    res, err := exec.Command("msiexec", args...).CombinedOutput()
    fmt.Print("OUT: ", string(res))
	return err
}
