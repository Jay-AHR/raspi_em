package raspi

import (
	"errors"
	"fmt"

	"github.com/Jay-AHR/raspi_em/gobot/sysfs"
)

// Removed readFile() as it was used to set i2c revision
// i2c revision will be set statically based on emulation

type RaspiAdaptor struct {
	name        string
	revision    string
	i2cLocation string
	digitalPins map[int]sysfs.DigitalPin
	pwmPins     []int
	i2cDevice   sysfs.I2cDevice
}

var pins = map[string]map[string]int{
	"3": map[string]int{
		"1": 0,
		"2": 2,
		"3": 2,
	},
	"5": map[string]int{
		"1": 1,
		"2": 3,
		"3": 3,
	},
	"7": map[string]int{
		"*": 4,
	},
	"8": map[string]int{
		"*": 14,
	},
	"10": map[string]int{
		"*": 15,
	},
	"11": map[string]int{
		"*": 17,
	},
	"12": map[string]int{
		"*": 18,
	},
	"13": map[string]int{
		"1": 21,
		"2": 27,
		"3": 27,
	},
	"15": map[string]int{
		"*": 22,
	},
	"16": map[string]int{
		"*": 23,
	},
	"18": map[string]int{
		"*": 24,
	},
	"19": map[string]int{
		"*": 10,
	},
	"21": map[string]int{
		"*": 9,
	},
	"22": map[string]int{
		"*": 25,
	},
	"23": map[string]int{
		"*": 11,
	},
	"24": map[string]int{
		"*": 8,
	},
	"26": map[string]int{
		"*": 7,
	},
	"29": map[string]int{
		"3": 5,
	},
	"31": map[string]int{
		"3": 6,
	},
	"32": map[string]int{
		"3": 12,
	},
	"33": map[string]int{
		"3": 13,
	},
	"35": map[string]int{
		"3": 19,
	},
	"36": map[string]int{
		"3": 16,
	},
	"37": map[string]int{
		"3": 26,
	},
	"38": map[string]int{
		"3": 20,
	},
	"40": map[string]int{
		"3": 21,
	},
}

func NewRaspiAdaptor(name string) *RaspiAdaptor {
	r := &RaspiAdaptor{
		name:        name,
		digitalPins: make(map[int]sysfs.DigitalPin),
		pwmPins:     []int{},
	}

	r.revision = "3"

	// Deleted code that called readFile() to obtain revision

	return r
}

func (r *RaspiAdaptor) Name() string { return r.name }

// Connect starts connection with board and creates
// digitalPins and pwmPins adaptor maps
// JER: This function appears to not be needed for RPi as for other boards
func (r *RaspiAdaptor) Connect() (errs []error) {
	return
}

// Finalize closes connection to board and pins
// JER: Removed contents
func (r *RaspiAdaptor) Finalize() (errs []error) {
	return
}

func (r *RaspiAdaptor) translatePin(pin string) (i int, err error) {
	// JER: check if per board-revision pins are valid
	if val, ok := pins[pin][r.revision]; ok {
		i = val
		// JER: check if pins that are present on all board revisions are valid
	} else if val, ok := pins[pin]["*"]; ok {
		i = val
	} else {
		err = errors.New("Not a valid pin")
		return
	}
	return
}

func (r *RaspiAdaptor) pwmPin(pin string) (i int, err error) {
	//  JER: check if a valid pin via the "pin" mapping
	i, err = r.translatePin(pin)
	if err != nil {
		return
	}

	newPin := true
	for _, pin := range r.pwmPins {
		if i == pin {
			newPin = false
			return
		}
	}

	if newPin {
		r.pwmPins = append(r.pwmPins, i)
	}

	return
}

// digitalPin returns matched digitalPin for specified values
func (r *RaspiAdaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
	//  JER: check if a valid pin via the "pin" mapping
	i, err := r.translatePin(pin)

	if err != nil {
		return
	}

	if r.digitalPins[i] == nil {
		r.digitalPins[i] = sysfs.NewDigitalPin(i)
		if err = r.digitalPins[i].Export(); err != nil {
			return
		}
	}

	if err = r.digitalPins[i].Direction(dir); err != nil {
		return
	}

	return r.digitalPins[i], nil
}

// DigitalRead reads digital value from pin
func (r *RaspiAdaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := r.digitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to specified pin
func (r *RaspiAdaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := r.digitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
}

// JER: Removed I2cStart, I2cWrite as QPID does not use these
// TODO: Add I2C methods back in, once QPID prototype is operational

// I2cRead returns value from i2c device using specified size
func (r *RaspiAdaptor) I2cRead(address int, size int) (t int, err error) {
	fmt.Println(address)
	fmt.Println(size)

	t, err = r.i2cDevice.Read()
	return t, err
}
