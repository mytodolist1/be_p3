package modul

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mytodolist1/be_p3/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func generateRandomFileName(originalFilename string) (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomFileName := fmt.Sprintf("%x%s", randomBytes, filepath.Ext(originalFilename))
	return randomFileName, nil
}

func SaveFileToGridFS(db *mongo.Database, col, filePath string) (model.GridFSFile, error) {
	// Open the local file
	file, err := os.Open(filePath)
	if err != nil {
		return model.GridFSFile{}, err
	}
	defer file.Close()

	// Generate a random filename
	randomFileName, err := generateRandomFileName(file.Name())
	if err != nil {
		return model.GridFSFile{}, err
	}

	// Create a GridFS bucket
	bucket, err := gridfs.NewBucket(
		db, options.GridFSBucket().SetName("mytodolistFiles"),
	)
	if err != nil {
		return model.GridFSFile{}, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return model.GridFSFile{}, err
	}

	fileSize := fileInfo.Size()

	// Upload the file to GridFS
	uploadStream, err := bucket.OpenUploadStream(randomFileName)
	if err != nil {
		return model.GridFSFile{}, err
	}
	defer uploadStream.Close()

	_, err = io.Copy(uploadStream, file)
	if err != nil {
		return model.GridFSFile{}, err
	}

	// Get the file's ID in GridFS
	fileID := uploadStream.FileID.(primitive.ObjectID)

	fmt.Printf("File %s uploaded successfully to GridFS with ID: %s\n", randomFileName, fileID.Hex())

	// Return a model.GridFSFile struct with relevant information
	return model.GridFSFile{
		ID:          fileID,
		FileName:    randomFileName,
		FileSize:    fileSize,
		ContentType: "application/octet-stream",
	}, nil
}
