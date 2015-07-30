package graphql

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"unicode/utf8"
)

var (
	ErrInvalidType error = errors.New("invalid type")
	validName      *regexp.Regexp
)

func init() {
	validName = regexp.MustCompile(`^[_A-Za-z][_0-9A-Za-z]*$`)
}

type ScalarTypePrimitive int

const (
	GQL_PRIMITIVE_INTEGER ScalarTypePrimitive = iota
	GQL_PRIMITIVE_FLOAT
	GWL_PRIMITIVE_STRING
)

func (s ScalarTypePrimitive) String() string {
	switch s {
	default:
		panic("ScalarTypePrimitive:string - unrecognised type")
	case GQL_PRIMITIVE_INTEGER:
		return fmt.Sprintf("Int")
	case GQL_PRIMITIVE_FLOAT:
		return fmt.Sprintf("Float")
	case GWL_PRIMITIVE_STRING:
		return fmt.Sprintf("String")
	}
}

type GraphQLType int

const (
	GQL_TYPE_SCALAR GraphQLType = iota
	GQL_TYPE_ENUM
	GQL_TYPE_OBJECT
)

type ObjectChild interface {
	GQLType() GraphQLType
}

type ObjectType struct {
	Name        string
	Description string
	Children    []ObjectChild
}

func (o ObjectType) GQLType() GraphQLType {
	return GQL_TYPE_OBJECT
}

func (o ObjectType) MarshalGraphQL(w io.Writer) error {

	buffer := bytes.Buffer{}
	buffer.WriteString("type ")
	buffer.WriteString(o.Name)
	buffer.WriteString(" ")

	if false {
		// TOTO: print the name of the interface this type implements
		// e.g.
		// buffer.WriteString(o.Implements.Name)
		// buffer.WriteString(' ')
	}

	for index, child := range o.Children {

		buffer.WriteString("\n\t")

		switch t := child.(type) {
		default:
			panic("unrecognised type")
		case ScalarType:
			buffer.WriteString(t.Name)
			buffer.WriteString(":\t")
			buffer.WriteString(t.Type.String())
		}

		if index == len(o.Children)-1 {
			buffer.WriteString("\n}")
		}
	}

	buffer.WriteRune('\n')

	if _, err := buffer.WriteTo(w); err != nil {
		return err
	}

	return nil

}

type ScalarType struct {
	Type        ScalarTypePrimitive
	Name        string
	Value       string
	Description string
}

func (s ScalarType) GQLType() GraphQLType {
	return GQL_TYPE_SCALAR
}

func (s ScalarType) IsValid() bool {
	switch s.Type {
	default:
		return false
	case GQL_PRIMITIVE_INTEGER:
		_, err := strconv.ParseInt(s.Value, 10, 32)
		return err == nil
	case GQL_PRIMITIVE_FLOAT:
		_, err := strconv.ParseFloat(s.Value, 64)
		return err == nil
	case GWL_PRIMITIVE_STRING:
		return utf8.ValidString(s.Value)
	}
}

type EnumType struct {
	Name        string
	Description string
	Values      []EnumValue
}

func (e EnumType) GQLType() GraphQLType {
	return GQL_TYPE_ENUM
}

type EnumValue struct {
	Name        string
	Description string
}

func (e EnumType) IsValid() bool {

	if !validName.MatchString(e.Name) {
		return false
	}

	for _, value := range e.Values {
		if !validName.MatchString(value.Name) {
			return false
		}
	}

	return true

}

func (e EnumType) MarshalGraphQL(w io.Writer) error {

	// not sure if this is the best place to be checking this
	if !e.IsValid() {
		return ErrInvalidType
	}

	buffer := bytes.Buffer{}
	buffer.WriteString("enum " + e.Name + " {")
	for index, value := range e.Values {

		buffer.WriteString("\n\t" + value.Name)

		if index == len(e.Values)-1 {
			buffer.WriteString("\n}")
		} else {
			buffer.WriteString(", ")
		}

	}

	buffer.WriteRune('\n')

	if _, err := buffer.WriteTo(w); err != nil {
		return err
	}

	return nil

}
