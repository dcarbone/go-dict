package dict_test

import (
	"testing"

	"github.com/dcarbone/go-dict"
)

func TestUnmarshalJSON(t *testing.T) {
	const wellDocumentedJson = `{
    "definitelyString": false,
    "docsSayInt": "200",
    "billIsSureItsABool": "true"
}`

	d, err := dict.UnmarshalJSON([]byte(wellDocumentedJson))
	if err != nil {
		t.Log(err.Error())
		t.Fail()
		return
	}

	// TODO: gosh this is lazy :D
	d.MustGetString("definitelyString")
	d.MustGetInt("docsSayInt")
	d.MustGetBool("billIsSureItsABool")
}
