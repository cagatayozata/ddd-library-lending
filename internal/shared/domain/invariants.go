package domain

import "fmt"

// Invariants — one place for domain precondition failures.
func Require(ok bool, msg string) error {
	if !ok {
		return fmt.Errorf("invariant: %s", msg)
	}
	return nil
}

func RequirePresent[T comparable](v T, name string) error {
	var zero T
	if v == zero {
		return fmt.Errorf("invariant: %s is required", name)
	}
	return nil
}

func RequireRange(v, min, max int, name string) error {
	if v < min || v > max {
		return fmt.Errorf("invariant: %s must be in [%d, %d], got %d", name, min, max, v)
	}
	return nil
}

func RequireText(s, name string, min, max int) error {
	n := len([]rune(s))
	if n < min || n > max {
		return fmt.Errorf("invariant: %s length must be in [%d, %d], got %d", name, min, max, n)
	}
	return nil
}
