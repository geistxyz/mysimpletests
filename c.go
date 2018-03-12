package main

 
 import (
	 "bufio"
	 "bytes"
	 "fmt"
     "io"
	 "io/ioutil"

	 "github.com/jung-kurt/gofpdf"
	 	"os"
	"path/filepath"
	"strings"
	"time"

	
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
 // since only core fonts are used (in this case Arial, a synonym for
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
	 pdf.SetFont("Arial", "B", 16)
	 pdf.Cell(40, 10, "Hello World!")
	 fileStr := Filename("basic")
	 err := pdf.OutputFileAndClose(fileStr)
	 Summary(err, fileStr)
	 // Output:
	 // Successfully generated pdf/basic.pdf
 }
 
// This example demonstrates word-wrapping, line justification and
 // page-breaking.
 func ExampleFpdf_MultiCell() {
	 pdf := gofpdf.New("P", "mm", "A4", "")
	 titleStr := "20000 Leagues Under the Seas"
	 pdf.SetTitle(titleStr, false)
	 pdf.SetAuthor("Jules Verne", false)
	 pdf.SetHeaderFunc(func() {
		 // Arial bold 15
		 pdf.SetFont("Arial", "B", 15)
		 // Calculate width of title and position
		 wd := pdf.GetStringWidth(titleStr) + 6
		 pdf.SetX((210 - wd) / 2)
		 // Colors of frame, background and text
		 pdf.SetDrawColor(0, 80, 180)
		 pdf.SetFillColor(230, 230, 0)
		 pdf.SetTextColor(220, 50, 50)
		 // Thickness of frame (1 mm)
		 pdf.SetLineWidth(1)
		 // Title
		 pdf.CellFormat(wd, 9, titleStr, "1", 1, "C", true, 0, "")
		 // Line break
		 pdf.Ln(10)
	 })
	 pdf.SetFooterFunc(func() {
		 // Position at 1.5 cm from bottom
		 pdf.SetY(-15)
		 // Arial italic 8
		 pdf.SetFont("Arial", "I", 8)
		 // Text color in gray
		 pdf.SetTextColor(128, 128, 128)
		 // Page number
		 pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			 "", 0, "C", false, 0, "")
	 })
	 chapterTitle := func(chapNum int, titleStr string) {
		 // 	// Arial 12
		 pdf.SetFont("Arial", "", 12)
		 // Background color
		 pdf.SetFillColor(200, 220, 255)
		 // Title
		 pdf.CellFormat(0, 6, fmt.Sprintf("Chapter %d : %s", chapNum, titleStr),
			 "", 1, "L", true, 0, "")
		 // Line break
		 pdf.Ln(4)
	 }
	 chapterBody := func(fileStr string) {
		 // Read text file
		 txtStr, err := ioutil.ReadFile(fileStr)
		 if err != nil {
			 pdf.SetError(err)
		 }
		 // Times 12
		 pdf.SetFont("Times", "", 12)
		 // Output justified text
		 pdf.MultiCell(0, 5, string(txtStr), "", "", false)
		 // Line break
		 pdf.Ln(-1)
		 // Mention in italics
		 pdf.SetFont("", "I", 0)
		 pdf.Cell(0, 5, "(end of excerpt)")
	 }
	 printChapter := func(chapNum int, titleStr, fileStr string) {
		 pdf.AddPage()
		 chapterTitle(chapNum, titleStr)
		 chapterBody(fileStr)
	 }
	 printChapter(1, "A RUNAWAY REEF", TextFile("20k_c1.txt"))
	 printChapter(2, "THE PROS AND CONS", TextFile("20k_c2.txt"))
	 fileStr := Filename("Fpdf_MultiCell")
	 err := pdf.OutputFileAndClose(fileStr)
	 Summary(err, fileStr)
	 // Output:
	 // Successfully generated pdf/Fpdf_MultiCell.pdf
 }
 // This example demonsrates the generation of headers, footers and page breaks.
 func ExampleFpdf_AddPage() {
	 pdf := gofpdf.New("P", "mm", "A4", "")
	 pdf.SetHeaderFunc(func() {
		 pdf.Image(ImageFile("logo.png"), 10, 6, 30, 0, false, "", 0, "")
		 pdf.SetY(5)
		 pdf.SetFont("Arial", "B", 15)
		 pdf.Cell(80, 0, "")
		 pdf.CellFormat(30, 10, "Title", "1", 0, "C", false, 0, "")
		 pdf.Ln(20)
	 })
	 pdf.SetFooterFunc(func() {
		 pdf.SetY(-15)
		 pdf.SetFont("Arial", "I", 8)
		 pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()),
			 "", 0, "C", false, 0, "")
	 })
	 pdf.AliasNbPages("")
	 pdf.AddPage()
	 pdf.SetFont("Times", "", 12)
	 for j := 1; j <= 40; j++ {
		 pdf.CellFormat(0, 10, fmt.Sprintf("Printing line number %d", j),
			 "", 1, "", false, 0, "")
	 }
	 fileStr := Filename("Fpdf_AddPage")
	 err := pdf.OutputFileAndClose(fileStr)
	 Summary(err, fileStr)
	 // Output:
	 // Successfully generated pdf/Fpdf_AddPage.pdf
 }
 
 // This example demonstrates the generation of a PDF document that has multiple
 // columns. This is accomplished with the SetLeftMargin() and Cell() methods.
 func ExampleFpdf_SetLeftMargin() {
	 var y0 float64
	 var crrntCol int
	 pdf := gofpdf.New("P", "mm", "A4", "")
	 pdf.SetDisplayMode("fullpage", "TwoColumnLeft")
	 titleStr := "20000 Leagues Under the Seas"
	 pdf.SetTitle(titleStr, false)
	 pdf.SetAuthor("Jules Verne", false)
	 setCol := func(col int) {
		 // Set position at a given column
		 crrntCol = col
		 x := 10.0 + float64(col)*65.0
		 pdf.SetLeftMargin(x)
		 pdf.SetX(x)
	 }
	 chapterTitle := func(chapNum int, titleStr string) {
		 // Arial 12
		 pdf.SetFont("Arial", "", 12)
		 // Background color
		 pdf.SetFillColor(200, 220, 255)
		 // Title
		 pdf.CellFormat(0, 6, fmt.Sprintf("Chapter %d : %s", chapNum, titleStr),
			 "", 1, "L", true, 0, "")
		 // Line break
		 pdf.Ln(4)
		 y0 = pdf.GetY()
	 }
	 chapterBody := func(fileStr string) {
		 // Read text file
		 txtStr, err := ioutil.ReadFile(fileStr)
		 if err != nil {
			 pdf.SetError(err)
		 }
		 // Font
		 pdf.SetFont("Times", "", 12)
		 // Output text in a 6 cm width column
		 pdf.MultiCell(60, 5, string(txtStr), "", "", false)
		 pdf.Ln(-1)
		 // Mention
		 pdf.SetFont("", "I", 0)
		 pdf.Cell(0, 5, "(end of excerpt)")
		 // Go back to first column
		 setCol(0)
	 }
	 printChapter := func(num int, titleStr, fileStr string) {
		 // Add chapter
		 pdf.AddPage()
		 chapterTitle(num, titleStr)
		 chapterBody(fileStr)
	 }
	 pdf.SetAcceptPageBreakFunc(func() bool {
		 // Method accepting or not automatic page break
		 if crrntCol < 2 {
			 // Go to next column
			 setCol(crrntCol + 1)
			 // Set ordinate to top
			 pdf.SetY(y0)
			 // Keep on page
			 return false
		 }
		 // Go back to first column
		 setCol(0)
		 // Page break
		 return true
	 })
	 pdf.SetHeaderFunc(func() {
		 // Arial bold 15
		 pdf.SetFont("Arial", "B", 15)
		 // Calculate width of title and position
		 wd := pdf.GetStringWidth(titleStr) + 6
		 pdf.SetX((210 - wd) / 2)
		 // Colors of frame, background and text
		 pdf.SetDrawColor(0, 80, 180)
		 pdf.SetFillColor(230, 230, 0)
		 pdf.SetTextColor(220, 50, 50)
		 // Thickness of frame (1 mm)
		 pdf.SetLineWidth(1)
		 // Title
		 pdf.CellFormat(wd, 9, titleStr, "1", 1, "C", true, 0, "")
		 // Line break
		 pdf.Ln(10)
		 // Save ordinate
		 y0 = pdf.GetY()
	 })
	 pdf.SetFooterFunc(func() {
		 // Position at 1.5 cm from bottom
		 pdf.SetY(-15)
		 // Arial italic 8
		 pdf.SetFont("Arial", "I", 8)
		 // Text color in gray
		 pdf.SetTextColor(128, 128, 128)
		 // Page number
		 pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
			 "", 0, "C", false, 0, "")
	 })
	 printChapter(1, "A RUNAWAY REEF", TextFile("20k_c1.txt"))
	 printChapter(2, "THE PROS AND CONS", TextFile("20k_c2.txt"))
	 fileStr := Filename("Fpdf_SetLeftMargin_multicolumn")
	 err := pdf.OutputFileAndClose(fileStr)
	 Summary(err, fileStr)
	 // Output:
	 // Successfully generated pdf/Fpdf_SetLeftMargin_multicolumn.pdf
 }
 
  func loremList() []string {
	 return []string{
		 "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod " +
			 "tempor incididunt ut labore et dolore magna aliqua.",
		 "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut " +
			 "aliquip ex ea commodo consequat.",
		 "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum " +
			 "dolore eu fugiat nulla pariatur.",
		 "Excepteur sint occaecat cupidatat non proident, sunt in culpa qui " +
			 "officia deserunt mollit anim id est laborum.",
	 }
 }
 
  // This example demonstrates word-wrapped table cells
 func ExampleFpdf_SplitLines_tables() {
	 const (
		 colCount = 3
		 colWd    = 60.0
		 marginH  = 15.0
		 lineHt   = 5.5
		 cellGap  = 2.0
	 )
	 // var colStrList [colCount]string
	 type cellType struct {
		 str  string
		 list [][]byte
		 ht   float64
	 }
	 var (
		 cellList [colCount]cellType
		 cell     cellType
	 )
 
	 pdf := gofpdf.New("P", "mm", "A4", "") // 210 x 297
	 header := [colCount]string{"Column A", "Column B", "Column C"}
	 alignList := [colCount]string{"L", "C", "R"}
	 strList := loremList()
	 pdf.SetMargins(marginH, 15, marginH)
	 pdf.SetFont("Arial", "", 14)
	 pdf.AddPage()
 
	 // Headers
	 pdf.SetTextColor(224, 224, 224)
	 pdf.SetFillColor(64, 64, 64)
	 for colJ := 0; colJ < colCount; colJ++ {
		 pdf.CellFormat(colWd, 10, header[colJ], "1", 0, "CM", true, 0, "")
	 }
	 pdf.Ln(-1)
	 pdf.SetTextColor(24, 24, 24)
	 pdf.SetFillColor(255, 255, 255)
 
	 // Rows
	 y := pdf.GetY()
	 count := 0
	 for rowJ := 0; rowJ < 2; rowJ++ {
		 maxHt := lineHt
		 // Cell height calculation loop
		 for colJ := 0; colJ < colCount; colJ++ {
			 count++
			 if count > len(strList) {
				 count = 1
			 }
			 cell.str = strings.Join(strList[0:count], " ")
			 cell.list = pdf.SplitLines([]byte(cell.str), colWd-cellGap-cellGap)
			 cell.ht = float64(len(cell.list)) * lineHt
			 if cell.ht > maxHt {
				 maxHt = cell.ht
			 }
			 cellList[colJ] = cell
		 }
		 // Cell render loop
		 x := marginH
		 for colJ := 0; colJ < colCount; colJ++ {
			 pdf.Rect(x, y, colWd, maxHt+cellGap+cellGap, "D")
			 cell = cellList[colJ]
			 cellY := y + cellGap + (maxHt-cell.ht)/2
			 for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				 pdf.SetXY(x+cellGap, cellY)
				 pdf.CellFormat(colWd-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0,
					 alignList[colJ], false, 0, "")
				 cellY += lineHt
			 }
			 x += colWd
		 }
		 y += maxHt + cellGap + cellGap
	 }
 
	 fileStr := Filename("Fpdf_SplitLines_tables")
	 err := pdf.OutputFileAndClose(fileStr)
	 Summary(err, fileStr)
	 // Output:
	 // Successfully generated pdf/Fpdf_SplitLines_tables.pdf
 }
 
 
  // Convert 'ABCDEFG' to, for example, 'A,BCD,EFG'
 func strDelimit(str string, sepstr string, sepcount int) string {
	 pos := len(str) - sepcount
	 for pos > 0 {
		 str = str[:pos] + sepstr + str[pos:]
		 pos = pos - sepcount
	 }
	 return str
 }
 
 
  // This example demonstrates various table styles.
 func ExampleFpdf_CellFormat_tables() {
	 pdf := gofpdf.New("P", "mm", "A4", "")
	 type countryType struct {
		 nameStr, capitalStr, areaStr, popStr string
	 }
	 countryList := make([]countryType, 0, 8)
	 header := []string{"Country", "Capital", "Area (sq km)", "Pop. (thousands)"}
	 loadData := func(fileStr string) {
		 fl, err := os.Open(fileStr)
		 if err == nil {
			 scanner := bufio.NewScanner(fl)
			 var c countryType
			 for scanner.Scan() {
				 // Austria;Vienna;83859;8075
				 lineStr := scanner.Text()
				 list := strings.Split(lineStr, ";")
				 if len(list) == 4 {
					 c.nameStr = list[0]
					 c.capitalStr = list[1]
					 c.areaStr = list[2]
					 c.popStr = list[3]
					 countryList = append(countryList, c)
				 } else {
					 err = fmt.Errorf("error tokenizing %s", lineStr)
				 }
			 }
			 fl.Close()
			 if len(countryList) == 0 {
				 err = fmt.Errorf("error loading data from %s", fileStr)
			 }
		 }
		 if err != nil {
			 pdf.SetError(err)
		 }
	 }
	 // Simple table
	 basicTable := func() {
		 for _, str := range header {
			 pdf.CellFormat(40, 7, str, "1", 0, "", false, 0, "")
		 }
		 pdf.Ln(-1)
		 for _, c := range countryList {
			 pdf.CellFormat(40, 6, c.nameStr, "1", 0, "", false, 0, "")
			 pdf.CellFormat(40, 6, c.capitalStr, "1", 0, "", false, 0, "")
			 pdf.CellFormat(40, 6, c.areaStr, "1", 0, "", false, 0, "")
			 pdf.CellFormat(40, 6, c.popStr, "1", 0, "", false, 0, "")
			 pdf.Ln(-1)
		 }
	 }
	 // Better table
	 improvedTable := func() {
		 // Column widths
		 w := []float64{40.0, 35.0, 40.0, 45.0}
		 wSum := 0.0
		 for _, v := range w {
			 wSum += v
		 }
		 // 	Header
		 for j, str := range header {
			 pdf.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
		 }
		 pdf.Ln(-1)
		 // Data
		 for _, c := range countryList {
			 pdf.CellFormat(w[0], 6, c.nameStr, "LR", 0, "", false, 0, "")
			 pdf.CellFormat(w[1], 6, c.capitalStr, "LR", 0, "", false, 0, "")
			 pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3),
				 "LR", 0, "R", false, 0, "")
			 pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
				 "LR", 0, "R", false, 0, "")
			 pdf.Ln(-1)
		 }
		 pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	 }
	 // Colored table
	 fancyTable := func() {
		 // Colors, line width and bold font
		 pdf.SetFillColor(255, 0, 0)
		 pdf.SetTextColor(255, 255, 255)
		 pdf.SetDrawColor(128, 0, 0)
		 pdf.SetLineWidth(.3)
		 pdf.SetFont("", "B", 0)
		 // 	Header
		 w := []float64{40, 35, 40, 45}
		 wSum := 0.0
		 for _, v := range w {
			 wSum += v
		 }
		 for j, str := range header {
			 pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
		 }
		 pdf.Ln(-1)
		 // Color and font restoration
		 pdf.SetFillColor(224, 235, 255)
		 pdf.SetTextColor(0, 0, 0)
		 pdf.SetFont("", "", 0)
		 // 	Data
		 fill := false
		 for _, c := range countryList {
			 pdf.CellFormat(w[0], 6, c.nameStr, "LR", 0, "", fill, 0, "")
			 pdf.CellFormat(w[1], 6, c.capitalStr, "LR", 0, "", fill, 0, "")
			 pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3),
				 "LR", 0, "R", fill, 0, "")
			 pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
				 "LR", 0, "R", fill, 0, "")
			 pdf.Ln(-1)
			 fill = !fill
		 }
		 pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	 }
	 loadData(TextFile("countries.txt"))
	 pdf.SetFont("Arial", "", 14)
	 pdf.AddPage()
	 basicTable()
	 pdf.AddPage()
	 improvedTable()
	 pdf.AddPage()
	 fancyTable()
	 fileStr := Filename("Fpdf_CellFormat_tables")
	 err := pdf.OutputFileAndClose(fileStr)
	 Summary(err, fileStr)
	 // Output:
	 // Successfully generated pdf/Fpdf_CellFormat_tables.pdf
 }
 
 
  func lorem() string {
	 return strings.Join(loremList(), " ")
 }
 
 // This examples demonstrates Landscape mode with images.
 func ExampleFpdf_SetAcceptPageBreakFunc() {
	 var y0 float64
	 var crrntCol int
	 loremStr := lorem()
	 pdf := gofpdf.New("L", "mm", "A4", "")
	 const (
		 pageWd = 297.0 // A4 210.0 x 297.0
		 margin = 10.0
		 gutter = 4
		 colNum = 3
		 colWd  = (pageWd - 2*margin - (colNum-1)*gutter) / colNum
	 )
	 setCol := func(col int) {
		 crrntCol = col
		 x := margin + float64(col)*(colWd+gutter)
		 pdf.SetLeftMargin(x)
		 pdf.SetX(x)
	 }
	 pdf.SetHeaderFunc(func() {
		 titleStr := "gofpdf"
		 pdf.SetFont("Helvetica", "B", 48)
		 wd := pdf.GetStringWidth(titleStr) + 6
		 pdf.SetX((pageWd - wd) / 2)
		 pdf.SetTextColor(128, 128, 160)
		 pdf.Write(12, titleStr[:2])
		 pdf.SetTextColor(128, 128, 128)
		 pdf.Write(12, titleStr[2:])
		 pdf.Ln(20)
		 y0 = pdf.GetY()
	 })
	 pdf.SetAcceptPageBreakFunc(func() bool {
		 if crrntCol < colNum-1 {
			 setCol(crrntCol + 1)
			 pdf.SetY(y0)
			 // Start new column, not new page
			 return false
		 }
		 setCol(0)
		 return true
	 })
	 pdf.AddPage()
	 pdf.SetFont("Times", "", 12)
	 for j := 0; j < 20; j++ {
		 if j == 1 {
			 pdf.Image(ImageFile("fpdf.png"), -1, 0, colWd, 0, true, "", 0, "")
		 } else if j == 5 {
			 pdf.Image(ImageFile("golang-gopher.png"),
				 -1, 0, colWd, 0, true, "", 0, "")
		 }
		 pdf.MultiCell(colWd, 5, loremStr, "", "", false)
		 pdf.Ln(-1)
	 }
	 fileStr := Filename("Fpdf_SetAcceptPageBreakFunc_landscape")
	 err := pdf.OutputFileAndClose(fileStr)
	 Summary(err, fileStr)
	 // Output:
	 // Successfully generated pdf/Fpdf_SetAcceptPageBreakFunc_landscape.pdf
 }
 
 
 
 
 // This example generates a PDF document with various page sizes.
 func ExampleFpdf_PageSize() {
	 pdf := gofpdf.NewCustom(&gofpdf.InitType{
		 UnitStr:    "in",
		 Size:       gofpdf.SizeType{Wd: 6, Ht: 6},
		 FontDirStr: FontDir(),
	 })
	 pdf.SetMargins(0.5, 1, 0.5)
	 pdf.SetFont("Times", "", 14)
	 pdf.AddPageFormat("L", gofpdf.SizeType{Wd: 3, Ht: 12})
	 pdf.SetXY(0.5, 1.5)
	 pdf.CellFormat(11, 0.2, "12 in x 3 in", "", 0, "C", false, 0, "")
	 pdf.AddPage() // Default size established in NewCustom()
	 pdf.SetXY(0.5, 3)
	 pdf.CellFormat(5, 0.2, "6 in x 6 in", "", 0, "C", false, 0, "")
	 pdf.AddPageFormat("P", gofpdf.SizeType{Wd: 3, Ht: 12})
	 pdf.SetXY(0.5, 6)
	 pdf.CellFormat(2, 0.2, "3 in x 12 in", "", 0, "C", false, 0, "")
	 for j := 0; j <= 3; j++ {
		 wd, ht, u := pdf.PageSize(j)
		 fmt.Printf("%d: %6.2f %s, %6.2f %s\n", j, wd, u, ht, u)
	 }
	 fileStr := Filename("Fpdf_PageSize")
	 err := pdf.OutputFileAndClose(fileStr)
	 Summary(err, fileStr)
	 // Output:
	 // 0:   6.00 in,   6.00 in
	 // 1:  12.00 in,   3.00 in
	 // 2:   6.00 in,   6.00 in
	 // 3:   3.00 in,  12.00 in
	 // Successfully generated pdf/Fpdf_PageSize.pdf
 }
 
 
 // This example demonstrate Clipped table cells
 func ExampleFpdf_ClipRect() {
	 marginCell := 2. // margin of top/bottom of cell
	 pdf := gofpdf.New("P", "mm", "A4", "")
	 pdf.SetFont("Arial", "", 12)
	 pdf.AddPage()
	 pagew, pageh := pdf.GetPageSize()
	 mleft, mright, _, mbottom := pdf.GetMargins()
 
	 cols := []float64{60, 100, pagew - mleft - mright - 100 - 60}
	 rows := [][]string{}
	 for i := 1; i <= 50; i++ {
		 word := fmt.Sprintf("%d:%s", i, strings.Repeat("A", i%100))
		 rows = append(rows, []string{word, word, word})
	 }
 
	 for _, row := range rows {
		 _, lineHt := pdf.GetFontSize()
		 height := lineHt + marginCell
 
		 x, y := pdf.GetXY()
		 // add a new page if the height of the row doesn't fit on the page
		 if y+height >= pageh-mbottom {
			 pdf.AddPage()
			 x, y = pdf.GetXY()
		 }
		 for i, txt := range row {
			 width := cols[i]
			 pdf.Rect(x, y, width, height, "")
			 pdf.ClipRect(x, y, width, height, false)
			 pdf.Cell(width, height, txt)
			 pdf.ClipEnd()
			 x += width
		 }
		 pdf.Ln(-1)
	 }
	 fileStr := Filename("Fpdf_ClippedTableCells")
	 err := pdf.OutputFileAndClose(fileStr)
	 Summary(err, fileStr)
	 // Output:
	 // Successfully generated pdf/Fpdf_ClippedTableCells.pdf
 }
 
 
 
 // This example demonstrate wrapped table cells
 func ExampleFpdf_Rect() {
	 marginCell := 2. // margin of top/bottom of cell
	 pdf := gofpdf.New("P", "mm", "A4", "")
	 pdf.SetFont("Arial", "", 12)
	 pdf.AddPage()
	 pagew, pageh := pdf.GetPageSize()
	 mleft, mright, _, mbottom := pdf.GetMargins()
 
	 cols := []float64{60, 100, pagew - mleft - mright - 100 - 60}
	 rows := [][]string{}
	 for i := 1; i <= 30; i++ {
		 word := fmt.Sprintf("%d:%s", i, strings.Repeat("A", i%100))
		 rows = append(rows, []string{word, word, word})
	 }
 
	 for _, row := range rows {
		 curx, y := pdf.GetXY()
		 x := curx
 
		 height := 0.
		 _, lineHt := pdf.GetFontSize()
 
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
	 fileStr := Filename("Fpdf_WrappedTableCells")
	 err := pdf.OutputFileAndClose(fileStr)
	 Summary(err, fileStr)
	 // Output:
	 // Successfully generated pdf/Fpdf_WrappedTableCells.pdf
 }

 
 
 
 
  // This example demonstrates the use of characters in the high range of the
 // Windows-1252 code page (gofdpf default). See the example for CellFormat (4)
 // for a way to do this automatically.
 func ExampleFpdf_CellFormat_codepageescape() {
	 pdf := gofpdf.New("P", "mm", "A4", "") // A4 210.0 x 297.0
	 fontSize := 16.0
	 pdf.SetFont("Helvetica", "", fontSize)
	 ht := pdf.PointConvert(fontSize)
	 write := func(str string) {
		 pdf.CellFormat(190, ht, str, "", 1, "C", false, 0, "")
		 pdf.Ln(ht)
	 }
	 pdf.AddPage()
	 htmlStr := `Until gofpdf supports UTF-8 encoded source text, source text needs ` +
		 `to be specified with all special characters escaped to match the code page ` +
		 `layout of the currently selected font. By default, gofdpf uses code page 1252.` +
		 ` See <a href="http://en.wikipedia.org/wiki/Windows-1252">Wikipedia</a> for ` +
		 `a table of this layout.`
	 html := pdf.HTMLBasicNew()
	 html.Write(ht, htmlStr)
	 pdf.Ln(2 * ht)
	 write("Voix ambigu\xeb d'un c\x9cur qui au z\xe9phyr pr\xe9f\xe8re les jattes de kiwi.")
	 write("Falsches \xdcben von Xylophonmusik qu\xe4lt jeden gr\xf6\xdferen Zwerg.")
	 write("Heiz\xf6lr\xfccksto\xdfabd\xe4mpfung")
	 write("For\xe5rsj\xe6vnd\xf8gn / Efter\xe5rsj\xe6vnd\xf8gn")
	 fileStr := Filename("Fpdf_CellFormat_codepageescape")
	 err := pdf.OutputFileAndClose(fileStr)
	 Summary(err, fileStr)
	 // Output:
	 // Successfully generated pdf/Fpdf_CellFormat_codepageescape.pdf
 }
  type fontResourceType struct {
 }

  func (f fontResourceType) Open(name string) (rdr io.Reader, err error) {
	 var buf []byte
	 buf, err = ioutil.ReadFile(FontFile(name))
	 if err == nil {
		 rdr = bytes.NewReader(buf)
		 fmt.Printf("Generalized font loader reading %s\n", name)
	 }
	 return
 }
 
  


  
func main() {
 ExampleFpdf_AddPage()
 //ExampleFpdf_MultiCell()
 //ExampleFpdf_SetLeftMargin ()
 //ExampleFpdf_SplitLines_tables ()
 //ExampleFpdf_CellFormat_tables()
 //ExampleFpdf_SetAcceptPageBreakFunc()
 //ExampleFpdf_PageSize()
 //ExampleFpdf_ClipRect()
 //ExampleFpdf_Rect()
 // No nos interesa : ExampleFpdf_SetLineJoinStyle()
 // +- ExampleFpdf_CellFormat_codepageescape()
 
 // no ExampleFpdf_CellFormat_align()
 // no ExampleFpdf_WriteAligned()
 // no ExampleFpdf_SetKeywords ()
 
 
 // +- ExampleFpdf_SetAlpha () 
 
 // no ExampleFpdf_LinearGradient()
 // +-- ExampleFpdf_ClipText() 
 // no ExampleFpdf_TransformBegin ()
 // +- ExampleFpdf_SplitLines() 
// ExampleFpdf_SVGBasicWrite ()
 
 fmt.Println("Hello world")
}