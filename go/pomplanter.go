package main

import (
	"fmt"
	"log"
	"os"
)

func BuildDiagram(poms []Pom) string {
	s := "@startuml\n\n"
	for _, pom := range poms {
		s = s + "[" + pom.ArtifactId + "]\n"
	}
	s = s + "\n"
	for _, pom := range poms {
		for _, p := range PointedBy(pom, poms) {
			s = s + "[" + p.ArtifactId + "] -DOWN-> [" + pom.ArtifactId + "]\n"
		}
	}
	s = s + "\n@enduml"
	return s
}

func main() {
	log.Println("Started POM Planter ================================================================================================")
	if len(os.Args) < 2 {
		log.Fatalln("Missing POM file name")
		return
	}
	var poms []Pom = make([]Pom, 0)
	poms = ReadTreePOM(os.Args[1], "root", poms)
	diagram := BuildDiagram(poms)
	fmt.Println(diagram)
	log.Println("Finished POM Planter ===============================================================================================")
}
