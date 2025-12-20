package models

// RangeValues represents a min-max range for a productivity level
type RangeValues struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// ProductivityRange represents the productivity thresholds with min-max ranges
// Used for classifying metrics into OK, Alert, or Critical levels
type ProductivityRange struct {
	Ok       RangeValues `json:"ok"`
	Alert    RangeValues `json:"alert"`
	Critical RangeValues `json:"critical"`
}

// ClassifyValue determines the productivity level for a given value based on the ranges
// Returns the ProductivityEnum (Ok, Alert, Critical) based on which range the value falls into
func (pr *ProductivityRange) ClassifyValue(value float64) ProductivityEnum {
	// Check if value falls within OK range
	if value >= pr.Ok.Min && value <= pr.Ok.Max {
		return ProductivityOk
	}
	// Check if value falls within Alert range
	if value >= pr.Alert.Min && value <= pr.Alert.Max {
		return ProductivityAlert
	}
	// Check if value falls within Critical range
	if value >= pr.Critical.Min && value <= pr.Critical.Max {
		return ProductivityCritical
	}
	// Default to Critical if outside all ranges
	return ProductivityCritical
}
