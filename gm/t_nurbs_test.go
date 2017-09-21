// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func checkDerivs(tst *testing.T, b *Nurbs, npts int, tol float64, verb bool) {
	δ := 1e-1 // used to avoid using central differences @ boundaries of t in [0,5]
	dana := make([]float64, 2)
	uu := make([]float64, 2)
	for _, u := range utl.LinSpace(b.b[0].tmin+δ, b.b[0].tmax-δ, npts) {
		for _, v := range utl.LinSpace(b.b[1].tmin+δ, b.b[1].tmax-δ, npts) {
			uu[0], uu[1] = u, v
			b.CalcBasisAndDerivs(uu)
			for i := 0; i < b.n[0]; i++ {
				for j := 0; j < b.n[1]; j++ {
					l := i + j*b.n[0]
					b.GetDerivL(dana, l)
					chk.DerivScaVec(tst, io.Sf("dR%d(%.3f,%.3f)", l, uu[0], uu[1]), tol, dana, uu, 1e-1, verb, func(x []float64) float64 {
						return b.RecursiveBasis(x, l)
					})
				}
			}
		}
	}
}

func TestNurbs01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Nurbs01. Elements, IndBasis, and ExtractSurfaces")

	// NURBS
	surf := FactoryNurbs.Surf2dExample1()
	elems := surf.Elements()
	nbasis := surf.GetElemNumBasis()
	io.Pforan("nbasis = %v\n", nbasis)
	chk.IntAssert(nbasis, 6) // orders := (2,1) => nbasis = (2+1)*(1+1) = 6

	// check basis and elements
	chk.Ints(tst, "elem[0]", elems[0], []int{2, 3, 1, 2})
	chk.Ints(tst, "elem[1]", elems[1], []int{3, 4, 1, 2})
	chk.Ints(tst, "elem[2]", elems[2], []int{4, 5, 1, 2})
	chk.Ints(tst, "ibasis0", surf.IndBasis(elems[0]), []int{0, 1, 2, 5, 6, 7})
	chk.Ints(tst, "ibasis1", surf.IndBasis(elems[1]), []int{1, 2, 3, 6, 7, 8})
	chk.Ints(tst, "ibasis2", surf.IndBasis(elems[2]), []int{2, 3, 4, 7, 8, 9})
	chk.IntAssert(surf.GetElemNumBasis(), len(surf.IndBasis(elems[0])))

	// check derivatives
	many := false
	if many {
		checkDerivs(tst, surf, 21, 1e-9, chk.Verbose) // many points
	} else {
		checkDerivs(tst, surf, 5, 1e-9, chk.Verbose)
	}

	// check spans
	solE := [][]int{{2, 3, 1, 2}, {3, 4, 1, 2}, {4, 5, 1, 2}}
	solL := [][]int{{0, 1, 2, 5, 6, 7}, {1, 2, 3, 6, 7, 8}, {2, 3, 4, 7, 8, 9}}
	for k, e := range elems {
		L := surf.IndBasis(e)
		io.Pforan("e=%v: L=%v\n", e, L)
		chk.Ints(tst, "span", e, solE[k])
		chk.Ints(tst, "L", L, solL[k])
	}

	// check indices along curve
	io.Pf("\n------------ indices along curve -------------\n")
	chk.Ints(tst, "l0s2a0", surf.IndsAlongCurve(0, 2, 0), []int{0, 1, 2})
	chk.Ints(tst, "l0s3a0", surf.IndsAlongCurve(0, 3, 0), []int{1, 2, 3})
	chk.Ints(tst, "l0s4a0", surf.IndsAlongCurve(0, 4, 0), []int{2, 3, 4})
	chk.Ints(tst, "l0s2a1", surf.IndsAlongCurve(0, 2, 1), []int{5, 6, 7})
	chk.Ints(tst, "l0s3a1", surf.IndsAlongCurve(0, 3, 1), []int{6, 7, 8})
	chk.Ints(tst, "l0s4a1", surf.IndsAlongCurve(0, 4, 1), []int{7, 8, 9})
	chk.Ints(tst, "l1s1a0", surf.IndsAlongCurve(1, 1, 0), []int{0, 5})
	chk.Ints(tst, "l1s1a1", surf.IndsAlongCurve(1, 1, 1), []int{1, 6})
	chk.Ints(tst, "l1s1a2", surf.IndsAlongCurve(1, 1, 2), []int{2, 7})
	chk.Ints(tst, "l1s1a3", surf.IndsAlongCurve(1, 1, 3), []int{3, 8})
	chk.Ints(tst, "l1s1a4", surf.IndsAlongCurve(1, 1, 4), []int{4, 9})

	// extract surfaces and check
	io.Pf("\n------------ extract surfaces -------------\n")
	surfs := surf.ExtractSurfaces()
	chk.Deep4(tst, "surf0: Q", 1e-15, surfs[0].Q, [][][][]float64{
		{{{0, 0, 0, 0.8}}},         // 0
		{{{0, 0.4 * 0.9, 0, 0.9}}}, // 5
	})
	chk.Deep4(tst, "surf1: Q", 1e-15, surfs[1].Q, [][][][]float64{
		{{{1.0 * 1.1, 0.1 * 1.1, 0, 1.1}}}, // 4
		{{{1.0 * 0.5, 0.5 * 0.5, 0, 0.5}}}, // 9
	})
	chk.Deep4(tst, "surf2: Q", 1e-15, surfs[2].Q, [][][][]float64{
		{{{0.00 * 0.80, 0.00 * 0.80, 0, 0.80}}}, // 0
		{{{0.25 * 1.00, 0.15 * 1.00, 0, 1.00}}}, // 1
		{{{0.50 * 0.70, 0.00 * 0.70, 0, 0.70}}}, // 2
		{{{0.75 * 1.20, 0.00 * 1.20, 0, 1.20}}}, // 3
		{{{1.00 * 1.10, 0.10 * 1.10, 0, 1.10}}}, // 4
	})
	chk.Deep4(tst, "surf3: Q", 1e-15, surfs[3].Q, [][][][]float64{
		{{{0.00 * 0.90, 0.40 * 0.90, 0, 0.90}}}, // 5
		{{{0.25 * 0.60, 0.55 * 0.60, 0, 0.60}}}, // 6
		{{{0.50 * 1.50, 0.40 * 1.50, 0, 1.50}}}, // 7
		{{{0.75 * 1.40, 0.40 * 1.40, 0, 1.40}}}, // 8
		{{{1.00 * 0.50, 0.50 * 0.50, 0, 0.50}}}, // 9
	})

	io.Pf("\n------------ elem bry local inds -----------\n")
	elembryinds := surf.ElemBryLocalInds()
	io.Pforan("elembryinds = %v\n", elembryinds)
	chk.IntDeep2(tst, "elembryinds", elembryinds, [][]int{
		{0, 1, 2},
		{2, 5},
		{3, 4, 5},
		{0, 3},
	})

	// refine NURBS
	c := surf.Krefine([][]float64{
		{0.5, 1.5, 2.5},
		{0.5},
	})
	elems = c.Elements()
	chk.IntAssert(c.GetElemNumBasis(), len(c.IndBasis(elems[0])))

	// check refined elements: round 1
	io.Pf("\n------------ refined -------------\n")
	chk.Ints(tst, "elem[0]", elems[0], []int{2, 3, 1, 2})
	chk.Ints(tst, "elem[1]", elems[1], []int{3, 4, 1, 2})
	chk.Ints(tst, "elem[2]", elems[2], []int{4, 5, 1, 2})
	chk.Ints(tst, "elem[3]", elems[3], []int{5, 6, 1, 2})
	chk.Ints(tst, "elem[4]", elems[4], []int{6, 7, 1, 2})
	chk.Ints(tst, "elem[5]", elems[5], []int{7, 8, 1, 2})

	// check refined elements: round 2
	chk.Ints(tst, "elem[ 6]", elems[6], []int{2, 3, 2, 3})
	chk.Ints(tst, "elem[ 7]", elems[7], []int{3, 4, 2, 3})
	chk.Ints(tst, "elem[ 8]", elems[8], []int{4, 5, 2, 3})
	chk.Ints(tst, "elem[ 9]", elems[9], []int{5, 6, 2, 3})
	chk.Ints(tst, "elem[10]", elems[10], []int{6, 7, 2, 3})
	chk.Ints(tst, "elem[11]", elems[11], []int{7, 8, 2, 3})

	// check refined basis: round 1
	chk.Ints(tst, "basis0", c.IndBasis(elems[0]), []int{0, 1, 2, 8, 9, 10})
	chk.Ints(tst, "basis1", c.IndBasis(elems[1]), []int{1, 2, 3, 9, 10, 11})
	chk.Ints(tst, "basis2", c.IndBasis(elems[2]), []int{2, 3, 4, 10, 11, 12})
	chk.Ints(tst, "basis3", c.IndBasis(elems[3]), []int{3, 4, 5, 11, 12, 13})
	chk.Ints(tst, "basis4", c.IndBasis(elems[4]), []int{4, 5, 6, 12, 13, 14})
	chk.Ints(tst, "basis5", c.IndBasis(elems[5]), []int{5, 6, 7, 13, 14, 15})

	// check refined basis: round 2
	chk.Ints(tst, "basis6", c.IndBasis(elems[6]), []int{8, 9, 10, 16, 17, 18})
	chk.Ints(tst, "basis7", c.IndBasis(elems[7]), []int{9, 10, 11, 17, 18, 19})
	chk.Ints(tst, "basis8", c.IndBasis(elems[8]), []int{10, 11, 12, 18, 19, 20})
	chk.Ints(tst, "basis9", c.IndBasis(elems[9]), []int{11, 12, 13, 19, 20, 21})
	chk.Ints(tst, "basis10", c.IndBasis(elems[10]), []int{12, 13, 14, 20, 21, 22})
	chk.Ints(tst, "basis11", c.IndBasis(elems[11]), []int{13, 14, 15, 21, 22, 23})

	io.Pf("\n------------ refined: inds along curve -------------\n")
	chk.Ints(tst, "l0s2a0", c.IndsAlongCurve(0, 2, 0), []int{0, 1, 2})
	chk.Ints(tst, "l0s7a0", c.IndsAlongCurve(0, 7, 0), []int{5, 6, 7})
	chk.Ints(tst, "l0s3a2", c.IndsAlongCurve(0, 3, 2), []int{17, 18, 19})
	chk.Ints(tst, "l0s7a2", c.IndsAlongCurve(0, 7, 2), []int{21, 22, 23})
	chk.Ints(tst, "l1s1a0", c.IndsAlongCurve(1, 1, 0), []int{0, 8})
	chk.Ints(tst, "l1s1a0", c.IndsAlongCurve(1, 2, 0), []int{8, 16})
	chk.Ints(tst, "l1s2a7", c.IndsAlongCurve(1, 1, 7), []int{7, 15})
	chk.Ints(tst, "l1s2a7", c.IndsAlongCurve(1, 2, 7), []int{15, 23})

	// plot
	if chk.Verbose {
		io.Pf("\n------------ plot -------------\n")
		ndim := 2
		plt.Reset(true, nil)
		PlotNurbs("/tmp/gosl/gm", "nurbs01a", surf, ndim, 41, true, true, nil, nil, nil, func() {
			plt.AxisOff()
			plt.Equal()
		})
		plt.Reset(true, nil)
		PlotNurbsBasis2d("/tmp/gosl/gm", "nurbs01b", surf, 0, 7, true, true, nil, nil, func(idx int) {
			plt.AxisOff()
			plt.Equal()
		})
		plt.Reset(true, &plt.A{Prop: 1.0})
		PlotNurbsDerivs2d("/tmp/gosl/gm", "nurbs01c", surf, 0, 7, false, false, nil, nil, func(idx int) {
			plt.AxisOff()
			plt.Equal()
		})
	}
}

