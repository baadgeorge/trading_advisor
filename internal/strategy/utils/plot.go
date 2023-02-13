package utils

import (
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
)

func CandlesToPlot(plots map[string][]CandleStruct) (*plot.Plot, error) {
	p := plot.New()
	//xticks := plot.TimeTicks{Format: "2006-01-02\n15:04"}
	//p.X.Tick.Marker = xticks

	p.X.Tick.Marker = plot.TimeTicks{
		Ticker: hplot.Ticks{N: 12},
		Format: "2006-01-02\n15:04",
	}
	p.Y.Tick.Marker = hplot.Ticks{N: 10}

	//p.X.Tick.Marker = hplot.Ticks{Format: "2006-01-02\n15:04", N: 10}
	//grid := plotter.NewGrid()
	p.Add(plotter.NewGrid())
	//p.Title.Text =.
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "Price"

	count := 0

	for name, candles := range plots {
		pts := make(plotter.XYs, len(candles))

		for k, v := range candles {
			pts[k].X = float64(v.Period.End.Unix())
			pts[k].Y = v.Value
		}

		line, points, err := plotter.NewLinePoints(pts)
		if err != nil {
			return nil, err
		}

		switch count {
		case 0:
			line.Color = color.RGBA{G: 255, A: 255}
		case 1:
			line.Color = color.RGBA{B: 255, A: 255}
		case 2:
			line.Color = color.RGBA{G: 100, B: 70, R: 20, A: 255}
		case 3:
			line.Color = color.RGBA{G: 200, B: 150, A: 255}
		}
		count++

		points.Shape = draw.CircleGlyph{}
		points.Radius = vg.Length(2)
		points.Color = color.RGBA{R: 255, A: 255}
		p.Add(line, points)

		//leg.Add(name, line)
		//p.Legend = leg

		p.Legend.Add(name, line)
		//p.Legend.Padding = vg.Length(5)

	}
	//p.Legend.Rectangle()
	p.Legend.XOffs = vg.Length(-5)
	//err = p.Save(10*vg.Centimeter, 5*vg.Centimeter, "timeseries.png")

	/*//TODO
	if err != nil {
		panic(err)

	}*/
	return p, nil
}
