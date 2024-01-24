package processutils

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Process interface {
	Start() error
	Kill() (string, error)
	CheckRunning() (string, error)
	GetProcessInfo() ([]*ProcessInfo, error)
}

type process struct {
	cmdName   string
	args      []string
	check     string
	checkArgs []string
}

var ErrNotFoundProcess = errors.New("not found process")

func NewProcess(cmdName string, args []string, checkArgs ...string) Process {
	r := &process{cmdName: cmdName, args: args, checkArgs: checkArgs}
	r.init()
	return r
}

func (p *process) Start() error {
	cmd := exec.Command(p.cmdName, p.args...)
	//cmd.Stdin = os.Stdin
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	cmdLine := p.cmdName
	count := len(p.args) / 2
	for i := 0; i < count; i++ {
		cmdLine = fmt.Sprintf("%s %s=%s", cmdLine, p.args[i*2], p.args[i*2+1])
	}

	logs.Infof(context.Background(), "", nil, cmdLine)
	_, err := p.CheckRunning()
	if err != nil && !errors.Is(err, ErrNotFoundProcess) {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return errors.New(fmt.Sprintf("start %s process fail, error %s", p.cmdName, err.Error()))
	}
	logs.Debugfmt(context.Background(), "", "dapr process id = %v ", cmd.Process.Pid)

	return err
}

func (p *process) init() {
	sb := strings.Builder{}
	sb.WriteString(`ps -ef|grep ` + p.cmdName)
	for _, s := range p.checkArgs {
		sb.WriteString("|grep ")
		sb.WriteString(s)
	}
	p.check = sb.String()
}

func (p *process) Kill() (pid string, err error) {
	defer func() {
		if err != nil {
			fmt.Println(fmt.Sprintf("Stop %s %s", p.cmdName, err.Error()))
		} else if pid == "" {
			fmt.Println(fmt.Sprintf("Stop %s not found", p.cmdName))
		} else {
			fmt.Println(fmt.Sprintf("Stop %s OK PID=%s", p.cmdName, pid))
		}
	}()

	pid, err = p.CheckRunning()
	if err != nil {
		return pid, err
	}
	err = p.kill(pid)
	if err != nil {
		return pid, err
	}
	return pid, err
}

func (p *process) kill(pid string) error {
	id, err := strconv.Atoi(pid)
	if err != nil {
		return err
	}
	proc, err := os.FindProcess(id)
	if err != nil {
		return err
	}
	return proc.Kill()
}

// CheckRunning 根据进程名判断进程是否运行
func (p *process) CheckRunning() (string, error) {
	pid, err := runCommand(p.cmdName, p.check)
	if err != nil {
		return "", err
	}
	return pid, nil
}

// GetPid 根据进程名称获取进程ID
func (p *process) GetPid() (pid int, err error) {
	var pidStr string
	if pidStr, err = runCommand(p.cmdName, p.check); err != nil {
		return
	}
	pid, err = strconv.Atoi(pidStr)
	return
}

// GetProcessInfo 根据进程名称获取进程ID
func (p *process) GetProcessInfo() ([]*ProcessInfo, error) {
	var result []byte
	var err error
	fmt.Println("  " + p.check)
	if runtime.GOOS == "windows" {
		result, err = exec.Command("cmd", "/c", p.check).Output()
	} else {
		result, err = exec.Command("/bin/sh", "-c", p.check).Output()
	}
	if err != nil {
		return nil, err
	}
	res := getPids(string(result), p.cmdName)
	return res, nil
}

func runCommand(cmdName string, cmdLine string) (string, error) {
	if runtime.GOOS == "windows" {
		return runInWindows(cmdName, cmdLine)
	}
	return runInLinux(cmdName, cmdLine)
}

func runInWindows(cmdName string, cmdLine string) (string, error) {
	result, err := exec.Command("cmd", "/c", cmdLine).Output()
	if err != nil {
		return "", err
	}
	pinfos := getPids(string(result), cmdName)
	if len(pinfos) > 0 {
		return pinfos[0].PID, nil
	}
	return "", ErrNotFoundProcess
}

func runInLinux(cmdName string, cmdLine string) (string, error) {
	result, err := exec.Command("/bin/sh", "-c", cmdLine).Output()
	if err != nil {
		return "", err
	}
	pinfos := getPids(string(result), cmdName)
	if len(pinfos) > 0 {
		return pinfos[0].PID, nil
	}
	return "", ErrNotFoundProcess
}

func getPids(outText string, cmdName string) []*ProcessInfo {
	outLines := strings.Split(outText, "\n")
	count := len(outLines) - 1
	var res []*ProcessInfo
	for i := 0; i < count; i++ {
		text := outLines[i]
		if len(text) != 0 {
			p := newProcessInfo(text)
			if p.CmdName == cmdName {
				fmt.Println(text)
				res = append(res, p)
			}
		}
	}
	return res
}

type ProcessInfo struct {
	UID       string
	PID       string
	PPID      string
	C         string
	StartTime string
	TTY       string
	Time      string
	CmdName   string
	CmdPath   string
}

func newProcessInfo(text string) *ProcessInfo {
	list := strings.Split(text, " ")
	index := 0
	p := &ProcessInfo{}
	for _, str := range list {
		v := strings.Trim(str, " ")
		if len(v) > 0 {
			index++
		} else {
			continue
		}
		switch index {
		case 1:
			p.UID = str
		case 2:
			p.PID = str
		case 3:
			p.PPID = str
		case 4:
			p.C = str
		case 5:
			p.StartTime = str
		case 6:
			p.TTY = str
		case 7:
			p.Time = str
		case 8:
			dir, file := filepath.Split(str)
			p.CmdPath = dir
			p.CmdName = file
		}
	}
	return p
}
