package display

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)

// ── ANSI helpers ─────────────────────────────────────────────────────────────

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visibleLen(s string) int {
	return len([]rune(ansiEscape.ReplaceAllString(s, "")))
}

func padRight(s string, width int) string {
	vl := visibleLen(s)
	if vl >= width {
		return s
	}
	return s + strings.Repeat(" ", width-vl)
}

func truncate(s string, maxLen int) string {
	runes := []rune(ansiEscape.ReplaceAllString(s, ""))
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}

// ── Colours ───────────────────────────────────────────────────────────────────

var (
	hiPrio   = color.New(color.FgRed, color.Bold)
	medPrio  = color.New(color.FgYellow)
	loPrio   = color.New(color.Faint)
	doneText = color.New(color.Faint)
	overdueC = color.New(color.FgRed)
	todayC   = color.New(color.FgYellow)
	inProgC  = color.New(color.FgCyan)
	successC = color.New(color.FgGreen, color.Bold)
	infoC    = color.New(color.FgCyan)
	mutedC   = color.New(color.Faint)
	boldC    = color.New(color.Bold)
)

// ── Terminal width ────────────────────────────────────────────────────────────

// termWidth returns the current terminal column count.
// Priority: ioctl syscall (term_unix.go) → $COLUMNS env var → 80 default.
func termWidth() int {
	if w := getTermWidth(); w > 0 {
		return w
	}
	if col := os.Getenv("COLUMNS"); col != "" {
		if n, err := strconv.Atoi(col); err == nil && n > 20 {
			return n
		}
	}
	return 80
}

// ── Responsive column layout ──────────────────────────────────────────────────

// column describes one rendered table column.
type column struct {
	header string
	width  int
}

// Layout constants (in visible characters).
const (
	colIndent = 2  // leading "  " before each row
	colSep    = 2  // "  " gap between columns
	colID     = 4
	colStatus = 2
	colPrio   = 5
	colDue    = 15

	minTitle = 12
	maxTitle = 52
	minTags  = 10
	maxTags  = 28
)

// fixedBase is the visual width consumed by the non-flexible columns and their
// separators: indent + (ID+sep) + (STATUS+sep) + (PRIO+sep) + (TITLE sep) + DUE
// = 2 + 6 + 4 + 7 + 2 + 15 = 36
const fixedBase = colIndent + (colID + colSep) + (colStatus + colSep) + (colPrio + colSep) + colSep + colDue

// computeColumns builds the column slice for the given terminal width,
// dynamically sizing TITLE and (optionally) TAGS.
func computeColumns(width int) []column {
	// Minimum width needed to show TAGS at all.
	withTagsMin := fixedBase + colSep + minTitle + minTags // +colSep = TAGS trailing sep

	if width < withTagsMin {
		// ── Compact mode: no TAGS column ──
		titleW := width - fixedBase
		if titleW < minTitle {
			titleW = minTitle
		}
		return []column{
			{"ID", colID},
			{"", colStatus},
			{"PRIO", colPrio},
			{"TITLE", titleW},
			{"DUE", colDue},
		}
	}

	// ── Normal / wide mode: TITLE + TAGS ──
	// Available space = total - fixedBase - TAGS's own colSep
	available := width - fixedBase - colSep
	// Give TITLE 65 %, TAGS 35 %, then clamp both.
	titleW := available * 65 / 100
	tagsW := available - titleW

	if titleW > maxTitle {
		tagsW += titleW - maxTitle
		titleW = maxTitle
	}
	if tagsW > maxTags {
		titleW += tagsW - maxTags
		tagsW = maxTags
	}
	if titleW > maxTitle {
		titleW = maxTitle
	}
	if titleW < minTitle {
		titleW = minTitle
	}
	if tagsW < minTags {
		tagsW = minTags
	}

	return []column{
		{"ID", colID},
		{"", colStatus},
		{"PRIO", colPrio},
		{"TITLE", titleW},
		{"TAGS", tagsW},
		{"DUE", colDue},
	}
}

