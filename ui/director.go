package ui

import (
	"log"

	"github.com/fogleman/nes/nes"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type View interface {
	Enter()
	Exit()
	Update(t, dt float64)
	Draw()
}

type Director struct {
	window    *glfw.Window
	audio     *Audio
	view      View
	timestamp float64
}

func NewDirector(window *glfw.Window, audio *Audio) *Director {
	director := Director{}
	director.window = window
	director.audio = audio
	return &director
}

func (d *Director) SetTitle(title string) {
	d.window.SetTitle(title)
}

func (d *Director) SetView(view View) {
	if d.view != nil {
		d.view.Exit()
	}
	d.view = view
	if d.view != nil {
		d.view.Enter()
	}
	d.timestamp = glfw.GetTime()
}

func (d *Director) Step() {
	now := glfw.GetTime()
	elapsed := now - d.timestamp
	d.timestamp = now
	d.view.Update(d.timestamp, elapsed)
	d.view.Draw()
}

func (d *Director) Run() {
	for !d.window.ShouldClose() {
		d.Step()
		d.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (d *Director) PlayROM(path string) {
	console, err := nes.NewConsole(path)
	if err != nil {
		log.Fatalln(err)
	}
	d.SetView(NewGameView(d, console, path))
}