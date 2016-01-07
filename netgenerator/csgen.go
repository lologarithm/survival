package main

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
)

func WriteCS(messages []Message, messageMap map[string]Message) {
	gobuf := &bytes.Buffer{}
	gobuf.WriteString("using System;\nusing System.IO;\nusing System.Text;\n\n")
	// 2. Generate go classes
	for _, msg := range messages {
		gobuf.WriteString("public class ")
		gobuf.WriteString(msg.Name)
		gobuf.WriteString(" {")
		for _, f := range msg.Fields {
			gobuf.WriteString("\n\tpublic ")
			gobuf.WriteString(strings.Replace(f.Type, "*", "", 0))
			gobuf.WriteString(" ")
			gobuf.WriteString(f.Name)
		}
		gobuf.WriteString("\n\n")

		gobuf.WriteString("\tpublic void Serialize(BinaryWriter writer) {\n")
		for _, f := range msg.Fields {
			WriteCSSerialize(f, 1, gobuf, messageMap)
		}
		gobuf.WriteString("\t}\n\n")
		gobuf.WriteString("\tpublic void Deserialize(BinaryReader buffer) {\n")
		for _, f := range msg.Fields {
			WriteCSDeserial(f, 1, gobuf, messageMap)
		}
		gobuf.WriteString("\t}\n}\n\n")

	}
	ioutil.WriteFile("../client/Scripts/messages/messages.cs", gobuf.Bytes(), 0775)
}

func WriteCSSerialize(f MessageField, scopeDepth int, buf *bytes.Buffer, messages map[string]Message) {
	for i := 0; i < scopeDepth+1; i++ {
		buf.WriteString("\t")
	}
	switch f.Type {
	case "byte", "int16", "int32", "int64", "uint16", "uint32", "uint64":
		buf.WriteString("writer.Write(")
		if scopeDepth == 1 {
			buf.WriteString("this.")
		}
		buf.WriteString(f.Name)
		buf.WriteString(");\n")
	case "string":
		buf.WriteString("writer.Write(")
		if scopeDepth == 1 {
			buf.WriteString("this.")
		}
		buf.WriteString(f.Name)
		buf.WriteString(".length);\n")
		for i := 0; i < scopeDepth+1; i++ {
			buf.WriteString("\t")
		}
		buf.WriteString("writer.Write(")
		if scopeDepth == 1 {
			buf.WriteString("this.")
		}
		buf.WriteString(f.Name)
		buf.WriteString(");\n")
	default:
		if f.Type[:2] == "[]" {
			// Array!
			buf.WriteString("writer.Write(")
			if scopeDepth == 1 {
				buf.WriteString("this.")
			}
			buf.WriteString(f.Name)
			buf.WriteString(".length);\n")
			for i := 0; i < scopeDepth+1; i++ {
				buf.WriteString("\t")
			}
			buf.WriteString("for ( int i=0; i < ")
			if scopeDepth == 1 {
				buf.WriteString("this.")
			}
			buf.WriteString(f.Name)
			buf.WriteString(".length; i++) {\n")
			fn := f.Name + "[i]"
			if scopeDepth == 1 {
				fn = "this." + fn
			}
			WriteCSSerialize(MessageField{Name: fn, Type: f.Type[2:], Order: f.Order}, scopeDepth+1, buf, messages)
			for i := 0; i < scopeDepth+1; i++ {
				buf.WriteString("\t")
			}
			buf.WriteString("}\n")
		} else {
			// Custom message deserial here.
			if scopeDepth == 1 {
				buf.WriteString("this.")
			}
			buf.WriteString(f.Name)
			buf.WriteString(".Serialize(buffer)\n")
		}
	}
}

