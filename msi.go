package tobuilds

import "fmt"

type MSI struct {
	ctx      *Ctx
	platform Platform
	name     string
}

func newMSI(ctx *Ctx, name string) *MSI {
	fmt.Println("INFO: registered new MSI", name)
	return &MSI{ctx, PlatformWindows, name}
}

type MsiExecOptions struct{}

func (m *MSI) Exec(opts MsiExecOptions) error {
	if !m.platform.isCurrent() {
		fmt.Printf("INFO: skipped run (%s) for different platform\n", m.name)
		return nil
	}
	f, err := m.ctx.GetFile(m.name)
	if err != nil {
		return err
	}
	fmt.Println("msi exec", f.Name())
	return nil
}
