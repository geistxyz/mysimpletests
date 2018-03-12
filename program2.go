package main

import (

	"fmt"
	"github.com/jung-kurt/gofpdf"
	"strings"
	"time"
	"os"
	"path/filepath"
)
var gofpdfDir string

func init() {
	//setRoot()
	gofpdfDir = "C:\\trinity_guard\\tests"
	gofpdf.SetDefaultCatalogSort(true)
	gofpdf.SetDefaultCreationDate(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
}

// Assign the relative path to the gofpdfDir directory based on current working
// directory
func setRoot() {
	wdStr, err := os.Getwd()
	if err == nil {
		gofpdfDir = ""
		list := strings.Split(filepath.ToSlash(wdStr), "/")
		for j := len(list) - 1; j >= 0 && list[j] != "gofpdf"; j-- {
			gofpdfDir = filepath.Join(gofpdfDir, "..")
		}
	} else {
		panic(err)
	}
}

// ImageFile returns a qualified filename in which the path to the image
// directory is prepended to the specified filename.
func ImageFile(fileStr string) string {
	return filepath.Join(gofpdfDir, "image", fileStr)
}

// FontDir returns the path to the font directory.
func FontDir() string {
	return filepath.Join(gofpdfDir, "font")
}

// FontFile returns a qualified filename in which the path to the font
// directory is prepended to the specified filename.
func FontFile(fileStr string) string {
	return filepath.Join(FontDir(), fileStr)
}

// TextFile returns a qualified filename in which the path to the text
// directory is prepended to the specified filename.
func TextFile(fileStr string) string {
	return filepath.Join(gofpdfDir, "text", fileStr)
}

// PdfDir returns the path to the PDF output directory.
func PdfDir() string {
	return filepath.Join(gofpdfDir, "pdf")
}

// PdfFile returns a qualified filename in which the path to the PDF output
// directory is prepended to the specified filename.
func PdfFile(fileStr string) string {
	return filepath.Join(PdfDir(), fileStr)
}

// Filename returns a qualified filename in which the example PDF directory
// path is prepended and the suffix ".pdf" is appended to the specified
// filename.
func Filename(baseStr string) string {
	return PdfFile(baseStr + ".pdf")
}

// referenceCompare compares the specified file with the file's reference copy
// located in the 'reference' subdirectory. All bytes of the two files are
// compared except for the value of the /CreationDate field in the PDF. This
// function succeeds if both files are equivalent except for their
// /CreationDate values or if the reference file does not exist.
func referenceCompare(fileStr string) (err error) {
	var refFileStr, refDirStr, dirStr, baseFileStr string
	dirStr, baseFileStr = filepath.Split(fileStr)
	refDirStr = filepath.Join(dirStr, "reference")
	err = os.MkdirAll(refDirStr, 0755)
	if err == nil {
		refFileStr = filepath.Join(refDirStr, baseFileStr)
		err = gofpdf.ComparePDFFiles(fileStr, refFileStr)
	}
	return
}

// Summary generates a predictable report for use by test examples. If the
// specified error is nil, the filename delimiters are normalized and the
// filename printed to standard output with a success message. If the specified
// error is not nil, its String() value is printed to standard output.
func Summary(err error, fileStr string) {
	if err == nil {
		err = referenceCompare(fileStr)
	}
	if err == nil {
		fileStr = filepath.ToSlash(fileStr)
		fmt.Printf("Successfully generated %s\n", fileStr)
	} else {
		fmt.Println(err)
	}
}
// This example demonstrates the generation of a simple PDF document. Note that
// since only core fonts are used (in this case Times, a synonym for
// Helvetica), an empty string can be specified for the font directory in the
// call to New(). Note also that the example.Filename() and example.Summary()
// functions belong to a separate, internal package and are not part of the
// gofpdf library. If an error occurs at some point during the construction of
// the document, subsequent method calls exit immediately and the error is
// finally retrieved with the output call where it can be handled by the
// application.
func Example() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 16)
	pdf.Cell(40, 10, "Hello World!")
	fileStr := Filename("basic")
	err := pdf.OutputFileAndClose(fileStr)
	Summary(err, fileStr)
	// Output:
	// Successfully generated pdf/basic.pdf
}
// This example demonstrate wrapped table cells
func ExampleFpdf_Rect(fontSize float64) {
	marginCell := 2. // margin of top/bottom of cell
	pdf := gofpdf.New("P", "mm", "A4", "")

	setHeaderAndFooter(pdf , "MI REPORTE MUY LINDO ", "SERVER 00000 " , "DATE.......")
	//pdf.SetFont("Times", "", 12)
	pdf.SetFont("Times", "", fontSize)
	pdf.AddPage()
	pagew, pageh := pdf.GetPageSize()
	mleft, mright, _, mbottom := pdf.GetMargins()
	_, lineHt := pdf.GetFontSize()

	//cols := []float64{pagew - mleft - mright - (6*30), 30, 30,30,30,30 , 30 }  // para size 12
	cols := []float64{ pagew - mleft - mright - (12*15), 15, 15, 15 , 15, 15, 15, 15, 15 , 15, 15, 15,15}  // para size 5

	rows := [][]string{}
	for i := 1; i <= 100; i++ {
		word := fmt.Sprintf("%d:%s", i, strings.Repeat("A", i%100))
		count := fmt.Sprintf("%d", i)
		// rows = append(rows, []string{count, word, word,word, word, word, word}) // size 12
		rows = append(rows, []string{count, word, word ,word, word, word,word, word, word ,word, word, word , word })  // size 5
	}

	for _, row := range rows {
		curx, y := pdf.GetXY()
		x := curx

		height := 0.

		for i, txt := range row {
			lines := pdf.SplitLines([]byte(txt), cols[i])
			h := float64(len(lines))*lineHt + marginCell*float64(len(lines))
			if h > height {
				height = h
			}
		}
		// add a new page if the height of the row doesn't fit on the page
		if pdf.GetY()+height > pageh-mbottom {
			pdf.AddPage()
			y = pdf.GetY()
		}
		for i, txt := range row {
			width := cols[i]
			pdf.Rect(x, y, width, height, "")
			pdf.MultiCell(width, lineHt+marginCell, txt, "", "", false)
			x += width
			pdf.SetXY(x, y)
		}
		pdf.SetXY(curx, y+height)
	}
	fileStr := Filename("Fpdf_WrappedTableCells_5")
	err := pdf.OutputFileAndClose(fileStr)
	Summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_WrappedTableCells.pdf
}

