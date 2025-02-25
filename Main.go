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

	container.parseConstantsPool(clazz)
	container.parseClassFlags(clazz)
	container.parseThisClass(clazz)
	container.parseSuperClass(clazz)
	container.parseInterfaces(clazz)
	container.parseFields(clazz)
	container.parseMethods(clazz)

	clazz.asJson()

}

// Util Funcs

func parseConstantPoolEntry(container *Container, tag int8) Constant_Pool_Type {
	tagValue := POOL_CONSTANTS[tag]

	switch tag {
	case CONSTANT_Methodref:
		return container.parseMethodref(tagValue)
	case CONSTANT_Class:
		return container.parseClass(tagValue)
	case CONSTANT_NameAndType:
		return container.parseNameAndType(tagValue)
	case CONSTANT_Utf8:
		return container.parseUtf8(tagValue)
	case CONSTANT_Fieldref:
		return container.parseFieldref(tagValue)
	case CONSTANT_String:
		return container.parseString(tagValue)
	case CONSTANT_InvokeDynamic:
		return container.parseInvokeDynamic(tagValue)
	case CONSTANT_Dynamic:
		return container.parseDynamic(tagValue)
	case CONSTANT_Double:
		return container.parseDouble(tagValue)
	case CONSTANT_Float:
		return container.parseFloat(tagValue)
	case CONSTANT_MethodHandle:
		return container.parseMethodHandle(tagValue)
	default:
		log.Printf("%s[%d] not yet implemented from the constant pool", POOL_CONSTANTS[tag], tag)
	}
	return nil
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

func to_uint16(data []byte) uint16 {
	var res uint16
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &res)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	return res
}

func to_int32(data []byte) int32 {
	var res int32
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

// Clazz Funcs : The class File Structure

type Clazz struct {
	Magic           string               `json:"magic"`
	Minor           int16                `json:"manor"`
	Major           int16                `json:"major"`
	ConstantsPool   []Constant_Pool_Type `json:"constants_pool"`
	AccessFlags     []string             `json:"access_flags"`
	ThisClass       Constant_Pool_Type   `json:"this_class"`
	SuperClass      Constant_Pool_Type   `json:"super_class"`
	InterfacesCount int16                `json:"interfaces_count"`
	FieldsCount     int16                `json:"fields_count"`
	Methods         []Method_Info        `json:"methods"`
}

func (c *Clazz) asJson() {
	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonData))

}

type Constant_Pool_Type interface {
}

type Methodref_Info struct {
	Tag                 string `json:"tag"`
	Class_index         int16  `json:"class_index"`
	Name_and_type_index int16  `json:"name_and_type_index"`
}

type Method_Info struct {
	Access_flags     []string         `json:"access_flags"`
	Name_index       int16            `json:"name_index"`
	Descriptor_index int16            `json:"descriptor_index"`
	Attributes       []Attribute_Info `json:"attributes"`
}

func (m *Method_Info) parseAccessFlags(flags uint16) {
	for bitmask, name := range MethodAccessFlags {
		if flags&bitmask != 0 {
			m.Access_flags = append(m.Access_flags, name)
		}
	}
}

func (m *Method_Info) parseAttributes(container *Container) {
	count := int(container.parse_u2())
	for i := 0; i < count; i++ {
		var attr Attribute_Info
		attr.Attribute_name_index = container.parse_u2()
		attrLength := to_int32(container.parse_u4())
		infobytes, _ := container.parse_(int(attrLength))
		attr.Infos = infobytes
		m.Attributes = append(m.Attributes, attr)
	}
}

