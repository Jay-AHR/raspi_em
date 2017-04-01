package raspi

import (
	"errors"
	"fmt"
	"os"
	//"github.com/hybridgroup/gobot"
	//"github.com/hybridgroup/gobot/sysfs"
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

// NewRaspiAdaptor creates a RaspiAdaptor with specified name and
// DigitalPin mapping (no PWM pins for RPi)
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
		// JER: check if pins present on all board revisions are valid
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

// I2cStart starts a i2c device in specified address
func (r *RaspiAdaptor) I2cStart(address int) (err error) {
	if r.i2cDevice == nil {
		r.i2cDevice, err = sysfs.NewI2cDevice(r.i2cLocation, address)
	}
	return err
}

// I2CWrite writes data to i2c device
func (r *RaspiAdaptor) I2cWrite(address int, data []byte) (err error) {
	if err = r.i2cDevice.SetAddress(address); err != nil {
		return
	}
	_, err = r.i2cDevice.Write(data)
	return
}

// I2cRead returns value from i2c device using specified size
func (r *RaspiAdaptor) I2cRead(address int, size int) (data []byte, err error) {
	if err = r.i2cDevice.SetAddress(address); err != nil {
		return
	}
	data = make([]byte, size)
	_, err = r.i2cDevice.Read(data)
	return
}

func (r *RaspiAdaptor) PwmWrite(pin string, val byte) (err error) {
	sysfsPin, err := r.pwmPin(pin)
	if err != nil {
		return err
	}
	return r.piBlaster(fmt.Sprintf("%v=%v\n", sysfsPin, gobot.FromScale(float64(val), 0, 255)))
}

func (r *RaspiAdaptor) ServoWrite(pin string, angle byte) (err error) {
	sysfsPin, err := r.pwmPin(pin)
	if err != nil {
		return err
	}

	val := (gobot.ToScale(gobot.FromScale(float64(angle), 0, 180), 0, 200) / 1000.0) + 0.05

	return r.piBlaster(fmt.Sprintf("%v=%v\n", sysfsPin, val))
}

func (r *RaspiAdaptor) piBlaster(data string) (err error) {
	fi, err := sysfs.OpenFile("/dev/pi-blaster", os.O_WRONLY|os.O_APPEND, 0644)
	defer fi.Close()

	if err != nil {
		return err
	}

	_, err = fi.WriteString(data)
	return
}
