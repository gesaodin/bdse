package util

import (
	"fmt"

	"github.com/extrame/xls"
)


func LeerMaticloXLS() {
  if xlFile, err := xls.Open("MAT-21-08-2018.xls", "utf-8"); err == nil {
      if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
          fmt.Print("Total Lines ", sheet1.MaxRow, sheet1.Name)

          for i := 0; i <= (int(sheet1.MaxRow)); i++ {
              row1 := sheet1.Row(i)
              col1 := row1.Col(5)
              col2 := row1.Col(9)
              col3 := row1.Col(10)
              // row1 := sheet1.Rows[uint16(i)]
              // col1 = row1.Cols[5]
              // col2 = row1.Cols[9]
              // fmt.Print("\n", col1.String(xlFile), ",", col2.String(xlFile))
              fmt.Print("\n", col1, ",", col2, ",", col3)
          }
      }


     }
}
