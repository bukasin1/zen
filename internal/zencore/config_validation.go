package zencore

import "fmt"

func (c Config) Validate() error {

	if c.AppName == "" {
		return fmt.Errorf("app name cannot be empty")
	}

	if c.HTTP.Addr == "" {
		return fmt.Errorf("http address cannot be empty")
	}

	if c.HTTP.ShutdownTimeout <= 0 {
		return fmt.Errorf(
			"http shutdown timeout must be greater than zero",
		)
	}

	return nil
}