func TestNurbs02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Nurbs02. Elements and IndBasis")

	// NURBS
	surf := FactoryNurbs.Surf2dQuarterPlateHole1()
	elems := surf.Elements()
	nbasis := surf.GetElemNumBasis()
	io.Pforan("nbasis = %v\n", nbasis)
	chk.IntAssert(nbasis, 9) // orders := (2,2) => nbasis = (2+1)*(2+1) = 9

	// check basis and elements
	chk.Ints(tst, "elem[0]", elems[0], []int{2, 3, 2, 3})
	chk.Ints(tst, "elem[1]", elems[1], []int{3, 4, 2, 3})
	chk.Ints(tst, "ibasis0", surf.IndBasis(elems[0]), []int{0, 1, 2, 4, 5, 6, 8, 9, 10})
	chk.Ints(tst, "ibasis1", surf.IndBasis(elems[1]), []int{1, 2, 3, 5, 6, 7, 9, 10, 11})
	chk.IntAssert(surf.GetElemNumBasis(), len(surf.IndBasis(elems[0])))

	// check derivatives
	many := false
	if many {
		checkDerivs(tst, surf, 21, 1e-5, chk.Verbose)
	} else {
		checkDerivs(tst, surf, 5, 1e-5, chk.Verbose)
	}

	// refine NURBS
	refined := surf.KrefineN(2, false)
	elems = refined.Elements()
	chk.IntAssert(refined.GetElemNumBasis(), len(refined.IndBasis(elems[0])))

	// check refined elements
	io.Pf("\n------------ refined -------------\n")
	chk.Ints(tst, "elem[0]", elems[0], []int{2, 3, 2, 3})
	chk.Ints(tst, "elem[1]", elems[1], []int{3, 4, 2, 3})
	chk.Ints(tst, "elem[2]", elems[2], []int{4, 5, 2, 3})
	chk.Ints(tst, "elem[3]", elems[3], []int{5, 6, 2, 3})
	chk.Ints(tst, "elem[4]", elems[4], []int{2, 3, 3, 4})
	chk.Ints(tst, "elem[5]", elems[5], []int{3, 4, 3, 4})
	chk.Ints(tst, "elem[6]", elems[6], []int{4, 5, 3, 4})
	chk.Ints(tst, "elem[7]", elems[7], []int{5, 6, 3, 4})

	// check refined basis
	chk.Ints(tst, "ibasis0", refined.IndBasis(elems[0]), []int{0, 1, 2, 6, 7, 8, 12, 13, 14})
	chk.Ints(tst, "ibasis1", refined.IndBasis(elems[1]), []int{1, 2, 3, 7, 8, 9, 13, 14, 15})
	chk.Ints(tst, "ibasis2", refined.IndBasis(elems[2]), []int{2, 3, 4, 8, 9, 10, 14, 15, 16})
	chk.Ints(tst, "ibasis3", refined.IndBasis(elems[3]), []int{3, 4, 5, 9, 10, 11, 15, 16, 17})
	chk.Ints(tst, "ibasis4", refined.IndBasis(elems[4]), []int{6, 7, 8, 12, 13, 14, 18, 19, 20})
	chk.Ints(tst, "ibasis5", refined.IndBasis(elems[5]), []int{7, 8, 9, 13, 14, 15, 19, 20, 21})
	chk.Ints(tst, "ibasis6", refined.IndBasis(elems[6]), []int{8, 9, 10, 14, 15, 16, 20, 21, 22})
	chk.Ints(tst, "ibasis7", refined.IndBasis(elems[7]), []int{9, 10, 11, 15, 16, 17, 21, 22, 23})

	// plot
	if chk.Verbose {
		io.Pf("\n------------ plot -------------\n")
		ndim := 2
		la := 0 + 0*surf.n[0]
		lb := 2 + 1*surf.n[0]
		plt.Reset(true, nil)
		PlotNurbs("/tmp/gosl/gm", "nurbs02a", surf, ndim, 41, true, true, nil, nil, nil, func() {
			plt.AxisOff()
			plt.Equal()
		})
		plt.Reset(true, &plt.A{Prop: 1.5})
		PlotNurbsBasis2d("/tmp/gosl/gm", "nurbs02b", surf, la, lb, false, false, nil, nil, func(idx int) {
			plt.AxisOff()
			plt.Equal()
		})
		plt.Reset(true, &plt.A{Prop: 1.7})
		PlotNurbsDerivs2d("/tmp/gosl/gm", "nurbs02c", surf, la, lb, false, false, nil, nil, func(idx int) {
			plt.AxisOff()
			plt.Equal()
		})
	}
}

