package terraform

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/cloudskiff/driftctl/mocks"
	terraformError "github.com/cloudskiff/driftctl/pkg/terraform/error"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

func TestProviderInstallerInstallDoesNotExist(t *testing.T) {

	assert := assert.New(t)
	fakeTmpHome := t.TempDir()

	expectedSubFolder := fmt.Sprintf("/.driftctl/plugins/%s_%s", runtime.GOOS, runtime.GOARCH)

	config := ProviderConfig{
		Key:     "aws",
		Version: "3.19.0",
	}

	mockDownloader := mocks.ProviderDownloaderInterface{}
	mockDownloader.On("Download", config.GetDownloadUrl(), path.Join(fakeTmpHome, expectedSubFolder)).Return(nil)

	installer := ProviderInstaller{
		downloader: &mockDownloader,
		config:     config,
		homeDir:    fakeTmpHome,
	}

	providerPath, err := installer.Install()
	mockDownloader.AssertExpectations(t)

	assert.Nil(err)
	assert.Equal(path.Join(fakeTmpHome, expectedSubFolder, config.GetBinaryName()), providerPath)

}

func TestProviderInstallerInstallAlreadyExist(t *testing.T) {

	assert := assert.New(t)
	fakeTmpHome := t.TempDir()
	expectedSubFolder := fmt.Sprintf("/.driftctl/plugins/%s_%s", runtime.GOOS, runtime.GOARCH)
	err := os.MkdirAll(path.Join(fakeTmpHome, expectedSubFolder), 0755)
	if err != nil {
		t.Error(err)
	}

	config := ProviderConfig{
		Key:     "aws",
		Version: "3.19.0",
	}

	_, err = os.Create(path.Join(fakeTmpHome, expectedSubFolder, config.GetBinaryName()))
	if err != nil {
		t.Error(err)
	}

	mockDownloader := mocks.ProviderDownloaderInterface{}

	installer := ProviderInstaller{
		downloader: &mockDownloader,
		config:     config,
		homeDir:    fakeTmpHome,
	}

	providerPath, err := installer.Install()
	mockDownloader.AssertExpectations(t)

	assert.Nil(err)
	assert.Equal(path.Join(fakeTmpHome, expectedSubFolder, config.GetBinaryName()), providerPath)

}

func TestProviderInstallerInstallAlreadyExistButIsDirectory(t *testing.T) {

	assert := assert.New(t)
	fakeTmpHome := t.TempDir()
	expectedSubFolder := fmt.Sprintf("/.driftctl/plugins/%s_%s", runtime.GOOS, runtime.GOARCH)

	config := ProviderConfig{
		Key:     "aws",
		Version: "3.19.0",
	}

	invalidDirPath := path.Join(fakeTmpHome, expectedSubFolder, config.GetBinaryName())
	err := os.MkdirAll(invalidDirPath, 0755)
	if err != nil {
		t.Error(err)
	}

	mockDownloader := mocks.ProviderDownloaderInterface{}

	installer := ProviderInstaller{
		downloader: &mockDownloader,
		config:     config,
		homeDir:    fakeTmpHome,
	}

	providerPath, err := installer.Install()
	mockDownloader.AssertExpectations(t)

	assert.Empty(providerPath)
	assert.NotNil(err)
	assert.Equal(
		fmt.Sprintf(
			"found directory instead of provider binary in %s",
			invalidDirPath,
		),
		err.Error(),
	)

}

// Ensure that if a provider exists with a postfix (_x5) we properly detect it
func TestProviderInstallerInstallPostfixIsHandler(t *testing.T) {

	assert := assert.New(t)
	fakeTmpHome := t.TempDir()
	expectedSubFolder := fmt.Sprintf("/.driftctl/plugins/%s_%s", runtime.GOOS, runtime.GOARCH)
	err := os.MkdirAll(path.Join(fakeTmpHome, expectedSubFolder), 0755)
	if err != nil {
		t.Error(err)
	}

	config := ProviderConfig{
		Key:     "aws",
		Version: "3.19.0",
	}

	_, err = os.Create(path.Join(fakeTmpHome, expectedSubFolder, config.GetBinaryName()+"_x5"))
	if err != nil {
		t.Fatal(err)
	}

	mockDownloader := mocks.ProviderDownloaderInterface{}

	installer := ProviderInstaller{
		downloader: &mockDownloader,
		config:     config,
		homeDir:    fakeTmpHome,
	}

	providerPath, err := installer.Install()
	mockDownloader.AssertExpectations(t)

	assert.Nil(err)
	assert.Equal(path.Join(fakeTmpHome, expectedSubFolder, config.GetBinaryName()+"_x5"), providerPath)

}

