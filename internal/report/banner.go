package report

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/layergrid/layergrid-cli/internal/version"
)

const LayergridWordmark = "    __                                    _     __\n" +
	"   / /   ____ ___  _____  _____/ ____/____(_)___/ /\n" +
	"  / /   / __ `/ / / / _ \\/ ___/ / __/ ___/ / __  / \n" +
	" / /___/ /_/ / /_/ /  __/ /  / /_/ / /  / / /_/ /  \n" +
	"/_____/\\__,_/\\__, /\\___/_/   \\____/_/  /_/\\__,_/   \n" +
	"            /____/"

type BannerOptions struct {
	NoColor   bool
	Width     int
	RuleCount int
}

func Banner(opts BannerOptions) string {
	r := humanRenderer{noColor: noColorEnabled(opts.NoColor), width: resolveWidth(opts.Width)}
	if r.width < 60 {
		return plaintextBanner(opts.RuleCount)
	}
	left := r.leftBannerPanel()
	quick := r.quickStartPanel(opts.RuleCount)
	links := r.linksPanel()
	right := lipgloss.JoinVertical(lipgloss.Left, quick, "", links)
	if r.width < designWidth {
		return lipgloss.JoinVertical(lipgloss.Left, left, "", quick, "", links) + "\n"
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right) + "\n"
}

func (r humanRenderer) leftBannerPanel() string {
	wordmark := lipgloss.NewStyle().
		Foreground(r.color(brandOrange)).
		Render(centerBlock(LayergridWordmark, 58))
	body := strings.Join([]string{
		"",
		wordmark,
		"",
		"   Secure every layer of your AI.",
		"",
		"                                Apache 2.0 · Go 1.22",
		"",
	}, "\n")
	return panel(r, "LayerGrid v"+version.Version, 60, body, true)
}

func (r humanRenderer) quickStartPanel(ruleCount int) string {
	bullet := "▸"
	if !r.noColor {
		bullet = lipgloss.NewStyle().Foreground(r.color(brandOrange)).Render(bullet)
	}
	body := fmt.Sprintf("\n  %s scan .\n    Scan the current dir\n\n  %s list-rules\n    Show all %d rules\n\n  %s explain <id>\n    Deep-dive on a rule\n\n  %s --help\n    Full reference\n", bullet, bullet, ruleCount, bullet, bullet)
	return panel(r, "Quick Start", 29, body, false)
}

func (r humanRenderer) linksPanel() string {
	body := "\n  Docs   layergrid.ai/docs\n  Score  trifecta.report\n         (coming soon)\n"
	return panel(r, "Links", 29, body, false)
}

func panel(r humanRenderer, label string, width int, body string, brandBorder bool) string {
	borderColor := r.subtleColor()
	if brandBorder {
		borderColor = r.color(brandOrange)
	}
	rendered := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Render(body)
	lines := strings.Split(rendered, "\n")
	if len(lines) == 0 {
		return rendered
	}
	top := legendTopLine(label, lipgloss.Width(lines[0]))
	if !r.noColor {
		top = lipgloss.NewStyle().Foreground(r.color(brandOrange)).Render(top)
	}
	lines[0] = top
	return strings.Join(lines, "\n")
}

func plaintextBanner(ruleCount int) string {
	return fmt.Sprintf("LayerGrid v%s\nSecure every layer of your AI.\n\nQuick Start\n  layergrid scan .\n  layergrid list-rules  (%d rules)\n  layergrid explain <id>\n", version.Version, ruleCount)
}

func legendTopLine(label string, width int) string {
	legend := "─ " + label + " "
	remaining := width - lipgloss.Width("╭╮") - lipgloss.Width(legend)
	if remaining < 0 {
		remaining = 0
	}
	return "╭" + legend + strings.Repeat("─", remaining) + "╮"
}

func centerBlock(s string, width int) string {
	lines := strings.Split(s, "\n")
	maxWidth := 0
	for _, line := range lines {
		maxWidth = maxInt(maxWidth, lipgloss.Width(line))
	}
	pad := maxInt(0, (width-maxWidth)/2)
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", pad) + line
	}
	return strings.Join(lines, "\n")
}
