package internal

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"tranzo/src/models"

	"github.com/tealeg/xlsx"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type headerMap map[string]int

func getHeaders(sheet *xlsx.Sheet, headers []string) headerMap {

	data := headerMap{}
	rows := sheet.Rows
	if len(rows) == 0 {
		return data
	}

	cells := rows[0].Cells

	for _, header := range headers {
		data[header] = -1
		for colIdx, cell := range cells {
			if strings.EqualFold(strings.TrimSpace(cell.Value), header) {
				data[header] = colIdx
				break
			}
		}
	}

	return data

}

func ReadExcel(dbSession *mgo.Session, sheet *xlsx.Sheet) {
	rows := sheet.Rows
	if len(rows) == 0 {
		return
	}

	headers := []string{"Name", "Age", "Gender"}
	headerData := getHeaders(sheet, headers)

	if colIdx, ok := headerData[headers[0]]; !ok || colIdx == -1 {
		return
	}

	for i, row := range rows[1:] {

		if len(row.Cells) == 0 {
			//If there occurs an empty row skip and continue
			continue
		}
		detIns := &models.Details{}

		if colIdx, ok := headerData[headers[0]]; ok && len(row.Cells) > colIdx {
			detIns.Name = strings.TrimSpace(row.Cells[colIdx].Value)
		}

		if colIdx, ok := headerData[headers[1]]; ok && colIdx > -1 && len(row.Cells) > colIdx {
			detIns.Age = row.Cells[colIdx].Value
		}

		if colIdx, ok := headerData[headers[2]]; ok && colIdx > -1 && len(row.Cells) > colIdx {
			detIns.Gender = row.Cells[colIdx].Value
		}

		if hasErr, err := detIns.Validate(); hasErr {
			fmt.Printf("%d th Row %v ", i, err["name"])
			continue
		}

		detIns.Id = bson.NewObjectId()
		err := dbSession.DB("tranzo").C("details").Insert(&detIns)
		if err != nil {
			log.Printf("ERROR: Insert Details  (%s) - %s", detIns.Name, err)
		}
	}

}

func FindSheet(sheets []*xlsx.Sheet, name string) (error, *xlsx.Sheet) {

	//find sheet
	for _, sheet := range sheets {
		if strings.EqualFold(strings.TrimSpace(sheet.Name), strings.TrimSpace(name)) {
			return nil, sheet
		}
	}

	return errors.New("Not found"), nil
}
