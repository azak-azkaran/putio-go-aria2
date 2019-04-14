package organize

import (
	"github.com/azak-azkaran/putio-go-aria2/utils"
	"io"
	"os"
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
	if output != 0 {
		t.Error("Compared failed")
	}
	putio.CRC32 = "cca7c6b3"
	output = CompareFiles("../output.json", putio)
	if output != 1 {
		t.Error("Compared failed")
	}

	err = os.Remove("../output.json")
	if err != nil {
		t.Error("File could not be removed")
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