// tableWidth returns the total visual width of a row for the given columns.
func tableWidth(cols []column) int {
	w := colIndent
	for i, c := range cols {
		w += c.width
		if i < len(cols)-1 {
			w += colSep
		}
	}
	return w
}

// hasTags reports whether the column set contains a TAGS column.
func hasTags(cols []column) bool {
	for _, c := range cols {
		if c.header == "TAGS" {
			return true
		}
	}
	return false
}

// ── Cell renderers ────────────────────────────────────────────────────────────

func cellStatus(t task.Task) string {
	switch t.Status {
	case task.StatusInProgress:
		return inProgC.Sprint("◐")
	case task.StatusDone:
		return doneText.Sprint("✓")
	default:
		return "○"
	}
}

func cellPriority(t task.Task) string {
	sym := t.PrioritySymbol()
	switch t.Priority {
	case 3:
		return hiPrio.Sprint(sym)
	case 2:
		return medPrio.Sprint(sym)
	default:
		return loPrio.Sprint(sym)
	}
}

func cellTitle(t task.Task, maxLen int) string {
	title := truncate(t.Title, maxLen)
	switch {
	case t.Status == task.StatusDone:
		return doneText.Sprint(title)
	case t.IsOverdue():
		return overdueC.Sprint(title)
	default:
		return title
	}
}

func cellDue(t task.Task) string {
	if t.Due == "" {
		return mutedC.Sprint("—")
	}
	due, err := time.Parse("2006-01-02", t.Due)
	if err != nil {
		return t.Due
	}
	formatted := due.Format("Jan 02")
	switch {
	case t.IsOverdue():
		return overdueC.Sprint(formatted + " OVERDUE")
	case t.IsDueToday():
		return todayC.Sprint("Today")
	default:
		return formatted
	}
}

func cellTags(t task.Task) string {
	if len(t.Tags) == 0 {
		return mutedC.Sprint("—")
	}
	parts := make([]string, len(t.Tags))
	for i, tg := range t.Tags {
		parts[i] = infoC.Sprint("#" + tg)
	}
	return strings.Join(parts, " ")
}

// ── Table renderer ────────────────────────────────────────────────────────────

func printTableHeader(cols []column) {
	fmt.Fprint(os.Stdout, "  ")
	for _, col := range cols {
		fmt.Fprint(os.Stdout, padRight(boldC.Sprint(col.header), col.width)+"  ")
	}
	fmt.Fprintln(os.Stdout)

	tw := tableWidth(cols)
	fmt.Fprintf(os.Stdout, "  %s\n", mutedC.Sprint(strings.Repeat("─", tw-colIndent)))
}

// PrintTable renders all tasks in an adaptive table.
func PrintTable(tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Println()
		mutedC.Println("  No tasks found.")
		mutedC.Println("  Add one with: tascii add \"Your task title\"")
		fmt.Println()
		return
	}

	cols := computeColumns(termWidth())
	showTags := hasTags(cols)

	// Column indices are always: 0=ID, 1=STATUS, 2=PRIO, 3=TITLE, then
	// 4=TAGS+5=DUE (wide) or 4=DUE (compact).
	const (
		idxID     = 0
		idxStatus = 1
		idxPrio   = 2
		idxTitle  = 3
	)

	fmt.Println()
	printTableHeader(cols)

	for _, t := range tasks {
		var cells []string
		if showTags {
			cells = []string{
				padRight(fmt.Sprintf("%d", t.ID), cols[idxID].width),
				padRight(cellStatus(t), cols[idxStatus].width),
				padRight(cellPriority(t), cols[idxPrio].width),
				padRight(cellTitle(t, cols[idxTitle].width), cols[idxTitle].width),
				padRight(cellTags(t), cols[4].width),
				cellDue(t),
			}
		} else {
			cells = []string{
				padRight(fmt.Sprintf("%d", t.ID), cols[idxID].width),
				padRight(cellStatus(t), cols[idxStatus].width),
				padRight(cellPriority(t), cols[idxPrio].width),
				padRight(cellTitle(t, cols[idxTitle].width), cols[idxTitle].width),
				cellDue(t),
			}
		}

		fmt.Fprint(os.Stdout, "  ")
		fmt.Fprintln(os.Stdout, strings.Join(cells, "  "))
	}

	tw := tableWidth(cols)
	fmt.Fprintf(os.Stdout, "  %s\n", mutedC.Sprint(strings.Repeat("─", tw-colIndent)))
	fmt.Printf("  ")
	mutedC.Printf("%d task(s)\n\n", len(tasks))
}

