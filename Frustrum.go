package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type PlaneEquations struct {
	a float32
	b float32
	c float32
	d float32
}
type FrustrumPlanes struct {
	left   PlaneEquations
	right  PlaneEquations
	top    PlaneEquations
	bottom PlaneEquations
	far    PlaneEquations
	near   PlaneEquations
}

func _getFrustrumSignedDistance(plane PlaneEquations, x, y, z float64) float64 {
	return float64(plane.a)*x + float64(plane.b)*y + float64(plane.c)*z + float64(plane.d)
}

func _inside(plane PlaneEquations, xmin, xmax, zmin, zmax int) bool {
	if _getFrustrumSignedDistance(plane, float64(xmax), 255, float64(zmax)) <= 0 &&
		_getFrustrumSignedDistance(plane, float64(xmax), 255, float64(zmin)) <= 0 &&
		_getFrustrumSignedDistance(plane, float64(xmin), 255, float64(zmin)) <= 0 &&
		_getFrustrumSignedDistance(plane, float64(xmin), 255, float64(zmax)) <= 0 &&
		_getFrustrumSignedDistance(plane, float64(xmax), 0, float64(zmax)) <= 0 &&
		_getFrustrumSignedDistance(plane, float64(xmax), 0, float64(zmin)) <= 0 &&
		_getFrustrumSignedDistance(plane, float64(xmin), 0, float64(zmin)) <= 0 &&
		_getFrustrumSignedDistance(plane, float64(xmin), 0, float64(zmax)) <= 0 {
		return false
	}
	return true
}

func InsideFrustrum(planes FrustrumPlanes, xmin, xmax, zmin, zmax int) bool {
	if (!_inside(planes.left, xmin, xmax, zmin, zmax) ||
		!_inside(planes.right, xmin, xmax, zmin, zmax)) ||
		(!_inside(planes.far, xmin, xmax, zmin, zmax) ||
			!_inside(planes.near, xmin, xmax, zmin, zmax)) {
		return false
	}
	return true
}

func ExtractViewFrustrumPlanes(projection mgl32.Mat4, view mgl32.Mat4) FrustrumPlanes {
	var frustrum FrustrumPlanes

	projection = projection.Mul4(view)

	frustrum.left.a = projection.At(3, 0) + projection.At(0, 0)
	frustrum.left.b = projection.At(3, 1) + projection.At(0, 1)
	frustrum.left.c = projection.At(3, 2) + projection.At(0, 2)
	frustrum.left.d = projection.At(3, 3) + projection.At(0, 3)

	frustrum.right.a = projection.At(3, 0) - projection.At(0, 0)
	frustrum.right.b = projection.At(3, 1) - projection.At(0, 1)
	frustrum.right.c = projection.At(3, 2) - projection.At(0, 2)
	frustrum.right.d = projection.At(3, 3) - projection.At(0, 3)

	frustrum.top.a = projection.At(3, 0) - projection.At(1, 0)
	frustrum.top.b = projection.At(3, 1) - projection.At(1, 1)
	frustrum.top.c = projection.At(3, 2) - projection.At(1, 2)
	frustrum.top.d = projection.At(3, 3) - projection.At(1, 3)

	frustrum.bottom.a = projection.At(3, 0) + projection.At(1, 0)
	frustrum.bottom.b = projection.At(3, 1) + projection.At(1, 1)
	frustrum.bottom.c = projection.At(3, 2) + projection.At(1, 2)
	frustrum.bottom.d = projection.At(3, 3) + projection.At(1, 3)

	frustrum.far.a = projection.At(3, 0) - projection.At(2, 0)
	frustrum.far.b = projection.At(3, 1) - projection.At(2, 1)
	frustrum.far.c = projection.At(3, 2) - projection.At(2, 2)
	frustrum.far.d = projection.At(3, 3) - projection.At(2, 3)

	frustrum.near.a = projection.At(3, 0) + projection.At(2, 0)
	frustrum.near.b = projection.At(3, 1) + projection.At(2, 1)
	frustrum.near.c = projection.At(3, 2) + projection.At(2, 2)
	frustrum.near.d = projection.At(3, 3) + projection.At(2, 3)

	return frustrum
}
