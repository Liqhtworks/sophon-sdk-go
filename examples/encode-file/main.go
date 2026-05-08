package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	sophon "github.com/Liqhtworks/sophon-sdk-go"
	"github.com/Liqhtworks/sophon-sdk-go/helpers"
	"github.com/google/uuid"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./examples/encode-file /path/to/video.mov")
		os.Exit(2)
	}

	apiKey := os.Getenv("SOPHON_API_KEY")
	if apiKey == "" {
		panic("SOPHON_API_KEY is required")
	}

	baseURL := os.Getenv("SOPHON_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.liqhtworks.xyz"
	}

	cfg := sophon.NewConfiguration()
	cfg.Servers = sophon.ServerConfigurations{{URL: baseURL}}
	cfg.AddDefaultHeader("Authorization", "Bearer "+apiKey)
	client := sophon.NewAPIClient(cfg)

	ctx := context.Background()
	uploads := helpers.NewUploadsClient(client.UploadsAPI)
	jobs := helpers.NewJobsClient(client.JobsAPI)

	inputPath := os.Args[1]
	reader, size, closer, err := helpers.OpenFileForUpload(inputPath)
	if err != nil {
		panic(err)
	}
	defer closer()

	mimeType := "video/mp4"
	if strings.EqualFold(filepath.Ext(inputPath), ".mov") {
		mimeType = "video/quicktime"
	}

	upload, err := helpers.UploadFile(
		ctx,
		uploads,
		reader,
		size,
		filepath.Base(inputPath),
		mimeType,
		helpers.UploadFileOptions{
			Concurrency: 4,
			OnProgress: func(p helpers.UploadProgress) {
				fmt.Printf("upload %d/%d\n", p.PartsDone, p.PartsTotal)
			},
		},
	)
	if err != nil {
		panic(err)
	}

	job, _, err := client.JobsAPI.CreateJob(ctx).
		IdempotencyKey(uuid.NewString()).
		CreateJobRequest(sophon.CreateJobRequest{
			Source:  helpers.JobSource.Upload(upload.UploadID),
			Profile: sophon.SOPHON_ESPRESSO,
		}).
		Execute()
	if err != nil {
		panic(err)
	}
	fmt.Println("created", job.GetId())

	final, err := helpers.WaitForJob(ctx, jobs, job.GetId(), helpers.WaitForJobOptions{
		Timeout: 30 * time.Minute,
		OnProgress: func(j *helpers.Job) {
			fmt.Printf("job %s: %s\n", j.ID, j.Status)
		},
	})
	if err != nil {
		panic(err)
	}
	if final.Status != "completed" {
		panic("job ended in " + final.Status)
	}

	if err := downloadOutput(baseURL, apiKey, final.ID, "sophon-output.mp4"); err != nil {
		panic(err)
	}
	fmt.Println("wrote sophon-output.mp4")
}

func downloadOutput(baseURL, apiKey, jobID, outputPath string) error {
	client := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", baseURL+"/v1/jobs/"+jobID+"/output", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusFound {
		return fmt.Errorf("expected output redirect, got %d", res.StatusCode)
	}

	location := res.Header.Get("Location")
	if location == "" {
		return fmt.Errorf("missing output redirect")
	}

	downloadURL, err := resolveURL(baseURL, location)
	if err != nil {
		return err
	}

	dl, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer dl.Body.Close()
	if dl.StatusCode < 200 || dl.StatusCode >= 300 {
		return fmt.Errorf("download failed: %d", dl.StatusCode)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, dl.Body)
	return err
}

func resolveURL(baseURL, location string) (string, error) {
	loc, err := url.Parse(location)
	if err != nil {
		return "", err
	}
	if loc.IsAbs() {
		return loc.String(), nil
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(loc).String(), nil
}