func TestNurbs03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Nurbs03. Elements and Krefine")

	// NURBS
	surf := FactoryNurbs.Curve2dExample1()
	elems := surf.Elements()
	nbasis := surf.GetElemNumBasis()
	io.Pforan("nbasis = %v\n", nbasis)
	chk.IntAssert(nbasis, 4) // orders := (3,) => nbasis = (3+1) = 4

	// check basis and elements
	chk.Ints(tst, "elem[0]", elems[0], []int{3, 4})
	chk.Ints(tst, "elem[1]", elems[1], []int{4, 5})
	chk.Ints(tst, "elem[2]", elems[2], []int{5, 6})
	chk.Ints(tst, "ibasis0", surf.IndBasis(elems[0]), []int{0, 1, 2, 3})
	chk.Ints(tst, "ibasis1", surf.IndBasis(elems[1]), []int{1, 2, 3, 4})
	chk.Ints(tst, "ibasis2", surf.IndBasis(elems[2]), []int{2, 3, 4, 5})

	// refine NURBS
	refined := surf.Krefine([][]float64{
		{0.15, 0.5, 0.85},
	})

	// plot
	if chk.Verbose {

		// geometry
		plt.Reset(true, &plt.A{WidthPt: 450})
		plotTwoNurbs2d("/tmp/gosl/gm", "nurbs03a", surf, refined, "original", "refined", func() {
			plt.AxisOff()
			plt.Equal()
		})

		// basis
		plt.Reset(true, &plt.A{Prop: 1.2})
		PlotNurbsBasis2d("/tmp/gosl/gm", "nurbs03b", surf, 0, 1, false, false, nil, nil, func(idx int) {
			plt.HideBorders(&plt.A{HideR: true, HideT: true})
		})
		plt.Reset(true, &plt.A{Prop: 1.2})
		plt.HideBorders(&plt.A{HideR: true, HideT: true})
		PlotNurbsDerivs2d("/tmp/gosl/gm", "nurbs03c", surf, 0, 1, false, false, nil, nil, func(idx int) {
			plt.HideBorders(&plt.A{HideR: true, HideT: true})
		})
	}
}

