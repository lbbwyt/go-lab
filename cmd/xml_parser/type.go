package main

import "encoding/xml"

type Exception struct {
	XMLName xml.Name `xml:"Exception"`

	ExceptionRecord      []Record         `xml:"ExceptionRecord"`
	AdditionalInfomation []AdditionalInfo `xml:"AdditionalInfomation"`

	Module ModuleType `xml:"Modules"'`
}

type AdditionalInfo struct {
	TargetProcCmd string `xml:"TargetProcCmd,attr"`
}

func (e *Exception) Reset() {
	e.ExceptionRecord = nil
	e.Module = ModuleType{}
}

type ExceptionRecord struct {
	ExceptionRecords []Record `xml:"FullPath,attr"`
}

type Record struct {
	ModuleName          string `xml:"ExceptionModuleName,attr"`
	ExceptionModuleName string `xml:"ExceptionEspReturnModuleName,attr"`
}

type ModuleType struct {
	ModuleType []Module `xml:"Module"`
}

type Module struct {
	DumpId         string
	TargetProcCmd  string
	FullPath       string `xml:"FullPath,attr"`
	BaseAddress    string `xml:"BaseAddress,attr"`
	Size           string `xml:"Size,attr"`
	TimeStamp      string `xml:"TimeStamp,attr"`
	FileVersion    string `xml:"FileVersion,attr"`
	ProductVersion string `xml:"ProductVersion,attr"`
	Signature      string `xml:"Signature,attr"`
}
