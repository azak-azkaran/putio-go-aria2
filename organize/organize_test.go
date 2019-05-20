package organize

import (
	"github.com/azak-azkaran/putio-go-aria2/utils"
	"github.com/orcaman/concurrent-map"
	"hash/crc32"
	"io"
	"os"
	"strconv"
	"testing"
)

func TestCreateFolder(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stdout)

	path := "../test"
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		t.Error("Folder is already there")
	}

	if !CreateFolder(path) {
		t.Error("Folder could not be created")
	}
	_, err = os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		t.Error("Folder is not created")
	}
	err = os.Remove(path)
	if err != nil {
		t.Error("Folder could not be removed")
	}
}
func TestCompareFiles(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stdout)

	err := Copy("../testdata/output.json", "../output.json")
	if err != nil {
		t.Error("testdata could not be copied")
	}

	file, err := os.Stat("../testdata/output.json")
	if err != nil {
		t.Error("testdata could not be copied")
	}

	var putio PutIoFiles
	putio.Folder = "../test/blub"
	putio.Name = "output.json"
	putio.PutIoID = 23
	putio.CRC32 = "16ec90d5"
	putio.Size = file.Size()

	output := CompareFiles("../output.json", putio)
	if !output {
		t.Error("Compared failed")
	}
	putio.CRC32 = "cca7c6b3"
	output = CompareFiles("../output.json", putio)
	if output {
		t.Error("Compared failed")
	}

	putio.CRC32 = "00188c02"
	output = CompareFiles("../output.json", putio)
	if output {
		t.Error("Compared failed even with padding")
	}

	err = os.Remove("../output.json")
	if err != nil {
		t.Error("File could not be removed")
	}
}

func TestOrganizeFolder(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stdout)
	folders := cmap.New()
	var conf Configuration

	path := "../testdata"
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		t.Error("Folder is already there")
	}
	files := GoOrganizeFolder(path, folders, conf)
	if len(files) == 0 {
		t.Error("No files found")
	}
}

func TestHandleFile(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stdout)

	var conf Configuration
	var putio PutIoFiles

	_, err := os.Stat("../test/output.json")
	if err != nil && !os.IsNotExist(err) {
		t.Error("File already moved")
	}

	file, err := os.Stat("../testdata/output.json")
	if err != nil {
		t.Error("testdata not available")
	}

	putio.Folder = "test/"
	putio.Name = "output.json"
	putio.PutIoID = 23
	putio.CRC32, _ = CreateCrc32("../testdata/output.json")
	putio.Size = file.Size()
	err = Copy("../testdata/output.json", "../output.json")
	if err != nil {
		t.Error("testdata could not be copied")
	}

	file, err = os.Stat("../output.json")
	if err != nil && os.IsNotExist(err) {
		t.Error("testdata was not copied")
	}
	HandleFile(putio, file, "../", conf, false)

	_, err = os.Stat("../test/output.json")
	if err != nil && !os.IsNotExist(err) {
		t.Error("File was not moved")
	}

	err = os.Remove("../test/output.json")
	if err != nil {
		t.Error("Error while removing File")
	}

	err = os.Remove("../test")
	if err != nil {
		t.Error("Error while removing Folder")
	}

}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}

func CreateCrc32(path string) (string, error) {
	offlineFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer offlineFile.Close()

	hash := crc32.NewIEEE()
	if _, err := io.Copy(hash, offlineFile); err != nil {
		return "", err
	}
	//Generate the hash
	hashInBytes := hash.Sum32()

	//Encode the hash to a string
	crc := strconv.FormatUint(uint64(hashInBytes), 16)
	return crc, nil
}
