package zraster

// zColor represents a candidate pixel together with its z depth.
type zColor struct {
	// Color premultiplied with rasterizer span alpha.
	r, g, b, a uint32
	// Depth info.
	z float64
}

type sortByZ []*zColor

func (byz sortByZ) Len() int      { return len(byz) }
func (byz sortByZ) Swap(i, j int) { byz[i], byz[j] = byz[j], byz[i] }

// We are looking from -z to +z, so a greater z value needs to be painted first.
func (byz sortByZ) Less(i, j int) bool { return byz[i].z > byz[j].z }
