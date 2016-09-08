package main

import (
	//	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
)

/*
mouse position updates use pixels, so the smallest drag motions will be
a jump of at least 1 pixel height.
the ratio of that height / LenOfVoid (bar representing the page size),
compared to the void/offscreen length of the text body,
gives us the jump size in scrolling through the text body
*/

type ScrollBar struct {
	PosX      float32
	PosY      float32
	LenOfBar  float32
	LenOfVoid float32 // length of the negative space representing the length of entire document
}

func (bar *ScrollBar) StartFrame(tp TextPanel) {
	hei /* height of panel */ := textRend.chHei * float32(tp.NumCharsY)

	if /* content smaller than screen */ len(tp.Body) <= tp.NumCharsY {
		// no need for scrollbar
		bar.LenOfBar = 0
		bar.LenOfVoid = hei
		code.LenOfOffscreenY = 0
	} else {
		bar.LenOfBar = float32(tp.NumCharsY) / float32(len(tp.Body)) * hei
		bar.LenOfVoid = hei - bar.LenOfBar
		code.LenOfOffscreenY = float32(len(tp.Body)-tp.NumCharsY) * textRend.chHei
	}
}

func (bar *ScrollBar) Scroll(tp *TextPanel, mousePixelDeltaY float64) {
	// y increment (for bar) in gl space
	yInc := float32(mousePixelDeltaY) * textRend.pixelHei

	bar.PosY -= yInc

	if bar.PosY < tp.Bottom+bar.LenOfBar {
		bar.PosY = tp.Bottom + bar.LenOfBar
	}
	if bar.PosY > tp.Top {
		bar.PosY = tp.Top
	}

	tp.OffsetY -= yInc / bar.LenOfVoid * tp.LenOfOffscreenY

	if tp.OffsetY > 0 {
		tp.OffsetY = 0
	}

	if tp.OffsetY < -tp.LenOfOffscreenY {
		tp.OffsetY = -tp.LenOfOffscreenY
	}
}

func (bar *ScrollBar) DrawVertical(atlasX, atlasY float32) {
	rad := textRend.ScreenRad
	sp := textRend.UvSpan
	u := float32(atlasX) * sp
	v := float32(atlasY) * sp

	top := bar.PosY                 //rad - 1
	bott := bar.PosY - bar.LenOfBar //-rad + 1

	gl.Normal3f(0, 0, 1)

	// bottom left   0, 1
	gl.TexCoord2f(u, v+sp)
	gl.Vertex3f(rad-textRend.chWid, bott, 0)

	// bottom right   1, 1
	gl.TexCoord2f(u+sp, v+sp)
	gl.Vertex3f(rad, bott, 0)

	// top right   1, 0
	gl.TexCoord2f(u+sp, v)
	gl.Vertex3f(rad, top, 0)

	// top left   0, 0
	gl.TexCoord2f(u, v)
	gl.Vertex3f(rad-textRend.chWid, top, 0)

	textRend.CurrX += textRend.chWid
}