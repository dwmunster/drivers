// Package lps22hb implements a driver for LPS22HB, a MEMS nano pressure sensor.
//
// Datasheet: https://www.st.com/resource/en/datasheet/dm00140895.pdf
//
package lps22hb

import (
	"tinygo.org/x/drivers"
)

// Device wraps an I2C connection to a HTS221 device.
type Device struct {
	bus     drivers.I2C
	Address uint8
}

// Connected returns whether LPS22HB has been found.
// It does a "who am I" request and checks the response.
func (d *Device) Connected() bool {
	data := []byte{0}
	d.bus.ReadRegister(d.Address, LPS22HB_WHO_AM_I_REG, data)
	return data[0] == 0xB1
}

// Configure sets up the LPS22HB device for communication.
func (d *Device) Configure() {
	// set to block update mode
	d.bus.WriteRegister(d.Address, LPS22HB_CTRL1_REG, []byte{0x02})
}

// ReadPressure returns the pressure in milli pascals (mPa).
func (d *Device) ReadPressure() (pressure int32, err error) {
	d.waitForOneShot()

	// read data
	data := []byte{0, 0, 0}
	d.bus.ReadRegister(d.Address, LPS22HB_PRESS_OUT_REG, data[:1])
	d.bus.ReadRegister(d.Address, LPS22HB_PRESS_OUT_REG+1, data[1:2])
	d.bus.ReadRegister(d.Address, LPS22HB_PRESS_OUT_REG+2, data[2:])
	pValue := float32(uint32(data[2])<<16|uint32(data[1])<<8|uint32(data[0])) / 4096.0

	return int32(pValue * 1000), nil
}

// ReadTemperature returns the temperature in celsius milli degrees (°C/1000).
func (d *Device) ReadTemperature() (temperature int32, err error) {
	d.waitForOneShot()

	// read data
	data := []byte{0, 0}
	d.bus.ReadRegister(d.Address, LPS22HB_TEMP_OUT_REG, data[:1])
	d.bus.ReadRegister(d.Address, LPS22HB_TEMP_OUT_REG+1, data[1:])
	tValue := float32(int16(uint16(data[1])<<8|uint16(data[0]))) / 100.0

	return int32(tValue * 1000), nil
}

// private functions

// wait and trigger one shot in block update
func (d *Device) waitForOneShot() {
	// trigger one shot
	d.bus.WriteRegister(d.Address, LPS22HB_CTRL2_REG, []byte{0x01})

	// wait until one shot is cleared
	data := []byte{1}
	for {
		d.bus.ReadRegister(d.Address, LPS22HB_CTRL2_REG, data)
		if data[0]&0x01 == 0 {
			break
		}
	}
}