func WriteCSDeserial(f MessageField, scopeDepth int, buf *bytes.Buffer, messages map[string]Message) {
	for i := 0; i < scopeDepth+1; i++ {
		buf.WriteString("\t")
	}
	switch f.Type {
	case "byte":
		if scopeDepth == 1 {
			buf.WriteString("this.")
		}
		buf.WriteString(f.Name)
		buf.WriteString(" = buffer.ReadByte()\n")
	case "int16", "int32", "int64", "uint16", "uint32", "uint64":
		if scopeDepth == 1 {
			buf.WriteString("this.")
		}
		buf.WriteString(f.Name)

		funcName := "Read"
		if f.Type[0] == 'u' {
			funcName += strings.ToUpper(f.Type[0:2]) + f.Type[2:]
		} else {
			funcName += strings.ToUpper(f.Type[0:1]) + f.Type[1:]
		}
		buf.WriteString(" = buffer.")
		buf.WriteString(funcName)
		buf.WriteString("();\n")

	case "string":
		lname := "l" + strconv.Itoa(f.Order) + "_" + strconv.Itoa(scopeDepth)
		buf.WriteString("var ")
		buf.WriteString(lname)
		buf.WriteString(" int\n")
		for i := 0; i < scopeDepth; i++ {
			buf.WriteString("\t")
		}
		buf.WriteString("binary.Read(buffer, binary.LittleEndian, &")
		buf.WriteString(lname)
		buf.WriteString(")\n")
		for i := 0; i < scopeDepth; i++ {
			buf.WriteString("\t")
		}
		tmpname := "temp" + strconv.Itoa(f.Order) + "_" + strconv.Itoa(scopeDepth)
		buf.WriteString(tmpname)
		buf.WriteString(" := make([]byte, ")
		buf.WriteString(lname)
		buf.WriteString(")\n")
		for i := 0; i < scopeDepth; i++ {
			buf.WriteString("\t")
		}
		buf.WriteString("buffer.Read(")
		buf.WriteString(tmpname)
		buf.WriteString(")\n")
		for i := 0; i < scopeDepth; i++ {
			buf.WriteString("\t")
		}
		if scopeDepth == 1 {
			buf.WriteString("this.")
		}
		buf.WriteString(f.Name)
		buf.WriteString(" = string(")
		buf.WriteString(tmpname)
		buf.WriteString(")\n")
	default:
		if f.Type[:2] == "[]" {
			// Get len of array
			lname := "l" + strconv.Itoa(f.Order) + "_" + strconv.Itoa(scopeDepth)
			buf.WriteString("var ")
			buf.WriteString(lname)
			buf.WriteString(" int\n")
			for i := 0; i < scopeDepth; i++ {
				buf.WriteString("\t")
			}
			buf.WriteString("binary.Read(buffer, binary.LittleEndian, &")
			buf.WriteString(lname)
			buf.WriteString(")\n")

			// Create array variable
			for i := 0; i < scopeDepth; i++ {
				buf.WriteString("\t")
			}
			if scopeDepth == 1 {
				buf.WriteString("this.")
			}
			buf.WriteString(f.Name)
			buf.WriteString(" = make([]")
			buf.WriteString(f.Type[2:])
			buf.WriteString(", ")
			buf.WriteString(lname)
			buf.WriteString(")\n")

			// Read each var into the array in loop
			for i := 0; i < scopeDepth; i++ {
				buf.WriteString("\t")
			}
			buf.WriteString("for i := 0; i < ")
			buf.WriteString(lname)
			buf.WriteString("; i++ {\n")
			fn := ""
			if scopeDepth == 1 {
				fn += "this."
			}
			fn += f.Name + "[i]"
			WriteCSDeserial(MessageField{Name: fn, Type: f.Type[2:]}, scopeDepth+1, buf, messages)
			for i := 0; i < scopeDepth; i++ {
				buf.WriteString("\t")
			}
			buf.WriteString("}\n")
		} else {
			// Custom message deserial here.
			if scopeDepth == 1 {
				buf.WriteString("this.")
			}
			buf.WriteString(f.Name)
			buf.WriteString(" = new(")
			buf.WriteString(f.Type[1:])
			buf.WriteString(")\n")

			for i := 0; i < scopeDepth; i++ {
				buf.WriteString("\t")
			}
			if scopeDepth == 1 {
				buf.WriteString("this.")
			}
			buf.WriteString(f.Name)
			buf.WriteString(".Deserialize(buffer)\n")
		}
	}

}
