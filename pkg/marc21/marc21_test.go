package marc21

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestXML(t *testing.T) {

	xmlFile := os.Getenv("TEST_MARCXML_FILE")

	doc, err := LoadXML(xmlFile)

	if err != nil {
		t.Errorf("LoadXML() failed: %q", err)
	}

	out, err := doc.AsXML()
	if err != nil {
		t.Errorf("Collection.AsXML() failed: %q", err)
	} else if out == "" {
		t.Errorf("Collection.AsXML() failed")
	}

}

func TestMARC(t *testing.T) {

	marcFile := os.Getenv("TEST_MARC_FILE")

	fi, err := os.Open(marcFile)
	if err != nil {
		t.Errorf("os.Open() failed: %q", err)
	}
	defer fi.Close()

	for {
		rec, err := ParseNextRecord(fi)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("ParseNextRecord() failed: %q", err)
			continue
		}

		// TODO: at some point this needs to move to specific test case
		// MARC records with specific expected outcomes. "Good" records
		// should have valid values while intentionally "bad" records
		// should have invalid or no values and should be evaluated
		// accordingly.
		rf := rec.RecordFormat()
		if rf != FmtUnknown {

			code, _ := rec.RecordType()
			if code == "" {
				t.Errorf("RecordType() failed")
			}
			code, _ = rec.CharacterCodingScheme()
			if code == "" {
				t.Errorf("CharacterCodingScheme() failed")
			}

			n := rec.RecordFormatName()
			if n == "" {
				t.Errorf("RecordFormatName() failed")
			}

		}
		if rf == Bibliography {
			code, _ := rec.BibliographyMaterialType()
			if code == "" {
				t.Errorf("BibliographyMaterialType() failed")
			}
		}

		out := fmt.Sprint(rec)
		if out == "" {
			t.Errorf("Record.Print() failed")
		}
	}
}

func TestMARC2XML(t *testing.T) {

	marcFile := os.Getenv("TEST_MARC_FILE")

	fi, err := os.Open(marcFile)
	if err != nil {
		t.Errorf("os.Open() failed: %q", err)
	}
	defer fi.Close()

	rec, err := ParseNextRecord(fi)
	if err != nil {
		t.Errorf("ParseNextRecord() failed: %q", err)
	}

	var doc Collection
	doc.Records = append(doc.Records, rec)

	out, err := doc.AsXML()
	if err != nil {
		t.Errorf("Collection.AsXML() failed: %q", err)
	} else if out == "" {
		t.Errorf("Collection.AsXML() failed")
	}
}

func TestXML2PrettyPrint(t *testing.T) {

	xmlFile := os.Getenv("TEST_MARCXML_FILE")

	doc, err := LoadXML(xmlFile)
	if err != nil {
		t.Errorf("LoadXML() failed: %q", err)
	}

	for _, rec := range doc.Records {
		out := fmt.Sprint(rec)
		if out == "" {
			t.Errorf("Print() failed")
		}
	}
}

func TestRecordAsMARC(t *testing.T) {

	marcFile := os.Getenv("TEST_MARC_FILE")

	fi, err := os.Open(marcFile)
	if err != nil {
		t.Errorf("os.Open() failed: %q", err)
	}
	defer fi.Close()

	rec, err := ParseNextRecord(fi)
	if err != nil {
		t.Errorf("ParseNextRecord() failed: %q", err)
	} else {
		_, err := rec.RecordAsMARC()
		if err != nil {
			t.Errorf("RecordAsMARC() failed: %q", err)
		}
	}
}

func TestLeaderPrint(t *testing.T) {

	marcFile := os.Getenv("TEST_MARC_FILE")

	fi, err := os.Open(marcFile)
	if err != nil {
		t.Errorf("os.Open() failed: %q", err)
	}
	defer fi.Close()

	rec, err := ParseNextRecord(fi)
	if err != nil {
		t.Errorf("ParseNextRecord() failed: %q", err)
	} else {
		out := fmt.Sprint(rec.Leader)
		// TODO: this really needs to test against a specific MARC record
		// and compare the output to the specific expected value.
		if out == "" {
			t.Errorf("Sprint(rec.Leader) failed")
		}

	}
}

func TestCustomString(t *testing.T) {
	rec := &Record{}
	rec.AddControlfield("005", "20230224185445.0")
	rec.AddDatafield(NewDatafield("020").SetInd2("0").AddSubfield("a", "test").AddSubfield("c", "9.99"))

	res := rec.CustomString(false, "\n", "|", "|", " ", " ")

	if res != "005     20230224185445.0\n020 |0 |a test |c 9.99 \n" {
		t.Errorf("CustomString returned wrong result\n%s", res)
	}
}
