package modul

import (
	"crypto/rand"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v56/github"
	"github.com/mytodolist1/be_p3/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
)

func SaveFileToGithub(usernameGhp, emailGhp, repoGhp, path string, r *http.Request) (string, error) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("error 1: %s", err)
	}
	defer file.Close()

	// Generate a random filename
	randomFileName, err := generateRandomFileName(handler.Filename)
	if err != nil {
		return "", fmt.Errorf("error 2: %s", err)
	}

	// Read the content of the file into a byte slice
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error 5: %s", err)
	}

	access_token := os.Getenv("GITHUB_ACCESS_TOKEN")
	if access_token == "" {
		return "", fmt.Errorf("error access token: %s", err)
	}

	// Initialize GitHub client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: access_token},
	)
	tc := oauth2.NewClient(r.Context(), ts)
	client := github.NewClient(tc)

	// Create a new repository file
	_, _, err = client.Repositories.CreateFile(r.Context(), usernameGhp, repoGhp, path+"/"+randomFileName, &github.RepositoryContentFileOptions{
		Message:   github.String("Add new file"),
		Content:   fileContent,
		Committer: &github.CommitAuthor{Name: github.String(usernameGhp), Email: github.String(emailGhp)},
	})
	if err != nil {
		return "", fmt.Errorf("error 6: %s", err)
	}

	imageUrl := "https://raw.githubusercontent.com/" + usernameGhp + "/" + repoGhp + "/main/" + path + "/" + randomFileName

	return imageUrl, nil
}

func generateRandomFileName(originalFilename string) (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomFileName := fmt.Sprintf("%x%s", randomBytes, filepath.Ext(originalFilename))
	return randomFileName, nil
}

func SaveFileToGridFS(db *mongo.Database, col string, r *http.Request) (model.GridFSFile, error) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		return model.GridFSFile{}, fmt.Errorf("error 1: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return model.GridFSFile{}, fmt.Errorf("error 2: %s", err)
	}

	fileSize := len(data)

	// Generate a random filename
	randomFileName, err := generateRandomFileName(handler.Filename)
	if err != nil {
		return model.GridFSFile{}, fmt.Errorf("error 3: %s", err)
	}

	// Create a new file in GridFS
	bucket, err := gridfs.NewBucket(
		db,
		options.GridFSBucket().SetName(col),
	)
	if err != nil {
		return model.GridFSFile{}, fmt.Errorf("error 4: %s", err)
	}

	uploadStream, err := bucket.OpenUploadStream(randomFileName)
	if err != nil {
		return model.GridFSFile{}, fmt.Errorf("error 5: %s", err)
	}

	defer uploadStream.Close()

	_, err = uploadStream.Write(data)
	if err != nil {
		return model.GridFSFile{}, fmt.Errorf("error 6: %s", err)
	}

	// Content Type based on file extension (modify as needed)
	contentType := mime.TypeByExtension(filepath.Ext(handler.Filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return model.GridFSFile{
		ID:          primitive.NewObjectID(),
		FileName:    randomFileName,
		FileSize:    int64(fileSize),
		ContentType: contentType,
	}, nil
}
