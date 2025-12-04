package models

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// VersionRange represents a version constraint
type VersionRange struct {
	Raw        string // Original version string
	Constraint *semver.Constraints
}

// NewVersionRange creates a new VersionRange from a constraint string
func NewVersionRange(constraint string) (*VersionRange, error) {
	if constraint == "" || constraint == "*" {
		// Accept any version
		constraint = "*"
	}

	// Normalize common patterns
	constraint = normalizeConstraint(constraint)

	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return &VersionRange{
			Raw:        constraint,
			Constraint: nil,
		}, nil // Return without error but with nil constraint for best-effort matching
	}

	return &VersionRange{
		Raw:        constraint,
		Constraint: c,
	}, nil
}

// normalizeConstraint converts common version patterns to semver constraints
func normalizeConstraint(constraint string) string {
	// Remove whitespace
	constraint = strings.TrimSpace(constraint)

	// Handle exact versions without operators
	if !strings.ContainsAny(constraint, "><=~^*") && constraint != "" {
		// Check if it looks like a version
		if _, err := semver.NewVersion(constraint); err == nil {
			return constraint
		}
	}

	// Handle ranges like [1.0,2.0) -> >=1.0 <2.0
	if strings.HasPrefix(constraint, "[") || strings.HasPrefix(constraint, "(") {
		return convertRangeBrackets(constraint)
	}

	return constraint
}

// convertRangeBrackets converts bracket notation to semver constraints
func convertRangeBrackets(constraint string) string {
	// Remove brackets
	constraint = strings.Trim(constraint, "[]() ")
	parts := strings.Split(constraint, ",")

	if len(parts) != 2 {
		return constraint
	}

	lower := strings.TrimSpace(parts[0])
	upper := strings.TrimSpace(parts[1])

	result := ""
	if lower != "" {
		result = ">=" + lower
	}
	if upper != "" {
		if result != "" {
			result += " "
		}
		result += "<" + upper
	}

	return result
}

// Check returns true if the given version satisfies this constraint
func (vr *VersionRange) Check(version string) bool {
	if vr.Constraint == nil {
		// No valid constraint, accept any version
		return true
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		// Can't parse version, be lenient
		return true
	}

	return vr.Constraint.Check(v)
}

// Intersect checks if two version ranges have any overlap
func (vr *VersionRange) Intersect(other *VersionRange) bool {
	// If either has no constraint, they intersect
	if vr.Constraint == nil || other.Constraint == nil {
		return true
	}

	// Try a range of common versions to see if both accept any
	testVersions := []string{
		"1.0.0", "1.1.0", "1.2.0", "1.5.0",
		"2.0.0", "2.1.0", "2.5.0",
		"3.0.0", "4.0.0", "5.0.0",
	}

	for _, tv := range testVersions {
		if vr.Check(tv) && other.Check(tv) {
			return true
		}
	}

	return false
}

// String returns the string representation
func (vr *VersionRange) String() string {
	if vr.Raw != "" {
		return vr.Raw
	}
	return "*"
}

// MarshalText implements encoding.TextMarshaler
func (vr *VersionRange) MarshalText() ([]byte, error) {
	return []byte(vr.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (vr *VersionRange) UnmarshalText(data []byte) error {
	constraint := string(data)
	newVR, err := NewVersionRange(constraint)
	if err != nil {
		return fmt.Errorf("failed to unmarshal version range: %w", err)
	}
	*vr = *newVR
	return nil
}
