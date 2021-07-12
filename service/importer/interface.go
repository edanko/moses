package importer

import "github.com/edanko/gen2dxf/pkg/wcog"

type UseCase interface {
	ImportGen(f string, wcog *wcog.WCOG) error
	ImportCsv(f string) error
	ImportTxt(f string) error
}
