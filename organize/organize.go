package organize

import (
	"bufio"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/azak-azkaran/putio-go-aria2/utils"
	cmap "github.com/orcaman/concurrent-map"
)

var markedFiles cmap.ConcurrentMap

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

func CompareFiles(path string, file PutIoFiles) bool {
	fileCrc, err := strconv.ParseInt(file.CRC32, 16, 64)
	if err != nil {
		utils.Error.Fatalln("Error while converting CRC from Putio: ", err)
		return false
	}
	var crc int64
	if markedFiles.Has(path) {
		val, b := markedFiles.Get(path)
		if b {
			crc = val.(int64)
		}
	} else {
		offlineFile, err := os.Open(path)
		if err != nil {
			utils.Error.Fatalln("Error while reading file: ", file.Name, "\tError: ", err)
			return false
		}
		hash := crc32.NewIEEE()
		//Copy the file in the interface
		if _, err := io.Copy(hash, offlineFile); err != nil {
			utils.Error.Fatalln("Error while copying file for CRC : ", file.Name, "\tError: ", err)
			return false
		}

		err = offlineFile.Close()
		if err != nil {
			utils.Error.Fatalln("Error while closing file: ", file.Name, "\tError: ", err)
			return false
		}
		//Generate the hash
		hashInBytes := hash.Sum32()

		//Encode the hash to a string
		crc = int64(hashInBytes)
	}
	utils.Info.Println("File: ", file.Name, "\tFolder: ", file.Folder)
	if fileCrc != crc {
		utils.Warning.Println("CRC values are different", "\nOnline CRC: ", strconv.FormatInt(fileCrc, 16), "\nOffline CRC: ", strconv.FormatInt(crc, 16))
		markedFiles.SetIfAbsent(path, crc)
		return false
	}
	if markedFiles.Has(path) {
		markedFiles.Remove(path)
	}
	return true
}

func RemoveOfflineFile(path string, stats os.FileInfo) bool {
	//if stats.Size() != file.Size {
	//	utils.Warning.Println("Size between files is different: ", "\nOnline: ", strconv.FormatInt(file.Size, 10), "\nOffline: ", strconv.FormatInt(stats.Size(), 10))
	//}

	//if stats.Size() == 0 {
	utils.Warning.Println("Trying to Remove: ", path)
	utils.Warning.Println("Lokal File size is: ", strconv.FormatInt(stats.Size(), 10), " removing file ")
	err := os.Remove(path)
	if err != nil {
		utils.Error.Fatalln("Error while removing offline file")
		return false
	}
	utils.Warning.Print("Offline File removed")
	//return false
	//}
	return true
}

func HandleFile(putFile PutIoFiles, file os.FileInfo, foldername string, conf Configuration, removeFile bool, moveFileToFolder string) {
	utils.Info.Println("Handling File: ", file.Name())
	completeFilepath := foldername + file.Name()
	completeFolderpath := moveFileToFolder + putFile.Folder
	compare := CompareFiles(completeFilepath, putFile)
	if compare {
		CreateFolder(completeFolderpath)
		if removeFile {
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
				err = os.Chown(completeFolderpath+"/"+putFile.Name, 1000, 1000)
				if err != nil {
					utils.Error.Fatalln("Error while changing Owner: ", putFile.Name, "\n", err)
					return
				}
			}
		}
	}
}

func GoOrganizeFolder(foldername string, folders cmap.ConcurrentMap, conf Configuration, moveFileToFolder string) []os.FileInfo {
	markedFiles = cmap.New()
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
		fi, err := os.Lstat(foldername + file.Name())
		if err != nil {
			utils.Error.Println("Error for Lstat: ", err)
			break
		}
		utils.Info.Println("Checking File: ", foldername+file.Name(),
			" IsDir", file.IsDir(),
			" ok:", ok,
			" Mode:", fi.Mode()&os.ModeSymlink == 0)

		if !file.IsDir() && ok && fi.Mode()&os.ModeSymlink == 0 {
			v := value.(PutIoFiles)
			HandleFile(v, file, foldername, conf, true, moveFileToFolder)
		}
	}
	utils.Info.Println("Removing corrupted files")
	for _, path := range markedFiles.Keys() {
		stats, _ := os.Stat(path)
		RemoveOfflineFile(path, stats)
	}
	return files
}
