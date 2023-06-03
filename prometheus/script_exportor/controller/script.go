package controller

import (
	"fmt"
	"io"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

func NewScriptCollector(homeDir string) *ScriptCollector {
	return &ScriptCollector{
		homeDir: homeDir,
		log:     zap.L().Named("script"),
	}
}

type ScriptCollector struct {
	// 脚本存放的目录
	homeDir string

	log logger.Logger
}

func (c *ScriptCollector) Exec(module, params string, dst io.Writer) error {
	// 校验参数
	if module == "" || params == "" {
		return fmt.Errorf("module or params<target> required")
	}

	// 找到脚本
	script, err := c.find(module)
	if err != nil {
		return err
	}
	c.log.Debugf("exec script %s, args: %s", script, params)

	// 定义执行命令
	var cmd *exec.Cmd
	ext := path.Ext(script)
	switch ext {
	case ".sh":
		cmd = exec.Command("bash", script, params)
	case ".py":
		cmd = exec.Command("python", script, params)
	default:
		cmd = exec.Command(script, params)
	}

	// 定义输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("exec module %s error, %s", module, err)
	}
	defer stdout.Close()

	// 运行命令
	if err := cmd.Start(); err != nil {
		msg := strings.ReplaceAll(err.Error(), script, module)
		return fmt.Errorf("run script %s error, %s, see log for detail path", module, msg)
	}

	// 获取结果
	_, err = io.Copy(dst, stdout)
	if err != nil {
		return err
	}

	return nil
}

func (c *ScriptCollector) find(module string) (string, error) {
	absPath, err := filepath.Abs(c.homeDir)
	if err != nil {
		return "", fmt.Errorf("find module %s abs path error %s", module, err)
	}

	// 防止路径中带有 ..
	if strings.Contains(module, "..") {
		return "", fmt.Errorf("module forbiden .. in module")
	}

	return path.Join(absPath, module), nil
}
