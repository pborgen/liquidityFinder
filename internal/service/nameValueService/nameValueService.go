package nameValueService

import (
	"errors"
	"strconv"

	"github.com/pborgen/liquidityFinder/internal/database/model/namevalue"
	"github.com/pborgen/liquidityFinder/internal/types"
)

const (
	DataTypeInt = 1
	DataTypeString = 2
	DataTypeBool = 3
)



func GetValueInt(name string) (int, error) {
	nameValue, err := namevalue.GetByName(name)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(nameValue.Value)
}

func GetValueString(name string) (string, error) {
	nameValue, err := namevalue.GetByName(name)
	if err != nil {
		return "", err
	}
	return nameValue.Value, nil
}

func GetValueBool(name string) (bool, error) {
	nameValue, err := namevalue.GetByName(name)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(nameValue.Value)
}

func UpdateValueString(name string, value string) (error) {

	if !namevalue.Exists(name) {
		return errors.New("NameValue does not exist")
	}

	nameValue, err := namevalue.GetByName(name)
	if err != nil {
		return err
	}

	if nameValue.DataType != DataTypeString {
		return errors.New("NameValue is not a string")
	}

	return namevalue.UpdateValue(name, value)
}

func UpdateValueInt(name string, value int) (error) {

	if !namevalue.Exists(name) {
		return errors.New("NameValue does not exist")
	}

	nameValue, err := namevalue.GetByName(name)
	if err != nil {
		return err
	}

	if nameValue.DataType != DataTypeInt {
		return errors.New("NameValue is not an int")
	}

	return namevalue.UpdateValue(name, strconv.Itoa(value))
}

func UpdateValueBool(name string, value bool) (error) {

	if !namevalue.Exists(name) {
		return errors.New("NameValue does not exist")
	}

	nameValue, err := namevalue.GetByName(name)
	if err != nil {
		return err
	}

	if nameValue.DataType != DataTypeBool {
		return errors.New("NameValue is not a bool")
	}

	return namevalue.UpdateValue(name, strconv.FormatBool(value))
}

func Insert(name string, value string, dataType int) (types.NameValue, error) {

	if dataType == DataTypeInt {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return types.NameValue{}, err
		}
		value = strconv.Itoa(intValue)
	}

	if dataType == DataTypeBool {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return types.NameValue{}, err
		}
		value = strconv.FormatBool(boolValue)
	}

	if namevalue.Exists(name) {
		return types.NameValue{}, errors.New("name already exists")
	}

	return namevalue.Insert(types.NameValue{Name: name, Value: value, DataType: dataType})
}

func Exists(name string) bool {
	return namevalue.Exists(name)
}