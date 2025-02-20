package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	container := InitContainerFromFile("App.class")

	clazz := &Clazz{
		Magic: hex.EncodeToString(container.parse_u4()),
		Minor: container.parse_u2(),
		Major: container.parse_u2(),
	}

	constPoolCount := int(container.parse_u2())
	for i := 0; i < constPoolCount-1; i++ {
		tag := container.parse_u1()
		cp := parseConstantPoolEntry(container, tag)

		if cp != nil {
			clazz.addContstantPool(cp)
		}
	}

	parseClassFlags(container, clazz)

	clazz.asJson()
}

func parseClassFlags(c *Container, cl *Clazz) {
	flags := c.parse_u2_u()

	for bitmask, name := range accessFlagsMap {
		if flags&bitmask != 0 {
			cl.AddFlagAccess(name)
		}
	}
}

func parseConstantPoolEntry(container *Container, tag int8) Constant_Pool_Type {
	tagValue := POOL_CONSTANTS[tag]

	switch tag {
	case CONSTANT_Methodref:
		return parseMethodref(container, tagValue)
	case CONSTANT_Class:
		return parseClass(container, tagValue)
	case CONSTANT_NameAndType:
		return parseNameAndType(container, tagValue)
	case CONSTANT_Utf8:
		return parseUtf8(container, tagValue)
	case CONSTANT_Fieldref:
		return parseFieldref(container, tagValue)
	case CONSTANT_String:
		return parseString(container, tagValue)
	default:
		log.Printf("%s[%d] not yet implemented from the constant pool", POOL_CONSTANTS[tag], tag)
		os.Exit(0)
	}
	return nil
}

func parseMethodref(container *Container, tagValue string) Methodref_Info {
	return Methodref_Info{
		Tag:                 tagValue,
		Class_index:         container.parse_u2(),
		Name_and_type_index: container.parse_u2(),
	}
}

func parseClass(container *Container, tagValue string) CONSTANT_Class_Info {
	return CONSTANT_Class_Info{
		Tag:        tagValue,
		Name_index: container.parse_u2(),
	}
}

func parseNameAndType(container *Container, tagValue string) CONSTANT_NameAndType_info {
	return CONSTANT_NameAndType_info{
		Tag:              tagValue,
		Name_index:       container.parse_u2(),
		Descriptor_index: container.parse_u2(),
	}
}

func parseUtf8(container *Container, tagValue string) CONSTANT_Utf8_Info {
	length := int(container.parse_u2())
	bytes, bytesAsString := container.parse_(length)
	return CONSTANT_Utf8_Info{
		Tag:         tagValue,
		Bytes:       bytes,
		StringBytes: bytesAsString,
	}
}

func parseFieldref(container *Container, tagValue string) CONSTANT_Fieldref_info {
	return CONSTANT_Fieldref_info{
		Tag:                 tagValue,
		Class_index:         container.parse_u2(),
		Name_and_type_index: container.parse_u2(),
	}
}

func parseString(container *Container, tagValue string) CONSTANT_String_info {
	return CONSTANT_String_info{
		Tag:          tagValue,
		String_index: container.parse_u2(),
	}
}

// Clazz Funcs : The class File Structe

type Clazz struct {
	Magic         string               `json:"magic"`
	Minor         int16                `json:"manor"`
	Major         int16                `json:"major"`
	ConstantsPool []Constant_Pool_Type `json:"constants_pool"`
	AccessFlags   []string             `json:"access_flags"`
}

func (c *Clazz) AddFlagAccess(f string) {
	c.AccessFlags = append(c.AccessFlags, f)
}

func (c *Clazz) addContstantPool(cp Constant_Pool_Type) {
	c.ConstantsPool = append(c.ConstantsPool, cp)
}

func (c *Clazz) asJson() {
	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

}

type Constant_Pool_Type interface{}

type Methodref_Info struct {
	Tag                 string `json:"tag"`
	Class_index         int16  `json:"class_index"`
	Name_and_type_index int16  `json:"name_and_type_index"`
}

type CONSTANT_Class_Info struct {
	Tag        string `json:"tag"`
	Name_index int16  `json:"name_index"`
}

type CONSTANT_NameAndType_info struct {
	Tag              string `json:"tag"`
	Name_index       int16  `json:"name_index"`
	Descriptor_index int16  `json:"descriptor_index"`
}

type CONSTANT_Utf8_Info struct {
	Tag         string `json:"tag"`
	Bytes       []byte `json:"bytes"`
	StringBytes string `json:"stringBytes"`
}

