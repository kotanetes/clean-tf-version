package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

func main() {

	region := flag.String("region", "us-east-1", "region")

	// identity version from locals.tf
	versions := readTfFiles("./master/template/locals.tf")

	entries, err := os.ReadDir("./" + *region)
	if err != nil {
		log.Fatal(err)
	}

	var res []TenantResult
	for _, e := range entries {

		if e.Type().IsDir() {

			locPath := "./" + *region + "/" + e.Name()

			fmt.Println(locPath)

			// get active release from tfvars.json
			tfJson := readTFVarsFromJson(locPath + "/local.tfvars.json")

			res = append(res, appendVersionForTenant(*region, e.Name(), versions, tfJson))
		}
	}

	file, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(*region+"-tenant-version.json", file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// appendVersionForTenant ...
func appendVersionForTenant(region string, tenant string, versionMap map[string]string, tfJson TFModelJson) TenantResult {

	tr := TenantResult{}

	tr.Region = region
	tr.Tenant = tenant

	v := versionMap[tfJson.SoftwareRelease]
	if val, ok := tfJson.ConfigVersion["messaging-services"]; ok {
		tr.Version = val
	} else {
		tr.Version = v
	}

	return tr
}

// readTfFiles ...
func readTfFiles(filename string) map[string]string {
	parser := hclparse.NewParser()
	f, parseDiags := parser.ParseHCLFile(filename)
	if parseDiags.HasErrors() {
		log.Fatal(parseDiags.Error())
	}

	versions := make(map[string]string, 0)
	var tfModel Local

	for _, v := range f.BlocksAtPos(hcl.Pos{1, 1, 1}) {

		decodeDiags := gohcl.DecodeBody(v.Body, nil, &tfModel)
		if decodeDiags.HasErrors() {
			log.Fatal(decodeDiags.Error())
		}

		hclAttr := tfModel.SoftwareReleases.(*hcl.Attribute)

		val, _ := hclAttr.Expr.Value(nil)

		d := val.AsValueMap()["2023.4"].AsValueMap()["config_versions"].AsValueMap()["messaging-services"].GoString()
		d2 := val.AsValueMap()["2023.5"].AsValueMap()["config_versions"].AsValueMap()["messaging-services"].GoString()

		re := regexp.MustCompile(`"[^"]+"`)

		versions["2023.4"] = strings.TrimRight(strings.TrimLeft(re.FindString(d), `"`), `"`)
		versions["2023.5"] = strings.TrimRight(strings.TrimLeft(re.FindString(d2), `"`), `"`)
	}

	return versions
}

func readTFVarsFromJson(filepath string) TFModelJson {

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("error reading json file, ", err.Error())
	}

	var tfJson TFModelJson

	err = json.Unmarshal(data, &tfJson)
	if err != nil {
		log.Fatalf("error unmarshalling json data, ", err.Error())
	}

	return tfJson

}
