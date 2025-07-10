package utils

import (
	"errors"
	gliderssh "github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh"
)

func ConnectMachine(session gliderssh.Session, machine string) error {
	config := &ssh.ClientConfig{
		User: "root", // 目标主机的用户名
		Auth: []ssh.AuthMethod{
			ssh.Password("mojory@1q2w3e4r"), // 替换为目标主机的真实密码
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

//
//func ConnectMachine(machine string) error {
//	config := &ssh.ClientConfig{
//		User: "root",
//		Auth: []ssh.AuthMethod{
//			ssh.Password("mojory@1q2w3e4r"),
//		},
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 跳过 host key 校验
//	}
//	config.SetDefaults()
//
//	addr := net.JoinHostPort(machine, "22")
//	client, err := ssh.Dial("tcp", addr, config)
//	if err != nil {
//		return err
//	}
//	defer func(client *ssh.Client) {
//		err := client.Close()
//		if err != nil {
//			return
//		}
//	}(client)
//
//	session, err := client.NewSession()
//	if err != nil {
//		return err
//	}
//	defer func(session *ssh.Session) {
//		err := session.Close()
//		if err != nil {
//			return
//		}
//	}(session)
//
//	// 设置终端模式
//	fd := int(os.Stdin.Fd())
//	oldState, err := term.MakeRaw(fd)
//	if err != nil {
//		return err
//	}
//	defer func(fd int, oldState *term.State) {
//		err := term.Restore(fd, oldState)
//		if err != nil {
//			panic(err)
//		}
//	}(fd, oldState)
//
//	session.Stdout = os.Stdout
//	session.Stderr = os.Stderr
//	session.Stdin = os.Stdin
//
//	// 请求伪终端
//	termWidth, termHeight, _ := term.GetSize(fd)
//	if err := session.RequestPty("xterm", termHeight, termWidth, ssh.TerminalModes{}); err != nil {
//		return fmt.Errorf("请求 pty 失败: %v", err)
//	}
//
//	// 启动 shell
//	if err := session.Shell(); err != nil {
//		return fmt.Errorf("启动 shell 失败: %v", err)
//	}
//
//	// 等待 shell 退出
//	return session.Wait()
//}
