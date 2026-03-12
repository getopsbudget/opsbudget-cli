package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// PrintTable prints a formatted table to stdout.
func PrintTable(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	for i, h := range headers {
		if i > 0 {
			fmt.Fprint(w, "\t")
		}
		fmt.Fprint(w, h)
	}
	fmt.Fprintln(w)

	for _, row := range rows {
		for i, col := range row {
			if i > 0 {
				fmt.Fprint(w, "\t")
			}
			fmt.Fprint(w, col)
		}
		fmt.Fprintln(w)
	}

	w.Flush()
}

// ColorStatus returns a color-coded status string for terminal output.
func ColorStatus(status string) string {
	switch status {
	case "up":
		return colorGreen + "up" + colorReset
	case "down":
		return colorRed + "down" + colorReset
	case "pending":
		return colorYellow + "pending" + colorReset
	default:
		return status
	}
}

// PrintJSON prints the value as formatted JSON to stdout.
func PrintJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
