package organize

import (
	"bufio"
	"github.com/azak-azkaran/putio-go-aria2/utils"
	cmap "github.com/orcaman/concurrent-map"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func Read(filename string) string {
	oauthToken := ""
	file, err := os.Open(filename)
	if err != nil {
		utils.Error.Fatalln("could not read file", err)
		panic(err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			oauthToken = text
		}
	}
	return oauthToken
}

func CreateFolder(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		utils.Info.Println("Creating folder: ", path)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			utils.Error.Fatalln("Error while creating folder", err)
			return false
		}
		return true
	}
	return true
}

func CompareFiles(path string, file PutIoFiles) int {
	offline_file, err := os.Open(path)
	if err != nil {
		utils.Error.Fatalln("Error while reading file: ", file.Name, "\tError: ", err)
		return -1
	}
	defer offline_file.Close()

	hash := crc32.NewIEEE()

	//Copy the file in the interface
	if _, err := io.Copy(hash, offline_file); err != nil {
		utils.Error.Fatalln("Error while reading file: ", file.Name, "\tError: ", err)
	}
	//Generate the hash
	hashInBytes := hash.Sum32()

	//Encode the hash to a string
	crc := strconv.FormatUint(uint64(hashInBytes), 16)

	utils.Info.Println("File: ", file.Name, "\tFolder: ", file.Folder)
	if strings.Compare(crc, file.CRC32) != 0 {
		utils.Warning.Println("CRC values are different", "\nCRC: ", file.CRC32, "\nCRC: ", crc)
		if strings.Contains(file.CRC32, crc) && len(file.CRC32)-1 == len(crc) {
			utils.Warning.Println("but crc value is contained")
			return 0
		}
		stats, _ := os.Stat(path)
		if stats.Size() != file.Size {
			utils.Warning.Println("Size between files is different: ", "\nOnline: ", strconv.FormatInt(file.Size, 10), "\nOffline: ", strconv.FormatInt(stats.Size(), 10))
			err := os.Remove(path)
			if err != nil {
				utils.Error.Fatalln("Error while removing offline file")
			}
			utils.Warning.Print("Offline File removed")
			return -1
		}
		return 1
	}
	return 0
}

func HandleFile(putFile PutIoFiles, file os.FileInfo, foldername string, conf Configuration) {
	complete_filepath := foldername + file.Name()
	complete_folderpath := foldername + putFile.Folder
	compare := CompareFiles(complete_filepath, putFile)
	if compare != -1 {
		CreateFolder(complete_folderpath)
		if compare == 0 {
			RemoveOnlineFile(conf, putFile)
		}

		newfolder, err := os.Stat(complete_folderpath)
		if err != nil {
			utils.Error.Fatalln("Error Folder missing will not move file: ", err)
			return
		} else {
			if len(putFile.Folder) != 0 && newfolder.IsDir() {
				utils.Info.Println("Moving to: ", complete_folderpath+"/"+putFile.Name)
				err := os.Rename(complete_filepath, complete_folderpath+"/"+putFile.Name)
				if err != nil {
					utils.Error.Fatalln("Error while moving File: ", putFile.Name, "\n", err)
					return
				}
			}
		}
	}
}

func OrganizeFolder(foldername string, folders cmap.ConcurrentMap, conf Configuration) {
	files, err := ioutil.ReadDir(foldername)
	if err != nil {
		utils.Error.Fatalln(err)
	}

	utils.Info.Println("Checking Files on Disk")
	for _, file := range files {
		value, ok := folders.Get(file.Name())
		if !file.IsDir() && ok {
			v := value.(PutIoFiles)
			HandleFile(v, file, foldername, conf)
		}
	}
}
