package organize

import (
	"hash/crc32"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/azak-azkaran/putio-go-aria2/utils"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/stretchr/testify/assert"
)

func TestCreateFolder(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stdout)

	path := "../test"
	_, err := os.Stat(path)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
	assert.True(t, CreateFolder(path))

	_, err = os.Stat(path)
	assert.NoError(t, err)
	assert.False(t, os.IsNotExist(err))

	err = os.Remove(path)
	assert.NoError(t, err)
}
func TestCompareFiles(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stdout)

	markedFiles = cmap.New()
	err := Copy("../testdata/output.json", "../output.json")
	assert.NoError(t, err)

	file, err := os.Stat("../testdata/output.json")
	assert.NoError(t, err)

	var putio PutIoFiles
	putio.Folder = "../test/blub"
	putio.Name = "output.json"
	putio.PutIoID = 23
	putio.CRC32 = "16ec90d5"
	putio.Size = file.Size()

	output := CompareFiles("../output.json", putio)
	assert.True(t, output)

	putio.CRC32 = "cca7c6b3"
	output = CompareFiles("../output.json", putio)
	assert.False(t, output)

	putio.CRC32 = "00188c02"
	output = CompareFiles("../output.json", putio)
	assert.False(t, output)

	err = os.Remove("../output.json")
	assert.NoError(t, err)
}

func TestOrganizeFolder(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stdout)
	folders := cmap.New()
	var conf Configuration

	path := "../testdata"
	_, err := os.Stat(path)
	assert.NoError(t, err)
	assert.False(t, os.IsNotExist(err))

	files := GoOrganizeFolder(path, folders, conf)
	assert.NotZero(t, files)
	for k, v := range files {
		utils.Info.Println("file:", k, "\t", v.Name())
		utils.Info.Println("Mode: ", v.Mode()&os.ModeSymlink == 0)
	}

}

func TestHandleFile(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stdout)

	markedFiles = cmap.New()
	var conf Configuration
	var putio PutIoFiles

	_, err := os.Stat("../test/output.json")
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))

	file, err := os.Stat("../testdata/output.json")
	assert.NoError(t, err)

	putio.Folder = "test/"
	putio.Name = "output.json"
	putio.PutIoID = 23
	putio.CRC32, _ = CreateCrc32("../testdata/output.json")
	putio.Size = file.Size()
	err = Copy("../testdata/output.json", "../output.json")
	assert.NoError(t, err)

	file, err = os.Stat("../output.json")
	assert.NoError(t, err)
	assert.False(t, os.IsNotExist(err))

	HandleFile(putio, file, "../", conf, false)

	_, err = os.Stat("../test/output.json")
	assert.NoError(t, err)
	assert.False(t, os.IsNotExist(err))

	err = os.Remove("../test/output.json")
	assert.NoError(t, err)

	err = os.Remove("../test")
	assert.NoError(t, err)
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
