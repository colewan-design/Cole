package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx        context.Context
	httpClient *http.Client
}

type FileInfo struct {
	Content  string `json:"content"`
	FilePath string `json:"filePath"`
	FileName string `json:"fileName"`
	DirName  string `json:"dirName"`
}

type DialogOptions struct {
	DefaultPath string `json:"defaultPath"`
}

type ParsePayload struct {
	Content    string `json:"content"`
	DateStart  string `json:"dateStart"`
	DateEnd    string `json:"dateEnd"`
	ReportType string `json:"reportType"`
}

type GenerateReportPayload struct {
	ParsedData   ParsedData `json:"parsedData"`
	UserInfo     UserInfo   `json:"userInfo"`
	EmployeeType string     `json:"employeeType"`
	OutputPath   string     `json:"outputPath"`
}

type ParsedData struct {
	ReportType   string        `json:"reportType"`
	DateStart    string        `json:"dateStart"`
	DateEnd      string        `json:"dateEnd"`
	Entries      []ReportEntry `json:"entries"`
	Summary      string        `json:"summary"`
	TotalHours   float64       `json:"totalHours"`
	OngoingTasks []OngoingTask `json:"ongoingTasks"`
}

type ReportEntry struct {
	Date        string   `json:"date"`
	Tasks       []string `json:"tasks"`
	HoursWorked float64  `json:"hoursWorked"`
	TimeIn      string   `json:"timeIn"`
	TimeOut     string   `json:"timeOut"`
}

type OngoingTask struct {
	Name            string   `json:"name"`
	PercentComplete float64  `json:"percentComplete"`
	Status          string   `json:"status"`
	NextSteps       []string `json:"nextSteps"`
}

type UserInfo struct {
	ReporterName    string `json:"reporterName"`
	ReporterNameSig string `json:"reporterNameSig"`
	Position        string `json:"position"`
	Office          string `json:"office"`
	SupervisorName  string `json:"supervisorName"`
	SupervisorPos1  string `json:"supervisorPos1"`
	SupervisorPos2  string `json:"supervisorPos2"`
}

func NewApp() *App {
	return &App{
		httpClient: &http.Client{Timeout: 2 * time.Minute},
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OpenFileDialog(options DialogOptions) (string, error) {
	dir := cleanDefaultDirectory(options.DefaultPath)
	if dir == "" {
		dir = a.getStoredPath("inputPath")
	}

	return wruntime.OpenFileDialog(a.ctx, wruntime.OpenDialogOptions{
		Title:            "Select Markdown file",
		DefaultDirectory: dir,
		Filters: []wruntime.FileFilter{
			{DisplayName: "Markdown Files (*.md)", Pattern: "*.md"},
		},
	})
}

func (a *App) ReadFile(filePath string) (*FileInfo, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, errors.New("no file path provided")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Content:  string(data),
		FilePath: filePath,
		FileName: filepath.Base(filePath),
		DirName:  filepath.Dir(filePath),
	}, nil
}

func (a *App) SaveFile(filePath string, encoded string) (string, error) {
	if strings.TrimSpace(filePath) == "" {
		return "", errors.New("no file path provided")
	}
	data, err := decodeFrontendBytes(encoded)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return "", err
	}
	return filePath, nil
}

func (a *App) ShowSaveDialog(options DialogOptions) (string, error) {
	defaultDir, defaultName := splitDefaultPath(options.DefaultPath)
	if defaultDir == "" {
		defaultDir = a.getStoredPath("outputPath")
	}

	return wruntime.SaveFileDialog(a.ctx, wruntime.SaveDialogOptions{
		Title:            "Save Word document",
		DefaultDirectory: defaultDir,
		DefaultFilename:  defaultName,
		Filters: []wruntime.FileFilter{
			{DisplayName: "Word Document (*.docx)", Pattern: "*.docx"},
		},
	})
}

func (a *App) OpenPath(target string) error {
	if strings.TrimSpace(target) == "" {
		return errors.New("no path provided")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	case "darwin":
		cmd = exec.Command("open", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}
	return cmd.Start()
}

func (a *App) OpenExternal(url string) error {
	if strings.TrimSpace(url) == "" {
		return errors.New("no url provided")
	}
	wruntime.BrowserOpenURL(a.ctx, url)
	return nil
}

func (a *App) GetSetting(key string) (any, error) {
	if isSecretKey(key) {
		return a.loadSecret(key)
	}
	settings, err := a.loadSettings()
	if err != nil {
		return nil, err
	}
	return settings[key], nil
}

func (a *App) SetSetting(key string, value any) error {
	if isSecretKey(key) {
		strValue, _ := value.(string)
		return a.saveSecret(key, strValue)
	}

	settings, err := a.loadSettings()
	if err != nil {
		return err
	}
	settings[key] = value
	return a.saveSettings(settings)
}

func (a *App) GetAllSettings() (map[string]any, error) {
	settings, err := a.loadSettings()
	if err != nil {
		return nil, err
	}
	apiKey, _ := a.loadSecret("apiKey")
	geminiKey, _ := a.loadSecret("geminiKey")
	settings["hasApiKey"] = strings.TrimSpace(apiKey) != ""
	settings["hasGeminiKey"] = strings.TrimSpace(geminiKey) != ""
	return settings, nil
}

func (a *App) ListOllamaModels(url string) ([]string, error) {
	endpoint := strings.TrimRight(strings.TrimSpace(url), "/")
	if endpoint == "" {
		endpoint = "http://127.0.0.1:11434"
	}

	req, err := http.NewRequestWithContext(a.ctx, http.MethodGet, endpoint+"/api/tags", nil)
	if err != nil {
		return []string{}, nil
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return []string{}, nil
	}
	defer resp.Body.Close()

	var payload struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return []string{}, nil
	}

	models := make([]string, 0, len(payload.Models))
	for _, model := range payload.Models {
		if strings.TrimSpace(model.Name) != "" {
			models = append(models, model.Name)
		}
	}
	return models, nil
}