func ExampleFpdf_Rect_10(fontSize float64, orientationStr string, title string ) {

	marginCell := 2. // margin of top/bottom of cell
	pdf := gofpdf.New(orientationStr, "mm", "A4", "")
	setHeaderAndFooter(pdf , title, "SERVER 00000 " , "DATE.......")
	pdf.SetFont("Times", "", fontSize)
	pdf.AddPage()
	pagew, pageh := pdf.GetPageSize()
	mleft, mright, _, mbottom := pdf.GetMargins()
	_, lineHt := pdf.GetFontSize()
	cols := getSizeOfColumns(pagew, mleft, mright) // para size 12
	//cols := []float64{ pagew - mleft - mright - (12*15), 15, 15, 15 , 15, 15, 15, 15, 15 , 15, 15, 15,15}  // para size 5

	rows := getRows()

	for _, row := range rows {
		curx, y := pdf.GetXY()
		x := curx
		height := splitLines(row, pdf, cols, lineHt, marginCell)
		y = AddPage(pdf, height, pageh, mbottom, y)
		createMultiCell(row, cols, pdf, x, y, height, lineHt, marginCell)
		pdf.SetXY(curx, y+height)
	}
	fileStr := Filename("Fpdf_WrappedTableCells_10_Refactorizado")
	err := pdf.OutputFileAndClose(fileStr)
	Summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_WrappedTableCells.pdf
}
func createMultiCell(row []string, cols []float64, pdf *gofpdf.Fpdf, x float64, y float64, height float64, lineHt float64, marginCell float64) {
	for i, txt := range row {
		width := cols[i]
		pdf.Rect(x, y, width, height, "")
		pdf.MultiCell(width, lineHt+marginCell, txt, "", "", false)
		x += width
		pdf.SetXY(x, y)
	}
}
func splitLines(row []string, pdf *gofpdf.Fpdf, cols []float64, lineHt float64, marginCell float64) float64 {
	var height = 0.
	for i, txt := range row {
		lines := pdf.SplitLines([]byte(txt), cols[i])
		h := float64(len(lines))*lineHt + marginCell*float64(len(lines))
		if h > height {
			height = h
		}
	}
	return height
}
func AddPage(pdf *gofpdf.Fpdf, height float64, pageh float64, mbottom float64, y float64) float64 {
	// add a new page if the height of the row doesn't fit on the page
	if pdf.GetY()+height > pageh-mbottom {
		pdf.AddPage()
		y = pdf.GetY()
	}
	return y
}
func getSizeOfColumns(pagew float64, mleft float64, mright float64) []float64 {
	return []float64{pagew - mleft - mright - (6 * 30), 30, 30, 30, 30, 30, 30}
}
func getRows() [][]string {
	rows := [][]string{}
	for i := 1; i <= 100; i++ {
		word := fmt.Sprintf("%d:%s", i, strings.Repeat("A", i%100))
		count := fmt.Sprintf("%d", i)
		rows = append(rows, []string{count, word, word, word, word, word, word}) // size 12
		//rows = append(rows, []string{count, word, word ,word, word, word,word, word, word ,word, word, word , word })  // size 5
	}
	return rows
}
func ExampleFpdf_Rect_LandScape() {
	marginCell := 2. // margin of top/bottom of cell
	pdf := gofpdf.New("L", "mm", "A4", "")
	setHeaderAndFooter(pdf , "MI REPORTE MUY LINDO ", "SERVER 00000 " , "DATE.......")

	//pdf.SetFont("Times", "", 10)
	/*pdf.SetHeaderFunc(func() {
		pdf.Image(ImageFile("logo.png"), 10, 6, 30, 0, false, "", 0, "")
		pdf.SetY(5)
		pdf.SetFont("Times", "B", 15)
		pdf.Cell(80, 0, "")
		pdf.CellFormat(30, 10, "Title", "", 0, "C", false, 0, "")
		pdf.Ln(20)
	})
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Times", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("")
	//pdf.SetFont("Times", "", 5)
	*/
	pdf.SetFont("Times", "", 5)
	pdf.SetFillColor(200, 220, 255)

	pdf.AddPage()
	pagew, pageh := pdf.GetPageSize()
	mleft, mright, _, mbottom := pdf.GetMargins()
	_, lineHt := pdf.GetFontSize()
	/*fmt.Println(abc)
	fmt.Println(lineHt)

	fmt.Println("===")

	fmt.Println(pagew)
	fmt.Println(pageh)
	fmt.Println(mleft)
	fmt.Println(mright)
	fmt.Println(def)
	fmt.Println(mbottom)*/
	// fmt.Println( pagew - mleft - mright - (12*15))

	//cols := []float64{pagew - mleft - mright - (9*30), 30, 30,30,30,30,30,30,30,30 }  // para size 12
	cols := []float64{ pagew - mleft - mright - (18*15), 15, 15, 15 , 15, 15, 15, 15, 15, 15, 15, 15,15, 15, 15 , 15, 15, 15,15}  // para size 5

	rows := [][]string{}
	for i := 1; i <= 100; i++ {
		word := fmt.Sprintf("%d:%s", i, strings.Repeat("A", i%100))
		count := fmt.Sprintf("%d", i)
		//rows = append(rows, []string{count, word, word,word, word, word, word, word, word, word}) // size 12
		rows = append(rows, []string{count, word, word ,word, word, word,word, word, word ,word, word, word , word , word, word ,word, word, word , word })  // size 5
	}

	for _, row := range rows {
		curx, y := pdf.GetXY()
		x := curx

		height := 0.

		for i, txt := range row {
			lines := pdf.SplitLines([]byte(txt), cols[i])
			h := float64(len(lines))*lineHt + marginCell*float64(len(lines))
			if h > height {
				height = h
			}
		}
		// add a new page if the height of the row doesn't fit on the page
		if pdf.GetY()+height > pageh-mbottom {
			pdf.AddPage()
			y = pdf.GetY()
		}
		for i, txt := range row {
			width := cols[i]
			pdf.Rect(x, y, width, height, "")
			//if i%2 == 0{
			pdf.MultiCell(width, lineHt+marginCell, txt, "", "", false)
			//} else {
			//	pdf.MultiCell(width, lineHt+marginCell, txt, "1", "", true)
			//}

			x += width
			pdf.SetXY(x, y)
		}
		pdf.SetXY(curx, y+height)
	}
	fileStr := Filename("Fpdf_WrappedTableCells_L_5")
	err := pdf.OutputFileAndClose(fileStr)
	Summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_WrappedTableCells.pdf
}
func ExampleFpdf_Rect_LandScape_10(fontSize float64) {
	marginCell := 2. // margin of top/bottom of cell
	pdf := gofpdf.New("L", "mm", "A4", "")
	setHeaderAndFooter(pdf , "MI REPORTE MUY LINDO ", "SERVER 00000 " , "DATE.......")

	//pdf.SetFont("Times", "", 5)
	pdf.SetFont("Times", "", fontSize)

	//pdf.SetFillColor(200, 220, 255)

	pdf.AddPage()
	pagew, pageh := pdf.GetPageSize()


	mleft, mright, _, mbottom := pdf.GetMargins()
	_, lineHt := pdf.GetFontSize()
	cols := []float64{pagew - mleft - mright - (9 * 30), 30, 30, 30, 30, 30, 30, 30, 30, 30} // para size 12
	//cols := []float64{ pagew - mleft - mright - (18*15), 15, 15, 15 , 15, 15, 15, 15, 15, 15, 15, 15,15, 15, 15 , 15, 15, 15,15}  // para size 5

	rows := [][]string{}
	for i := 1; i <= 100; i++ {
		word := fmt.Sprintf("%d:%s", i, strings.Repeat("A", i%100))
		count := fmt.Sprintf("%d", i)
		rows = append(rows, []string{count, word, word, word, word, word, word, word, word, word}) // size 12
		//rows = append(rows, []string{count, word, word ,word, word, word,word, word, word ,word, word, word , word , word, word ,word, word, word , word })  // size 5
	}

	for _, row := range rows {
		curx, y := pdf.GetXY()
		x := curx

		height := 0.

		for i, txt := range row {
			lines := pdf.SplitLines([]byte(txt), cols[i])
			h := float64(len(lines))*lineHt + marginCell*float64(len(lines))
			if h > height {
				height = h
			}
		}
		// add a new page if the height of the row doesn't fit on the page
		if pdf.GetY()+height > pageh-mbottom {
			pdf.AddPage()
			y = pdf.GetY()
		}
		for i, txt := range row {
			width := cols[i]
			pdf.Rect(x, y, width, height, "")
			//if i%2 == 0{
			pdf.MultiCell(width, lineHt+marginCell, txt, "", "", false)
			//} else {
			//	pdf.MultiCell(width, lineHt+marginCell, txt, "1", "", true)
			//}

			x += width
			pdf.SetXY(x, y)
		}
		pdf.SetXY(curx, y+height)
	}
	fileStr := Filename("Fpdf_WrappedTableCells_L_10")
	err := pdf.OutputFileAndClose(fileStr)
	Summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_WrappedTableCells.pdf
}
func setHeaderAndFooter(pdf *gofpdf.Fpdf , title string , server string , date_str string) {

	pdf.SetHeaderFunc(func() {
		pdf.Image(ImageFile("logo.png"), 10, 6, 30, 0, false, "", 0, "")
		pdf.SetY(5)
		//pdf.SetX(50)
		pdf.SetFont("Times", "B", 15)
		pdf.Cell(60, 0, "")
		lines := pdf.SplitLines([]byte(title) , 100)

		pdf.MultiCell(100, 5, title, "", "C", false)

		pdf.SetFont("Times", "B", 10)
		wd := pdf.GetStringWidth(title) + 6
		pdf.SetX((210 - wd) / 2)
		var sizel = len(lines) * 5 + 5
		pdf.SetY(float64(sizel)  )
		pdf.Ln(-1)
		pdf.SetFont("Times", "B", 10)
		pdf.CellFormat(105, 8, server, "", 0, "C", false, 0, "")
		pdf.CellFormat(105, 8, date_str, "", 0, "C", false, 0, "")
		sizel = sizel + 16
		pdf.Ln(10)

	})

	pdf.SetFooterFunc(func() {
		pdf.SetY(-10)
		pdf.SetFont("Times", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("")
}



func ExampleFpdf_Rect_LandScape_L_15_10() {
	marginCell := 2. // margin of top/bottom of cell
	pdf := gofpdf.New("L", "mm", "A4", "")
	setHeaderAndFooter(pdf , "MI REPORTE MUY LINDO ", "SERVER 00000 " , "DATE.......")


	pdf.SetFont("Times", "", 10)
	pdf.SetFillColor(200, 220, 255)

	pdf.AddPage()
	pagew, pageh := pdf.GetPageSize()
	mleft, mright, _, mbottom := pdf.GetMargins()
	_, lineHt := pdf.GetFontSize()

	cols := []float64{ pagew - mleft - mright - (18*15), 15, 15, 15 , 15, 15, 15, 15, 15, 15, 15, 15,15, 15, 15 , 15, 15, 15,15}  // para size 5

	rows := [][]string{}
	for i := 1; i <= 100; i++ {
		word := fmt.Sprintf("%d:%s", i, strings.Repeat("A", i%100))
		count := fmt.Sprintf("%d", i)
		//rows = append(rows, []string{count, word, word,word, word, word, word, word, word, word}) // size 12
		rows = append(rows, []string{count, word, word ,word, word, word,word, word, word ,word, word, word , word , word, word ,word, word, word , word })  // size 5
	}

	for _, row := range rows {
		curx, y := pdf.GetXY()
		x := curx

		height := 0.

		for i, txt := range row {
			lines := pdf.SplitLines([]byte(txt), cols[i])
			h := float64(len(lines))*lineHt + marginCell*float64(len(lines))
			if h > height {
				height = h
			}
		}
		// add a new page if the height of the row doesn't fit on the page
		if pdf.GetY()+height > pageh-mbottom {
			pdf.AddPage()
			y = pdf.GetY()
		}
		for i, txt := range row {
			width := cols[i]
			pdf.Rect(x, y, width, height, "")
			//if i%2 == 0{
			pdf.MultiCell(width, lineHt+marginCell, txt, "", "", false)
			//} else {
			//	pdf.MultiCell(width, lineHt+marginCell, txt, "1", "", true)
			//}

			x += width
			pdf.SetXY(x, y)
		}
		pdf.SetXY(curx, y+height)
	}
	fileStr := Filename("Fpdf_WrappedTableCells_L_15_10")
	err := pdf.OutputFileAndClose(fileStr)
	Summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_WrappedTableCells.pdf
}
func ExampleFpdf_Rect_15_10(fontSize float64) {
	marginCell := 2. // margin of top/bottom of cell
	pdf := gofpdf.New("P", "mm", "A4", "")

	setHeaderAndFooter(pdf , "MI REPORTE MUY LINDO ", "SERVER 00000 " , "DATE.......")
	//pdf.SetFont("Times", "", 12)
	pdf.SetFont("Times", "", fontSize)
	pdf.AddPage()
	pagew, pageh := pdf.GetPageSize()
	mleft, mright, _, mbottom := pdf.GetMargins()
	_, lineHt := pdf.GetFontSize()

	//cols := []float64{pagew - mleft - mright - (6*30), 30, 30,30,30,30 , 30 }  // para size 12
	cols := []float64{ pagew - mleft - mright - (12*15), 15, 15, 15 , 15, 15, 15, 15, 15 , 15, 15, 15,15}  // para size 5

	rows := [][]string{}
	for i := 1; i <= 100; i++ {
		word := fmt.Sprintf("%d:%s", i, strings.Repeat("A", i%100))
		count := fmt.Sprintf("%d", i)
		// rows = append(rows, []string{count, word, word,word, word, word, word}) // size 12
		rows = append(rows, []string{count, word, word ,word, word, word,word, word, word ,word, word, word , word })  // size 5
	}

	for _, row := range rows {
		curx, y := pdf.GetXY()
		x := curx

		height := 0.

		for i, txt := range row {
			lines := pdf.SplitLines([]byte(txt), cols[i])
			h := float64(len(lines))*lineHt + marginCell*float64(len(lines))
			if h > height {
				height = h
			}
		}
		// add a new page if the height of the row doesn't fit on the page
		if pdf.GetY()+height > pageh-mbottom {
			pdf.AddPage()
			y = pdf.GetY()
		}
		for i, txt := range row {
			width := cols[i]
			pdf.Rect(x, y, width, height, "")
			pdf.MultiCell(width, lineHt+marginCell, txt, "", "", false)
			x += width
			pdf.SetXY(x, y)
		}
		pdf.SetXY(curx, y+height)
	}
	fileStr := Filename("Fpdf_WrappedTableCells_P_15_10")
	err := pdf.OutputFileAndClose(fileStr)
	Summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_WrappedTableCells.pdf
}

func main() {

	//ExampleFpdf_Rect(5)
	//ExampleFpdf_Rect_LandScape()
	//ExampleFpdf_Rect_LandScape_L_15_10()
	//ExampleFpdf_Rect_15_10(10)
	//ExampleFpdf_Rect_LandScape_10(10)
	ExampleFpdf_Rect_10(5, "P" , "Asertividad: Expresar tus verdaderos sentimientos y defender tus derechos puede ser maravillosamente reconfortante. Cuando dices lo que quieres, independientemente de si lo consigues o no, logras vivir de forma mas autentica y feliz.")
	fmt.Println("Hello world")
}
