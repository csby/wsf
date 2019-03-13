package types

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ArgsParser interface {
	Parse(key, value string)
}

type Args struct {
	moduleVersion Version
	moduleType    string
	moduleName    string
	moduleRemark  string
	modulePath    string
}

func (s *Args) Parse(args []string, moduleType, moduleName, moduleVersion, moduleRemark string, parser ArgsParser) {
	s.moduleType = moduleType
	s.moduleName = moduleName
	s.moduleRemark = moduleRemark
	s.moduleVersion.Parse(moduleVersion)

	if nil == args {
		return
	}
	argsLength := len(args)
	if argsLength < 1 {
		return
	}

	s.modulePath = args[0]
	absModulePath, err := filepath.Abs(args[0])
	if err == nil {
		s.modulePath = absModulePath
	}

	for argsIndex := 1; argsIndex < argsLength; argsIndex++ {
		arg := args[argsIndex]
		splitIndex := strings.Index(arg, "=")
		if splitIndex < 1 {
			if strings.ToLower("--type") == strings.ToLower(arg) { // 获取模块类型
				fmt.Fprint(os.Stdout, moduleType)
				os.Exit(0)
			} else if strings.ToLower("--module") == strings.ToLower(arg) { // 获取模块名称
				fmt.Fprint(os.Stdout, moduleName)
				os.Exit(0)
			} else if strings.ToLower("--version") == strings.ToLower(arg) { // 获取模块版本号
				fmt.Fprint(os.Stdout, moduleVersion)
				os.Exit(0)
			} else if strings.ToLower("--remark") == strings.ToLower(arg) { // 获取模块备注说明
				fmt.Fprint(os.Stdout, moduleRemark)
				os.Exit(0)
			} else {
				if nil != parser {
					parser.Parse(strings.ToLower(arg), "")
				}
			}

			continue
		}

		if nil == parser {
			continue
		}
		if splitIndex >= len(arg)-1 {
			continue
		}
		key := strings.ToLower(arg[0:splitIndex])
		val := strings.Trim(arg[splitIndex+1:], "\"")
		parser.Parse(strings.ToLower(key), val)
	}
}

func (s *Args) ModulePath() string {
	return s.modulePath
}

func (s *Args) ModuleFolder() string {
	return filepath.Dir(s.modulePath)
}

func (s *Args) ModuleType() string {
	return s.moduleType
}

func (s *Args) ModuleName() string {
	return s.moduleName
}

func (s *Args) ModuleRemark() string {
	return s.moduleRemark
}

func (s *Args) ModuleVersion() *Version {
	return &s.moduleVersion
}

func (s *Args) ParseNew(path string) (*Args, error) {
	out, err := exec.Command(path, "--version").Output()
	if err != nil {
		return nil, err
	}
	args := &Args{modulePath: path}
	if !args.moduleVersion.Parse(string(out[:])) {
		return nil, fmt.Errorf("invalid version: %s", string(out[:]))
	}

	out, err = exec.Command(path, "--type").Output()
	if err != nil {
		return nil, err
	}
	args.moduleType = string(out[:])

	out, err = exec.Command(path, "--module").Output()
	if err != nil {
		return nil, err
	}
	args.moduleName = string(out[:])

	out, err = exec.Command(path, "--remark").Output()
	if err != nil {
		return nil, err
	}
	args.moduleRemark = string(out[:])

	return args, nil
}
