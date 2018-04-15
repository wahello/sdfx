//-----------------------------------------------------------------------------
/*

Axoloti Board Mounting Kit

*/
//-----------------------------------------------------------------------------

package main

import . "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

var front_panel_thickness = 3.0
var front_panel_length = 170.0
var front_panel_height = 50.0
var front_panel_y_offset = 15.0

var base_width = 50.0
var base_length = 170.0
var base_thickness = 3.0

var base_foot_width = 10.0
var base_foot_corner_radius = 3.0

var pcb_width = 50.0
var pcb_length = 160.0

var pillar_height = 12.0

//-----------------------------------------------------------------------------

// multiple standoffs
func standoffs() SDF3 {

	k := &StandoffParms{
		PillarHeight:   pillar_height,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.0,
	}

	z_ofs := 0.5 * (pillar_height + base_thickness)

	// from the board mechanicals
	positions := V3Set{
		{3.5, 10.0, z_ofs},   // H1
		{3.5, 40.0, z_ofs},   // H2
		{54.0, 40.0, z_ofs},  // H3
		{156.5, 10.0, z_ofs}, // H4
		//{54.0, 10.0, z_ofs},  // H5
		{156.5, 40.0, z_ofs}, // H6
		{44.0, 10.0, z_ofs},  // H7
		{116.0, 10.0, z_ofs}, // H8
	}

	return Standoffs3D(k, positions)
}

//-----------------------------------------------------------------------------

func base() SDF3 {
	// base
	pp := &PanelParms{
		Size:         V2{base_length, base_width},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{7.0, 20.0, 7.0, 20.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	s0 := Panel2D(pp)

	// cutout
	l := base_length - (2.0 * base_foot_width)
	w := 18.0
	s1 := Box2D(V2{l, w}, base_foot_corner_radius)
	y_ofs := 0.5 * (base_width - pcb_width)
	s1 = Transform2D(s1, Translate2d(V2{0, y_ofs}))

	s2 := Extrude3D(Difference2D(s0, s1), base_thickness)
	x_ofs := 0.5 * pcb_length
	y_ofs = pcb_width - (0.5 * base_width)
	s2 = Transform3D(s2, Translate3d(V3{x_ofs, y_ofs, 0}))

	// standoffs
	s3 := standoffs()

	s4 := Union3D(s2, s3)
	s4.(*UnionSDF3).SetMin(PolyMin(3.0))

	return s4
}

//-----------------------------------------------------------------------------

type PanelHole struct {
	center V2   // center of hole
	hole   SDF2 // 2d hole
}

func front_panel() SDF3 {

	s_midi := Circle2D(0.5 * 15.5)
	s_jack := Circle2D(0.5 * 11.5)
	s_led := Box2D(V2{1.6, 1.6}, 0)

	fb := &FingerButtonParms{
		Width:  3.5,
		Gap:    0.5,
		Length: 20.0,
	}
	s_button := Transform2D(FingerButton2D(fb), Rotate2d(DtoR(-90)))

	jack_x := 123.0
	midi_x := 18.2
	led_x := 62.7
	pb_x := 52.8

	holes := []PanelHole{
		{V2{midi_x, 9.3}, s_midi},                // MIDI DIN Jack
		{V2{midi_x + 20.32, 9.3}, s_midi},        // MIDI DIN Jack
		{V2{jack_x, 8.14}, s_jack},               // 1/4" Stereo Jack
		{V2{jack_x + 19.5, 8.14}, s_jack},        // 1/4" Stereo Jack
		{V2{107.4, 2.3}, Circle2D(0.5 * 5.5)},    // 3.5 mm Headphone Jack
		{V2{led_x, 0.5}, s_led},                  // LED
		{V2{led_x + 3.635, 0.5}, s_led},          // LED
		{V2{pb_x, 0.8}, s_button},                // Push Button
		{V2{pb_x + 5.334, 0.8}, s_button},        // Push Button
		{V2{84.1, 1.3}, Box2D(V2{16.0, 8.0}, 0)}, // micro SD card
		{V2{96.7, 1.3}, Box2D(V2{11.0, 8.0}, 0)}, // micro USB connector
		{V2{73.1, 7.1}, Box2D(V2{7.5, 15.0}, 0)}, // fullsize USB connector
	}

	s := make([]SDF2, len(holes))
	for i, k := range holes {
		s[i] = Transform2D(k.hole, Translate2d(k.center))
	}
	cutouts := Union2D(s...)

	// overall panel
	pp := &PanelParms{
		Size:         V2{front_panel_length, front_panel_height},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	panel := Panel2D(pp)

	x_ofs := 0.5 * pcb_length
	y_ofs := (0.5 * front_panel_height) - front_panel_y_offset
	panel = Transform2D(panel, Translate2d(V2{x_ofs, y_ofs}))

	return Extrude3D(Difference2D(panel, cutouts), front_panel_thickness)
}

//-----------------------------------------------------------------------------

func main() {
	s0 := front_panel()
	s1 := base()
	RenderSTL(s0, 400, "fp.stl")
	RenderSTL(s1, 400, "base.stl")
	s0 = Transform3D(s0, Translate3d(V3{0, 80, 0}))
	RenderSTL(Union3D(s0, s1), 400, "both.stl")
}

//-----------------------------------------------------------------------------