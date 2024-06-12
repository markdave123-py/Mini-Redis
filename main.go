package main

import (
	"Database/lib"
	"fmt"
	"net"
)


func main(){
	fmt.Println("server listening on port : 6379")

	l, err := net.Listen("tcp", ":6379")

	if err != nil{
		fmt.Println(err)
		return
	}

	conn, err := l.Accept()

	if err != nil{
		fmt.Println(err)
		return
	}

	defer conn.Close() // close connection once finished

	for{
		resp := lib.NewResp(conn)

		value, err := resp.Read()

		if err != nil{
			fmt.Println(err)
			return
		}
		_ = value

		writer := lib.NewWriter(conn)

		writer.Write(lib.Value{Typ: "string", Str: "OK"})



	}

}

		// _, err =conn.Read(buf)
		// if err != nil{

		// 	if err == io.EOF{
		// 		break
		// 	}

		// 	fmt.Println("error reading from client: ", err.Error())
		// 	os.Exit(1)

		// }
		// conn.Write([]byte("+OK\r\n"))


	// input := "$5\r\nAhmed\r\n"

	// reader := bufio.NewReader(strings.NewReader(input))

	// b, _  := reader.ReadByte()

	// if b != '$' {
	// 	fmt.Println("invalid type expecting bulk type only")
	// 	os.Exit(1)
	// }

	// size , _ := reader.ReadByte()

	// strSize, _ := strconv.ParseInt(string(size), 10, 64)

	// reader.ReadByte()
	// reader.ReadByte()

	// name :=make([]byte, strSize)
	// reader.Read(name)

	// fmt.Println(string(name))
