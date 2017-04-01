package sysfs

import (
	"strconv"
)

const (
	// IN gpio direction
	IN = "in"
	// OUT gpio direction
	OUT = "out"
	// HIGH gpio level
	HIGH = 1
	// LOW gpio level
	LOW = 0
)

// DigitalPin is the interface for sysfs gpio interactions
type DigitalPin interface {
	// Unexport unexports the pin and releases the pin from the operating system
	Unexport() error
	// Export exports the pin for use by the operating system
	Export() error
	// Read reads the current value of the pin
	Read() (int, error)
	// Direction sets the direction for the pin
	Direction(string) error
	// Write writes to the pin
	Write(int) error
}

type digitalPin struct {
	pin   string
	label string

	// JER: For emulation, changed from original type "File" to types shown
	value     int
	direction string
}

// NewDigitalPin returns a DigitalPin given the pin number and an optional sysfs pin label.
// If no label is supplied the default label will prepend "gpio" to the pin number,
// eg. a pin number of 10 will have a label of "gpio10"
func NewDigitalPin(pin int, v ...string) DigitalPin {
	d := &digitalPin{pin: strconv.Itoa(pin)}
	if len(v) > 0 {
		d.label = v[0]
	} else {
		d.label = "gpio" + d.pin
	}

	return d
}

// JER: Maintained DigitalPin interface with Direction(), and simplified for emulation
func (d *digitalPin) Direction(dir string) error {
	d.direction = dir
	return nil
}

// JER: Maintained DigitalPin interface with Write() and simplified for emulation
// TODO: Use channel to signal changes for "connection" to pin
func (d *digitalPin) Write(b int) error {
	d.value = b
	return nil
}

// JER: Maintained DigitalPin interface with Read() and simplified for emulation
func (d *digitalPin) Read() (n int, err error) {
	return d.value, nil
}

// JER: Export not needed. Removed contents of function and return
// <nil> to preserve signature and DigitalPin interface
func (d *digitalPin) Export() error {
	return nil
}

// JER: Unexport not needed. Removed contents of function and return
// <nil> to preserve signature and DigitalPin interface
func (d *digitalPin) Unexport() error {
	return nil
}

// JER: removed writeFile and readFile as not needed for emulation