type Attribute_Info struct {
	Attribute_name_index int16  `json:"attribute_name_index"`
	Infos                []byte `josn:"infos"`
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

type CONSTANT_Dynamic_info struct {
	Tag                         string `json:"tag"`
	Bootstrap_method_attr_index int16  `json:"bootstrap_method_attr_index"`
	Name_and_type_index         int16  `json:"name_and_type_index"`
}

type CONSTANT_InvokeDynamic_info struct {
	Tag                         string `json:"tag"`
	Bootstrap_method_attr_index int16  `json:"bootstrap_method_attr_index"`
	Name_and_type_index         int16  `json:"name_and_type_index"`
}

type CONSTANT_Double_info struct {
	Tag       string `json:"tag"`
	HighBytes []byte `json:"high_bytes"`
	LowBytes  []byte `json:"low_bytes"`
}

type CONSTANT_Float_info struct {
	Tag   string `json:"tag"`
	Bytes []byte `json:"bytes"`
}

type CONSTANT_MethodHandle_info struct {
	Tag            string `json:"tag"`
	ReferenceKind  string `json:"reference_kind"`
	ReferenceIndex int16  `json:"reference_index"`
}

type CONSTANT_String_info struct {
	Tag          string `json:"tag"`
	String_index int16  `json:"string_index"`
}

// CONTAINER Funcs (the file contentes 'bytes', cursor 'current index to parse', parse funcs)

type Container struct {
	Content []byte
	Cursor  int
}

func (container *Container) parseConstantsPool(cl *Clazz) {
	constPoolCount := int(container.parse_u2())
	for i := 0; i < constPoolCount-1; i++ {
		tag := container.parse_u1()
		cp := parseConstantPoolEntry(container, tag)

		if cp != nil {
			cl.ConstantsPool = append(cl.ConstantsPool, cp)
		}
	}
}

func (container *Container) parseClassFlags(cl *Clazz) {
	flags := container.parse_u2_u()

	for bitmask, name := range classAccessFlagsMap {
		if flags&bitmask != 0 {
			cl.AccessFlags = append(cl.AccessFlags, name)
		}
	}
}

func (container *Container) getConstantPoolInfos(cl *Clazz, index int16) (Constant_Pool_Type, error) {
	if index == 0 || int(index) > len(cl.ConstantsPool) {
		return nil, fmt.Errorf("invalid constant pool index: %d", index)
	}

	cp, ok := cl.ConstantsPool[index-1].(CONSTANT_Class_Info)
	if !ok {
		return nil, fmt.Errorf("unexpected type in constant pool at index %d", index)
	}

	return cl.ConstantsPool[cp.Name_index], nil
}

func (container *Container) parseThisClass(cl *Clazz) {
	thisClassIndex := container.parse_u2()
	cp, err := container.getConstantPoolInfos(cl, thisClassIndex)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	cl.ThisClass = cp
}

func (container *Container) parseSuperClass(cl *Clazz) {
	superClassIndex := container.parse_u2()
	cp, err := container.getConstantPoolInfos(cl, superClassIndex)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	cl.SuperClass = cp
}

func (container *Container) parseInterfaces(cl *Clazz) {
	count := container.parse_u2()
	cl.InterfacesCount = count
	if count > 0 {
		log.Fatal("interface parser is not yet implemented")
	}
}

func (container *Container) parseFields(cl *Clazz) {
	count := container.parse_u2()
	cl.FieldsCount = count
	if count > 0 {
		log.Fatal("Fields parser is not yet implemented")
	}
}

func (container *Container) parseMethods(cl *Clazz) {
	count := container.parse_u2()
	for i := 0; i < int(count); i++ {
		var m Method_Info
		m.parseAccessFlags(container.parse_u2_u())
		m.Name_index = container.parse_u2()
		m.Descriptor_index = container.parse_u2()
		m.parseAttributes(container)
		cl.Methods = append(cl.Methods, m)
	}
}

func (container *Container) parseMethodref(tagValue string) Methodref_Info {
	return Methodref_Info{
		Tag:                 tagValue,
		Class_index:         container.parse_u2(),
		Name_and_type_index: container.parse_u2(),
	}
}

func (container *Container) parseClass(tagValue string) CONSTANT_Class_Info {
	return CONSTANT_Class_Info{
		Tag:        tagValue,
		Name_index: container.parse_u2() - 1,
	}
}

func (container *Container) parseNameAndType(tagValue string) CONSTANT_NameAndType_info {
	return CONSTANT_NameAndType_info{
		Tag:              tagValue,
		Name_index:       container.parse_u2(),
		Descriptor_index: container.parse_u2(),
	}
}

func (container *Container) parseUtf8(tagValue string) CONSTANT_Utf8_Info {
	length := int(container.parse_u2())
	bytes, bytesAsString := container.parse_(length)
	return CONSTANT_Utf8_Info{
		Tag:         tagValue,
		Bytes:       bytes,
		StringBytes: bytesAsString,
	}
}

func (container *Container) parseFieldref(tagValue string) CONSTANT_Fieldref_info {
	return CONSTANT_Fieldref_info{
		Tag:                 tagValue,
		Class_index:         container.parse_u2(),
		Name_and_type_index: container.parse_u2(),
	}
}

func (container *Container) parseString(tagValue string) CONSTANT_String_info {
	return CONSTANT_String_info{
		Tag:          tagValue,
		String_index: container.parse_u2(),
	}
}

func (c *Container) parseInvokeDynamic(tagValue string) CONSTANT_InvokeDynamic_info {
	return CONSTANT_InvokeDynamic_info{
		Tag:                         tagValue,
		Bootstrap_method_attr_index: c.parse_u2(),
		Name_and_type_index:         c.parse_u2(),
	}
}

func (c *Container) parseDynamic(tagValue string) CONSTANT_Dynamic_info {
	return CONSTANT_Dynamic_info{
		Tag:                         tagValue,
		Bootstrap_method_attr_index: c.parse_u2(),
		Name_and_type_index:         c.parse_u2(),
	}
}

func (c *Container) parseDouble(tagValue string) CONSTANT_Double_info {
	return CONSTANT_Double_info{
		Tag:       tagValue,
		HighBytes: c.parse_u4(),
		LowBytes:  c.parse_u4(),
	}
}

func (c *Container) parseFloat(tagValue string) CONSTANT_Float_info {
	return CONSTANT_Float_info{
		Tag:   tagValue,
		Bytes: c.parse_u4(),
	}
}

func (c *Container) parseMethodHandle(tagValue string) CONSTANT_MethodHandle_info {
	return CONSTANT_MethodHandle_info{
		Tag:            tagValue,
		ReferenceKind:  MethodHandlers[int8(c.parse_u1())],
		ReferenceIndex: c.parse_u2(),
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

// Define constant pool tags

const (
	CONSTANT_Utf8               int8 = 1
	CONSTANT_Integer            int8 = 3
	CONSTANT_Float              int8 = 4
	CONSTANT_Long               int8 = 5
	CONSTANT_Double             int8 = 6
	CONSTANT_Class              int8 = 7
	CONSTANT_String             int8 = 8
	CONSTANT_Fieldref           int8 = 9
	CONSTANT_Methodref          int8 = 10
	CONSTANT_InterfaceMethodref int8 = 11
	CONSTANT_NameAndType        int8 = 12
	CONSTANT_MethodHandle       int8 = 15
	CONSTANT_MethodType         int8 = 16
	CONSTANT_Dynamic            int8 = 17
	CONSTANT_InvokeDynamic      int8 = 18
	CONSTANT_Module             int8 = 19
	CONSTANT_Package            int8 = 20
)

var POOL_CONSTANTS = map[int8]string{
	CONSTANT_Utf8:               "CONSTANT_Utf8",
	CONSTANT_Integer:            "CONSTANT_Integer",
	CONSTANT_Float:              "CONSTANT_Float",
	CONSTANT_Long:               "CONSTANT_Long",
	CONSTANT_Double:             "CONSTANT_Double",
	CONSTANT_Class:              "CONSTANT_Class",
	CONSTANT_String:             "CONSTANT_String",
	CONSTANT_Fieldref:           "CONSTANT_Fieldref",
	CONSTANT_Methodref:          "CONSTANT_Methodref",
	CONSTANT_InterfaceMethodref: "CONSTANT_InterfaceMethodref",
	CONSTANT_NameAndType:        "CONSTANT_NameAndType",
	CONSTANT_MethodHandle:       "CONSTANT_MethodHandle",
	CONSTANT_MethodType:         "CONSTANT_MethodType",
	CONSTANT_Dynamic:            "CONSTANT_Dynamic",
	CONSTANT_InvokeDynamic:      "CONSTANT_InvokeDynamic",
	CONSTANT_Module:             "CONSTANT_Module",
	CONSTANT_Package:            "CONSTANT_Package",
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

var classAccessFlagsMap = map[uint16]string{
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

const (
	M_ACC_PUBLIC     = 0x0001
	ACC_PRIVATE      = 0x0002
	ACC_PROTECTED    = 0x0004
	ACC_STATIC       = 0x0008
	M_ACC_FINAL      = 0x0010
	ACC_SYNCHRONIZED = 0x0020
	ACC_BRIDGE       = 0x0040
	ACC_VARARGS      = 0x0080
	ACC_NATIVE       = 0x0100
	M_ACC_ABSTRACT   = 0x0400
	ACC_STRICT       = 0x0800
	M_ACC_SYNTHETIC  = 0x1000
)

var MethodAccessFlags = map[uint16]string{
	M_ACC_PUBLIC:     "ACC_PUBLIC",
	ACC_PRIVATE:      "ACC_PRIVATE",
	ACC_PROTECTED:    "ACC_PROTECTED",
	ACC_STATIC:       "ACC_STATIC",
	M_ACC_FINAL:      "ACC_FINAL",
	ACC_SYNCHRONIZED: "ACC_SYNCHRONIZED",
	ACC_BRIDGE:       "ACC_BRIDGE",
	ACC_VARARGS:      "ACC_VARARGS",
	ACC_NATIVE:       "ACC_NATIVE",
	M_ACC_ABSTRACT:   "ACC_ABSTRACT",
	ACC_STRICT:       "ACC_STRICT",
	M_ACC_SYNTHETIC:  "ACC_SYNTHETIC",
}

var MethodHandlers = map[int8]string{
	1: "REF_getField  getfield C.f:T",
	2: "REF_getStatic getstatic C.f:T",
	3: "REF_putField  putfield C.f:T",
	4: "REF_putStatic putstatic C.f:T",
	5: "REF_invokeVirtual invokevirtual C.m:(A*)T",
	6: "REF_invokeStatic invokestatic C.m:(A*)T",
	7: "REF_invokeSpecial invokespecial C.m:(A*)T",
	8: "REF_newInvokeSpecial new C; dup; invokespecial C.<init>:(A*)V",
	9: "REF_invokeInterface",
}
