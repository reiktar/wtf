package gcal

import (
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

type Widget struct {
	view.ScrollableWidget

	calEvents []*CalEvent
	err       error
	settings  *Settings
	tviewApp  *tview.Application
}

func NewWidget(tviewApp *tview.Application, redrawChan chan bool,pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.Common),

		tviewApp: tviewApp,
		settings: settings,
	}

	widget.SetRenderFunction(widget.display)
	widget.initializeKeyboardControls()

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Disable() {
	widget.TextWidget.Disable()
}

func (widget *Widget) Refresh() {
	if isAuthenticated(widget.settings.email) {
		widget.fetchAndDisplayEvents()
		return
	}

	widget.tviewApp.Suspend(widget.authenticate)
	widget.Refresh()
}

func (widget *Widget) Open() {
	widget.GetSelected()
	calEvent := widget.calEvents[widget.GetSelected()]

	link := calEvent.MeetingLink()
	if link != "" {
		utils.OpenFile(link)
	}
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) fetchAndDisplayEvents() {
	calEvents, err := widget.Fetch()
	if err != nil {
		widget.err = err
		widget.calEvents = []*CalEvent{}
	} else {
		widget.err = nil
		widget.calEvents = calEvents
	}

	widget.display()
}
