package jenkins

import (
	"fmt"
	"regexp"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/wtf"
)

const HelpText = `
 Keyboard commands for Jenkins:

   /: Show/hide this help window
   j: Select the next job in the list
   k: Select the previous job in the list
   r: Refresh the data

   arrow down: Select the next job in the list
   arrow up:   Select the previous job in the list

   return: Open the selected job in a browser
`

type Widget struct {
	wtf.HelpfulWidget
	wtf.KeyboardWidget
	wtf.ScrollableWidget

	settings *Settings
	view     *View
}

func NewWidget(app *tview.Application, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		HelpfulWidget:    wtf.NewHelpfulWidget(app, pages, HelpText),
		KeyboardWidget:   wtf.NewKeyboardWidget(),
		ScrollableWidget: wtf.NewScrollableWidget(app, settings.common, true),

		settings: settings,
	}

	widget.SetRenderFunction(widget.Render)
	widget.initializeKeyboardControls()
	widget.View.SetInputCapture(widget.InputCapture)

	widget.HelpfulWidget.SetView(widget.View)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	view, err := widget.Create(
		widget.settings.url,
		widget.settings.user,
		widget.settings.apiKey,
	)
	widget.view = view

	if err != nil {
		widget.Redraw(widget.CommonSettings.Title, err.Error(), true)
		return
	}

	widget.SetItemCount(len(widget.view.Jobs))

	widget.Render()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) Render() {
	if widget.view == nil {
		return
	}

	title := fmt.Sprintf("%s: [red]%s", widget.CommonSettings.Title, widget.view.Name)

	widget.Redraw(title, widget.contentFrom(widget.view), false)
}

func (widget *Widget) contentFrom(view *View) string {
	var str string
	for idx, job := range view.Jobs {
		var validID = regexp.MustCompile(widget.settings.jobNameRegex)

		if validID.MatchString(job.Name) {
			str = str + fmt.Sprintf(
				`["%d"][%s] [%s]%-6s[white][""]`,
				idx,
				widget.RowColor(idx),
				widget.jobColor(&job),
				job.Name,
			)

			str = str + "\n"
		}
	}

	return str
}

func (widget *Widget) jobColor(job *Job) string {
	switch job.Color {
	case "blue":
		// Override color if successBallColor boolean param provided in config
		return widget.settings.successBallColor
	case "red":
		return "red"
	default:
		return "white"
	}
}

func (widget *Widget) openJob() {
	sel := widget.GetSelected()
	if sel >= 0 && widget.view != nil && sel < len(widget.view.Jobs) {
		job := &widget.view.Jobs[sel]
		wtf.OpenFile(job.Url)
	}
}
