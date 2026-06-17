package main

import (
	"archive/zip"
	"bytes"
	"embed"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

//go:embed resources/bsu-logo.png
var embeddedResources embed.FS

var reportMonths = []string{
	"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

var daysOfWeek = []string{
	"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
}

func buildReportDocx(parsedData ParsedData, userInfo UserInfo, employeeType string) ([]byte, error) {
	start, err := time.Parse("2006-01-02", parsedData.DateStart)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", parsedData.DateEnd)
	if err != nil {
		return nil, err
	}

	logoBytes, _ := embeddedResources.ReadFile("resources/bsu-logo.png")
	period := formatReportPeriod(start, end)

	documentXML := buildDocumentXML(parsedData, userInfo, employeeType, period)
	headerXML := buildHeaderXML(employeeType, len(logoBytes) > 0)

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := map[string]string{
		"[Content_Types].xml":             contentTypesXML(employeeType, len(logoBytes) > 0),
		"_rels/.rels":                     rootRelsXML(),
		"docProps/app.xml":                appPropsXML(),
		"docProps/core.xml":               corePropsXML(),
		"word/document.xml":               documentXML,
		"word/_rels/document.xml.rels":    documentRelsXML(employeeType, len(logoBytes) > 0),
		"word/styles.xml":                 stylesXML(),
		"word/settings.xml":               settingsXML(),
		"word/fontTable.xml":              fontTableXML(),
		"word/theme/theme1.xml":           themeXML(),
		"word/webSettings.xml":            webSettingsXML(),
		"word/header1.xml":                headerXML,
		"word/_rels/header1.xml.rels":     headerRelsXML(len(logoBytes) > 0),
	}

	for name, content := range files {
		if err := writeZipTextFile(zw, name, content); err != nil {
			return nil, err
		}
	}

	if len(logoBytes) > 0 {
		if err := writeZipBinaryFile(zw, "word/media/bsu-logo.png", logoBytes); err != nil {
			return nil, err
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeZipTextFile(zw *zip.Writer, name, content string) error {
	file, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(content))
	return err
}

func writeZipBinaryFile(zw *zip.Writer, name string, data []byte) error {
	file, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

func buildDocumentXML(parsedData ParsedData, userInfo UserInfo, employeeType, period string) string {
	var body strings.Builder
	body.WriteString(xmlDecl())
	body.WriteString(`<w:document xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas" xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing" xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:w10="urn:schemas-microsoft-com:office:word" xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml" xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml" xmlns:wpg="http://schemas.microsoft.com/office/word/2010/wordprocessingGroup" xmlns:wpi="http://schemas.microsoft.com/office/word/2010/wordprocessingInk" xmlns:wne="http://schemas.microsoft.com/office/2006/wordml" xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape" mc:Ignorable="w14 wp14"><w:body>`)

	body.WriteString(titleBoxXML(employeeType))
	body.WriteString(blankParagraphXML())
	body.WriteString(centerParagraphXML(fmt.Sprintf("For the Period Covered: %s", period), "Calibri", 20, true, false))
	body.WriteString(namePositionTableXML(userInfo))
	body.WriteString(leftParagraphXML("(Family Name, First Name, Middle Name)", "Calibri", 18, false, true, 720))
	body.WriteString(leftParagraphXML(fmt.Sprintf("College/Office Assignment: %s", fallbackUnderscore(userInfo.Office)), "Calibri", 20, true, true, 0))
	body.WriteString(blankParagraphXML())
	body.WriteString(mainTaskTableXML(parsedData))
	body.WriteString(blankParagraphXML())
	body.WriteString(signatureTableXML(userInfo))

	body.WriteString(`<w:sectPr>`)
	if employeeType != "external" {
		body.WriteString(`<w:headerReference w:type="default" r:id="rId4"/>`)
	}
	body.WriteString(`<w:pgSz w:w="11906" w:h="16838"/><w:pgMar w:top="1440" w:right="720" w:bottom="720" w:left="720" w:header="360" w:footer="720" w:gutter="0"/><w:cols w:space="720"/><w:docGrid w:linePitch="360"/></w:sectPr>`)
	body.WriteString(`</w:body></w:document>`)
	return body.String()
}

func buildHeaderXML(employeeType string, withLogo bool) string {
	if employeeType == "external" {
		return xmlDecl() + `<w:hdr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"/>`
	}

	var b strings.Builder
	b.WriteString(xmlDecl())
	b.WriteString(`<w:hdr xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">`)
	b.WriteString(`<w:tbl><w:tblPr><w:tblW w:w="10588" w:type="dxa"/><w:tblBorders><w:top w:val="nil"/><w:left w:val="nil"/><w:bottom w:val="nil"/><w:right w:val="nil"/><w:insideH w:val="nil"/><w:insideV w:val="nil"/></w:tblBorders></w:tblPr><w:tblGrid><w:gridCol w:w="1600"/><w:gridCol w:w="6500"/><w:gridCol w:w="400"/><w:gridCol w:w="2088"/></w:tblGrid><w:tr>`)
	b.WriteString(`<w:tc><w:tcPr><w:tcW w:w="1600" w:type="dxa"/><w:vAlign w:val="center"/></w:tcPr>`)
	if withLogo {
		b.WriteString(`<w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:drawing><wp:inline distT="0" distB="0" distL="0" distR="0"><wp:extent cx="666750" cy="666750"/><wp:docPr id="1" name="BSU Logo"/><a:graphic><a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture"><pic:pic><pic:nvPicPr><pic:cNvPr id="0" name="bsu-logo.png"/><pic:cNvPicPr/></pic:nvPicPr><pic:blipFill><a:blip r:embed="rId1"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill><pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="666750" cy="666750"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr></pic:pic></a:graphicData></a:graphic></wp:inline></w:drawing></w:r></w:p>`)
	} else {
		b.WriteString(blankParagraphXML())
	}
	b.WriteString(`</w:tc>`)
	b.WriteString(`<w:tc><w:tcPr><w:tcW w:w="6500" w:type="dxa"/><w:vAlign w:val="center"/></w:tcPr>`)
	b.WriteString(centerParagraphXML("Republic of the Philippines", "Calibri", 18, false, false))
	b.WriteString(centerParagraphXML("Benguet State University", "Old English Text MT", 36, false, false))
	b.WriteString(centerParagraphXML("La Trinidad, Benguet, 2601", "Calibri", 18, false, false))
	b.WriteString(`</w:tc>`)
	b.WriteString(`<w:tc><w:tcPr><w:tcW w:w="400" w:type="dxa"/></w:tcPr>`)
	b.WriteString(blankParagraphXML())
	b.WriteString(`</w:tc>`)
	b.WriteString(`<w:tc><w:tcPr><w:tcW w:w="2088" w:type="dxa"/><w:vAlign w:val="center"/></w:tcPr><w:tbl><w:tblPr><w:tblW w:w="2048" w:type="dxa"/><w:tblBorders><w:top w:val="dotted" w:sz="4" w:color="000000"/><w:left w:val="dotted" w:sz="4" w:color="000000"/><w:bottom w:val="dotted" w:sz="4" w:color="000000"/><w:right w:val="dotted" w:sz="4" w:color="000000"/></w:tblBorders></w:tblPr><w:tr><w:tc><w:tcPr><w:tcMar><w:top w:w="60" w:type="dxa"/><w:bottom w:w="60" w:type="dxa"/><w:left w:w="80" w:type="dxa"/><w:right w:w="80" w:type="dxa"/></w:tcMar></w:tcPr>`)
	b.WriteString(centerParagraphXML("AR-HRMO form 1.0", "Brush Script MT", 18, false, true))
	b.WriteString(`</w:tc></w:tr></w:tbl></w:tc>`)
	b.WriteString(`</w:tr></w:tbl></w:hdr>`)
	return b.String()
}

func titleBoxXML(employeeType string) string {
	title := map[string]string{
		"casual":    "Accomplishment Report for casual personnel",
		"plantilla": "Accomplishment Report for plantilla personnel",
		"cos":       "Accomplishment Report for COS personnel",
		"external":  "Accomplishment Report",
	}[employeeType]
	if title == "" {
		title = "Accomplishment Report for casual personnel"
	}

	return `<w:tbl><w:tblPr><w:tblW w:w="10588" w:type="dxa"/><w:tblBorders><w:top w:val="nil"/><w:left w:val="nil"/><w:bottom w:val="nil"/><w:right w:val="nil"/><w:insideH w:val="nil"/><w:insideV w:val="nil"/></w:tblBorders></w:tblPr><w:tblGrid><w:gridCol w:w="2300"/><w:gridCol w:w="5988"/><w:gridCol w:w="2300"/></w:tblGrid><w:tr><w:tc><w:tcPr><w:tcW w:w="2300" w:type="dxa"/></w:tcPr><w:p/></w:tc><w:tc><w:tcPr><w:tcW w:w="5988" w:type="dxa"/><w:tcBorders><w:top w:val="single" w:sz="6" w:color="000000"/><w:left w:val="single" w:sz="6" w:color="000000"/><w:bottom w:val="single" w:sz="6" w:color="000000"/><w:right w:val="single" w:sz="6" w:color="000000"/></w:tcBorders></w:tcPr>` +
		centerParagraphXML(title, "Copperplate Gothic Bold", 20, false, false) +
		`</w:tc><w:tc><w:tcPr><w:tcW w:w="2300" w:type="dxa"/></w:tcPr><w:p/></w:tc></w:tr></w:tbl>`
}

func namePositionTableXML(userInfo UserInfo) string {
	name := fallbackUnderscore(userInfo.ReporterName)
	position := fallbackUnderscore(userInfo.Position)

	left := paragraphWithRuns("left", []runSpec{
		{Text: "NAME:  ", Font: "Calibri", Size: 20},
		{Text: name, Font: "Calibri", Size: 20, Bold: true},
	}, 0)
	right := paragraphWithRuns("left", []runSpec{
		{Text: "Position:  ", Font: "Calibri", Size: 20},
		{Text: position, Font: "Calibri", Size: 20},
	}, 0)

	return `<w:tbl><w:tblPr><w:tblW w:w="10588" w:type="dxa"/><w:tblBorders><w:top w:val="nil"/><w:left w:val="nil"/><w:bottom w:val="nil"/><w:right w:val="nil"/><w:insideH w:val="nil"/><w:insideV w:val="nil"/></w:tblBorders></w:tblPr><w:tblGrid><w:gridCol w:w="5294"/><w:gridCol w:w="5294"/></w:tblGrid><w:tr><w:tc><w:tcPr><w:tcW w:w="5294" w:type="dxa"/></w:tcPr>` +
		left +
		`</w:tc><w:tc><w:tcPr><w:tcW w:w="5294" w:type="dxa"/></w:tcPr>` +
		right +
		`</w:tc></w:tr></w:tbl>`
}

func mainTaskTableXML(parsedData ParsedData) string {
	entryMap := map[string]ReportEntry{}
	for _, entry := range parsedData.Entries {
		entryMap[entry.Date] = entry
	}

	start, _ := time.Parse("2006-01-02", parsedData.DateStart)
	end, _ := time.Parse("2006-01-02", parsedData.DateEnd)

	var rows strings.Builder
	rows.WriteString(`<w:tbl><w:tblPr><w:tblW w:w="10588" w:type="dxa"/><w:tblBorders><w:top w:val="single" w:sz="4" w:color="000000"/><w:left w:val="single" w:sz="4" w:color="000000"/><w:bottom w:val="single" w:sz="4" w:color="000000"/><w:right w:val="single" w:sz="4" w:color="000000"/><w:insideH w:val="single" w:sz="4" w:color="000000"/><w:insideV w:val="single" w:sz="4" w:color="000000"/></w:tblBorders></w:tblPr><w:tblGrid><w:gridCol w:w="729"/><w:gridCol w:w="9859"/></w:tblGrid>`)

	for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {
		entry := entryMap[day.Format("2006-01-02")]
		taskParagraphs := []string{}

		if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
			taskParagraphs = append(taskParagraphs, leftParagraphXML(daysOfWeek[int(day.Weekday())], "Arial Narrow", 20, false, true, 0))
		} else if len(entry.Tasks) > 0 {
			for _, task := range entry.Tasks {
				taskParagraphs = append(taskParagraphs, leftParagraphXML(task, "Arial Narrow", 20, false, false, 0))
			}
		} else {
			taskParagraphs = append(taskParagraphs, blankParagraphXML())
		}

		rows.WriteString(`<w:tr><w:tc><w:tcPr><w:tcW w:w="729" w:type="dxa"/><w:vAlign w:val="center"/></w:tcPr>`)
		rows.WriteString(centerParagraphXML(fmt.Sprintf("%d", day.Day()), "Cambria", 20, false, false))
		rows.WriteString(`</w:tc><w:tc><w:tcPr><w:tcW w:w="9859" w:type="dxa"/><w:vAlign w:val="center"/></w:tcPr>`)
		for _, paragraph := range taskParagraphs {
			rows.WriteString(paragraph)
		}
		rows.WriteString(`</w:tc></w:tr>`)
	}

	rows.WriteString(`</w:tbl>`)
	return rows.String()
}

func signatureTableXML(userInfo UserInfo) string {
	reporter := fallbackUnderscore(firstNonEmpty(userInfo.ReporterNameSig, userInfo.ReporterName))
	supervisor := fallbackUnderscore(userInfo.SupervisorName)

	leftLabel := centerParagraphXML("Reported by:", "Cambria", 20, false, false)
	rightLabel := centerParagraphXML("Certified Correct:", "Cambria", 20, false, false)
	leftName := centerParagraphXML(reporter, "Calibri", 20, true, false) + centerParagraphXML("Signature Over Printed Name", "Calibri", 20, false, false)
	rightName := centerParagraphXML(supervisor, "Calibri", 20, true, false) +
		centerParagraphXML(userInfo.SupervisorPos1, "Calibri", 20, false, false) +
		centerParagraphXML(userInfo.SupervisorPos2, "Calibri", 20, false, false)

	return `<w:tbl><w:tblPr><w:tblW w:w="10588" w:type="dxa"/><w:tblBorders><w:top w:val="nil"/><w:left w:val="nil"/><w:bottom w:val="nil"/><w:right w:val="nil"/><w:insideH w:val="nil"/><w:insideV w:val="nil"/></w:tblBorders></w:tblPr><w:tblGrid><w:gridCol w:w="5294"/><w:gridCol w:w="5294"/></w:tblGrid><w:tr><w:tc><w:tcPr><w:tcW w:w="5294" w:type="dxa"/></w:tcPr>` +
		leftLabel +
		`</w:tc><w:tc><w:tcPr><w:tcW w:w="5294" w:type="dxa"/></w:tcPr>` +
		rightLabel +
		`</w:tc></w:tr><w:tr><w:tc><w:tcPr><w:tcW w:w="5294" w:type="dxa"/></w:tcPr><w:p><w:pPr><w:spacing w:after="600"/></w:pPr></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="5294" w:type="dxa"/></w:tcPr><w:p><w:pPr><w:spacing w:after="600"/></w:pPr></w:p></w:tc></w:tr><w:tr><w:tc><w:tcPr><w:tcW w:w="5294" w:type="dxa"/></w:tcPr>` +
		leftName +
		`</w:tc><w:tc><w:tcPr><w:tcW w:w="5294" w:type="dxa"/></w:tcPr>` +
		rightName +
		`</w:tc></w:tr></w:tbl>`
}

type runSpec struct {
	Text   string
	Font   string
	Size   int
	Bold   bool
	Italic bool
}

func paragraphWithRuns(alignment string, runs []runSpec, indentLeft int) string {
	var b strings.Builder
	b.WriteString(`<w:p><w:pPr>`)
	if alignment == "center" {
		b.WriteString(`<w:jc w:val="center"/>`)
	}
	if indentLeft > 0 {
		b.WriteString(fmt.Sprintf(`<w:ind w:left="%d"/>`, indentLeft))
	}
	b.WriteString(`</w:pPr>`)
	for _, run := range runs {
		b.WriteString(runXML(run.Text, run.Font, run.Size, run.Bold, run.Italic))
	}
	b.WriteString(`</w:p>`)
	return b.String()
}

func centerParagraphXML(text, font string, size int, bold, italic bool) string {
	return paragraphWithRuns("center", []runSpec{{Text: text, Font: font, Size: size, Bold: bold, Italic: italic}}, 0)
}

func leftParagraphXML(text, font string, size int, bold, italic bool, indentLeft int) string {
	return paragraphWithRuns("left", []runSpec{{Text: text, Font: font, Size: size, Bold: bold, Italic: italic}}, indentLeft)
}

func blankParagraphXML() string {
	return `<w:p/>`
}

func runXML(text, font string, size int, bold, italic bool) string {
	if font == "" {
		font = "Calibri"
	}
	if size == 0 {
		size = 20
	}

	var props strings.Builder
	props.WriteString(`<w:rPr>`)
	props.WriteString(fmt.Sprintf(`<w:rFonts w:ascii="%s" w:hAnsi="%s"/>`, xmlEscape(font), xmlEscape(font)))
	props.WriteString(fmt.Sprintf(`<w:sz w:val="%d"/><w:szCs w:val="%d"/>`, size, size))
	if bold {
		props.WriteString(`<w:b/>`)
	}
	if italic {
		props.WriteString(`<w:i/>`)
	}
	props.WriteString(`</w:rPr>`)

	preserve := ""
	if strings.HasPrefix(text, " ") || strings.HasSuffix(text, " ") {
		preserve = ` xml:space="preserve"`
	}
	return `<w:r>` + props.String() + `<w:t` + preserve + `>` + xmlEscape(text) + `</w:t></w:r>`
}

func xmlDecl() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`
}

func xmlEscape(value string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(value))
	return buf.String()
}

func fallbackUnderscore(value string) string {
	if strings.TrimSpace(value) == "" {
		return "______________________________"
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func formatReportPeriod(start, end time.Time) string {
	startMonth := reportMonths[int(start.Month())-1]
	endMonth := reportMonths[int(end.Month())-1]
	if start.Month() == end.Month() {
		return fmt.Sprintf("%s %d-%d, %d", startMonth, start.Day(), end.Day(), end.Year())
	}
	return fmt.Sprintf("%s %d - %s %d, %d", startMonth, start.Day(), endMonth, end.Day(), end.Year())
}

func contentTypesXML(employeeType string, withLogo bool) string {
	parts := []string{
		xmlDecl(),
		`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`,
		`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>`,
		`<Default Extension="xml" ContentType="application/xml"/>`,
		`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>`,
		`<Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>`,
		`<Override PartName="/word/settings.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.settings+xml"/>`,
		`<Override PartName="/word/fontTable.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.fontTable+xml"/>`,
		`<Override PartName="/word/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>`,
		`<Override PartName="/word/webSettings.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.webSettings+xml"/>`,
		`<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>`,
		`<Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>`,
	}
	if employeeType != "external" {
		parts = append(parts, `<Override PartName="/word/header1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/>`)
	}
	if withLogo {
		parts = append(parts, `<Default Extension="png" ContentType="image/png"/>`)
	}
	parts = append(parts, `</Types>`)
	return strings.Join(parts, "")
}

func rootRelsXML() string {
	return xmlDecl() + `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/><Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/></Relationships>`
}

func documentRelsXML(employeeType string, _ bool) string {
	relationships := []string{
		xmlDecl(),
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`,
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>`,
		`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings" Target="settings.xml"/>`,
		`<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/fontTable" Target="fontTable.xml"/>`,
		`<Relationship Id="rId5" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="theme/theme1.xml"/>`,
	}
	if employeeType != "external" {
		relationships = append(relationships, `<Relationship Id="rId4" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/header" Target="header1.xml"/>`)
	}
	relationships = append(relationships, `</Relationships>`)
	return strings.Join(relationships, "")
}

func headerRelsXML(withLogo bool) string {
	if !withLogo {
		return xmlDecl() + `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`
	}
	return xmlDecl() + `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/bsu-logo.png"/></Relationships>`
}

func stylesXML() string {
	return xmlDecl() + `<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:docDefaults><w:rPrDefault><w:rPr><w:rFonts w:ascii="Calibri" w:hAnsi="Calibri"/><w:sz w:val="20"/><w:szCs w:val="20"/></w:rPr></w:rPrDefault></w:docDefaults></w:styles>`
}

func settingsXML() string {
	return xmlDecl() + `<w:settings xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:zoom w:percent="100"/><w:defaultTabStop w:val="720"/></w:settings>`
}

func fontTableXML() string {
	return xmlDecl() + `<w:fonts xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:font w:name="Calibri"/><w:font w:name="Cambria"/><w:font w:name="Arial Narrow"/><w:font w:name="Copperplate Gothic Bold"/><w:font w:name="Brush Script MT"/><w:font w:name="Old English Text MT"/></w:fonts>`
}

func themeXML() string {
	return xmlDecl() + `<a:theme xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" name="Office Theme"><a:themeElements><a:clrScheme name="Office"><a:dk1><a:sysClr val="windowText" lastClr="000000"/></a:dk1><a:lt1><a:sysClr val="window" lastClr="FFFFFF"/></a:lt1><a:accent1><a:srgbClr val="4F81BD"/></a:accent1></a:clrScheme><a:fontScheme name="Office"><a:majorFont><a:latin typeface="Calibri"/></a:majorFont><a:minorFont><a:latin typeface="Calibri"/></a:minorFont></a:fontScheme><a:fmtScheme name="Office"/></a:themeElements></a:theme>`
}

func webSettingsXML() string {
	return xmlDecl() + `<w:webSettings xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"/>`
}

func appPropsXML() string {
	return xmlDecl() + `<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes"><Application>Cole</Application></Properties>`
}

func corePropsXML() string {
	now := time.Now().UTC().Format(time.RFC3339)
	return xmlDecl() + `<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><dc:title>Cole Report</dc:title><dc:creator>Cole</dc:creator><cp:lastModifiedBy>Cole</cp:lastModifiedBy><dcterms:created xsi:type="dcterms:W3CDTF">` + now + `</dcterms:created><dcterms:modified xsi:type="dcterms:W3CDTF">` + now + `</dcterms:modified></cp:coreProperties>`
}
