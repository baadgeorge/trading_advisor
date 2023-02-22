package utils

import (
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"time"
)

/*type MyTicks struct{}

func (MyTicks) Ticks(min, max float64) []plot.Tick {
	if max <= min {
		panic("illegal range")
	}
	var ticks []plot.Tick
	k := (max - min) / 2
	for i := 0; i <= 12; i++ {
		tmp := min + k*float64(i+1)
		ticks = append(ticks, plot.Tick{Value: tmp, Label: "jddj"})
	}
	return ticks
}*/

func CustomXAxis(names []string, p *plot.Plot) *plot.Plot {
	p.X.Tick.Width = 0.5
	p.X.Tick.Length = 6
	p.X.Width = 0.5

	//p.Y.Padding = p.X.Tick.Label.Width(names[0]) / 2

	var numbers int
	l := len(names)
	if l < 10 {
		numbers = l
	} else if l < 20 {
		numbers = l / 2
	} else if l < 40 {
		numbers = l / 3
	} else if l < 100 {
		numbers = l / 4
	} else {
		numbers = l / 5
	}

	ticks := make([]plot.Tick, len(names))
	step := len(names) / numbers
	for i, name := range names {
		if i%step == 0 {
			ticks[i] = plot.Tick{float64(i), name}
		} else {
			ticks[i] = plot.Tick{float64(i), ""}
			ticks[i].IsMinor()
		}
	}
	p.X.Tick.Marker = plot.ConstantTicks(ticks)
	return p
}

func CandlesToPlot(plots map[string][]CandleStruct) (*plot.Plot, error) {
	p := plot.New()
	//xticks := plot.TimeTicks{Format: "2006-01-02\n15:04"}
	//Format: "2006-01-02\n15:04"}
	//Time:   plot.UnixTimeIn(time.Local)}

	//p.X.Tick.Marker = xticks
	/*{
		Ticker: hplot.Ticks{N: 12},
		Format: "2006-01-02\n15:04",
	}*/
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
			//pts[k].X = float64(v.Period.End.Unix())
			pts[k].X = float64(k)
			pts[k].Y = v.Value
		}

		line, points, err := plotter.NewLinePoints(pts)
		if err != nil {
			return nil, err
		}

		switch count {
		case 0:
			line.Color = color.RGBA{G: 255, A: 255}
			var xTime []string
			for _, v := range candles {

				xTime = append(xTime, v.Period.End.In(time.Local).Format("2006-01-02\n15:04"))
			}
			CustomXAxis(xTime, p)
			//p.NominalX(xTime...)
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
