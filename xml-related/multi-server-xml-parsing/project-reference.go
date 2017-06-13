package xmlresources

import "time"

type Project struct {
	ProjectName  *ProjectName
	AppVer       string
	XMLVer       string
	DbExportDate string
	DbExportTime string
	Areas        []*Area   `xml:">Area"`
	Timestamp    time.Time `xml:"-"`
}

func (p *Project) Process(vis Visitor) {
	vis.VisitProject(p)

	if p.ProjectName != nil {
		p.ProjectName.Process(vis)
	}

	for _, area := range p.Areas {
		area.Process(vis)
	}
}

type ProjectName struct {
	ProjectName string `xml:",attr"`
	UUID        uint32 `xml:",attr"`
}

func (p *ProjectName) Process(v Visitor) {
	// no-op
}

type Area struct {
	UUID         uint32         `xml:",attr"`
	Areas        []*Area        `xml:">Area"`
	DeviceGroups []*DeviceGroup `xml:">DeviceGroup"`
	Outputs      []*Output      `xml:">Output"`
}

func (a *Area) Process(vis Visitor) {
	for _, area := range a.Areas {
		area.Process(vis)
	}

	for _, deviceGroup := range a.DeviceGroups {
		deviceGroup.Process(vis)
	}
	for _, output := range a.Outputs {
		output.Process(vis)
	}
}
