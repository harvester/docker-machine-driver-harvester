package harvester

type NodeDriverError string

func (e NodeDriverError) Error() string {
	return string(e)
}

const (
	ParseLabelsSyntaxErr = NodeDriverError("parse labels syntax error")
	FormatLabelValueErr  = NodeDriverError("failed for format label value")
)
