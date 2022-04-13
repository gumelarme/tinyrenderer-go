package main

import (
	"image/color"
	// "math/rand"

	"github.com/gumelarme/tinyrenderer-go/pkg/targa"
	m "github.com/gumelarme/tinyrenderer-go/pkg/model"
)

var (
	white  = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	red    = color.RGBA{0xFF, 0, 0, 0xFF}
	yellow = color.RGBA{0xFF, 0xFF, 0, 0xFF}
	black  = color.RGBA{0, 0, 0, 0xFF}
)

func main() {
	width, height := 800, 800
	image := targa.NewImage(width, height, targa.UncompressedRGB, 24)
	image.FillRGB(black)
	model, err := m.NewModel("obj/african_head.obj")
	if err != nil {
		panic(err)
	}

	light_dir := m.Vec3f{0, 0, -120}
	for f := 0; f < model.FacesCount(); f++ {
		face := model.GetFace(f)
		screenCoord := make([]m.Vec2, 3)
		worldCoord := make([]m.Vec3f, 3)
		for j := 0; j < len(face); j++ {
			v, _ := model.GetVertex(face[j][0])
			screenCoord[j] = m.Vec2{
				X: (int((v.X + 1) * float32(width) / 2)),
				Y: int((v.Y + 1) * float32(height) / 2),
			}
			worldCoord[j] = v.ToVec3f()
		}

		// col := color.RGBA{
		// 	uint8(rand.Intn(256)),
		// 	uint8(rand.Intn(256)),
		// 	uint8(rand.Intn(256)),
		// 	1,
		// }

		n := worldCoord[2].Subtract(worldCoord[0])
		n = n.Power(worldCoord[1].Subtract(worldCoord[0]))
		n.Normalize()

		intensity := light_dir.MultiplyVec(n)
		if intensity > 0 {
			shade := uint8(intensity * 2)
			BarycentricTriangle(screenCoord, image, color.RGBA{shade, shade, shade, 255})
			// BarycentricTriangle(screenCoord, image, col) // colorful
		}

	}
	image.WriteToFile("output.tga")
}

func wireframe() {
	width, height := 800, 800
	image := targa.NewImage(width, height, targa.UncompressedRGB, 24)
	image.FillRGB(black)
	m, err := m.NewModel("obj/african_head.obj")
	if err != nil {
		panic(err)
	}

	for f := 0; f < m.FacesCount(); f++ {
		face := m.GetFace(f)
		for v := 0; v < len(face); v++ {
			// get vertex only
			v0, err := m.GetVertex(face[v][0])

			if err != nil {
				panic(err)
			}

			v1, err := m.GetVertex(face[(v+1)%3][0])

			if err != nil {
				panic(err)
			}

			var shift, scale float32 = 1, 2

			x0 := int((v0.X + shift) * float32(width) / scale)
			y0 := int((v0.Y + shift) * float32(height) / scale)

			x1 := int((v1.X + shift) * float32(width) / scale)
			y1 := int((v1.Y + shift) * float32(height) / scale)

			// fmt.Printf("Drawing from (%d, %d) to (%d, %d)\n", x0, y0, x1, y1)
			DrawLine(x0, y0, x1, y1, image, yellow)
		}
	}
	image.WriteToFile("output.tga")
}

func DrawLineVec2(t0, t1 m.Vec2, image *targa.TGAImage, col color.RGBA) {
	DrawLine(t0.X, t0.Y, t1.X, t1.Y, image, col)
}

