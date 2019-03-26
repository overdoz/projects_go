package xmlutil

import (
	"fmt"
	"github.com/beevik/etree"
	"strings"
	"unicode"
)
type VDPdata struct {
	Action string
	Subaction string
	VIN string
}


// vdpdata models.VDPdata
func (vdpdata VDPdata) GenerateXML(result string) string {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8" standalone="yes"`)
	doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)

	vdpRes := doc.CreateElement("ns1:vdpResponse")
	header := vdpRes.CreateElement("ns1:header")

	action := header.CreateElement("ns1:action")
	action.CreateAttr("ns1:name", vdpdata.Action)

	subaction := action.CreateElement("ns1:subaction")
	subaction.CreateAttr("ns1:name", vdpdata.Subaction)

	eeSnapshot := header.CreateElement("ns1:eeSnapshot")
	eesUid := eeSnapshot.CreateElement("ns1:eesUid")
	eesUid.CreateText(vdpdata.VIN+";2014-10-24T12:29:49.650+0200")
	eesVersion := eeSnapshot.CreateElement("ns1:eesVersion")
	eesVersion.CreateText("2.2.0")

	trackingID := header.CreateElement("ns1:trackingID")
	trackingID.CreateText(vdpdata.VIN+"_Tassi2")
	sourceID := header.CreateElement("ns1:trackingID")
	sourceID.CreateText("Test")

	responseState := vdpRes.CreateElement("ns1:responseState")
	responseState.CreateAttr("ns1:state", "0")
	messageItem := responseState.CreateElement("ns1:messageItem")
	messageSource := messageItem.CreateElement("ns1:messageSource")
	messageSource.CreateText("DDA")

	requestHeader := vdpRes.CreateElement("ns1:requestHeader")

	action2 := requestHeader.CreateElement("ns1:action")
	action2.CreateAttr("ns1:name", "DDA")
	subaction2 := action2.CreateElement("ns1:subaction")
	subaction2.CreateAttr("ns1:name", "RAR")

	eeSnapshot2 := requestHeader.CreateElement("ns1:eeSnapshot")
	eesUid2 := eeSnapshot2.CreateElement("ns1:eesUid")
	eesUid2.CreateText(vdpdata.VIN)
	eesVersion2 := eeSnapshot2.CreateElement("ns1:eesVersion")
	eesVersion2.CreateText("2.2.0")

	trackingID2 := requestHeader.CreateElement("ns1:trackingID")
	trackingID2.CreateText("WDDVP9AB1EJ001221_Tassi2")
	sourceID2 := requestHeader.CreateElement("ns1:trackingID")
	sourceID2.CreateText("Test")
	market := requestHeader.CreateElement("ns1:market")
	market.CreateText("Test")
	tenant := requestHeader.CreateElement("ns1:tenant")
	tenant.CreateText("Test")
	user := requestHeader.CreateElement("ns1:user")
	user.CreateText("Test")

	language := requestHeader.CreateElement("ns1:language")
	language.CreateText("Test")

	payload := vdpRes.CreateElement("ns1:payload")

	data := etree.NewDocument()

	formattedXML := []byte(formatString(result))
	err := data.ReadFromBytes(formattedXML)
	if err != nil {
		panic(err)
	}

	payload.AddChild(data)

	// create temporary XML to formate
	tempExport, err := doc.WriteToString()
	if err != nil {
		panic(err)
	}

	// create final XML to export
	export := etree.NewDocument()


	//parse temporary XML to fix syntax
	err = export.ReadFromString(deleteEmptyTags(tempExport))
	if err != nil {
		panic(err)
	}

	//export final XML
	err = export.WriteToFile("./EES2.xml")
	if err != nil {
		fmt.Println(err.Error())
	}

	return deleteEmptyTags(tempExport)
}

func parseData(input string) (string, string, string) {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	stringArray := strings.FieldsFunc(input, f)

	trackingID := ""
	action := ""
	subaction := ""

	for i, v := range stringArray {
		if v == "jobID" {
			if trackingID == "" {
				trackingID = stringArray[i+1]
			}
		}
		if v == "jobID" {
			if action == "" {
				action = stringArray[i+1]
			}
		}
		if v == "jobID" {
			if subaction == "" {
				subaction = stringArray[i+1]
			}
		}
	}
	return trackingID, action, subaction

}

func formatString(input string) string {
	// delete first line of received XML
	subs := strings.SplitAfterN(input, "?>", 2)
	return subs[1]

}

func deleteEmptyTags(input string) string {
	emptyTags := strings.Replace(input, "<>", "", -1)
	emptyCloseTags := strings.Replace(emptyTags, "</>", "", -1)
	return emptyCloseTags
}
