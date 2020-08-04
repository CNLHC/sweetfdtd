package fdtd

import (
	"encoding/xml"
	"errors"
	"strings"
)

var (
	ErrorLicense error = errors.New("license Error")
)

const (
	LineTypeUpdate   = "Update"
	LineTypePlain    = "Plain"
	LineTypeStatus   = "Status"
	LineTypeComplete = "Complete"
)

type FDTDLine struct {
	Type       string
	Update     LineProgress
	Status     string
	Unknown    string
	PlainError error
}

type LineProgress struct {
	XMLName  xml.Name `xml:"update"`
	Progress float32  `xml:"progress"`
	Time     string   `xml:"time"`
	Shutoff  float32  `xml:"shutoff"`
}
type LineStatus struct {
	XMLName xml.Name `xml:"status"`
}
type LineComplete struct {
	XMLName xml.Name `xml:"complete"`
}
type LineUnknownXml struct {
	XMLName xml.Name
}

func ParseStdOutLine(line string) FDTDLine {
	var progress LineProgress
	if err := xml.Unmarshal([]byte(line), &progress); err == nil {
		if progress.XMLName.Local != "" {
			return FDTDLine{
				Type:   LineTypeUpdate,
				Update: progress,
			}
		}
	}

	var status LineStatus
	if err := xml.Unmarshal([]byte(line), &status); err == nil {
		var s string
		xml.Unmarshal([]byte(line), &s)

		if status.XMLName.Local != "" {
			return FDTDLine{
				Type:   LineTypeStatus,
				Status: s,
			}
		}
	}

	var complete LineComplete
	if err := xml.Unmarshal([]byte(line), &complete); err == nil {
		return FDTDLine{Type: LineTypeComplete}
	}

	return FDTDLine{
		Type:       LineTypePlain,
		PlainError: FilterPlainError(line),
	}
}

func FilterPlainError(line string) error {
	if strings.Contains(line, "The flexNet error code") {
		return ErrorLicense
	}
	if strings.Contains(line, "there was a failure with the license") {
		return ErrorLicense
	}
	return nil
}