func (a *App) ParseWithClaude(payload ParsePayload) (map[string]any, error) {
	raw, err := a.parseReport(payload)
	if err != nil {
		return nil, err
	}

	cleaned := strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(raw), "```"), "```json"))
	var result map[string]any
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		snippet := cleaned
		if len(snippet) > 300 {
			snippet = snippet[:300] + "..."
		}
		return nil, fmt.Errorf("LLM returned invalid JSON. Try a narrower date range.\n\nRaw response:\n%s", snippet)
	}

	return result, nil
}

func (a *App) GenerateReport(payload GenerateReportPayload) (string, error) {
	if strings.TrimSpace(payload.OutputPath) == "" {
		return "", errors.New("no output path provided")
	}
	if strings.TrimSpace(payload.ParsedData.DateStart) == "" || strings.TrimSpace(payload.ParsedData.DateEnd) == "" {
		return "", errors.New("parsed data is missing date range")
	}

	docBytes, err := buildReportDocx(payload.ParsedData, payload.UserInfo, payload.EmployeeType)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(payload.OutputPath, docBytes, 0o644); err != nil {
		return "", err
	}
	return payload.OutputPath, nil
}

func (a *App) RenderReport(payload GenerateReportPayload) (string, error) {
	if strings.TrimSpace(payload.ParsedData.DateStart) == "" || strings.TrimSpace(payload.ParsedData.DateEnd) == "" {
		return "", errors.New("parsed data is missing date range")
	}

	docBytes, err := buildReportDocx(payload.ParsedData, payload.UserInfo, payload.EmployeeType)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(docBytes), nil
}

func (a *App) getStoredPath(key string) string {
	settings, err := a.loadSettings()
	if err != nil {
		return ""
	}
	value, _ := settings[key].(string)
	return value
}

func (a *App) storageDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "Cole")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func (a *App) settingsPath() (string, error) {
	dir, err := a.storageDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.json"), nil
}

func (a *App) secretsPath() (string, error) {
	dir, err := a.storageDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "secrets.json"), nil
}

func (a *App) loadSettings() (map[string]any, error) {
	path, err := a.settingsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]any{}, nil
		}
		return nil, err
	}

	settings := map[string]any{}
	if len(data) == 0 {
		return settings, nil
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		return map[string]any{}, nil
	}
	return settings, nil
}

func (a *App) saveSettings(settings map[string]any) error {
	path, err := a.settingsPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (a *App) loadSecret(key string) (string, error) {
	secrets, err := a.loadSecrets()
	if err != nil {
		return "", err
	}
	encoded := secrets[key]
	if strings.TrimSpace(encoded) == "" {
		return "", nil
	}
	return decryptSecret(encoded)
}

func (a *App) saveSecret(key, value string) error {
	secrets, err := a.loadSecrets()
	if err != nil {
		return err
	}
	if strings.TrimSpace(value) == "" {
		delete(secrets, key)
		return a.saveSecrets(secrets)
	}

	encrypted, err := encryptSecret(value)
	if err != nil {
		return err
	}
	secrets[key] = encrypted
	return a.saveSecrets(secrets)
}

func (a *App) loadSecrets() (map[string]string, error) {
	path, err := a.secretsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, nil
		}
		return nil, err
	}

	secrets := map[string]string{}
	if len(data) == 0 {
		return secrets, nil
	}
	if err := json.Unmarshal(data, &secrets); err != nil {
		return map[string]string{}, nil
	}
	return secrets, nil
}

func (a *App) saveSecrets(secrets map[string]string) error {
	path, err := a.secretsPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func cleanDefaultDirectory(defaultPath string) string {
	if strings.TrimSpace(defaultPath) == "" {
		return ""
	}
	info, err := os.Stat(defaultPath)
	if err == nil && info.IsDir() {
		return defaultPath
	}
	return filepath.Dir(defaultPath)
}

func splitDefaultPath(defaultPath string) (string, string) {
	if strings.TrimSpace(defaultPath) == "" {
		return "", ""
	}
	return filepath.Dir(defaultPath), filepath.Base(defaultPath)
}

func isSecretKey(key string) bool {
	return key == "apiKey" || key == "geminiKey"
}

func decodeFrontendBytes(value string) ([]byte, error) {
	if strings.TrimSpace(value) == "" {
		return []byte{}, nil
	}
	return base64.StdEncoding.DecodeString(value)
}

func readAllBody(body io.ReadCloser) string {
	defer body.Close()
	data, _ := io.ReadAll(body)
	return string(data)
}