type CONSTANT_Fieldref_info struct {
	Tag                 string `json:"tag"`
	Class_index         int16  `json:"class_index"`
	Name_and_type_index int16  `json:"name_and_type_index"`
}

type CONSTANT_String_info struct {
	Tag          string `json:"tag"`
	String_index int16  `json:"string_index"`
}

// CONTAINER Funcs

type Container struct {
	Content []byte
	Cursor  int
}

func InitContainerFromFile(file string) *Container {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	return &Container{
		Content: data,
		Cursor:  0,
	}
}

func (c *Container) parse_u1() int8 {
	bytes := c.Content[c.Cursor : c.Cursor+1]
	c.Cursor++
	return to_int8(bytes)
}

func (c *Container) parse_u2() int16 {
	bytes := c.Content[c.Cursor : c.Cursor+2]
	c.Cursor += 2
	return to_int16(bytes)
}

func (c *Container) parse_u2_u() uint16 {
	bytes := c.Content[c.Cursor : c.Cursor+2]
	c.Cursor += 2
	return to_uint16(bytes)
}

func (c *Container) parse_(steps int) ([]byte, string) {
	bytes := c.Content[c.Cursor : c.Cursor+steps]
	c.Cursor += steps
	return bytes, string(bytes)
}

func (c *Container) parse_u4() []byte {
	bytes := c.Content[c.Cursor : c.Cursor+4]
	c.Cursor += 4
	return bytes
}

func to_uint16(data []byte) uint16 {
	var res uint16
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &res)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	return res
}

func to_int16(data []byte) int16 {
	var res int16
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &res)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	return res
}

func to_int8(data []byte) int8 {
	var res int8
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &res)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	return res
}

// CONSTANT POOL TAGS

var CONSTANT_Class int8 = 7
var CONSTANT_Fieldref int8 = 9
var CONSTANT_Methodref int8 = 10
var CONSTANT_InterfaceMethodref int8 = 11
var CONSTANT_String int8 = 8
var CONSTANT_Integer int8 = 3
var CONSTANT_Float int8 = 4
var CONSTANT_Long int8 = 5
var CONSTANT_Double int8 = 6
var CONSTANT_NameAndType int8 = 12
var CONSTANT_Utf8 int8 = 1
var CONSTANT_MethodHandle int8 = 15
var CONSTANT_MethodType int8 = 16
var CONSTANT_Dynamic int8 = 17
var CONSTANT_InvokeDynamic int8 = 18
var CONSTANT_Module int8 = 19
var CONSTANT_Package int8 = 20

var POOL_CONSTANTS = map[int8]string{
	1:  "CONSTANT_Utf8",
	3:  "CONSTANT_Integer",
	4:  "CONSTANT_Float",
	5:  "CONSTANT_Long",
	6:  "CONSTANT_Double",
	7:  "CONSTANT_Class",
	8:  "CONSTANT_String",
	9:  "CONSTANT_Fieldref",
	10: "CONSTANT_Methodref",
	11: "CONSTANT_InterfaceMethodref",
	12: "CONSTANT_NameAndType",
	15: "CONSTANT_MethodHandle",
	16: "CONSTANT_MethodType",
	17: "CONSTANT_Dynamic",
	18: "CONSTANT_InvokeDynamic",
	19: "CONSTANT_Module",
	20: "CONSTANT_Package",
}

// Class flags constants

const (
	ACC_PUBLIC     = 0x0001
	ACC_FINAL      = 0x0010
	ACC_SUPER      = 0x0020
	ACC_INTERFACE  = 0x0200
	ACC_ABSTRACT   = 0x0400
	ACC_SYNTHETIC  = 0x1000
	ACC_ANNOTATION = 0x2000
	ACC_ENUM       = 0x4000
	ACC_MODULE     = 0x8000
)

var accessFlagsMap = map[uint16]string{
	ACC_PUBLIC:     "ACC_PUBLIC",
	ACC_FINAL:      "ACC_FINAL",
	ACC_SUPER:      "ACC_SUPER",
	ACC_INTERFACE:  "ACC_INTERFACE",
	ACC_ABSTRACT:   "ACC_ABSTRACT",
	ACC_SYNTHETIC:  "ACC_SYNTHETIC",
	ACC_ANNOTATION: "ACC_ANNOTATION",
	ACC_ENUM:       "ACC_ENUM",
	ACC_MODULE:     "ACC_MODULE",
}
