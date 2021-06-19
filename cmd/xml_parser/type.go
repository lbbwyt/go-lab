package main

import "encoding/xml"

type Exception struct {
	XMLName xml.Name `xml:"Exception"`

	ExceptionRecord []Record `xml:"ExceptionRecord"`

	Module ModuleType `xml:"Modules"'`
}

type ExceptionRecord struct {
	ExceptionRecords []Record `xml:"FullPath,attr"`
}

type Record struct {
	ModuleName          string `xml:"ModuleName,attr"`
	ExceptionModuleName string `xml:"ExceptionModuleName,attr"`
}

type ModuleType struct {
	ModuleType []Module `xml:"Module"`
}

type Module struct {
	FullPath       string `xml:"FullPath,attr"`
	BaseAddress    string `xml:"BaseAddress,attr"`
	Size           string `xml:"Size,attr"`
	TimeStamp      string `xml:"TimeStamp,attr"`
	FileVersion    string `xml:"FileVersion,attr"`
	ProductVersion string `xml:"ProductVersion,attr"`
	Signature      string `xml:"Signature,attr"`
}
