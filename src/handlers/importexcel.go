package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"tranzo/src/internal"

	"github.com/tealeg/xlsx"
)

const details = "sheet1"

func (p *Provider) ImportExcelFile(w http.ResponseWriter, r *http.Request) {

	dbSession := p.db.Copy()
	defer dbSession.Close()

	res := make(map[string]interface{})
	buf := &bytes.Buffer{}

	file, err := os.Open("/Users/purushotham/Desktop/tranzo.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	io.Copy(buf, file)

	xLSXFile, err := xlsx.OpenBinary(buf.Bytes())
	if err != nil {
		log.Println(err.Error())
		res["message"] = "Error while reading file"
		res["error"] = err.Error()
		renderJson(w, http.StatusBadRequest, res)
		return
	}

	if err, sheet := internal.FindSheet(xLSXFile.Sheets, details); err == nil {
		internal.ReadExcel(dbSession, sheet)
	}

	res["message"] = "Details Imported successfully..."
	renderJson(w, http.StatusOK, res)

}
