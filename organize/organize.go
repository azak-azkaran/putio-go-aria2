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
	offlineFile, err := os.Open(path)
	if err != nil {
		utils.Error.Fatalln("Error while reading file: ", file.Name, "\tError: ", err)
		return -1
	}
	hash := crc32.NewIEEE()
	//Copy the file in the interface
	if _, err := io.Copy(hash, offlineFile); err != nil {
		utils.Error.Fatalln("Error while copying file for CRC : ", file.Name, "\tError: ", err)
		return -1
	}

	err = offlineFile.Close()
	if err != nil {
		utils.Error.Fatalln("Error while closing file: ", file.Name, "\tError: ", err)
		return -1
	}
	//Generate the hash
	hashInBytes := hash.Sum32()

	//Encode the hash to a string
	crc := int64(hashInBytes)

	fileCrc, err := strconv.ParseInt(file.CRC32, 16, 64)
	if err != nil {
		utils.Error.Fatalln("Error while converting CRC from Putio: ", err)
		return -1
	}
	utils.Info.Println("File: ", file.Name, "\tFolder: ", file.Folder)
	if fileCrc != crc {
		utils.Warning.Println("CRC values are different", "\nOnline CRC: ", strconv.FormatInt(fileCrc, 16), "\nOffline CRC: ", strconv.FormatInt(crc, 16))

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

func HandleFile(putFile PutIoFiles, file os.FileInfo, foldername string, conf Configuration, removeFile bool) {
	completeFilepath := foldername + file.Name()
	completeFolderpath := foldername + putFile.Folder
	compare := CompareFiles(completeFilepath, putFile)
	if compare != -1 {
		CreateFolder(completeFolderpath)
		if compare == 0 && removeFile {
			RemoveOnlineFile(conf, putFile)
		}

		newfolder, err := os.Stat(completeFolderpath)
		if err != nil {
			utils.Error.Fatalln("Error Folder missing will not move file: ", err)
			return
		} else {
			if len(putFile.Folder) != 0 && newfolder.IsDir() {
				utils.Info.Println("Moving to: ", completeFolderpath+"/"+putFile.Name)
				err := os.Rename(completeFilepath, completeFolderpath+"/"+putFile.Name)
				if err != nil {
					utils.Error.Fatalln("Error while moving File: ", putFile.Name, "\n", err)
					return
				}
			}
		}
	}
}

func GoOrganizeFolder(foldername string, folders cmap.ConcurrentMap, conf Configuration) []os.FileInfo {
	files, err := ioutil.ReadDir(foldername)
	if err != nil {
		utils.Error.Fatalln(err)
		return nil
	}

	if !strings.HasSuffix(foldername, "/") {
		foldername = foldername + "/"
	}

	utils.Info.Println("Checking Files on Disk")
	for _, file := range files {
		value, ok := folders.Get(file.Name())
		if !file.IsDir() && ok {
			v := value.(PutIoFiles)
			HandleFile(v, file, foldername, conf, true)
		}
	}
	return files
}
