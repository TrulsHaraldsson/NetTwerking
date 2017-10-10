package d7024e


import (
	"crypto/sha1"
	"io/ioutil"
	"os"
	"reflect"
	//"fmt"
)

/*
* A list of stored files, change accordingly later.
 */
type Storage struct {
	Files []file
}

func NewStorage() Storage {
	return Storage{Files: []file{}}
}

type file struct {
	Name []uint8
	Text []byte
}

/*
* When a file misses the update, it is no longer prioritized and
* moved to Memory from RAM.
 */
func (storage *Storage) MoveToMemory(name []byte) {
	file := file{}
	//compareName := storage.HashFile(name)
	for i, v := range storage.Files {
		if reflect.DeepEqual(v.Name, name) {
			file.Name = v.Name
			file.Text = v.Text

			//Delete file in Files
			storage.Files = append(storage.Files[:i], storage.Files[i+1:]...)

			//Insert into Memory
			storage.Memory(name, file.Text)
			//break out of for loop.
			break
		}
	}
}

func (storage *Storage) deleteFromRam(name string) bool {
	for i, v := range storage.Files {
		if reflect.DeepEqual(v.Name, name) {
			//Delete file in Files
			storage.Files = append(storage.Files[:i], storage.Files[i+1:]...)
			return true
		}
	}
	return false
}

/*
* Look if the RAM storage include a certain file, if so return file, else
check Memory if it's there and return.
*/
func (storage *Storage) Search(name []byte) *file {
	returnedFile := storage.ReadRAM(name)
	if returnedFile == nil {
		// Check if memory has file.
		returnedFile = storage.ReadMemory(name)
	}
	return returnedFile
}

/*
* Check if file is in RAM
 */
func (storage *Storage) ReadRAM(name []byte) *file {
	file := file{}
	//compareName := storage.HashFile(name)
	for _, v := range storage.Files {
		//fmt.Println("\n",string(v.Name),"\n", string(name),"\n")
		if reflect.DeepEqual(v.Name, name) {
			file.Name = v.Name
			file.Text = v.Text
			return &file
		}
	}
	return nil
}

/*
* Read Memory and see if a file is there, if so, add file to RAM
and return it.
*/
func (storage *Storage) ReadMemory(name []byte) *file {
	filename := "../newfiles/" + string(name)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	//fmt.Printf("File contents: %s", content,"\n")
	returnedFile := &file{name, []byte(content)}

	/*
	 * When a file is retrieved from Memory, add it from Memory to RAM for "caching"
	 */
	storage.RAM(name, []byte(content))
	return returnedFile
}

/*
* Store a file into RAM, does not return anything.

 */
func (storage *Storage) RAM(name []byte, text []byte) {
	//fileName := storage.HashFile(name)
	newFile := file{name, text}
	//fmt.Println("Name : ", name, "\nString(Name)", string(name))
	storage.Files = append(storage.Files, newFile)
	//fmt.Println("Files : \n",storage.Files,"\n")
	return
}

/*
* Store a file into Memory, does not return anything.
 */
func (storage *Storage) Memory(name []byte, text []byte) {
	filename := "../newfiles/" + string(name)
	err2 := ioutil.WriteFile(filename, text, 0644)
	if err2 != nil {
		panic(err2)
	}
	// Check dir after creation to confirm correctness!
	/*
	   path = "./../newfiles"
	   files, err = ioutil.ReadDir(path)
	   if err != nil {
	     fmt.Println(err)
	   }
	   fmt.Println("\n")
	   for _, f := range files {
	     fmt.Println(f.Name())
	   }
	   fmt.Println("After new file !")
	*/
}

/*
* Convert a file name of type []byte into []uint8 (SHA-1)!
 */
func (storage *Storage) HashFile(name []byte) []uint8 {
	hashing := sha1.New()
	return hashing.Sum(name)
}

func (storage *Storage) DeleteFile(name string) {
	storage.deleteFromRam(name)
	path := "../newfiles/" + string(name)
	os.Remove(path) // clean up
}
