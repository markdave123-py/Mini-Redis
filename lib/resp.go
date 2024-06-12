package lib

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING= '+'
	ERROR = '-'
	INTEGER = ':'
	BULK = '$'
	ARRAY = '*'
)


type Value struct {
		Typ string
		Str string
		Num int
		Bulk string
		Array []Value

}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {

	return &Resp{reader: bufio.NewReader(rd)}
}


func (v Value) Marshall() []byte{
	switch v.Typ{
	case "array":
		return v.MarshallArray()

	case "bulk":
		return v.MarshallBulk()

	case "string":
		return v.MarshallString()

	case "null":
		return v.MarshallNull()

	case "error":
		return v.MarshallError()

	default:
		return []byte{}
	}
}

func (v Value) MarshallString() []byte{
	var bytes []byte

	bytes = append(bytes, STRING)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) MarshallBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) MarshallArray() []byte {

	len := len(v.Array)
	var bytes []byte

	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r','\n')

	for i:= 0; i < len; i++{
		bytes = append(bytes, v.Array[i].Marshall()...)
	}

	return bytes
}


func (v Value) MarshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) MarshallNull() []byte {
	return []byte("$-1\r\n")
}

func ( r *Resp) readLine() (line []byte, n int, err error){
	for{
		b, err  := r.reader.ReadByte()

		if err != nil{
			return nil, 0 , err
		}
		n += 1

		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r'{
			break
		}
	}

	return line[:len(line)-2], n, nil
}

func (r * Resp) readInteger() (x int, n int, err error){


		line, n , err := r.readLine()

		if err != nil{
			return 0,0,err
		}

		i64,err := strconv.ParseInt(string(line), 10, 64)

		if err != nil{
			return 0,n,err
		}

		return int(i64), n ,nil

}

func (r * Resp) Read() (Value, error){
	_type, err := r.reader.ReadByte()

	if err != nil{
		return Value{}, err
	}

	switch _type{
	case ARRAY:
		return r.readArray()

	case BULK:
		return r.readBulk()

	default:
		fmt.Printf("unknown type: %v", string(_type))
		return Value{}, nil
	}
}

func (r *Resp) readArray() (Value, error){

	v := Value{}

	v.Typ = "array"

	arrLen, _, err := r.readInteger()

	if err != nil{
		return v, err
	}

	v.Array = make([]Value, 0)

	for i := 0; i < arrLen; i++{
		val, err := r.Read()

		if err != nil{
			return v, err
		}

		v.Array = append(v.Array, val)
	}

	return v, nil
}


func (r *Resp) readBulk() (Value, error){

	v := Value{}

	v.Typ = "bulk"

	len, _, err := r.readInteger()

	if err != nil{
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.Bulk = string(bulk)

	r.readLine()

	return v, nil
}

type Writer struct {

	writer io.Writer
}

func NewWriter (w io.Writer) *Writer{
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshall()

	_, err := w.writer.Write(bytes)

	if err != nil {
		return err
	}

	return nil
}
