package intelite

import (
	"time"
)

const (
	pref_pulse     = 3450
	pref_space     = 1650
	pulse          = 450
	zero_space     = 450
	one_space      = 1100
	packet_size    = 14
	pattern_size   = 2 + 14*8*2 + 2
	b_header       = 0b01110101
	min_brightness = 1
	max_brightness = 10
	min_color      = 1
	max_color      = 10
)

type Mode byte

const (
	ModeOff    Mode = 0b00000000
	ModeNormal Mode = 0b00000001
	ModeNight  Mode = 0b00000010
)

func DefaultSignal() *signal {
	return &signal{
		mode_flags: ModeNormal,
		color:      5,
		brightness: 6,
		night_mode: 1,
	}
}

func NewSignal() *signal {
	return &signal{}
}

type signal struct {
	// b1
	imitation  bool
	off_flag   bool
	on_flag    bool
	sleep_mode byte
	mode_flags Mode
	// b5-b6
	imitation_hour   int
	imitation_minute int
	// b7-b10
	on_hour    int
	on_minute  int
	off_hour   int
	off_minute int
	// b11
	color      int
	brightness int
	// b12
	eco_mode   bool
	max_mode   bool
	night_mode int
}

func (s *signal) Mode(mode Mode) *signal {
	s.mode_flags = mode
	s.max_mode = false
	return s
}

func (s *signal) NightBrightness(v int) *signal {
	if v < 0 {
		s.night_mode = 0
	} else if v > 3 {
		s.night_mode = 3
	} else {
		s.night_mode = v
	}
	return s
}

func (s *signal) SleepMode(v int) *signal {
	if v < 0 {
		s.sleep_mode = 0
	} else if v > 3 {
		s.sleep_mode = 3
	} else {
		s.sleep_mode = byte(v)
	}
	return s
}

func (s *signal) AutoOnEnabled(v bool) *signal {
	s.on_flag = v
	return s
}

func (s *signal) AutoOnTime(t time.Time) *signal {
	s.on_hour = t.Hour()
	s.on_minute = t.Minute()
	return s
}

func (s *signal) AutoOffEnabled(v bool) *signal {
	s.off_flag = v
	return s
}

func (s *signal) AutoOffTime(t time.Time) *signal {
	s.off_hour = t.Hour()
	s.off_minute = t.Minute()
	return s
}

func (s *signal) ImitationEnabled(v bool) *signal {
	s.imitation = v
	return s
}

func (s *signal) ImitationTime(t time.Time) *signal {
	s.imitation_hour = t.Hour()
	s.imitation_minute = t.Minute()
	return s
}

func (s *signal) Brightness(v int) *signal {
	if v < min_brightness {
		s.brightness = min_brightness
	} else if v > max_brightness {
		s.brightness = max_brightness
	} else {
		s.brightness = v
	}
	return s
}

func (s *signal) Color(v int) *signal {
	if v < min_color {
		s.color = min_color
	} else if v > max_color {
		s.color = max_color
	} else {
		s.color = v
	}
	return s
}

func (s *signal) EcoModeEnabled(v bool) *signal {
	s.eco_mode = v
	return s

}

func (s *signal) MaxMode(v bool) *signal {
	s.max_mode = v
	return s
}

func (s *signal) GenPattern() []int {

	packet := make([]int, packet_size)

	// Header byte
	packet[0] = b_header

	// Flags
	packet[1] = 0
	if s.imitation {
		packet[1] |= 0b10000000
	}
	if s.off_flag {
		packet[1] |= 0b01000000
	}
	if s.on_flag {
		packet[1] |= 0b00100000
	}
	packet[1] |= int(s.sleep_mode&0b00000011) << 3
	packet[1] |= int(s.mode_flags)

	// Current time
	t := time.Now()
	packet[2] = t.Hour()
	packet[3] = t.Minute()
	packet[4] = t.Second()

	// Home imitation time
	packet[5] = s.imitation_hour
	packet[6] = s.imitation_minute

	// Auto on
	packet[7] = s.on_hour
	packet[8] = s.on_minute

	// Auto off
	packet[9] = s.off_hour
	packet[10] = s.off_minute

	// Color/brightness
	if s.max_mode {
		packet[11] = 0b01011010
	} else {
		packet[11] = ((s.color & 0b00001111) << 4) | (s.brightness & 0b00001111)
	}

	// Eco/max/night mode
	packet[12] = s.night_mode & 0b00001111
	if s.eco_mode {
		packet[12] |= 0b10000000
	}
	if s.max_mode {
		packet[12] |= 0b01000000
	}

	// Calc checksum
	checksum := 0
	for i := 0; i < packet_size-1; i++ {
		checksum = (checksum + packet[i]) & 0xff
	}
	packet[13] = (checksum + 0x55) & 0xff

	// Generate IR pattern
	pattern := make([]int, pattern_size)
	p := 0
	pattern[p] = pref_pulse
	p++
	pattern[p] = pref_space
	p++
	for i := 0; i < packet_size; i++ {
		p = fillByte(pattern, p, packet[i])
	}
	pattern[p] = pulse
	p++
	pattern[p] = zero_space
	return pattern
}

func fillByte(pattern []int, pos int, b int) int {
	for i := 1; i <= 128; i *= 2 {
		pattern[pos] = pulse
		pos++
		if (b & i) != 0 {
			pattern[pos] = one_space
		} else {
			pattern[pos] = zero_space
		}
		pos++
	}
	return pos
}
