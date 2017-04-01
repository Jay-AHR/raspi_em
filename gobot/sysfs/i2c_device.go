package sysfs

type i2cDevice struct{}

// NewI2cDevice returns an io.ReadWriteCloser with the proper ioctrl given
// an i2c bus location.
func NewI2cDevice(location string) (d *i2cDevice, err error) {
	d = &i2cDevice{}
	return
}

// JER: Read() is only method called from QPID
// returns a token integer for now
// TODO: Add implementation features
func (d *i2cDevice) Read() (n int, err error) {
	return 120, nil
}
