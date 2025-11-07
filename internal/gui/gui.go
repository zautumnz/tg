package gui

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"time"

	"github.com/zautumnz/tg/internal/font"
	"github.com/zautumnz/tg/internal/gui/popup"
	"github.com/zautumnz/tg/internal/termutil"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

type GUI struct {
	mouseStateLeft   MouseState
	mouseStateRight  MouseState
	mouseStateMiddle MouseState
	mouseDrag        bool
	size             image.Point // pixels
	terminal         *termutil.Terminal
	updateChan       chan struct{}
	lastClick        time.Time
	clickCount       int
	fontManager      *font.Manager
	mousePos         termutil.Position
	popupMessages    []popup.Message
	startupFuncs     []func(g *GUI)
	keyState         *keyState
	cursorImage      *ebiten.Image
}

type MouseState uint8

const (
	MouseStateNone MouseState = iota
	MouseStatePressed
)

func New(terminal *termutil.Terminal, options ...Option) (*GUI, error) {

	g := &GUI{
		terminal:    terminal,
		size:        image.Point{80, 30},
		updateChan:  make(chan struct{}),
		fontManager: font.NewManager(),
		keyState:    newKeyState(),
	}

	for _, option := range options {
		if err := option(g); err != nil {
			return nil, err
		}
	}

	terminal.SetWindowManipulator(NewManipulator(g))

	return g, nil
}

func (g *GUI) Run() error {

	go func() {
		if err := g.terminal.Run(g.updateChan, uint16(g.size.X), uint16(g.size.Y)); err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	ebiten.SetScreenTransparent(true)
	ebiten.SetScreenClearedEveryFrame(true)
	ebiten.SetWindowResizable(true)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)

	for _, f := range g.startupFuncs {
		go f(g)
	}

	go g.watchForUpdate()

	return ebiten.RunGame(g)
}

func (g *GUI) watchForUpdate() {
	for range g.updateChan {
		ebiten.ScheduleFrame()
		if g.keyState.AnythingPressed() {
			go func() {
				time.Sleep(time.Millisecond * 10)
				ebiten.ScheduleFrame()
			}()
		}
	}
}
