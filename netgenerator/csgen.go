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
			gobuf.WriteString(goTypeToCS(f.Type))
			gobuf.WriteString(" ")
			gobuf.WriteString(f.Name)
			gobuf.WriteString(";")
		}
		gobuf.WriteString("\n\n")

		gobuf.WriteString("\tpublic void Serialize(BinaryWriter buffer) {\n")
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
	ioutil.WriteFile("../client/Assets/Scripts/messages/messages.cs", gobuf.Bytes(), 0775)
}

func goTypeToCS(tn string) string {
	if len(tn) > 3 {
		for tn[:2] == "[]" {
			tn = tn[2:] + tn[:2]
		}
		if tn[:3] == "uin" {
			tn = strings.ToUpper(tn[0:2]) + tn[2:]
		} else if tn[:3] == "int" {
			tn = strings.ToUpper(tn[0:1]) + tn[1:]
		}
	}
	tn = strings.Replace(tn, "*", "", -1)

	return tn
}

func WriteCSSerialize(f MessageField, scopeDepth int, buf *bytes.Buffer, messages map[string]Message) {
	for i := 0; i < scopeDepth+1; i++ {
		buf.WriteString("\t")
	}
	switch f.Type {
	case "byte", "int16", "int32", "int64", "uint16", "uint32", "uint64":
		buf.WriteString("buffer.Write(")
		if scopeDepth == 1 {
			buf.WriteString("this.")
		}
		buf.WriteString(f.Name)
		buf.WriteString(");\n")
	case "string":
		buf.WriteString("buffer.Write((Int32)")
		if scopeDepth == 1 {
			buf.WriteString("this.")
		}
		buf.WriteString(f.Name)
		buf.WriteString(".Length);\n")
		for i := 0; i < scopeDepth+1; i++ {
			buf.WriteString("\t")
		}
		buf.WriteString("buffer.Write(")
		if scopeDepth == 1 {
			buf.WriteString("this.")
		}
		buf.WriteString(f.Name)
		buf.WriteString(");\n")
	default:
		if f.Type[:2] == "[]" {
			// Array!
			buf.WriteString("buffer.Write((Int32)")
			if scopeDepth == 1 {
				buf.WriteString("this.")
			}
			buf.WriteString(f.Name)
			buf.WriteString(".Length);\n")
			for i := 0; i < scopeDepth+1; i++ {
				buf.WriteString("\t")
			}

			loopvar := "v" + strconv.Itoa(scopeDepth+1)
			buf.WriteString("for (int ")
			buf.WriteString(loopvar)
			buf.WriteString(" = 0; ")
			buf.WriteString(loopvar)
			buf.WriteString(" < ")
			if scopeDepth == 1 {
				buf.WriteString("this.")
			}
			buf.WriteString(f.Name)
			buf.WriteString(".Length; ")
			buf.WriteString(loopvar)
			buf.WriteString("++) {\n")
			fn := f.Name + "[" + loopvar + "]"
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
			buf.WriteString(".Serialize(buffer);\n")
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
		buf.WriteString(" = buffer.ReadByte();\n")
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
		buf.WriteString("int ")
		buf.WriteString(lname)
		buf.WriteString(" = buffer.ReadInt32();\n")
		for i := 0; i < scopeDepth+1; i++ {
			buf.WriteString("\t")
		}
		buf.WriteString("byte[] ")
		tmpname := "temp" + strconv.Itoa(f.Order) + "_" + strconv.Itoa(scopeDepth)
		buf.WriteString(tmpname)
		buf.WriteString(" = buffer.ReadBytes(")
		buf.WriteString(lname)
		buf.WriteString(");\n")
		for i := 0; i < scopeDepth+1; i++ {
			buf.WriteString("\t")
		}
		if scopeDepth == 1 {
			buf.WriteString("this.")
		}
		buf.WriteString(f.Name)
		buf.WriteString(" = System.Text.Encoding.UTF8.GetString(")
		buf.WriteString(tmpname)
		buf.WriteString(");\n")
	default:
		if f.Type[:2] == "[]" {
			// Get len of array
			lname := "l" + strconv.Itoa(f.Order) + "_" + strconv.Itoa(scopeDepth)
			buf.WriteString("int ")
			buf.WriteString(lname)
			buf.WriteString(" = buffer.ReadInt32();\n")
			for i := 0; i < scopeDepth+1; i++ {
				buf.WriteString("\t")
			}

			// Create array variable
			if scopeDepth == 1 {
				buf.WriteString("this.")
			}
			buf.WriteString(f.Name)
			buf.WriteString(" = new ")
			t := goTypeToCS(f.Type)
			numdim := 0
			for t[len(t)-2:] == "[]" {
				t = t[:len(t)-2]
				numdim++
			}
			buf.WriteString(t)
			buf.WriteString("[")
			buf.WriteString(lname)
			buf.WriteString("]")
			for i := 0; i < numdim-1; i++ {
				buf.WriteString("[]")
			}
			buf.WriteString(";\n")

			// Read each var into the array in loop
			for i := 0; i < scopeDepth+1; i++ {
				buf.WriteString("\t")
			}
			loopvar := "v" + strconv.Itoa(scopeDepth+1)
			buf.WriteString("for (int ")
			buf.WriteString(loopvar)
			buf.WriteString(" = 0; ")
			buf.WriteString(loopvar)
			buf.WriteString(" < ")
			buf.WriteString(lname)
			buf.WriteString("; ")
			buf.WriteString(loopvar)
			buf.WriteString("++) {\n")
			fn := ""
			if scopeDepth == 1 {
				fn += "this."
			}
			fn += f.Name + "[" + loopvar + "]"
			WriteCSDeserial(MessageField{Name: fn, Type: f.Type[2:]}, scopeDepth+1, buf, messages)
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
			buf.WriteString(" = new ")
			buf.WriteString(f.Type[1:])
			buf.WriteString("();\n")

			for i := 0; i < scopeDepth+1; i++ {
				buf.WriteString("\t")
			}
			if scopeDepth == 1 {
				buf.WriteString("this.")
			}
			buf.WriteString(f.Name)
			buf.WriteString(".Deserialize(buffer);\n")
		}
	}

}
