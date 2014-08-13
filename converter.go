/*
Copyright (c) 2012, Jan Schlicht <jan.schlicht@gmail.com>

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted, provided that the above copyright notice
and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND
FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS
OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF
THIS SOFTWARE.
*/

package resize

import "image"

// Keep value in [0,255] range.
func clampUint8(in int32) uint8 {
	if in < 0 {
		return 0
	}
	if in > 255 {
		return 255
	}
	return uint8(in)
}

// Keep value in [0,65535] range.
func clampUint16(in int64) uint16 {
	if in < 0 {
		return 0
	}
	if in > 65535 {
		return 65535
	}
	return uint16(in)
}

func resizeGeneric(in image.Image, out *image.RGBA64, scale float64, coeffs []int32, offset []int, filterLength int) {
	oldBounds := image.Rect(0, 0, in.Bounds().Dx(), in.Bounds().Dy())
	newBounds := out.Bounds()

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int64
			var sum int64
			start := offset[y]
			ci := (y - newBounds.Min.Y) * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(oldBounds.Max.X):
						break
					case xi >= oldBounds.Max.X:
						xi = oldBounds.Min.X
					default:
						xi = oldBounds.Max.X - 1
					}
					r, g, b, a := in.At(xi, x).RGBA()
					rgba[0] += int64(coeff) * int64(r)
					rgba[1] += int64(coeff) * int64(g)
					rgba[2] += int64(coeff) * int64(b)
					rgba[3] += int64(coeff) * int64(a)
					sum += int64(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := clampUint16(rgba[0] / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
			value = clampUint16(rgba[1] / sum)
			out.Pix[offset+2] = uint8(value >> 8)
			out.Pix[offset+3] = uint8(value)
			value = clampUint16(rgba[2] / sum)
			out.Pix[offset+4] = uint8(value >> 8)
			out.Pix[offset+5] = uint8(value)
			value = clampUint16(rgba[3] / sum)
			out.Pix[offset+6] = uint8(value >> 8)
			out.Pix[offset+7] = uint8(value)
		}
	}
}

func resizeRGBA(in *image.RGBA, out *image.RGBA, scale float64, coeffs []int16, offset []int, filterLength int) {
	oldBounds := image.Rect(0, 0, in.Rect.Dx(), in.Rect.Dy())
	newBounds := out.Bounds()
	minX := oldBounds.Min.X * 4
	maxX := (oldBounds.Max.X - oldBounds.Min.X - 1) * 4

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[(x-oldBounds.Min.Y)*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int32
			var sum int32
			start := offset[y]
			ci := (y - newBounds.Min.Y) * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(oldBounds.Max.X):
						xi *= 4
					case xi >= oldBounds.Max.X:
						xi = maxX
					default:
						xi = minX
					}
					rgba[0] += int32(coeff) * int32(row[xi+0])
					rgba[1] += int32(coeff) * int32(row[xi+1])
					rgba[2] += int32(coeff) * int32(row[xi+2])
					rgba[3] += int32(coeff) * int32(row[xi+3])
					sum += int32(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*4
			out.Pix[xo+0] = clampUint8(rgba[0] / sum)
			out.Pix[xo+1] = clampUint8(rgba[1] / sum)
			out.Pix[xo+2] = clampUint8(rgba[2] / sum)
			out.Pix[xo+3] = clampUint8(rgba[3] / sum)
		}
	}
}

func resizeRGBA64(in *image.RGBA64, out *image.RGBA64, scale float64, coeffs []int32, offset []int, filterLength int) {
	oldBounds := image.Rect(0, 0, in.Rect.Dx(), in.Rect.Dy())
	newBounds := out.Bounds()
	minX := oldBounds.Min.X * 8
	maxX := (oldBounds.Max.X - oldBounds.Min.X - 1) * 8

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[(x-oldBounds.Min.Y)*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]int64
			var sum int64
			start := offset[y]
			ci := (y - newBounds.Min.Y) * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(oldBounds.Max.X):
						xi *= 8
					case xi >= oldBounds.Max.X:
						xi = maxX
					default:
						xi = minX
					}
					rgba[0] += int64(coeff) * int64(uint16(row[xi+0])<<8|uint16(row[xi+1]))
					rgba[1] += int64(coeff) * int64(uint16(row[xi+2])<<8|uint16(row[xi+3]))
					rgba[2] += int64(coeff) * int64(uint16(row[xi+4])<<8|uint16(row[xi+5]))
					rgba[3] += int64(coeff) * int64(uint16(row[xi+6])<<8|uint16(row[xi+7]))
					sum += int64(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := clampUint16(rgba[0] / sum)
			out.Pix[xo+0] = uint8(value >> 8)
			out.Pix[xo+1] = uint8(value)
			value = clampUint16(rgba[1] / sum)
			out.Pix[xo+2] = uint8(value >> 8)
			out.Pix[xo+3] = uint8(value)
			value = clampUint16(rgba[2] / sum)
			out.Pix[xo+4] = uint8(value >> 8)
			out.Pix[xo+5] = uint8(value)
			value = clampUint16(rgba[3] / sum)
			out.Pix[xo+6] = uint8(value >> 8)
			out.Pix[xo+7] = uint8(value)
		}
	}
}

