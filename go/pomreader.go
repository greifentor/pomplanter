package main

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type Pom struct {
	XMLName      xml.Name     `xml:"project"`
	ArtifactId   string       `xml:"artifactId"`
	GroupId      string       `xml:"groupId"`
	Version      string       `xml:"version"`
	Properties   Properties   `xml:"properties"`
	Dependencies Dependencies `xml:"dependencies"`
	Modules      Modules      `xml:"modules"`
	ModuleName   string
}

type Properties map[string]string

type Dependencies struct {
	XMLName    xml.Name     `xml:"dependencies"`
	Dependency []Dependency `xml:"dependency"`
}

type Dependency struct {
	XMLName    xml.Name `xml:"dependency"`
	ArtifactId string   `xml:"artifactId"`
	GroupId    string   `xml:"groupId"`
	Scope      string   `xml:"scope"`
	Version    string   `xml:"version"`
}

type Modules struct {
	XMLName xml.Name `xml:"modules"`
	Module  []string `xml:"module"`
}

func (i *Properties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*i = make(Properties)
	for {
		tok, err := d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if se, ok := tok.(xml.StartElement); ok {
			tok, err = d.Token()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			if data, ok := tok.(xml.CharData); ok {
				(*i)[se.Name.Local] = string(data)
			}
		}
	}
}

func toValue(element string, pom Pom) string {
	element = strings.ReplaceAll(element, "${project.artifactId}", pom.ArtifactId)
	element = strings.ReplaceAll(element, "${project.version}", pom.Version)
	for key, e := range pom.Properties {
		element = strings.ReplaceAll(element, "${"+key+"}", e)
	}
	return element
}

func ReadPOM(fileName string) *Pom {
	pom := new(Pom)
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	data := string(content)
	if err := xml.Unmarshal([]byte(data), pom); err != nil {
		panic(err)
	}
	return pom
}

func ReadTreePOM(fileName string, moduleName string, poms []Pom) []Pom {
	log.Println("Reading:", fileName)
	pom := ReadPOM(fileName)
	pom.ModuleName = moduleName
	poms = append(poms, *pom)
	fileName = strings.ReplaceAll(fileName, "\\pom.xml", "")
	for _, module := range pom.Modules.Module {
		poms = ReadTreePOM(fileName+"\\"+module+"\\pom.xml", module, poms)
	}
	return poms
}

func ToMap(poms []Pom) map[string]Pom {
	m := make(map[string]Pom)
	for _, pom := range poms {
		m[pom.ArtifactId] = pom
	}
	return m
}

func PointedBy(pom Pom, poms []Pom) []Pom {
	r := make([]Pom, 0)
	for _, p := range poms {
		if isPointingTo(pom, p) {
			r = append(r, p)
		}
	}
	return r
}

func isPointingTo(pom Pom, pomToCheck Pom) bool {
	for _, dependency := range pomToCheck.Dependencies.Dependency {
		if (dependency.ArtifactId == pom.ArtifactId) && (dependency.Scope != "test") {
			// && (dependency.GroupId == pom.GroupId) {
			return true
		}
	}
	return false
}

func LogPOM(pom Pom) {
	log.Println("artifactId:", pom.ArtifactId)
	log.Println("version:   ", pom.Version)
	log.Println("properties:")
	for key, element := range pom.Properties {
		log.Println("           ", key, "->", toValue(element, pom))
	}
	log.Println("dependencies:")
	for i, d := range pom.Dependencies.Dependency {
		log.Println("           ", i, "->", d.ArtifactId, d.GroupId, toValue(d.Version, pom))
	}
	log.Println("modules:")
	for i, m := range pom.Modules.Module {
		log.Println("           ", i, "->", m)
	}
}
