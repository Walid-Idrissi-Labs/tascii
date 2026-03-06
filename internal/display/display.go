package display

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)


var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)


func visibleLen(s string) int {
	stripped := ansiEscape.ReplaceAllString(s, "")
	return len([]rune(stripped))
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

func cellTitle(t task.Task) string {
	title := truncate(t.Title, 38)
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


type column struct {
	header string
	width  int
}


var columns = []column{
	{"ID", 4},
	{"", 2},     
	{"PRIO", 5},
	{"TITLE", 38},
	{"TAGS", 22},
	{"DUE", 15},
}


func printTableHeader() {
	fmt.Fprint(os.Stdout, "  ")
	for _, col := range columns {
		fmt.Fprint(os.Stdout, padRight(boldC.Sprint(col.header), col.width)+"  ")
	}
	fmt.Fprintln(os.Stdout)


	total := 0
	for _, col := range columns {
		total += col.width + 2
	}
	fmt.Fprintf(os.Stdout, "  %s\n", mutedC.Sprint(strings.Repeat("─", total)))
}


func PrintTable(tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Println()
		mutedC.Println("  No tasks found.")
		mutedC.Println("  Add one with: tascii add \"Your task title\"")
		fmt.Println()
		return
	}

	fmt.Println()
	printTableHeader()

	for _, t := range tasks {

		cells := []string{
			padRight(fmt.Sprintf("%d", t.ID), columns[0].width),
			padRight(cellStatus(t), columns[1].width),
			padRight(cellPriority(t), columns[2].width),
			padRight(cellTitle(t), columns[3].width),
			padRight(cellTags(t), columns[4].width),
			cellDue(t), 
		}

		fmt.Fprint(os.Stdout, "  ")
		fmt.Fprintln(os.Stdout, strings.Join(cells, "  "))
	}


	total := 0
	for _, col := range columns {
		total += col.width + 2
	}
	fmt.Fprintf(os.Stdout, "  %s\n", mutedC.Sprint(strings.Repeat("─", total)))
	fmt.Printf("  ")
	mutedC.Printf("%d task(s)\n\n", len(tasks))
}


func PrintDetail(t task.Task) {
	fmt.Println()


	fmt.Printf("  %s  %s\n", cellPriority(t), boldC.Sprint(t.Title))
	fmt.Printf("  %s\n", mutedC.Sprint(strings.Repeat("─", 50)))
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

	fmt.Println()
	fmt.Printf("   ")
	boldC.Printf("%d task(s)", len(tasks))
	fmt.Printf("  —  ")
	fmt.Printf("%s todo  ", mutedC.Sprintf("%d", todo))
	fmt.Printf("%s in-progress  ", inProgC.Sprintf("%d", inProgress))
	fmt.Printf("%s done", successC.Sprintf("%d", done))
	if overdue > 0 {
		fmt.Printf("  %s", overdueC.Sprintf("%d overdue !", overdue))
	}
	fmt.Println()
	fmt.Println()
}


func PrintSuccess(msg string) {
	successC.Printf("  ✓ %s\n", msg)
}


func PrintInfo(msg string) {
	infoC.Printf("  ℹ %s\n", msg)
}
