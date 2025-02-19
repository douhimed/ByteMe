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
		Magic:         hex.EncodeToString(container.parse_u4()),
		Minor:         container.parse_u2(),
		Major:         container.parse_u2(),
		ConstantsPool: make([]Constant_Pool_Type, 0),
	}

	constPoolCount := int(container.parse_u2())
	for i := 0; i < constPoolCount-1; i++ {
		tag := container.parse_u1()

		var cp Constant_Pool_Type
		if tag == CONSTANT_Methodref {
			cp = Methodref_Info{
				Tag:                 POOL_CONSTANTS[tag],
				Class_index:         container.parse_u2(),
				Name_and_type_index: container.parse_u2(),
			}
		} else if tag == CONSTANT_Class {
			cp = CONSTANT_Class_Info{
				Tag:        POOL_CONSTANTS[tag],
				Name_index: container.parse_u2(),
			}
		} else if tag == CONSTANT_NameAndType {
			cp = CONSTANT_NameAndType_info{
				Tag:              POOL_CONSTANTS[tag],
				Name_index:       container.parse_u2(),
				Descriptor_index: container.parse_u2(),
			}
		} else if tag == CONSTANT_Utf8 {
			length := int(container.parse_u2())
			value := container.parse_(length)
			str := fmt.Sprintf("%s", value)
			cp = CONSTANT_Utf8_Info{
				Tag:   POOL_CONSTANTS[tag],
				Bytes: str,
			}
		} else if tag == CONSTANT_Fieldref {
			cp = CONSTANT_Fieldref_info{
				Tag:                 POOL_CONSTANTS[tag],
				Class_index:         container.parse_u2(),
				Name_and_type_index: container.parse_u2(),
			}
		} else if tag == CONSTANT_String {
			cp = CONSTANT_String_info{
				Tag:          POOL_CONSTANTS[tag],
				String_index: container.parse_u2(),
			}
		}

		if cp != nil {
			clazz.addContstantPool(cp)
		}
	}

	clazz.asJson()
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

// Clazz Funcs : The class File Structe

type Clazz struct {
	Magic         string               `json:"magic"`
	Minor         int16                `json:"manor"`
	Major         int16                `json:"major"`
	ConstantsPool []Constant_Pool_Type `json:"constants_pool"`
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
	Tag   string `json:"tag"`
	Bytes string `json:"bytesAsString"`
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

func (c *Container) parse_(steps int) []byte {
	bytes := c.Content[c.Cursor : c.Cursor+steps]
	c.Cursor += steps
	return bytes
}

func (c *Container) parse_u4() []byte {
	bytes := c.Content[c.Cursor : c.Cursor+4]
	c.Cursor += 4
	return bytes
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
