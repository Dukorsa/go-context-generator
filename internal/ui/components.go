package ui

import (
	"image"
	"image/color"

	"gioui.org/f32" // For f32.Point
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// Card provides a styled surface for content.
type Card struct {
	Color        color.NRGBA // Background color of the card
	CornerRadius unit.Dp     // Radius for all corners
	Elevation    unit.Dp     // Controls shadow depth; 0 for no shadow
	BorderWidth  unit.Dp     // Optional border width
	BorderColor  color.NRGBA // Optional border color
	Hoverable    bool        // If true, card can visually respond to hover (requires external Clickable)
	hovered      bool        // Internal state, set externally if Hoverable
}

// SetHovered allows external logic (like a Clickable) to inform the card it's being hovered.
func (c *Card) SetHovered(hovered bool) {
	if c.Hoverable {
		c.hovered = hovered
	}
}

// Layout draws the card and then the provided widget w within it.
func (c Card) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	// --- Shadow Layer ---
	if c.Elevation > 0 {
		// Slightly larger offset for a more prominent shadow, but still soft
		// Elevation is used to scale the shadow properties
		elevationPx := gtx.Dp(c.Elevation)

		// Shadow color - soft black, more transparent for larger blurs
		shadowAlpha := uint8(max(10, 30-elevationPx)) // Alpha decreases with elevation, min 10
		shadowColor := color.NRGBA{A: shadowAlpha}

		// The shadow will be slightly larger than the card to simulate a blur
		// We'll draw the shadow as a rounded rect slightly offset and larger
		// Note: True blur is complex in Gio's immediate mode. This is an approximation.
		// A simple offset shadow is more straightforward:
		offset := image.Pt(gtx.Dp(c.Elevation)/3, gtx.Dp(c.Elevation)/2)
		shadowOps := op.Record(gtx.Ops)

		// Original widget dimensions will be determined first
		dims := w(gtx) // Call widget function once to get dimensions

		// Shadow Rectangle
		shadowRect := image.Rectangle{Max: dims.Size}
		shadowRadius := gtx.Dp(c.CornerRadius)

		// Apply offset for the shadow
		op.Offset(offset).Add(gtx.Ops)

		cl := clip.UniformRRect(shadowRect, shadowRadius).Push(gtx.Ops)
		paint.Fill(gtx.Ops, shadowColor) // Use the computed shadow color
		cl.Pop()

		call := shadowOps.Stop()
		call.Add(gtx.Ops) // Add the recorded shadow operations

		// Redraw the actual widget on top (content)
		// This is a common pattern: measure, draw effects, draw content
		// However, the original Card drew content first, then background.
		// Let's stick to: record content, draw bg/shadow, play content.
	}

	// --- Main Card Surface & Content ---
	macro := op.Record(gtx.Ops)
	dims := w(gtx) // Content is laid out here
	call := macro.Stop()

	// Card background color
	bgColor := c.Color
	if c.Hoverable && c.hovered {
		// Subtle hover effect: slightly lighten or change background
		// This is a placeholder; actual color depends on your theme.
		// Example: make it a bit lighter if ColorSurface
		if bgColor == ColorSurface {
			bgColor = color.NRGBA{R: 250, G: 250, B: 252, A: 255} // Very light blue/gray tint
		} else {
			bgColor.A = 230 // Slightly more transparent or lighter
		}
	}

	cardRect := image.Rectangle{Max: dims.Size}
	radiusPx := gtx.Dp(c.CornerRadius)

	// Draw background
	bgClip := clip.UniformRRect(cardRect, radiusPx).Push(gtx.Ops)
	paint.Fill(gtx.Ops, bgColor)
	bgClip.Pop()

	// --- Border Layer (Optional) ---
	if c.BorderWidth > 0 {
		borderWidthPx := float32(gtx.Dp(c.BorderWidth))
		borderColor := c.BorderColor
		if (borderColor == color.NRGBA{}) { // Default border color if not specified
			borderColor = ColorBorder // From your theme.go
		}

		// To draw a border, we stroke a path inset by half the border width
		// so the border is drawn on the edge of the card background.
		// Or, more simply, stroke the same rrect path with the desired width.
		// The stroke is centered on the path, so it will extend inwards and outwards.
		// For an "inner" border, you'd inset the shape first.
		// For an "outer" border, you'd outset it.
		// For a centered border, use the same shape.

		// This is tricky. Simpler: draw the border with clip.Stroke on the original cardRect

		clip.Stroke{
			Path:  clip.RRect{Rect: cardRect, SE: radiusPx, SW: radiusPx, NW: radiusPx, NE: radiusPx}.Path(gtx.Ops),
			Width: borderWidthPx,
		}.Op().Push(gtx.Ops)

		paint.ColorOp{Color: borderColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		bgClip.Pop()
	}

	// Add the recorded content operations (drawn on top)
	call.Add(gtx.Ops)

	return dims
}

// ProgressBar displays a progress visually.
type ProgressBar struct {
	Progress   float32     // Value between 0.0 and 1.0
	Color      color.NRGBA // Main color for the progress fill (e.g., theme.ColorPrimary)
	TrackColor color.NRGBA // Color for the background track of the bar
	Height     unit.Dp     // Height of the progress bar
	Radius     unit.Dp     // Corner radius for the bar and fill
}

// Layout draws the progress bar.
func (p ProgressBar) Layout(gtx layout.Context) layout.Dimensions {
	heightPx := gtx.Dp(p.Height)
	widthPx := gtx.Constraints.Max.X // Fill available width

	// Use a default track color if not specified, from theme or a light gray
	trackColor := p.TrackColor
	if (trackColor == color.NRGBA{}) {
		trackColor = ColorBorderLight // Assuming ColorBorderLight is defined in theme.go (e.g., a very light gray)
		// If not, use a hardcoded default: color.NRGBA{R: 230, G: 230, B: 230, A: 255}
	}

	// --- Draw Track ---
	trackRect := image.Rectangle{
		Max: image.Point{X: widthPx, Y: heightPx},
	}
	radiusPx := gtx.Dp(p.Radius)
	if radiusPx > heightPx/2 { // Cap radius to half height
		radiusPx = heightPx / 2
	}

	// Clip for the track
	trackClip := clip.UniformRRect(trackRect, radiusPx).Push(gtx.Ops)
	paint.Fill(gtx.Ops, trackColor)
	trackClip.Pop()

	// --- Draw Progress Fill ---
	if p.Progress > 0 {
		fillWidthPx := int(float32(widthPx) * p.Progress)
		if fillWidthPx < 2*radiusPx && fillWidthPx > 0 { // Ensure minimum width for visible rounding
			// If the fill width is less than twice the radius,
			// it might look odd or disappear. Adjust if needed or ensure progress is significant.
		}
		if fillWidthPx <= 0 { // No fill needed if width is zero or less
			return layout.Dimensions{Size: image.Pt(widthPx, heightPx)}
		}

		fillRect := image.Rectangle{
			Max: image.Point{X: fillWidthPx, Y: heightPx},
		}

		// Gradient for the fill (subtle)
		// Make the end color slightly lighter or a different hue of p.Color
		// This is a simple way to get a lighter version.
		// For more control, you might define StartColor and EndColor in the struct.
		gradientEndColor := p.Color
		gradientEndColor.R = uint8(min(255, int(p.Color.R)+20))
		gradientEndColor.G = uint8(min(255, int(p.Color.G)+20))
		gradientEndColor.B = uint8(min(255, int(p.Color.B)+20))

		fillClip := clip.UniformRRect(fillRect, radiusPx).Push(gtx.Ops)

		paint.LinearGradientOp{
			Stop1:  f32.Point{X: 0, Y: float32(heightPx) / 2.0}, // Gradient direction (e.g., left to right)
			Stop2:  f32.Point{X: float32(fillWidthPx), Y: float32(heightPx) / 2.0},
			Color1: p.Color,
			Color2: gradientEndColor, // A slightly lighter/different shade of p.Color
		}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		fillClip.Pop()
	}

	return layout.Dimensions{
		Size: image.Point{X: widthPx, Y: heightPx},
	}
}

// Helper for min/max if not available in this Go version's stdlib for uint8
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// func min(a, b int) int { // Not used here, but good to have
// 	if a < b {
// 		return a
// 	}
// 	return b
// }
