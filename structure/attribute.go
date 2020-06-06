package structure

// Attribute column attribute
type Attribute string

// Attributes
const (
	Unsigned      Attribute = "unsigned"
	Nullable      Attribute = "nullable"
	AutoIncrement Attribute = "auto_increment"
	Stored        Attribute = "stored"
)
