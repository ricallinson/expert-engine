//
// Copyright 2017, Yahoo Inc.
// Copyrights licensed under the New BSD License.
// See the accompanying LICENSE file for terms.
//

package engine

type ToggleSwitch struct {
	*PinInput
}

// Returns a new instance of ToggleSwitch.
// The value of `pin` must be in the range of 1-25 mapping to the Raspberry Pi GPIO pins.
func NewToggleSwitch(pin int) *ToggleSwitch {
	this := &ToggleSwitch{
		NewPinInput(pin),
	}
	return this
}
