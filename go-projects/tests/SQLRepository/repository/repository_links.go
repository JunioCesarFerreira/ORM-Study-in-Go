package repository

type BaseLinks struct {
	TableName     string
	MasterColName string
	LinkColName   string
}

func NewBaseLinks(tableName, masterCol, linkCol string) *BaseLinks {
	return &BaseLinks{
		TableName:     tableName,
		MasterColName: masterCol,
		LinkColName:   linkCol,
	}
}

type Links struct {
	BaseLinks
	MasterId int
	LinksIds []int
}

func (b BaseLinks) NewLinks(masterId int, linksIds []int) Links {
	return Links{
		BaseLinks: b,
		MasterId:  masterId,
		LinksIds:  linksIds,
	}
}

func (b BaseLinks) NewSelectLinks(masterId int) Links {
	return Links{
		BaseLinks: b,
		MasterId:  masterId,
	}
}
