package main

import (
	"context"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts/candlesticks"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	candlesticksclient "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

type CandlesticksView struct {
	client           candlesticksclient.Client
	updateInProgress bool
	canvas           *charts.Canvas
	chart            *candlesticks.Chart
	windowSize       tea.WindowSizeMsg

	program *tea.Program
}

func NewCandlesticksView(program *tea.Program) *CandlesticksView {
	candlesticksClient, err := candlesticksclient.NewClient(config.LoadNATS())
	if err != nil {
		log.Fatal(err)
	}

	cv := &CandlesticksView{
		client:  candlesticksClient,
		program: program,
	}

	cv.chart = candlesticks.NewChart(&candlestick.List{}, period.H1)

	cv.canvas = charts.NewCanvas(utils.Must(time.Parse(time.RFC3339, "2022-12-01T01:00:00Z")), time.Hour)
	cv.canvas.AddChart(cv.chart)

	return cv
}

func (cv *CandlesticksView) moveCount() int {
	return cv.windowSize.Width * 2 / 3
}

func (cv *CandlesticksView) Update(message tea.Msg) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, candlestickKeyLeft):
			for i := 0; i < cv.moveCount(); i++ {
				cv.canvas.MoveLeft()
			}
		case key.Matches(msg, candlestickKeyRight):
			for i := 0; i < cv.moveCount(); i++ {
				cv.canvas.MoveRight()
			}
		}

	case tea.WindowSizeMsg:
		cv.windowSize = msg
	}

	cv.updateMissingCandlesticks()
}

func (cv *CandlesticksView) updateMissingCandlesticks() {
	first, last := cv.chart.MissingData(cv.windowSize.Width)
	if first != nil && last != nil {
		go func() {
			if cv.updateInProgress {
				return
			}
			cv.updateInProgress = true
			defer func() { cv.updateInProgress = false }()

			delta := time.Duration(cv.windowSize.Width)
			first = utils.ToReference(first.Add(-time.Hour * delta))
			last = utils.ToReference(last.Add(time.Hour * delta))

			list, err := cv.client.Read(context.TODO(), client.ReadCandlesticksPayload{
				Exchange: "binance",
				Pair:     "ETH-USDT",
				Period:   period.H1,
				Start:    first,
				End:      last,
			})
			if err != nil {
				return
			}
			_ = cv.chart.UpsertData(list)

			// Send the main program an update
			cv.program.Send(dataUpdate{})
		}()
	}
}

func (cv *CandlesticksView) View(xPad, yPad int) string {
	cv.canvas.SetHeight(cv.windowSize.Height - yPad)
	cv.canvas.SetWidth(cv.windowSize.Width - xPad)

	return cv.canvas.View()
}

var (
	candlestickKeyLeft = key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	)

	candlestickKeyRight = key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	)
)

func (cv CandlesticksView) Keys() []key.Binding {
	return []key.Binding{
		candlestickKeyLeft, candlestickKeyRight,
	}
}
