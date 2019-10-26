package identicon

import (
	"crypto/md5"
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/llgcode/draw2d/draw2dimg"
)

type Point struct {
	x, y int
}

type DrawingPoint struct {
	topLeft     Point
	bottomRight Point
}

type GridPoint struct {
	value byte
	index int
}

type Identicon struct {
	Name       string
	hash       [16]byte
	color      [3]byte
	grid       []byte
	gridPoints []GridPoint
	pixelMap   []DrawingPoint
}

// WriteTo writes the identicon image to the given writer
func (i Identicon) WriteImage(w io.Writer) error {
	var img = image.NewRGBA(image.Rect(0, 0, 250, 250))
	col := color.RGBA{R: i.color[0], G: i.color[1], B: i.color[2], A: 255}

	for _, pixel := range i.pixelMap {
		rect(img, col, float64(pixel.topLeft.x), float64(pixel.topLeft.y), float64(pixel.bottomRight.x), float64(pixel.bottomRight.y))
	}

	return png.Encode(w, img)
}

type ApplyFunc func(Identicon) Identicon

func Generate(input []byte) Identicon {
	identiconPipe := []ApplyFunc{
		pickColor, buildGrid, filterOddSquares, buildPixelMap,
	}
	identicon := hashInput(input)
	for _, applyFunc := range identiconPipe {
		identicon = applyFunc(identicon)
	}
	return identicon
}

func hashInput(input []byte) Identicon {
	checkSum := md5.Sum(input)
	return Identicon{
		Name: string(input),
		hash: checkSum,
	}
}

func pickColor(identicon Identicon) Identicon {
	rgb := [3]byte{}
	copy(rgb[:], identicon.hash[:3])
	identicon.color = rgb
	return identicon
}

func buildGrid(identicon Identicon) Identicon {
	var grid []byte
	for i := 0; i < len(identicon.hash) && i+3 <= len(identicon.hash)-1; i += 3 {
		chunk := make([]byte, 5)
		copy(chunk, identicon.hash[i:i+3])
		chunk[3] = chunk[1]
		chunk[4] = chunk[0]
		grid = append(grid, chunk...)

	}
	identicon.grid = grid
	return identicon
}

func filterOddSquares(identicon Identicon) Identicon {
	var grid []GridPoint
	for i, code := range identicon.grid {
		if code%2 == 0 {
			point := GridPoint{
				value: code,
				index: i,
			}
			grid = append(grid, point)
		}
	}
	identicon.gridPoints = grid
	return identicon
}

func rect(img *image.RGBA, col color.Color, x1, y1, x2, y2 float64) {
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(col)
	gc.MoveTo(x1, y1)
	gc.LineTo(x1, y1)
	gc.LineTo(x1, y2)
	gc.MoveTo(x2, y1)
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	gc.SetLineWidth(0)
	gc.FillStroke()
}

func buildPixelMap(identicon Identicon) Identicon {
	var drawingPoints []DrawingPoint

	pixelFunc := func(p GridPoint) DrawingPoint {
		horizontal := (p.index % 5) * 50
		vertical := (p.index / 5) * 50
		topLeft := Point{horizontal, vertical}
		bottomRight := Point{horizontal + 50, vertical + 50}

		return DrawingPoint{
			topLeft,
			bottomRight,
		}
	}

	for _, gridPoint := range identicon.gridPoints {
		drawingPoints = append(drawingPoints, pixelFunc(gridPoint))
	}
	identicon.pixelMap = drawingPoints
	return identicon
}