func TestProviderInstallerVersionDoesNotExist(t *testing.T) {

	assert := assert.New(t)

	config := ProviderConfig{
		Key:     "aws",
		Version: "666.666.666",
	}

	mockDownloader := mocks.ProviderDownloaderInterface{}
	mockDownloader.On("Download", mock.Anything, mock.Anything).Return(terraformError.ProviderNotFoundError{})

	installer := ProviderInstaller{
		downloader: &mockDownloader,
		config:     config,
	}

	_, err := installer.Install()

	assert.Equal("Provider version 666.666.666 does not exist", err.Error())
}

func TestProviderInstallerWithConfigDirectory(t *testing.T) {

	assert := assert.New(t)
	fakeTmpHome := t.TempDir()

	expectedSubFolder := fmt.Sprintf("/.driftctl/plugins/%s_%s", runtime.GOOS, runtime.GOARCH)

	config := ProviderConfig{
		Key:       "aws",
		Version:   "3.19.0",
		ConfigDir: fakeTmpHome,
	}

	mockDownloader := mocks.ProviderDownloaderInterface{}
	mockDownloader.On("Download", config.GetDownloadUrl(), path.Join(fakeTmpHome, expectedSubFolder)).Return(nil)

	installer, _ := NewProviderInstaller(config)
	installer.downloader = &mockDownloader

	providerPath, err := installer.Install()
	mockDownloader.AssertExpectations(t)

	assert.Nil(err)
	assert.Equal(path.Join(fakeTmpHome, expectedSubFolder, config.GetBinaryName()), providerPath)

}

func TestProviderInstallerWithConfigDirectoryCantCreate(t *testing.T) {

	assert := assert.New(t)
	fakeTmpHome := t.TempDir()

	expectedSubFolder := fmt.Sprintf("/.driftctl/plugins/%s_%s", runtime.GOOS, runtime.GOARCH)

	config := ProviderConfig{
		Key:       "aws",
		Version:   "3.19.0",
		ConfigDir: fakeTmpHome,
	}

	_ = os.Chmod(fakeTmpHome, 0555)

	mockDownloader := mocks.ProviderDownloaderInterface{}
	mockDownloader.On("Download", config.GetDownloadUrl(), path.Join(fakeTmpHome, expectedSubFolder)).Return(nil)
	mockInstaller := mocks.HomeDirInterface{}
	mockInstaller.On("getProviderDirectory").Return("")

	installer, _ := NewProviderInstaller(config)
	installer.downloader = &mockDownloader

	providerPath, err := installer.Install()

	assert.Equal("", providerPath)
	assert.True(err != nil)

}

func TestProviderInstallerWithConfigDirectoryCantWrite(t *testing.T) {

	assert := assert.New(t)
	fakeTmpHome := t.TempDir()

	expectedSubFolder := fmt.Sprintf("/.driftctl/plugins/%s_%s", runtime.GOOS, runtime.GOARCH)

	config := ProviderConfig{
		Key:       "aws",
		Version:   "3.19.0",
		ConfigDir: fakeTmpHome,
	}

	mockDownloader := mocks.ProviderDownloaderInterface{}
	mockDownloader.On("Download", config.GetDownloadUrl(), path.Join(fakeTmpHome, expectedSubFolder)).Return(nil)

	installer, _ := NewProviderInstaller(config)
	installer.downloader = &mockDownloader

	joinedPath := path.Join(installer.getProviderDirectory())
	_ = os.MkdirAll(joinedPath, 0000)

	providerPath, err := installer.Install()

	assert.Equal("", providerPath)
	assert.EqualErrorf(err, fmt.Sprintf("can't write in configuration directory: open %s/.is_directory_writable: permission denied", joinedPath), "")
}
