package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const axisMargin = 30
const scalingFactor = 30

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	points := []image.Point{convertPoint(image.Point{0, 0}, img.Bounds()), convertPoint(image.Point{40, 400}, img.Bounds()), convertPoint(image.Point{250, 200}, img.Bounds())}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: basicfont.Face7x13,
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)
	drawVerticalAxis(img, img.Bounds(), color.White)
	drawHorizontalAxis(img, img.Bounds(), color.White)
	drawVerticalPegs(img, img.Bounds(), color.White)
	drawHorizontalPegs(img, img.Bounds(), color.White)
	drawHorizontalLabels(img, d)
	drawVerticalLabels(img, d)
	drawTitle(img, d, "My Title")
	drawPoints(img, points, color.White)
	encodeImageToFile(img)
}

func drawVerticalAxis(img draw.Image, rect image.Rectangle, color color.Color) {
	cutoff := calculateAxisCutoff(rect.Max.Y)
	for i := rect.Max.Y - axisMargin; i >= cutoff; i-- {
		img.Set(axisMargin, i, color)
	}
}

func drawHorizontalAxis(img draw.Image, rect image.Rectangle, color color.Color) {
	cutoff := calculateAxisCutoff(rect.Max.X)
	for i := axisMargin; i <= rect.Max.X-cutoff; i++ {
		img.Set(i, rect.Max.Y-axisMargin, color)
	}
}

func drawVerticalPegs(img draw.Image, rect image.Rectangle, color color.Color) {
	for i := rect.Max.Y - axisMargin - scalingFactor; i > 0; i -= scalingFactor {
		img.Set(axisMargin-2, i, color)
		img.Set(axisMargin-1, i, color)
		img.Set(axisMargin+1, i, color)
		img.Set(axisMargin+2, i, color)
	}
}

func drawHorizontalPegs(img draw.Image, rect image.Rectangle, color color.Color) {
	for i := axisMargin + scalingFactor; i < rect.Max.X; i += scalingFactor {
		img.Set(i, rect.Max.Y-axisMargin-2, color)
		img.Set(i, rect.Max.Y-axisMargin-1, color)
		img.Set(i, rect.Max.Y-axisMargin+1, color)
		img.Set(i, rect.Max.Y-axisMargin+2, color)
	}
}

func drawVerticalLabels(img *image.RGBA, d *font.Drawer) {
	padding := 5
	labelNum := scalingFactor
	rect := img.Bounds()
	for i := rect.Max.Y - axisMargin - scalingFactor; i > 0; i -= scalingFactor {
		label := strconv.Itoa(labelNum)
		x := fixed.I(axisMargin) - d.MeasureString(label) - fixed.I(padding)
		y := fixed.I(i + padding)
		point := fixed.Point26_6{
			X: x,
			Y: y,
		}
		drawText(d, point, label)
		labelNum += scalingFactor
	}
}

func drawHorizontalLabels(img *image.RGBA, d *font.Drawer) {
	padding := 5
	labelNum := scalingFactor
	rect := img.Bounds()
	for i := rect.Min.X + axisMargin + scalingFactor; i < rect.Max.X; i += scalingFactor {
		label := strconv.Itoa(labelNum)
		x := fixed.I(i - padding)
		y := fixed.I(rect.Max.Y - axisMargin + (3 * padding))
		point := fixed.Point26_6{
			X: x,
			Y: y,
		}
		drawText(d, point, label)
		labelNum += scalingFactor
	}
}

func drawTitle(img *image.RGBA, d *font.Drawer, title string) {
	rect := img.Bounds()
	x := fixed.I(rect.Max.X/2) - (d.MeasureString(title) / 2)
	y := fixed.I(rect.Min.Y + 10)
	point := fixed.Point26_6{
		X: x,
		Y: y,
	}
	drawText(d, point, title)
}

func drawText(d *font.Drawer, point fixed.Point26_6, label string) {
	d.Dot = point
	d.DrawString(label)
}

func drawPoints(img draw.Image, points []image.Point, color color.Color) {
	for i := 1; i < len(points); i++ {
		drawLine(img, color, points[i-1].X, points[i-1].Y, points[i].X, points[i].Y)
	}
}

func drawLine(img draw.Image, color color.Color, x0, y0, x1, y1 int) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	var e2 int
	for {
		img.Set(x0, y0, color)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 = 2 * err
		if e2 > -dy {
			err = err - dy
			x0 = x0 + sx
		}
		if e2 < dx {
			err = err + dx
			y0 = y0 + sy
		}
	}
}

func encodeImageToFile(img image.Image) {
	f, err := os.Create("graph.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	opts := &jpeg.Options{Quality: 100}
	encodeErr := jpeg.Encode(f, img, opts)
	if encodeErr != nil {
		panic(encodeErr)
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func convertPoint(pt image.Point, rect image.Rectangle) image.Point {
	return image.Point{pt.X + axisMargin, (rect.Max.Y - pt.Y - axisMargin)}
}

func calculateAxisCutoff(sideLength int) int {
	lengthRemaining := sideLength - axisMargin
	for {
		if lengthRemaining < scalingFactor {
			return lengthRemaining
		}
		lengthRemaining -= scalingFactor
	}
}