func TestNurbs04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Nurbs04. KrefineN and file read-write")

	// NURBS
	a := FactoryNurbs.Surf2dQuarterPlateHole1()
	b := a.KrefineN(2, false)
	c := a.KrefineN(4, false)

	// plot
	if chk.Verbose {
		ndim := 2
		plt.Reset(true, nil)
		PlotNurbs("/tmp/gosl/gm", "nurbs04a", b, ndim, 41, true, true, nil, nil, nil, func() {
			plt.AxisOff()
			plt.Equal()
		})
		plt.Reset(true, nil)
		plt.Reset(true, nil)
		plotTwoNurbs2d("/tmp/gosl/gm", "nurbs04c", a, b, "original", "refined", func() {
			plt.AxisOff()
			plt.Equal()
		})
		plt.Reset(true, nil)
		plotTwoNurbs2d("/tmp/gosl/gm", "nurbs04d", a, c, "original", "refined", func() {
			plt.AxisOff()
			plt.Equal()
		})
	}
}

func TestNurbs05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Nurbs05. PointDeriv. Curve")

	// NURBS
	verts := [][]float64{
		{0, 0, 0, 1},
		{1, 3, 0, 1},
		{4, 3, 0, 1},
		{5, 0, 0, 1},
		{6, 2, 0, 1},
	}
	knots := [][]float64{
		{0, 0, 0, 2.0 / 5.0, 3.0 / 5.0, 1, 1, 1},
	}
	curve := NewNurbs(1, []int{2}, knots)
	curve.SetControl(verts, utl.IntRange(len(verts)))

	// derivatives
	dCduA := la.NewMatrix(2, curve.Gnd())
	dCduB := la.NewMatrix(2, curve.Gnd())
	dCduC := la.NewMatrix(2, curve.Gnd())
	dCduD := la.NewMatrix(2, curve.Gnd())
	dCduE := la.NewMatrix(2, curve.Gnd())
	cA := la.NewVector(2)
	cB := la.NewVector(2)
	cC := la.NewVector(2)
	cD := la.NewVector(2)
	cE := la.NewVector(2)

	// check
	curve.PointDeriv(dCduA, cA, []float64{0}, 2)
	chk.Array(tst, "dCdu @ u=0  ", 1e-17, GetVec2d(dCduA, 0, false), []float64{5, 15})

	curve.PointDeriv(dCduB, cB, []float64{2.0 / 5.0}, 2)
	chk.Array(tst, "dCdu @ u=2/5", 1e-17, GetVec2d(dCduB, 0, false), []float64{10, 0})

	curve.PointDeriv(dCduC, cC, []float64{3.0 / 5.0}, 2)
	chk.Array(tst, "dCdu @ u=3/5", 1e-17, GetVec2d(dCduC, 0, false), []float64{10.0 / 3.0, -10})

	curve.PointDeriv(dCduD, cD, []float64{4.0 / 5.0}, 2)
	chk.Array(tst, "dCdu @ u=4/5", 1e-15, GetVec2d(dCduD, 0, true), []float64{1, 0})

	curve.PointDeriv(dCduE, cE, []float64{1}, 2)
	chk.Array(tst, "dCdu @ u=1  ", 1e-17, GetVec2d(dCduE, 0, false), []float64{5, 10})

	// plot
	if chk.Verbose {
		PlotNurbs("/tmp/gosl/gm", "nurbs05", curve, 2, 41, false, true, nil, nil, nil, func() {
			DrawArrow2dM(cA, dCduA, 0, true, 1, nil)
			DrawArrow2dM(cB, dCduB, 0, true, 1, nil)
			DrawArrow2dM(cC, dCduC, 0, true, 1, nil)
			DrawArrow2dM(cD, dCduD, 0, true, 1, nil)
			DrawArrow2dM(cE, dCduE, 0, true, 1, nil)
			curve.PointDeriv(dCduA, cA, []float64{1.0 / 5.0}, 2)
			DrawArrow2dM(cA, dCduA, 0, true, 1, nil)
			plt.Gll("x", "y", nil)
			plt.HideTRborders()
			plt.Equal()
		})
	}
}

