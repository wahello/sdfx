//-----------------------------------------------------------------------------
/*

Gyroid Cubes
Gyroid Teapot

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"
	"os"

	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

func gyroidCube() (sdf.SDF3, error) {

	l := 100.0   // cube side
	k := l * 0.1 // 10 cycles per side

	gyroid, err := sdf.Gyroid3D(v3.Vec{k, k, k})
	if err != nil {
		return nil, err
	}

	box, err := sdf.Box3D(v3.Vec{l, l, l}, 0)
	if err != nil {
		return nil, err
	}

	return sdf.Intersect3D(box, gyroid), nil
}

//-----------------------------------------------------------------------------

func gyroidSurface() (sdf.SDF3, error) {

	l := 60.0    // cube side
	k := l * 0.5 // 2 cycles per side

	s, err := sdf.Gyroid3D(v3.Vec{k, k, k})
	if err != nil {
		return nil, err
	}

	s, err = sdf.Shell3D(s, k*0.025)
	if err != nil {
		return nil, err
	}

	box, err := sdf.Box3D(v3.Vec{l, l, l}, 0)
	if err != nil {
		return nil, err
	}

	s = sdf.Intersect3D(box, s)

	// remove the isolated bits on the cube corners
	sphere, err := sdf.Sphere3D(k * 0.15)
	if err != nil {
		return nil, err
	}
	d := l * 0.5
	s0 := sdf.Transform3D(sphere, sdf.Translate3d(v3.Vec{d, d, d}))
	s1 := sdf.Transform3D(sphere, sdf.Translate3d(v3.Vec{-d, -d, -d}))

	return sdf.Difference3D(s, sdf.Union3D(s0, s1)), nil
}

//-----------------------------------------------------------------------------

func gyroidTeapot() (sdf.SDF3, error) {

	stl := "../../files/teapot.stl"

	// read the stl file.
	file, err := os.OpenFile(stl, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}

	// create the SDF from the STL mesh
	teapot, err := obj.ImportSTL(file, 20, 3, 5)
	if err != nil {
		return nil, err
	}

	min := teapot.BoundingBox().Min
	max := teapot.BoundingBox().Max
	dimX := max.X - min.X
	dimY := max.Y - min.Y
	dimZ := max.Z - min.Z
	dimMin := math.Min(dimX, math.Min(dimY, dimZ)) // teapot shortest side

	k := dimMin * 0.1 // 10 cycles per shortest side

	gyroid, err := sdf.Gyroid3D(v3.Vec{X: k, Y: k, Z: k})
	if err != nil {
		return nil, err
	}

	return sdf.Intersect3D(teapot, gyroid), nil
}

//-----------------------------------------------------------------------------

func main() {

	s0, err := gyroidCube()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s0, "gyroid_cube.stl", render.NewMarchingCubesUniform(300))

	s1, err := gyroidSurface()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s1, "gyroid_surface.stl", render.NewMarchingCubesUniform(150))

	s2, err := gyroidTeapot()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s2, "gyroid_teapot.stl", render.NewMarchingCubesUniform(50))

}

//-----------------------------------------------------------------------------
