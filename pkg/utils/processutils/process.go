package processutils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

type Process interface {
	Start() error
	Kill() (string, error)
	CheckRunning() (pid string, err error, isFound bool)
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

func logInfo(format string, args ...interface{}) {
	logs.Infof(nil, nil, format, args)
}

func logError(format string, args ...interface{}) {
	logs.Errorf(nil, nil, format, args)
}

func (p *process) Start() error {
	var cmd *exec.Cmd
	line := p.cmdName + " " + strings.Join(p.args, " ")

	cmdLine := p.cmdName
	count := len(p.args)
	for i := 0; i < count; i++ {
		cmdLine = fmt.Sprintf("%s %s ", cmdLine, p.args[i])
	}

	logInfo(cmdLine)
	_, err, _ := p.CheckRunning()
	if err != nil && !errors.Is(err, ErrNotFoundProcess) {
		fmt.Printf("dapr start() error = %s \n", err.Error())
		return err
	}

	cmd, err = runCmd(p.cmdName, line)
	_, stderr := getCmdOuts(cmd)
	if err != nil {
		if stderr.Len() > 0 {
			logError("dapr start() error = %s \n", stderr.String())
		}
		return errors.New(fmt.Sprintf("start %s process fail, error %s", p.cmdName, err.Error()))
	}
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

	pid, err, _ = p.CheckRunning()
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
func (p *process) CheckRunning() (pid string, err error, isFound bool) {
	pid, err = runCmdPid(p.cmdName, p.check)
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return "", errors.New(string(exit.Stderr)), false
		}
		return "", err, false
	}
	return pid, nil, true
}

// GetPid 根据进程名称获取进程ID
func (p *process) GetPid() (pid int, err error) {
	var pidStr string
	if pidStr, err = runCmdPid(p.cmdName, p.check); err != nil {
		return
	}
	pid, err = strconv.Atoi(pidStr)
	return
}

// GetProcessInfo 根据进程名称获取进程ID
func (p *process) GetProcessInfo() ([]*ProcessInfo, error) {
	var err error
	logInfo(p.check)
	cmd, err := runCmd(p.cmdName, p.check)
	if err != nil {

		return nil, err
	}
	stderr, _ := getCmdOuts(cmd)
	res := getPids(stderr.String(), p.cmdName)
	return res, nil
}

func getCmdOuts(cmd *exec.Cmd) (stdout *bytes.Buffer, stderr *bytes.Buffer) {
	stdout = cmd.Stdout.(*bytes.Buffer)
	stderr = cmd.Stderr.(*bytes.Buffer)
	return stdout, stderr
}

func runCmdPid(cmdName string, cmdLine string) (string, error) {
	var stdout *bytes.Buffer
	cmd, err := runCmd(cmdName, cmdLine)
	if err != nil {
		return "", err
	}
	stdout, _ = getCmdOuts(cmd)

	pInfos := getPids(stdout.String(), cmdName)
	if len(pInfos) > 0 {
		return pInfos[0].PID, nil
	}
	return "", ErrNotFoundProcess
}

func runCmd(cmdName string, cmdLine string) (*exec.Cmd, error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd := newCmd(cmdLine)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		logError("runCmd() error = %s \n", err.Error())
		return nil, err
	}
	return cmd, nil
}

func newCmd(cmdLine string) *exec.Cmd {
	var cmd *exec.Cmd
	if isWindowsOS() {
		cmd = exec.Command("powershell", "-NoProfile", "-Command", cmdLine)
	} else {
		cmd = exec.Command("/bin/sh", "-c", cmdLine)
	}

	return cmd
}

func isWindowsOS() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
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
