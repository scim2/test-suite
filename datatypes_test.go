package testsuite

import (
	"fmt"
	"testing"
)

func ExampleIsDateTime() {
	fmt.Println(IsDateTime("2008-01-23T04:56:22Z"))
	// Output:
	// true
}

func TestIsBase64(t *testing.T) {
	for _, example := range []string{
		"",
		"Zg==",
		"Zm9v",
		"Zm9vYg==",
		"Zm9vYmE=",
		"Zm9vYmFy",
		"MY======", // base 32
		"MZXQ====", // base 32
	} {
		if !IsBase64(example) {
			t.Error(example)
		}
	}
}

var exampleSCIMResource = `{
  "schemas": [
      "urn:ietf:params:scim:schemas:core:2.0:User",
      "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"
  ],
  "id": "2819c223-7f76-453a-413861904646",
  "externalId": "701984",
  "userName": "bjensen@example.com",
  "name": {
    "formatted": "Ms. Barbara J Jensen, III",
    "familyName": "Jensen",
    "givenName": "Barbara",
    "middleName": "Jane",
    "honorificPrefix": "Ms.",
    "honorificSuffix": "III"
  },
  "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": {
    "employeeNumber": "701984",
    "costCenter": "4130"
  },
  "meta": {
    "resourceType": "User",
    "created": "2010-01-23T04:56:22Z",
    "lastModified": "2011-05-13T04:42:34Z",
    "version": "W\/\"3694e05e9dff591\"",
    "location": "https://example.com/v2/Users/2819c223-7f76-453a-413861904646"
  }
}`

func TestIsObject(t *testing.T) {
	if !IsObject(exampleSCIMResource) {
		t.Error()
	}
}
