package graphix

type Screen struct {
	width  float64
	height float64
	x0     float64
	y0     float64
	xscale float64
	yscale float64
}

// NewScreen creates a Screen with dimension width and height, mapping into the rectangle [x0,x1)×(y0,y1].
func NewScreen(width, height, x0, y0, x1, y1 float64) *Screen {
	return &Screen{
		width:  width,
		height: height,
		x0:     x0,
		y0:     y0,
		xscale: width / (x1 - x0),
		yscale: height / (y1 - y0),
	}
}

// Map maps a projection q (world coordinate: right for +x, up for +y) into
// p (screen coordinate: right for +x, down for +y) and returns p.
func (sc *Screen) Map(p *Projection, q *Projection) *Projection {
	p[0] = (q[0] - sc.x0) * sc.xscale
	p[1] = sc.height - (q[1]-sc.y0)*sc.yscale
	return p
}
