package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// Card representa um cartão com sombra e bordas arredondadas
type Card struct {
	Color        color.NRGBA
	CornerRadius unit.Dp
	Elevation    unit.Dp
}

// ProgressBar representa uma barra de progresso moderna
type ProgressBar struct {
	Progress float32
	Color    color.NRGBA
	Height   unit.Dp
	Radius   unit.Dp
}

func (c Card) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	// Shadow
	if c.Elevation > 0 {
		shadowOffset := float32(c.Elevation) * 0.5
		shadowColor := color.NRGBA{A: 20}

		shadowMacro := op.Record(gtx.Ops)
		op.Offset(image.Pt(int(shadowOffset), int(shadowOffset))).Add(gtx.Ops)

		rect := image.Rectangle{
			Max: image.Point{X: dims.Size.X, Y: dims.Size.Y},
		}
		radius := gtx.Dp(c.CornerRadius)

		stack := clip.UniformRRect(rect, radius).Push(gtx.Ops)
		paint.Fill(gtx.Ops, shadowColor)
		stack.Pop()

		shadowCall := shadowMacro.Stop()
		shadowCall.Add(gtx.Ops)
	}

	// Card background
	rect := image.Rectangle{
		Max: image.Point{X: dims.Size.X, Y: dims.Size.Y},
	}
	radius := gtx.Dp(c.CornerRadius)

	stack := clip.UniformRRect(rect, radius).Push(gtx.Ops)
	paint.Fill(gtx.Ops, c.Color)
	stack.Pop()

	call.Add(gtx.Ops)
	return dims
}

func (p ProgressBar) Layout(gtx layout.Context) layout.Dimensions {
	height := gtx.Dp(p.Height)
	width := gtx.Constraints.Max.X

	// Background
	rect := image.Rectangle{
		Max: image.Point{X: width, Y: height},
	}
	radius := gtx.Dp(p.Radius)

	stack := clip.UniformRRect(rect, radius).Push(gtx.Ops)
	paint.Fill(gtx.Ops, ColorBorder)
	stack.Pop()

	// Progress
	if p.Progress > 0 {
		progressWidth := int(float32(width) * p.Progress)
		progressRect := image.Rectangle{
			Max: image.Point{X: progressWidth, Y: height},
		}
		stack := clip.UniformRRect(progressRect, radius).Push(gtx.Ops)
		paint.Fill(gtx.Ops, p.Color)
		stack.Pop()
	}

	return layout.Dimensions{
		Size: image.Pt(width, height),
	}
}