func DrawLine(x0, y0, x1, y1 int, image *targa.TGAImage, col color.RGBA) {
	steep := false
	if Abs(x1-x0) < Abs(y1-y0) {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
		steep = true
	}

	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	dx := x1 - x0
	dy := y1 - y0

	derror := float64(Abs(dy) * 2)
	err := 0.0
	y := y0

	for x := x0; x <= x1; x++ {
		if steep {
			// skip if out of bound
			if e := image.SetPixelRGB(y, x, col); e != nil {
				continue
			}
		} else {
			// skip if out of bound
			if e := image.SetPixelRGB(x, y, col); e != nil {
				continue
			}
		}

		err += derror
		if err >= 0.5 {
			if y1 > y0 {
				y += 1
			} else {
				y -= 1
			}
			err -= float64(dx) * 2
		}

	}
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Abs(val int) int {
	if val < 0 {
		val = -val
	}
	return val
}

func Absf(val float64) float64 {
	if val < 0 {
		val = -val
	}
	return val
}

func Triangle(t0, t1, t2 m.Vec2, image *targa.TGAImage, col color.RGBA) {
	DrawLineVec2(t0, t1, image, col)
	DrawLineVec2(t0, t2, image, col)
	DrawLineVec2(t1, t2, image, col)
	// Sort'em
	if t0.Y > t1.Y {
		t0, t1 = t1, t0
	}

	if t0.Y > t2.Y {
		t0, t2 = t2, t0
	}

	if t1.Y > t2.Y {
		t1, t2 = t2, t1
	}

	totalHeight := t2.Y - t0.Y
	segmentHeight := t1.Y - t0.Y + 1
	for y := t0.Y; y < t1.Y; y++ {
		alpha := float32(y-t0.Y) / float32(totalHeight)
		beta := float32(y-t0.Y) / float32(segmentHeight)
		a := t0.Add(t2.Subtract(t0).MultiplyFloat32(alpha))
		b := t0.Add(t1.Subtract(t0).MultiplyFloat32(beta))
		if a.X > b.X {
			a, b = b, a
		}

		for j := a.X; j < b.X; j++ {
			image.SetPixelRGB(j, y, col)
		}
	}

	segmentHeight = t2.Y - t1.Y + 1
	for y := t1.Y; y < t2.Y; y++ {
		alpha := float32(y-t0.Y) / float32(totalHeight)
		beta := float32(y-t1.Y) / float32(segmentHeight)
		a := t0.Add(t2.Subtract(t0).MultiplyFloat32(alpha))
		b := t1.Add(t2.Subtract(t1).MultiplyFloat32(beta))
		if a.X > b.X {
			a, b = b, a
		}

		for j := a.X; j < b.X; j++ {
			image.SetPixelRGB(j, y, col)
		}
	}
}

func Barycentric(points []m.Vec2, p m.Vec2) m.Vec3f {
	u := m.Cross(
		m.Vec3f{
			X: float64(points[2].X - points[0].X),
			Y: float64(points[1].X - points[0].X),
			Z: float64(points[0].X - p.X),
		},

		m.Vec3f{
			X: float64(points[2].Y - points[0].Y),
			Y: float64(points[1].Y - points[0].Y),
			Z: float64(points[0].Y - p.Y),
		},
	)

	if Absf(u.Z) < 1 {
		return m.Vec3f{X: -1, Y: 1, Z: 1}
	}

	return m.Vec3f{
		X: 1.0 - (u.X+u.Y)/u.Z,
		Y: u.Y / u.Z,
		Z: u.X / u.Z,
	}
}

func BarycentricTriangle(points []m.Vec2, image *targa.TGAImage, col color.RGBA) {
	bboxMin := m.Vec2{X: image.Width - 1, Y: image.Height - 1}
	clamp := bboxMin
	bboxMax := m.Vec2{X: 0, Y: 0}
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			bboxMin.Set(j, Max(0, Min(bboxMin.Get(j), points[i].Get(j))))
			bboxMax.Set(j, Min(clamp.Get(j), Max(bboxMax.Get(j), points[i].Get(j))))
		}
	}

	var P m.Vec2

	for P.X = bboxMin.X; P.X <= bboxMax.X; P.X++ {
		for P.Y = bboxMin.Y; P.Y <= bboxMax.Y; P.Y++ {
			screen := Barycentric(points, P)
			// fmt.Println(screen)
			if screen.X < 0 || screen.Y < 0 || screen.Z < 0 {
				continue
			}
			image.SetPixelRGB(P.X, P.Y, col)
		}
	}
}
