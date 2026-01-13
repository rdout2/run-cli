package domainmapping

import (
	"fmt"
	"strings"

	api_domainmapping "github.com/JulienBreux/run-cli/internal/run/api/domainmapping"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_domainmapping "github.com/JulienBreux/run-cli/internal/run/model/domainmapping"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/header"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/table"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	listHeaders = []string{
		"DOMAIN",
		"SERVICE",
		"RECORD TYPE",
		"STATE",
		"REGION",
		"CREATED",
		"UPDATED"}

	listExpansions = []int{
		2, // DOMAIN
		2, // SERVICE
		1, // RECORD TYPE
		1, // STATE
		1, // REGION
		2, // CREATED
		2, // UPDATED
	}

	listTable      *table.Table
	domainMappings []model_domainmapping.DomainMapping
)

const (
	LIST_PAGE_TITLE     = "Domain Mappings"
	LIST_PAGE_ID        = "domainmappings-list"
	LIST_PAGE_SHORTCUT  = tcell.KeyCtrlD
	MODAL_PAGE_ID       = "modal-dns-records"
)

var listDomainMappingsFunc = api_domainmapping.List

// List returns a list of domain mappings.
func List(app *tview.Application) *table.Table {
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	app.SetFocus(listTable.Table)

	return listTable
}

// Load populates the table with the provided list of domain mappings.
func Load(newDomainMappings []model_domainmapping.DomainMapping) {
	domainMappings = newDomainMappings
	render(domainMappings)
}

func ListReload(app *tview.Application, currentInfo info.Info, onResult func(error)) {
	listTable.Table.Clear()
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
	listTable.Table.SetTitle(fmt.Sprintf(" %s loading ", LIST_PAGE_TITLE))

	app.SetFocus(listTable.Table)

	go func() {
		// Fetch real data
		var err error
		domainMappings, err = listDomainMappingsFunc(currentInfo.Project, currentInfo.Region)

		app.QueueUpdateDraw(func() {
			defer func() {
				if len(domainMappings) == 0 {
					listTable.Table.Clear()
					listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
				}
				onResult(err)
			}()

			if err != nil {
				return
			}

			render(domainMappings)
		})
	}()
}

func render(dms []model_domainmapping.DomainMapping) {
	listTable.Table.Clear()
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	for i, dm := range dms {
		row := i + 1 // +1 for header row

		recordType := "n/a"
		if len(dm.Records) > 0 {
			recordType = dm.Records[0].Type
		}

		state := "Unknown"
		for _, c := range dm.Conditions {
			if c.Type == "Ready" {
				if c.State == "True" {
					state = "Ready"
				} else {
					state = c.Message
				}
				break
			}
		}

		listTable.Table.SetCell(row, 0, tview.NewTableCell(dm.Name))
		listTable.Table.SetCell(row, 1, tview.NewTableCell(dm.RouteName))
		listTable.Table.SetCell(row, 2, tview.NewTableCell(recordType))
		listTable.Table.SetCell(row, 3, tview.NewTableCell(state))
		listTable.Table.SetCell(row, 4, tview.NewTableCell(dm.Region))
		listTable.Table.SetCell(row, 5, tview.NewTableCell(humanize.Time(dm.CreateTime)))
		listTable.Table.SetCell(row, 6, tview.NewTableCell(humanize.Time(dm.UpdateTime)))
	}

	// Refresh title
	listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", LIST_PAGE_TITLE, len(dms)))
}

// GetSelectedDomainMappingFull returns the full domain mapping object for the selected row.
func GetSelectedDomainMappingFull() *model_domainmapping.DomainMapping {
	row, _ := listTable.Table.GetSelection()
	if row < 1 || len(domainMappings) == 0 {
		return nil
	}
	return &domainMappings[row-1]
}

func Shortcuts() {
	header.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<r> [white]Refresh
[dodgerblue]<enter> [white]DNS Records`
	header.ContextShortcutView.SetText(shortcuts)
}

// DNSRecordsModal creates a modal to display DNS records for the domain mapping.
func DNSRecordsModal(app *tview.Application, dm *model_domainmapping.DomainMapping, closeFunc func()) tview.Primitive {
	var sb strings.Builder
	for _, r := range dm.Records {
		fmt.Fprintf(&sb, "[yellow::b]Type:[white] %s\n", r.Type)
		fmt.Fprintf(&sb, "[yellow::b]Name:[white] %s\n", r.Name)
		fmt.Fprintf(&sb, "[yellow::b]Data:[white] %s\n", r.RRData)
		// TTL is not in the resource record struct currently, but usually present. 
		// For now we omit or hardcode if we had it.
		fmt.Fprintf(&sb, "\n")
	}

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(sb.String()).
		SetWrap(true).
		SetScrollable(true)

	textView.SetBorder(true).SetTitle(" DNS Records ")

	// Instructions
	instructions := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Press 'q' or 'esc' to close")

	// Button
	btnOk := tview.NewButton("Ok").SetSelectedFunc(closeFunc)
	
	// Layout
	content := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, false).
		AddItem(instructions, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(btnOk, 10, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), 1, 0, true)

	content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Rune() == 'q' {
			closeFunc()
			return nil
		}
		return event
	})
	
	// Create a Grid to center the modal
	grid := tview.NewGrid().
		SetColumns(0, 80, 0).
		SetRows(0, 20, 0).
		AddItem(content, 1, 1, 1, 1, 0, 0, true)

	return grid
}