func resizeGray(in *image.Gray, out *image.Gray, scale float64, coeffs []int16, offset []int, filterLength int) {
	oldBounds := image.Rect(0, 0, in.Rect.Dx(), in.Rect.Dy())
	newBounds := out.Bounds()
	minX := oldBounds.Min.X
	maxX := (oldBounds.Max.X - oldBounds.Min.X - 1)

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[(x-oldBounds.Min.Y)*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray int32
			var sum int32
			start := offset[y]
			ci := (y - newBounds.Min.Y) * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(oldBounds.Max.X):
						break
					case xi >= oldBounds.Max.X:
						xi = maxX
					default:
						xi = minX
					}
					gray += int32(coeff) * int32(row[xi])
					sum += int32(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x - newBounds.Min.X)
			out.Pix[offset] = clampUint8(gray / sum)
		}
	}
}

func resizeGray16(in *image.Gray16, out *image.Gray16, scale float64, coeffs []int32, offset []int, filterLength int) {
	oldBounds := image.Rect(0, 0, in.Rect.Dx(), in.Rect.Dy())
	newBounds := out.Bounds()
	minX := oldBounds.Min.X * 2
	maxX := (oldBounds.Max.X - oldBounds.Min.X - 1) * 2

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[(x-oldBounds.Min.Y)*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var gray int64
			var sum int64
			start := offset[y]
			ci := (y - newBounds.Min.Y) * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(oldBounds.Max.X):
						xi *= 2
					case xi >= oldBounds.Max.X:
						xi = maxX
					default:
						xi = minX
					}
					gray += int64(coeff) * int64(uint16(row[xi+0])<<8|uint16(row[xi+1]))
					sum += int64(coeff)
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*2
			value := clampUint16(gray / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
		}
	}
}

func resizeYCbCr(in *ycc, out *ycc, scale float64, coeffs []int16, offset []int, filterLength int) {
	oldBounds := image.Rect(0, 0, in.Rect.Dx(), in.Rect.Dy())
	newBounds := out.Bounds()
	minX := oldBounds.Min.X * 3
	maxX := (oldBounds.Max.X - oldBounds.Min.X - 1) * 3

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[(x-oldBounds.Min.Y)*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var p [3]int32
			var sum int32
			start := offset[y]
			ci := (y - newBounds.Min.Y) * filterLength
			for i := 0; i < filterLength; i++ {
				coeff := coeffs[ci+i]
				if coeff != 0 {
					xi := start + i
					switch {
					case uint(xi) < uint(oldBounds.Max.X):
						xi *= 3
					case xi >= oldBounds.Max.X:
						xi = maxX
					default:
						xi = minX
					}
					p[0] += int32(coeff) * int32(row[xi+0])
					p[1] += int32(coeff) * int32(row[xi+1])
					p[2] += int32(coeff) * int32(row[xi+2])
					sum += int32(coeff)
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*3
			out.Pix[xo+0] = clampUint8(p[0] / sum)
			out.Pix[xo+1] = clampUint8(p[1] / sum)
			out.Pix[xo+2] = clampUint8(p[2] / sum)
		}
	}
}

func nearestYCbCr(in *ycc, out *ycc, scale float64, coeffs []bool, offset []int, filterLength int) {
	oldBounds := image.Rect(0, 0, in.Rect.Dx(), in.Rect.Dy())
	newBounds := out.Bounds()
	minX := oldBounds.Min.X * 3
	maxX := (oldBounds.Max.X - oldBounds.Min.X - 1) * 3

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[(x-oldBounds.Min.Y)*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var p [3]float32
			var sum float32
			start := offset[y]
			ci := (y - newBounds.Min.Y) * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case uint(xi) < uint(oldBounds.Max.X):
						xi *= 3
					case xi >= oldBounds.Max.X:
						xi = maxX
					default:
						xi = minX
					}
					p[0] += float32(row[xi+0])
					p[1] += float32(row[xi+1])
					p[2] += float32(row[xi+2])
					sum++
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*3
			out.Pix[xo+0] = floatToUint8(p[0] / sum)
			out.Pix[xo+1] = floatToUint8(p[1] / sum)
			out.Pix[xo+2] = floatToUint8(p[2] / sum)
		}
	}
}