func TestNurbs06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Nurbs06. PointDeriv. Curve")

	// NURBS
	verts := [][]float64{
		{1, 0, 0, 1},
		{1, 1, 0, 1},
		{0, 1, 0, 2},
	}
	knots := [][]float64{
		{0, 0, 0, 1, 1, 1},
	}
	curve := NewNurbs(1, []int{2}, knots)
	curve.SetControl(verts, utl.IntRange(len(verts)))

	// derivatives
	dCduA := la.NewMatrix(2, curve.Gnd())
	dCduB := la.NewMatrix(2, curve.Gnd())
	dCduC := la.NewMatrix(2, curve.Gnd())
	dCduD := la.NewMatrix(2, curve.Gnd())
	cA := la.NewVector(2)
	cB := la.NewVector(2)
	cC := la.NewVector(2)
	cD := la.NewVector(2)

	// check
	curve.PointDeriv(dCduA, cA, []float64{0}, 2)
	chk.Array(tst, "dCdu @ u=0  ", 1e-17, GetVec2d(dCduA, 0, false), []float64{0, 2})

	curve.PointDeriv(dCduB, cB, []float64{2.0 / 5.0}, 2)

	curve.PointDeriv(dCduC, cC, []float64{3.0 / 5.0}, 2)

	curve.PointDeriv(dCduD, cD, []float64{1}, 2)
	chk.Array(tst, "dCdu @ u=1  ", 1e-17, GetVec2d(dCduD, 0, false), []float64{-1, 0})

	// plot
	if chk.Verbose {
		PlotNurbs("/tmp/gosl/gm", "nurbs06", curve, 2, 41, false, true, nil, nil, nil, func() {
			DrawArrow2dM(cA, dCduA, 0, false, 1, nil)
			DrawArrow2dM(cB, dCduB, 0, false, 1, nil)
			DrawArrow2dM(cC, dCduC, 0, false, 1, nil)
			DrawArrow2dM(cD, dCduD, 0, false, 1, nil)
			plt.Gll("x", "y", nil)
			plt.HideTRborders()
			plt.Equal()
		})
	}
}

