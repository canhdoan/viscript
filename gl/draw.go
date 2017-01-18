package gl

import (
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/cGfx"
	"github.com/corpusc/viscript/gfx"
	"github.com/go-gl/gl/v2.1/gl"
)

func SetColor(newColor []float32) {
	cGfx.PrevColor = cGfx.CurrColor
	cGfx.CurrColor = newColor
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &newColor[0])
}

func DrawQuad(atlasX, atlasY int, r *app.Rectangle) {
	sp /* span */ := cGfx.UvSpan
	u := float32(atlasX) * sp
	v := float32(atlasY) * sp

	gl.Normal3f(0, 0, 1)

	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(r.Left, r.Bottom, 0)

	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(r.Right, r.Bottom, 0)

	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(r.Right, r.Top, 0)

	gl.TexCoord2f(u, v)
	gl.Vertex3f(r.Left, r.Top, 0)
}

func Draw9Sliced(r *app.PicRectangle) {
	// // skip invisible or inverted rects
	// if w <= 0 || h <= 0 {
	// 	return
	// }

	/*gl.Normal3f(0, 0, 1)

	for iX := 0; iX < 3; iX++ {
		for iY := 0; iY < 3; iY++ {
			gl.TexCoord2f(uSpots[iX], vSpots[iY+1]) // left bottom
			gl.Vertex3f(xSpots[iX], ySpots[iY+1], 0)

			gl.TexCoord2f(uSpots[iX+1], vSpots[iY+1]) // right bottom
			gl.Vertex3f(xSpots[iX+1], ySpots[iY+1], 0)

			gl.TexCoord2f(uSpots[iX+1], vSpots[iY]) // right top
			gl.Vertex3f(xSpots[iX+1], ySpots[iY], 0)

			gl.TexCoord2f(uSpots[iX], vSpots[iY]) // left top
			gl.Vertex3f(xSpots[iX], ySpots[iY], 0)
		}
	}*/
}

func drawAll() { // ATM ONLY draws 9slices, but without 9 slicing them
	gfx.DrawAll()
	cGfx.DrawAll()

	// skip 0 so we can use it as a code for being uninitialized
	for i := 1; i < len(cGfx.Rects); i++ {
		if cGfx.Rects[i].State == app.RectState_Active {
			//gfx.SetColor(gfx.Rects[i].Color)

			if cGfx.Rects[i].Type == app.RectType_9Slice {
				Draw9Sliced(cGfx.Rects[i])
				DrawQuad(cGfx.Pic_GradientBorder.X, cGfx.Pic_GradientBorder.Y, cGfx.Rects[i].Rectangle)
			} else {
				DrawQuad(cGfx.Pic_GradientBorder.X, cGfx.Pic_GradientBorder.Y, cGfx.Rects[i].Rectangle)
			}
		}
	}
}
