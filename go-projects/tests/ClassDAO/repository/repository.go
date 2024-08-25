package repository

import (
	"database/sql"
	"m/tests/ClassDAO/dao"
	"m/tests/ClassDAO/entities"
)

func InsertClass(db *sql.DB, name string) (int, error) {
	class := &entities.Class{Name: name}
	daoClass := dao.NewDAO(db)
	return daoClass.Create("CLASSES", class)
}

func UpdateClass(db *sql.DB, class *entities.Class) error {
	daoClass := dao.NewDAO(db)
	return daoClass.Update("CLASSES", class)
}

func DeleteClass(db *sql.DB, classID int) error {
	daoClass := dao.NewDAO(db)
	return daoClass.Delete("CLASSES", classID)
}

func ReadClass(db *sql.DB, classID int) (*entities.Class, error) {
	class := &entities.Class{}
	daoClass := dao.NewDAO(db)

	err := daoClass.Read("CLASSES", classID, class)
	if err != nil {
		return class, err
	}

	condition := "CLASS_ID = $1"
	args := []interface{}{classID}
	model := &entities.Object{}

	objects, err := daoClass.ReadMultiple("OBJECTS", condition, args, model)
	if err != nil {
		return class, err
	}

	for _, obj := range objects {
		tmpObj := *obj.(*entities.Object)
		condition := "OBJECT_ID = $1"
		args := []interface{}{tmpObj.Id}
		model := &entities.Item{}

		items, err := daoClass.ReadMultiple("ITEMS_BY_OBJECT_VIEW", condition, args, model)
		if err != nil {
			return class, err
		}

		for _, item := range items {
			tmpObj.Items = append(tmpObj.Items, *item.(*entities.Item))
		}

		class.Objects = append(class.Objects, tmpObj)
	}

	return class, nil
}