func TestNurbs07(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Nurbs07. PointDeriv. Surface")

	// surface
	xc, yc, zc, r, R := 0.0, 0.0, 0.0, 2.0, 4.0
	surf := FactoryNurbs.Surf3dTorus(xc, yc, zc, r, R)

	cA := la.NewVector(3)
	dCduA := la.NewMatrix(3, surf.Gnd())
	surf.PointDeriv(dCduA, cA, []float64{1, 0}, 3)

	cB := la.NewVector(3)
	dCduB := la.NewMatrix(3, surf.Gnd())
	surf.PointDeriv(dCduB, cB, []float64{2, 0}, 3)

	cC := la.NewVector(3)
	dCduC := la.NewMatrix(3, surf.Gnd())
	surf.PointDeriv(dCduC, cC, []float64{1, 1}, 3)

	cD := la.NewVector(3)
	dCduD := la.NewMatrix(3, surf.Gnd())
	surf.PointDeriv(dCduD, cD, []float64{2, 1}, 3)

	chk.Array(tst, "cA", 1e-17, cA, []float64{6, 0, 0})
	chk.Array(tst, "cB", 1e-17, cB, []float64{4, 0, 2})
	chk.Array(tst, "cC", 1e-17, cC, []float64{0, 6, 0})
	chk.Array(tst, "cD", 1e-17, cD, []float64{0, 4, 2})
	chk.Array(tst, "dCduA_u", 1e-15, GetVec3d(dCduA, 0, true), []float64{0, 0, 1})
	chk.Array(tst, "dCduA_v", 1e-15, GetVec3d(dCduA, 1, true), []float64{0, 1, 0})
	chk.Array(tst, "dCduB_u", 1e-15, GetVec3d(dCduB, 0, true), []float64{-1, 0, 0})
	chk.Array(tst, "dCduB_v", 1e-15, GetVec3d(dCduB, 1, true), []float64{0, 1, 0})
	chk.Array(tst, "dCduC_u", 1e-15, GetVec3d(dCduC, 0, true), []float64{0, 0, 1})
	chk.Array(tst, "dCduC_v", 1e-15, GetVec3d(dCduC, 1, true), []float64{-1, 0, 0})
	chk.Array(tst, "dCduD_u", 1e-15, GetVec3d(dCduD, 0, true), []float64{0, -1, 0})
	chk.Array(tst, "dCduD_v", 1e-15, GetVec3d(dCduD, 1, true), []float64{-1, 0, 0})

	// plot
	if chk.Verbose {
		nu, nv := 18, 41
		plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
		//surf.DrawCtrl(3, false, &plt.A{C: "grey", Lw: 0.5}, nil)
		surf.DrawSurface(3, nu, nv, true, false, &plt.A{C: plt.C(2, 9), Rstride: 1, Cstride: 1}, &plt.A{C: "#2782c8", Lw: 0.5})
		//surf.DrawSurface(3, nu, nv, false, true, &plt.A{C: plt.C(0, 9), Rstride: 1, Cstride: 1}, &plt.A{C: "#2782c8", Lw: 0.5})
		plt.Sphere(cA, 0.5, 11, 11, &plt.A{C: "r", Surf: true})
		plt.Sphere(cB, 0.5, 11, 11, &plt.A{C: "r", Surf: true})

		sf := 3.0
		DrawArrow3dM(cA, dCduA, 0, true, sf, &plt.A{C: plt.C(0, 0)})
		DrawArrow3dM(cA, dCduA, 1, true, sf, &plt.A{C: plt.C(1, 0)})

		DrawArrow3dM(cB, dCduB, 0, true, sf, &plt.A{C: plt.C(0, 0)})
		DrawArrow3dM(cB, dCduB, 1, true, sf, &plt.A{C: plt.C(1, 0)})

		DrawArrow3dM(cC, dCduC, 0, true, sf, &plt.A{C: plt.C(0, 0)})
		DrawArrow3dM(cC, dCduC, 1, true, sf, &plt.A{C: plt.C(1, 0)})

		DrawArrow3dM(cD, dCduD, 0, true, sf, &plt.A{C: plt.C(0, 0)})
		DrawArrow3dM(cD, dCduD, 1, true, sf, &plt.A{C: plt.C(1, 0)})

		plt.Default3dView(-6.1, 6.1, -6.1, 6.1, -6.1, 6.1, true)
		plt.Save("/tmp/gosl/gm", "nurbs07")
		//plt.ShowSave("/tmp/gosl/gm", "nurbs07")
	}
}

