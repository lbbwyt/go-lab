package main

import (
	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o xmlParser.exe .

var xmlParser = NewKXmlParser(getFileName())

func getFileName() string {
	dir, _ := os.Getwd()
	log.Infof("解析根目录：%s", dir)
	return dir + "/exception.xlsx"
}

var ExceptionPool sync.Pool

func main() {
	defer log.Info("解析完成")
	dir, _ := os.Getwd()
	log.Infof("解析根目录：%s", dir)

	//清空文件并自动写入表头
	CleanExcel()

	//并发的形式追加内容
	res := GetAllFileFullPath(dir, ".xml")
	wg := new(sync.WaitGroup)
	for _, v := range res {
		wg.Add(1)
		go func(v string, group *sync.WaitGroup) {
			parserXml(group, v)
		}(v, wg)
	}
	wg.Wait()
}

func CleanExcel() {
	log.Info("clean excel file: " + xmlParser.fileName)
	file, err := excelize.OpenFile(xmlParser.fileName)
	if err != nil {
		log.WithError(err).Error("OpenFile err" + xmlParser.fileName)
		panic(err)
	}
	sheet_name := "Sheet1"
	//获取流式写入器
	streamWriter, err := file.NewStreamWriter(sheet_name)
	if err != nil {
		panic(err)
	}

	//将新加modules写进流式写入器
	row := make([]interface{}, 0, 8)
	row = append(row, "DumpId")
	row = append(row, "FullPath")
	row = append(row, "BaseAddress")
	row = append(row, "Signature")
	row = append(row, "FileVersion")
	row = append(row, "ProductVersion")
	row = append(row, "Size")
	row = append(row, "TimeStamp")
	cell, _ := excelize.CoordinatesToCellName(1, 1) //决定写入的位置
	if err := streamWriter.SetRow(cell, row); err != nil {
		log.WithError(err).Error("SetRow err")
		panic(err)
	}

	//结束流式写入过程
	if err := streamWriter.Flush(); err != nil {
		log.WithError(err).Error("Flush err")
		panic(err)
	}
	//保存工作簿
	if err := file.SaveAs(xmlParser.fileName); err != nil {
		log.WithError(err).Error("save err")
		panic(err)
	}

}

func parserXml(wg *sync.WaitGroup, xmlName string) {
	defer wg.Done()
	log.Infof("解析文件：%s ", xmlName)

	exception := xmlParser.ExceptionPool.Get().(*Exception)
	defer xmlParser.ExceptionPool.Put(exception)
	xmlParser.OpenFile(xmlName, exception)

	//将解析文件的内容追加到指定excel、
	xmlParser.lock.Lock()
	defer xmlParser.lock.Unlock()

	if exception.ExceptionRecord == nil || len(exception.ExceptionRecord) == 0 {
		return
	}

	er := exception.ExceptionRecord[0]

	moduleName := er.ModuleName
	if moduleName == "" {
		moduleName = er.ExceptionModuleName
	}
	if moduleName == "" {
		return
	}
	if len(exception.Module.ModuleType) == 0 {
		return
	}

	modules := make([]Module, 0)
	for _, v := range exception.Module.ModuleType {
		if v.FullPath == moduleName {
			modules = append(modules, v)
		}
	}

	//批量写入modules
	file, err := excelize.OpenFile(xmlParser.fileName)
	if err != nil {
		log.WithError(err).Error("OpenFile err" + xmlParser.fileName)
		//excelize.NewFile(xmlParser.fileName)
	}
	sheet_name := "Sheet1"
	//获取流式写入器
	streamWriter, err := file.NewStreamWriter(sheet_name)
	if err != nil {
		panic(err)
	}
	rows, _ := file.GetRows(sheet_name) //获取行内容
	cols, _ := file.GetCols(sheet_name) //获取列内容
	//将源文件内容先写入excel
	for rowid, row_pre := range rows {
		row_p := make([]interface{}, len(cols))
		for colID_p := 0; colID_p < len(cols); colID_p++ {
			if row_pre == nil {
				row_p[colID_p] = nil
			} else {
				row_p[colID_p] = row_pre[colID_p]
			}
		}
		cell_pre, _ := excelize.CoordinatesToCellName(1, rowid+1)
		if err := streamWriter.SetRow(cell_pre, row_p); err != nil {
			log.WithError(err).Error("First SetRow err")
			panic(err)
		}
	}

	//解析dunpId
	var dumpId = getDumpId(xmlName)

	//将新加modules写进流式写入器
	for rowID := 0; rowID < len(modules); rowID++ {
		m := modules[rowID]
		row := make([]interface{}, 0, 8)
		row = append(row, dumpId)
		row = append(row, m.FullPath)
		row = append(row, m.BaseAddress)
		row = append(row, m.Signature)
		row = append(row, m.FileVersion)
		row = append(row, m.ProductVersion)
		row = append(row, m.Size)
		row = append(row, m.TimeStamp)
		cell, _ := excelize.CoordinatesToCellName(1, rowID+len(rows)+1) //决定写入的位置
		if err := streamWriter.SetRow(cell, row); err != nil {
			log.WithError(err).Error("SetRow err")
			panic(err)
		}
	}
	//结束流式写入过程
	if err := streamWriter.Flush(); err != nil {
		log.WithError(err).Error("Flush err")
		panic(err)
	}
	//保存工作簿
	if err := file.SaveAs(xmlParser.fileName); err != nil {
		log.WithError(err).Error("save err")
		panic(err)
	}

}

func getDumpId(path string) string {
	index := strings.LastIndex(path, "/")
	sub := path[0:index]
	index = strings.LastIndex(sub, "/")
	last := path[index:len(sub)]
	return strings.Trim(last, "/")
}

func GetAllFileFullPath(curFullPath string, suffix string) []string {
	var res = make([]string, 0)
	filepath.Walk(curFullPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, suffix) {
			res = append(res, path)
		}
		return nil
	})
	return res
}