// ── Detail view ───────────────────────────────────────────────────────────────

// PrintDetail renders one task's full detail, with a separator that matches
// the terminal width (capped for readability).
func PrintDetail(t task.Task) {
	tw := termWidth()
	sepLen := tw - 4 // 2 indent + 2 extra padding
	if sepLen > 70 {
		sepLen = 70
	}
	if sepLen < 20 {
		sepLen = 20
	}

	fmt.Println()
	fmt.Printf("  %s  %s\n", cellPriority(t), boldC.Sprint(t.Title))
	fmt.Printf("  %s\n", mutedC.Sprint(strings.Repeat("─", sepLen)))
	fmt.Println()

	labelRow := func(label, value string) {
		mutedC.Printf("  %-14s", label)
		fmt.Println(value)
	}

	labelRow("ID:", fmt.Sprintf("%d", t.ID))
	labelRow("Status:", fmt.Sprintf("%s  %s", cellStatus(t), t.Status.String()))
	labelRow("Priority:", cellPriority(t))
	labelRow("Due:", cellDue(t))

	if len(t.Tags) > 0 {
		labelRow("Tags:", cellTags(t))
	}
	labelRow("Created:", t.CreatedAt.Format("Jan 02, 2006 15:04"))

	if t.Note != "" {
		fmt.Println()
		mutedC.Println("  Note:")
		for _, line := range strings.Split(t.Note, "\n") {
			fmt.Printf("    %s\n", line)
		}
	}

	fmt.Println()
}

// ── Today view ────────────────────────────────────────────────────────────────

func PrintToday(tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Println()
		successC.Println("  ✓ Nothing due today.")
		fmt.Println()
		return
	}

	fmt.Println()
	boldC.Println("  Due Today & Overdue")
	PrintTable(tasks)
}

// ── Summary view ──────────────────────────────────────────────────────────────

// PrintSummary shows task counts. On narrow terminals (< 72 cols) the stats
// wrap to a second line to avoid horizontal overflow.
func PrintSummary(tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Println()
		mutedC.Println("  No tasks added. Start with: tascii add \"Your task\"")
		fmt.Println()
		return
	}

	var todo, inProgress, done, overdue int
	for _, t := range tasks {
		switch t.Status {
		case task.StatusTodo:
			todo++
		case task.StatusInProgress:
			inProgress++
		case task.StatusDone:
			done++
		}
		if t.IsOverdue() {
			overdue++
		}
	}

	tw := termWidth()
	fmt.Println()
	fmt.Printf("   ")
	boldC.Printf("%d task(s)", len(tasks))

	if tw < 72 {
		// Narrow: wrap stats onto next line.
		fmt.Println()
		fmt.Printf("   ")
	} else {
		fmt.Printf("  —  ")
	}

	fmt.Printf("%s todo  ", mutedC.Sprintf("%d", todo))
	fmt.Printf("%s in-progress  ", inProgC.Sprintf("%d", inProgress))
	fmt.Printf("%s done", successC.Sprintf("%d", done))
	if overdue > 0 {
		if tw < 72 {
			fmt.Printf("\n   %s", overdueC.Sprintf("%d overdue !", overdue))
		} else {
			fmt.Printf("  %s", overdueC.Sprintf("%d overdue !", overdue))
		}
	}
	fmt.Println()
	fmt.Println()
}

// ── Misc helpers ──────────────────────────────────────────────────────────────

func PrintSuccess(msg string) {
	successC.Printf("  ✓ %s\n", msg)
}

func PrintInfo(msg string) {
	infoC.Printf("  ℹ %s\n", msg)
}

func PrintMuted(msg string) {
	mutedC.Println(msg)
}