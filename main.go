package main

import (
	"bastion/utils"
	"fmt"
	"github.com/muesli/termenv"
	"io"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	gliderssh "github.com/gliderlabs/ssh"
)

func main() {
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

				chosen = strings.TrimSpace(strings.Split(chosen, "-")[1])
				_, _ = fmt.Fprintf(session, "\n✅ 你选择了：%s\n", chosen)
				if err = utils.ConnectMachine(session, chosen); err != nil {
					_, _ = session.Write([]byte("连接失败: " + err.Error() + "\n"))
					//_ = session.Exit(1)
				}
			}
		},
		PasswordHandler: func(ctx gliderssh.Context, pass string) bool {
			return ctx.User() == "admin" && pass == "123456"
		},
	}
	log.Printf("starting gliderssh server at %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(0).Foreground(lipgloss.Color("#1890ff")).Bold(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("#FF00FF"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("👉 " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
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
			i, ok := m.list.SelectedItem().(item)
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
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}

var items = []list.Item{
	item("devflow_web - 192.168.1.113"),
	item("devflow - 192.168.1.112"),
	item("cddjob - 192.168.1.123"),
	item("tmp - 192.168.1.4"),
	item("Currywurst - 192.168.1.5"),
	item("Okonomiyaki - 192.168.1.6"),
	item("Pasta - 192.168.1.7"),
	item("Fillet Mignon - 192.168.1.8"),
	item("Caviar - 192.168.1.9)"),
	item("Text - 192.168.1.10"),
	item("Text - 192.168.1.11"),
	item("Text - 192.168.1.12"),
	item("Text - 192.168.1.13"),
	item("Text - 192.168.1.14"),
	item("Text - 192.168.1.15"),
	item("Text - 192.168.1.16"),
	item("Text - 192.168.1.17"),
	item("Text - 192.168.1.18"),
	item("Text - 192.168.1.19"),
	item("Text - 192.168.1.20"),
	item("Text - 192.168.1.21"),
	item("Text - 192.168.1.22"),
	item("Text - 192.168.1.23"),
	item("Text - 192.168.1.24"),
	item("Text - 192.168.1.25"),
	item("Text - 192.168.1.26"),
	item("Text - 192.168.1.27"),
	item("Text - 192.168.1.28"),
	item("Text - 192.168.1.29"),
	item("Text - 192.168.1.30"),
	item("Text - 192.168.1.31"),
	item("Text - 192.168.1.32"),
	item("Text - 192.168.1.33"),
	item("Text - 192.168.1.34"),
	item("Text - 192.168.1.35"),
	item("Text - 192.168.1.36"),
	item("Text - 192.168.1.37"),
	item("Text - 192.168.1.38"),
	item("Text - 192.168.1.39"),
	item("Text - 192.168.1.40"),
	item("Text - 192.168.1.41"),
	item("Text - 192.168.1.42"),
	item("Text - 192.168.1.43"),
	item("Text - 192.168.1.44"),
	item("Text - 192.168.1.45"),
	item("Text - 192.168.1.46"),
	item("Text - 192.168.1.47"),
	item("Text - 192.168.1.48"),
	item("Text - 192.168.1.49"),
	item("Text - 192.168.1.50"),
	item("Text - 192.168.1.51"),
	item("Text - 192.168.1.52"),
	item("Text - 192.168.1.53"),
	item("Text - 192.168.1.54"),
	item("Text - 192.168.1.55"),
	item("Text - 192.168.1.56"),
	item("Text - 192.168.1.57"),
	item("Text - 192.168.1.58"),
	item("Text - 192.168.1.59"),
	item("Text - 192.168.1.60"),
	item("Text - 192.168.1.61"),
	item("Text - 192.168.1.62"),
	item("Text - 192.168.1.63"),
	item("Text - 192.168.1.64"),
	item("Text - 192.168.1.65"),
	item("Text - 192.168.1.66"),
	item("Text - 192.168.1.67"),
	item("Text - 192.168.1.68"),
	item("Text - 192.168.1.69"),
	item("Text - 192.168.1.70"),
}

func initialModel() model {
	l := list.New(items, itemDelegate{}, 20, 0)
	l.Title = "请选择要连接的机器："
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return model{list: l}
}
