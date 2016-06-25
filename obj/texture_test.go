package obj

import (
	"bytes"
	"testing"
)

var tNullIndex = int64(-1)

var textureReadTests = []struct {
	Items   stringList
	Error   string
	Texture TextureCoord
}{
	{stringList{"1", "1", "1" /*---------------------*/}, "" /*-------------------------------*/, TextureCoord{tNullIndex, 1, 1, 1}},
	{stringList{"1", "1" /*--------------------------*/}, "" /*-------------------------------*/, TextureCoord{tNullIndex, 1, 1, 0}},
	{stringList{"1" /*------------------------------*/}, "TextureCoord: item length is incorrect" /**/, TextureCoord{tNullIndex, 0, 0, 0}},
	{stringList{"1.000", "-1.000", "-1.000" /*-------*/}, "" /*-------------------------------*/, TextureCoord{tNullIndex, 1, -1, -1}},
	{stringList{"0.999", "-1.000", "-1.001" /*-------*/}, "" /*-------------------------------*/, TextureCoord{tNullIndex, 0.999, -1, -1.001}},
	{stringList{"x", "-1.000", "-1.001" /*-----------*/}, "TextureCoord: unable to parse U coordinate" /*---*/, TextureCoord{tNullIndex, 0, 0, 0}},
	{stringList{"1.000", "y", "-1.001" /*------------*/}, "TextureCoord: unable to parse V coordinate" /*---*/, TextureCoord{tNullIndex, 1, 0, 0}},
	{stringList{"1.000", "1", "z" /*-----------------*/}, "TextureCoord: unable to parse W coordinate" /*---*/, TextureCoord{tNullIndex, 1, 1, 0}},
}

func TestReadTexture(t *testing.T) {

	for _, test := range textureReadTests {
		n, err := parseTextCoord(test.Items.ToByteList())

		failed := false

		if test.Error == "" && err != nil {
			failed = true
		} else if err != nil && test.Error != err.Error() {
			failed = true
		}

		if n.U != test.Texture.U || n.V != test.Texture.V || n.W != test.Texture.W {
			failed = true
		}

		if failed {
			t.Errorf("parseTextCoord(%s) => %v, '%v', expected %v, '%v'", test.Items, n, err, test.Texture, test.Error)
		}
	}
}

var textureWriteTests = []struct {
	Texture TextureCoord
	Output  string
	Error   string
}{
	{TextureCoord{tNullIndex, 1, 1, 1}, "1.000 1.000 1.000", ""},
	{TextureCoord{tNullIndex, -1, 1, 1}, "-1.000 1.000 1.000", ""},
	{TextureCoord{tNullIndex, -1.001, 0.999, 1}, "-1.001 0.999 1.000", ""},
}

func TestWriteTexture(t *testing.T) {

	for _, test := range textureWriteTests {
		var buf bytes.Buffer
		err := writeTextCoord(&test.Texture, &buf)

		failed := false

		body := string(buf.Bytes())
		if test.Output != body {
			failed = true
		}

		if test.Error == "" && err != nil {
			failed = true
		} else if err != nil && test.Error != err.Error() {
			failed = true
		}

		if failed {
			t.Errorf("writeTextCoord(%v, wr) => '%v', '%v', expected '%v', '%v'",
				test.Texture, body, err, test.Output, test.Error)
		}
	}

}
