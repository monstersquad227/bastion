package utils

import (
	"bastion/repository"
	"bastion/service"
	"errors"
	gliderssh "github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh"
)

func ConnectMachine(session gliderssh.Session, machine string) error {
	bastionSvc := &service.BastionService{
		BastionRepository: &repository.BastionRepository{},
	}
	encryptText, err := bastionSvc.GetPasswordByIp(machine)
	if err != nil {
		return err
	}
	plainPassword, err := DecryptPassword(encryptText)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: "root", // 目标主机的用户名
		Auth: []ssh.AuthMethod{
			ssh.Password(plainPassword), // 替换为目标主机的真实密码
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产环境应验证主机密钥
	}

	// 建立到目标主机的连接
	targetAddr := machine + ":22"
	client, err := ssh.Dial("tcp", targetAddr, config)
	if err != nil {
		return err
	}
	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil {
			return
		}
	}(client)

	// 在目标主机上创建会话
	targetSession, err := client.NewSession()
	if err != nil {
		return err
	}
	defer func(targetSession *ssh.Session) {
		err := targetSession.Close()
		if err != nil {
			return
		}
	}(targetSession)

	// 将当前会话的IO重定向到目标会话
	targetSession.Stdin = session
	targetSession.Stdout = session
	targetSession.Stderr = session

	// 设置终端模式
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // 启用回显
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	}

	// 获取当前PTY
	if pty, winCh, ok := session.Pty(); ok {
		// 请求伪终端
		if err := targetSession.RequestPty(pty.Term, pty.Window.Height, pty.Window.Width, modes); err != nil {
			return err
		}

		// 监听窗口大小变化
		go func() {
			for win := range winCh {
				_ = targetSession.WindowChange(win.Height, win.Width)
			}
		}()
	}

	// 启动远程shell
	if err := targetSession.Shell(); err != nil {
		return err
	}

	// 等待会话结束
	if err := targetSession.Wait(); err != nil {
		var exitErr *ssh.ExitError
		if errors.As(err, &exitErr) {
			_ = session.Exit(exitErr.ExitStatus())
		}
		return err
	}
	return nil
}
