package d7024e

/*
* TODO: Remove all fmt's!
*/

import(
  "crypto/sha1"
  "log"
  "reflect"
  "os"
  "io/ioutil"
  //"strconv"
)

type Storage struct{}

/*
* A list of stored files, change accordingly later.
*/
var Files []file

type file struct{
  Name []uint8
  Text []byte
}

/*
* When a file misses the update, it is no longer prioriticed and
* moved to Memory from RAM.
*/
func (storage *Storage) MoveToMemory (name []byte){
  file := file{}
  compareName := storage.HashFile(name)
  for i, v := range Files {
    if reflect.DeepEqual(v.Name, compareName) {
      file.Name = v.Name
      file.Text = v.Text

      //Delete file in Files
      Files = append(Files[:i], Files[i+1:]...)

      //Insert into Memory
      storage.Memory(name, file.Text)
      //break out of for loop.
      break
    }
  }
}

/*
* Look if the RAM storage include a certain file, if so return file, else
check Memory if it's there and return.
* TODO: Files that are requested within a timer, are refreshed in main Memory, rest are stored somewhere else.
* TODO: Include timers for each file within main Memory such that they are discarded from main Memory when the timer runs out. Then purge files that are not used for "very long" time.
*/
func (storage *Storage) Search(name []byte) file{
  returnedFile := file{}
  compareName := storage.HashFile(name)

  //Check local list if file exist, else check Memory.
  for _, v := range Files {
    if reflect.DeepEqual(v.Name, compareName) {
      returnedFile.Name = v.Name
      returnedFile.Text = v.Text
      //break out of for loop.
      break
    }
  }
  if returnedFile.Name == nil{
    // Check if memory has file.
    returnedFile = storage.ReadMemory(name)
  }
  return returnedFile
}

/*
* Read Memory and see if a file is there, if so, move file to RAM
and return it.
*/
func (storage *Storage) ReadMemory(name []byte) file {
  content, err := ioutil.ReadFile("/tmp/" + string(name))
	if err != nil {
		log.Fatal(err)
	}
  hashedName := storage.HashFile(name)
  returnedFile := file{hashedName, []byte(content)}

  /*
  * When a file is retrieved from Memory, move it from Memory to RAM.
  */
  storage.RAM(name,[]byte(content))
  path := "/tmp/" + string(name)
  os.Remove(path) // clean up temp

  return returnedFile
}

/*
* Store a file into RAM, does not return anything.
*/
func (storage *Storage) RAM(name []byte, text []byte){
  fileName := storage.HashFile(name)
  newFile := file{fileName, text}
  Files = append(Files, newFile)
  return
}

/*
* Store a file into Memory, does not return anything.
*/
func (storage *Storage) Memory(name []byte, text []byte) {
  file := "/tmp/" + string(name)
  err := ioutil.WriteFile(file, text, 0644)
  if err != nil{
    panic(err)
  }
}

/*
* Convert a file name of type []byte into []uint8 (SHA-1)!
*/
func (storage *Storage) HashFile(name []byte) []uint8 {
  hashing := sha1.New()
  return hashing.Sum(name)
}
