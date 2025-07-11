package main

import (
	bastionModel "bastion/model"
	"bastion/repository"
	"bastion/service"
	"bastion/utils"
	"fmt"
	"github.com/muesli/termenv"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	gliderssh "github.com/gliderlabs/ssh"
)

func main() {
	repository.InitMysql()
	svc := &service.BastionService{
		BastionRepository: &repository.BastionRepository{},
	}

	server := &gliderssh.Server{
		Addr: ":42678",
		Handler: func(session gliderssh.Session) {
			pty, winCh, _ := session.Pty()
			out := termenv.NewOutput(session)
			// 设置环境变量
			_ = os.Setenv("CLICOLOR_FORCE", "1")
			//_ = os.Setenv("TERM", "xterm-256color")

			sizeMsg := tea.WindowSizeMsg{
				Width:  pty.Window.Width,
				Height: pty.Window.Height,
			}
			sizeCh := make(chan tea.Msg)
			go func() {
				for win := range winCh {
					sizeCh <- tea.WindowSizeMsg{Width: win.Width, Height: win.Height}
				}
			}()

			for {
				prog := tea.NewProgram(
					initialModel(),
					tea.WithAltScreen(),
					tea.WithInput(session),
					tea.WithOutput(out),
					tea.WithMouseCellMotion(),
				)
				// 处理窗口大小变化
				go func() {
					prog.Send(sizeMsg)
					for msg := range sizeCh {
						prog.Send(msg)
					}
				}()
				//go func() {
				//	for win := range winCh {
				//		prog.Send(tea.WindowSizeMsg{
				//			Width:  win.Width,
				//			Height: win.Height,
				//		})
				//	}
				//}()

				finalModel, err := prog.Run()
				if err != nil {
					_, _ = fmt.Fprintf(session, "Error running program: %v\n", err)
					_ = session.Exit(1)
				}
				chosen := finalModel.(model).choice
				if chosen == "" {
					_, _ = fmt.Fprintln(session, "❌ 未选择主机，退出")
					_ = session.Exit(0)
					return
				}
				parts := strings.Split(chosen, "-")
				chosen = strings.TrimSpace(parts[len(parts)-1])
				_, _ = fmt.Fprintf(session, "\n✅ 你选择了：%s\n", chosen)
				if err = utils.ConnectMachine(session, chosen); err != nil {
					_, _ = session.Write([]byte("连接失败: " + err.Error() + "\n"))
					//_ = session.Exit(1)
				}
			}
		},
		PasswordHandler: func(ctx gliderssh.Context, pass string) bool {
			return utils.ValidateUser(ctx.User(), pass, svc)
		},
	}
	log.Printf("starting gliderssh server at %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(bastionModel.Item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return bastionModel.QuitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return bastionModel.QuitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}

func initialModel() model {
	l := list.New(bastionModel.DefaultVMList, bastionModel.ItemDelegate{}, 20, 0)
	l.Title = "请选择要连接的机器："
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = bastionModel.TitleStyle
	l.Styles.PaginationStyle = bastionModel.PaginationStyle
	l.Styles.HelpStyle = bastionModel.HelpStyle

	return model{list: l}
}
