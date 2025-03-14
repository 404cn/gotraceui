package main

import (
	"image/color"

	mycolor "honnef.co/go/gotraceui/color"
	"honnef.co/go/gotraceui/trace/ptrace"
)

var colors [colorLast]color.NRGBA
var colorsOklch [colorLast]mycolor.Oklch

func init() {
	// Our base lightness is 56.51, and our base chroma is 0.122
	const l = 58.51
	const c = 0.122
	const lStep1 = 15
	const lStep2 = 10

	colorsOklch[colorStateActive] = oklch(l, c, 143.74) // Manually chosen
	colorsOklch[colorStateStack] = oklchDelta(colorsOklch[colorStateActive], lStep1, -0.01, 0)
	colorsOklch[colorStateCPUSample] = oklchDelta(colorsOklch[colorStateStack], lStep2, -0.01, 0)

	colorsOklch[colorStateReady] = oklch(l, c, 206.35) // Manually chosen
	colorsOklch[colorStateInactive] = oklch(l, 0, 0)

	colorsOklch[colorStateUserRegion] = oklch(l+lStep1+lStep2, c, 331.18) // Manually chosen

	// Manually chosen. This is the rarest blocked state, so we darken it to have more range for the other states.
	colorsOklch[colorStateBlocked] = oklch(l-5, c, 23.89)
	colorsOklch[colorStateBlockedSyscall] = oklch(l, c, 23.89)
	colorsOklch[colorStateBlockedNet] = oklch(l+6, c-0.01, 23.89)
	colorsOklch[colorStateBlockedHappensBefore] = oklch(l+lStep2, c, 23.89)
	colorsOklch[colorStateBlockedGC] = oklch(l, c, 0) // a blend of colorStateGC and red

	colorsOklch[colorStateGC] = oklch(l, c, 302.36)
	colorsOklch[colorStateSTW] = oklch(l, c+0.072, 23.89) // STW is the most severe form of blocking, hence the increased chroma

	colorsOklch[colorTimelineLabel] = oklch(62.68, 0, 0)
	colorsOklch[colorTimelineBorder] = oklch(89.75, 0, 0)

	// 	// TODO(dh): find a nice color for this
	// We don't use the l constant for thse colors because they're independent from the span colors
	colorsOklch[colorSpanHighlightedPrimaryOutline] = oklch(70.71, 0.322, 328.36)
	colorsOklch[colorSpanHighlightedSecondaryOutline] = oklch(88.44, 0.27, 137.68)

	colorsOklch[colorStateMerged] = oklch(l+lStep1, c, 109.91) // Manually chosen, made brighter so it stands out in gradients

	colorsOklch[colorStateStuck] = oklch(0, 0, 0)
	colorsOklch[colorStateDone] = oklch(0, 0, 0)

	for i, c := range colorsOklch {
		colors[i] = c.NRGBA()
	}

	colors[colorStateUnknown] = rgba(0xFFFF00FF)
	colors[colorStatePlaceholderStackSpan] = rgba(0xe8e8d5ff)
}

type colorIndex uint8

const (
	colorStateUnknown colorIndex = iota

	colorStateInactive
	colorStateActive

	colorStateBlocked
	colorStateBlockedHappensBefore
	colorStateBlockedNet
	colorStateBlockedGC
	colorStateBlockedSyscall
	colorStateGC
	colorStateSTW

	colorStateReady
	colorStateStuck
	colorStateMerged
	colorStateUserRegion
	colorStateStack
	colorStateCPUSample
	colorStateDone
	colorStatePlaceholderStackSpan

	colorStateLast

	colorTimelineLabel
	colorTimelineBorder

	colorSpanHighlightedPrimaryOutline
	colorSpanHighlightedSecondaryOutline

	colorLast
)

var stateColors = [256]colorIndex{
	// per-G states
	ptrace.StateInactive:                colorStateInactive,
	ptrace.StateActive:                  colorStateActive,
	ptrace.StateBlocked:                 colorStateBlocked,
	ptrace.StateBlockedSend:             colorStateBlockedHappensBefore,
	ptrace.StateBlockedRecv:             colorStateBlockedHappensBefore,
	ptrace.StateBlockedSelect:           colorStateBlockedHappensBefore,
	ptrace.StateBlockedSync:             colorStateBlockedHappensBefore,
	ptrace.StateBlockedCond:             colorStateBlockedHappensBefore,
	ptrace.StateBlockedNet:              colorStateBlockedNet,
	ptrace.StateBlockedGC:               colorStateBlockedGC,
	ptrace.StateBlockedSyscall:          colorStateBlockedSyscall,
	ptrace.StateStuck:                   colorStateStuck,
	ptrace.StateReady:                   colorStateReady,
	ptrace.StateCreated:                 colorStateReady,
	ptrace.StateGCMarkAssist:            colorStateGC,
	ptrace.StateGCSweep:                 colorStateGC,
	ptrace.StateGCIdle:                  colorStateGC,
	ptrace.StateGCDedicated:             colorStateGC,
	ptrace.StateGCFractional:            colorStateGC,
	ptrace.StateBlockedSyncOnce:         colorStateBlockedHappensBefore,
	ptrace.StateBlockedSyncTriggeringGC: colorStateGC,
	ptrace.StateUserRegion:              colorStateUserRegion,
	ptrace.StateStack:                   colorStateStack,
	ptrace.StateCPUSample:               colorStateCPUSample,
	ptrace.StateDone:                    colorStateDone,

	// per-P states
	ptrace.StateRunningG: colorStateActive,

	// per-M states
	ptrace.StateRunningP: colorStateActive,
}

func oklch(l, c, h float32) mycolor.Oklch {
	return mycolor.Oklch{L: l / 100, C: c, H: h, Alpha: 1}
}

func oklchDelta(b mycolor.Oklch, l, c, h float32) mycolor.Oklch {
	b.L += l / 100
	b.C += c
	b.H += h
	if b.L < 0 {
		b.L = 0
	}
	if b.L > 1 {
		b.L = 1
	}
	if b.C < 0 {
		b.C = 0
	}
	return b
}

func rgba(c uint32) color.NRGBA {
	// XXX does endianness matter?
	return color.NRGBA{
		A: uint8(c & 0xFF),
		B: uint8(c >> 8 & 0xFF),
		G: uint8(c >> 16 & 0xFF),
		R: uint8(c >> 24 & 0xFF),
	}
}
