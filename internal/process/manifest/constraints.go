package manifest

import "github.com/Masterminds/semver/v3"

type Constraints struct {
	semver.Constraints

	isReference bool
}

func NewConstraints(text string) (*Constraints, error) {
	constraints := &Constraints{}
	if err := constraints.UnmarshalText([]byte(text)); err != nil {
		return nil, err
	}

	return constraints, nil
}

func NewReferenceConstraintsWithVersion(text string) (*Constraints, error) {
	constraints := &Constraints{isReference: true}
	if err := constraints.Constraints.UnmarshalText([]byte(text)); err != nil {
		return nil, err
	}

	return constraints, nil
}

func (constraints *Constraints) IsReference() bool {
	return constraints.isReference
}

func (constraints *Constraints) MarshalText() ([]byte, error) {
	if constraints.isReference {
		return []byte("reference"), nil
	}

	return constraints.Constraints.MarshalText()
}

func (constraints *Constraints) UnmarshalText(text []byte) error {
	if string(text) == "reference" {
		constraints.isReference = true
		return nil
	}

	return constraints.Constraints.UnmarshalText(text)
}
