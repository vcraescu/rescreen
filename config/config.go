package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const layoutCols = 4
const layoutRows = 4

// Config object of config file
type Config struct {
	Layout   layoutConfig
	Monitors MonitorsConfig `json:"monitors"`
}

// MonitorConfig object of config file
type monitorConfig struct {
	Scale   float32 `json:"scale"`
	Primary bool    `json:"primary"`
}

// MonitorsConfig map of config file
type MonitorsConfig map[string]monitorConfig

// LayoutConfig slice of config file
type layoutConfig []string

func (c Config) validate() error {
	err := c.Layout.validate()
	if err != nil {
		return fmt.Errorf("invalid layout config: %s", err)
	}

	return nil
}

// Size returns layout matrix size
func (l layoutConfig) Size() (int, int) {
	return layoutRows, layoutCols
}

func (l layoutConfig) validate() error {
	u := make(map[string]bool)
	for _, ID := range l {
		if ID == "" {
			continue
		}

		_, ok := u[ID]
		if ok {
			return fmt.Errorf(`"%s" appears twice`, ID)
		}

		u[ID] = true
	}

	if len(u) == 0 {
		return fmt.Errorf("empty layout")
	}

	return nil
}

// Row returns N row from layout matrix
func (l layoutConfig) Row(index int) ([]string, error) {
	rowsCount, colsCount := l.Size()

	if index >= rowsCount || index < 0 {
		return nil, fmt.Errorf("index must be between 0 and %d", rowsCount-1)
	}

	var rc int
	row := make([]string, colsCount)
	for i, v := range l {
		ci := i % colsCount
		row[ci] = v
		if ci != colsCount-1 {
			continue
		}

		if rc == index {
			return row, nil
		}

		rc++
		row = make([]string, colsCount)
	}

	return row, nil
}

// RowsCount returns the numbers of rows from layout matrix
func (l layoutConfig) RowsCount() int {
	c := len(l)
	if c == 0 {
		return 0
	}

	return c / layoutCols
}

// Matrix returns layout matrix
func (l layoutConfig) Matrix() [][]string {
	matrix := make([][]string, l.RowsCount())
	for i := 0; i < l.RowsCount(); i++ {
		row, err := l.Row(i)
		if err != nil {
			panic(err)
		}

		matrix[i] = row
	}

	return matrix
}

// IsScaled returns true is the monitor should be rescaled or not
func (mc monitorConfig) IsScaled() bool {
	return mc.Scale > 0
}

// Scaling returns the scaling factor or 1 if no rescale is necessary
func (mc monitorConfig) Scaling() float32 {
	if !mc.IsScaled() {
		return 1
	}

	return mc.Scale
}

// LoadFile load config file and unserialize it
func LoadFile(filename string) (*Config, error) {
	var cfg Config
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return &cfg, err
	}

	if err = json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	if err = cfg.validate(); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