func TestNurbs08(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Nurbs08. Quarter of Hemisphere")

	surf := FactoryNurbs.Surf3dQuarterHemisphere(1, 2, 3, 4)

	cA := la.NewVector(3)
	dCduA := la.NewMatrix(3, surf.Gnd())
	surf.PointDeriv(dCduA, cA, []float64{0, 1}, 3)
	io.Pf("dCduA =\n%v\n", dCduA.Print("%23.15e"))

	cB := la.NewVector(3)
	dCduB := la.NewMatrix(3, surf.Gnd())
	surf.PointDeriv(dCduB, cB, []float64{0.25, 1}, 3)
	io.Pf("dCduB =\n%v\n", dCduB.Print("%23.15e"))

	cC := la.NewVector(3)
	dCduC := la.NewMatrix(3, surf.Gnd())
	surf.PointDeriv(dCduC, cC, []float64{0.5, 1}, 3)
	io.Pf("dCduC =\n%v\n", dCduC.Print("%23.15e"))

	cD := la.NewVector(3)
	dCduD := la.NewMatrix(3, surf.Gnd())
	surf.PointDeriv(dCduD, cD, []float64{1, 1}, 3)
	io.Pf("dCduD =\n%v\n", dCduD.Print("%23.15e"))

	chk.Array(tst, "cA", 1e-15, cA, []float64{1, 2, 7})
	chk.Array(tst, "cB", 1e-15, cB, []float64{1, 2, 7})
	chk.Array(tst, "cC", 1e-15, cC, []float64{1, 2, 7})
	chk.Array(tst, "cD", 1e-15, cD, []float64{1, 2, 7})
	chk.Array(tst, "dCduA_u", 1e-15, GetVec3d(dCduA, 0, false), []float64{0, 0, 0})
	chk.Array(tst, "dCduB_u", 1e-15, GetVec3d(dCduB, 0, false), []float64{0, 0, 0})
	chk.Array(tst, "dCduC_u", 1e-15, GetVec3d(dCduC, 0, false), []float64{0, 0, 0})
	chk.Array(tst, "dCduD_u", 1e-15, GetVec3d(dCduD, 0, false), []float64{0, 0, 0})

	if chk.Verbose {
		nu, nv := 21, 21
		plt.Reset(true, &plt.A{WidthPt: 500, Dpi: 150})
		surf.DrawCtrl(3, false, &plt.A{C: "grey", Lw: 0.5}, nil)
		surf.DrawSurface(3, nu, nv, true, false, &plt.A{C: plt.C(2, 9), Rstride: 2, Cstride: 2}, &plt.A{C: "#2782c8", Lw: 0.5})
		surf.DrawVectors3d(5, 5, 1, nil, nil)

		sf := 1.0
		DrawArrow3dM(cA, dCduA, 0, true, sf, &plt.A{C: plt.C(0, 2), Lw: 4})
		DrawArrow3dM(cA, dCduA, 1, true, sf, &plt.A{C: plt.C(1, 2), Lw: 4})

		plt.Default3dView(-1, 5, -1, 5, -1, 5, true)
		plt.Save("/tmp/gosl/gm", "nurbs07")
		//plt.ShowSave("/tmp/gosl/gm", "nurbs07")
	}
}
