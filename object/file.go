package object

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type FileState int

const (
	UnknownFileState FileState = iota
	ReadFileState
	ReadlineFileState
	WriteFileState
)

type File struct {
	Value   *os.File
	scanner *bufio.Scanner
	State   FileState
}

func (f *File) Type() Type      { return FileType }
func (f *File) Inspect() string { return fmt.Sprintf("File[%s]", f.Value.Name()) }

func (f *File) Read() Object {
	if f.State == UnknownFileState {
		// not opened in write mode, and no read function has been called yet
		// so set the file to _read
		f.State = ReadFileState
	} else if f.State != ReadFileState {
		return &Error{Message: "file not in 'read' state"}
	}
	data, err := io.ReadAll(f.Value)
	if err != nil {
		return ErrorFromGo(err)
	}
	defer f.Value.Close()
	return &String{Value: string(data)}
}

func (f *File) Readline() Object {
	if f.State == UnknownFileState {
		// not opened in write mode, and no read function has been called yet
		// so set the file to _readline
		f.State = ReadlineFileState
	} else if f.State != ReadlineFileState {
		return &Error{Message: "file not in 'readline' state"}
	}
	if f.scanner == nil {
		f.scanner = bufio.NewScanner(f.Value)
	}
	if !f.scanner.Scan() {
		return nil
	}
	return &String{Value: f.scanner.Text()}
}

func (f *File) Write(obj Object) Object {
	if f.State != WriteFileState {
		return &Error{Message: "file not in 'write' state"}
	}
	n, err := fmt.Fprint(f.Value, obj.Inspect())
	if err != nil {
		return ErrorFromGo(err)
	}
	return &Integer{Value: int64(n)}
}

func (f *File) Close() Object {
	if err := f.Value.Close(); err != nil {
		return ErrorFromGo(err)
	}
	return nil
}
