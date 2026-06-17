package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

const parserSystemPrompt = "You are a government report assistant. Extract work entries from the provided Obsidian daily work log. Return ONLY valid JSON - no markdown fences, no explanation, no extra text. Just the raw JSON object."

var monthNames = []string{
	"january", "february", "march", "april", "may", "june",
	"july", "august", "september", "october", "november", "december",
}

func (a *App) parseReport(payload ParsePayload) (string, error) {
	providerSettings, err := a.providerSettings()
	if err != nil {
		return "", err
	}

	prompt := buildPrompt(payload.ReportType, payload.DateStart, payload.DateEnd, payload.Content)

	switch providerSettings.Provider {
	case "ollama":
		return a.parseWithOllama(providerSettings.OllamaURL, providerSettings.OllamaModel, prompt)
	case "gemini":
		key, err := a.loadSecret("geminiKey")
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(key) == "" {
			return "", fmt.Errorf("no Gemini API key found. Add it in Settings first")
		}
		return a.parseWithGemini(key, prompt)
	default:
		key, err := a.loadSecret("apiKey")
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(key) == "" {
			return "", fmt.Errorf("no API key found. Add your Claude API key in Settings first")
		}
		return a.parseWithAnthropic(key, prompt)
	}
}

type providerSettings struct {
	Provider    string
	OllamaURL   string
	OllamaModel string
}

func (a *App) providerSettings() (providerSettings, error) {
	settings, err := a.loadSettings()
	if err != nil {
		return providerSettings{}, err
	}

	getString := func(key, fallback string) string {
		value, _ := settings[key].(string)
		if strings.TrimSpace(value) == "" {
			return fallback
		}
		return value
	}

	return providerSettings{
		Provider:    getString("provider", "claude"),
		OllamaURL:   getString("ollamaUrl", "http://127.0.0.1:11434"),
		OllamaModel: getString("ollamaModel", "llama3"),
	}, nil
}

func (a *App) parseWithAnthropic(apiKey, prompt string) (string, error) {
	requestBody := map[string]any{
		"model":      "claude-sonnet-4-6",
		"max_tokens": 8192,
		"system":     parserSystemPrompt,
		"messages": []map[string]any{
			{"role": "user", "content": prompt},
		},
	}

	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequestWithContext(a.ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Claude API error %d: %s", resp.StatusCode, readAllBody(resp.Body))
	}
	defer resp.Body.Close()

	var payload struct {
		StopReason string `json:"stop_reason"`
		Content    []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if payload.StopReason == "max_tokens" {
		return "", fmt.Errorf("response was cut off - your notes may be too large. Try a narrower date range")
	}
	for _, part := range payload.Content {
		if part.Type == "text" {
			return strings.TrimSpace(part.Text), nil
		}
	}
	return "", fmt.Errorf("Claude returned no text content")
}

func (a *App) parseWithGemini(apiKey, prompt string) (string, error) {
	requestBody := map[string]any{
		"systemInstruction": map[string]any{
			"parts": []map[string]string{{"text": parserSystemPrompt}},
		},
		"contents": []map[string]any{
			{
				"parts": []map[string]string{{"text": prompt}},
			},
		},
	}

	body, _ := json.Marshal(requestBody)
	endpoint := "https://generativelanguage.googleapis.com/v1beta/models/gemini-3.1-flash-lite:generateContent?key=" + apiKey
	req, err := http.NewRequestWithContext(a.ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Gemini API error %d: %s", resp.StatusCode, readAllBody(resp.Body))
	}
	defer resp.Body.Close()

	var payload struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if len(payload.Candidates) == 0 || len(payload.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("Gemini returned no content")
	}
	return strings.TrimSpace(payload.Candidates[0].Content.Parts[0].Text), nil
}

func (a *App) parseWithOllama(ollamaURL, ollamaModel, prompt string) (string, error) {
	requestBody := map[string]any{
		"model": ollamaModel,
		"messages": []map[string]string{
			{"role": "system", "content": parserSystemPrompt},
			{"role": "user", "content": prompt},
		},
		"stream": false,
		"options": map[string]any{
			"num_ctx":     16384,
			"num_predict": 8192,
		},
	}

	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequestWithContext(a.ctx, http.MethodPost, strings.TrimRight(ollamaURL, "/")+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Ollama error %d: %s", resp.StatusCode, readAllBody(resp.Body))
	}
	defer resp.Body.Close()

	var payload struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	return strings.TrimSpace(payload.Message.Content), nil
}

func buildPrompt(reportType, dateStart, dateEnd, content string) string {
	filtered := filterContentByDateRange(content, dateStart, dateEnd)
	base := fmt.Sprintf("Date range: %s to %s\n\nWork log:\n%s\n\n", dateStart, dateEnd, filtered)

	switch reportType {
	case "AR":
		return base + fmt.Sprintf(`You are writing an official Accomplishment Report for a government IT employee. The audience is HR and supervisors - non-technical readers.

Rules for writing task descriptions:
- Rewrite each task as a clear, plain-language sentence that a non-technical reader can understand.
- Describe WHAT was done and WHY it matters, not HOW it was coded.
- Avoid: component names, file paths, CSS properties, hex colors, variable names, code syntax, or technical jargon.
- Good example: "Redesigned the employee payroll edit form with reorganized deduction headers for better usability."
- Bad example: "Replaced edlg-net-box, grouped GSIS/HDMF headers via scoped CSS flex layout."
- If multiple tasks happened on the same day, list each as a separate string in the "tasks" array.
- If the log says WFH, begin the sentence with "WFH: ".
- Leaves, holidays, and weekends should have an empty tasks array [].

Extract all work entries within the date range and return this exact JSON:
{
  "reportType": "AR",
  "dateStart": "%s",
  "dateEnd": "%s",
  "entries": [
    { "date": "YYYY-MM-DD", "tasks": ["Plain sentence describing what was accomplished."], "hoursWorked": 8 }
  ],
  "summary": "A 1-2 sentence plain-language summary of the period's accomplishments."
}`, dateStart, dateEnd)
	case "DTR":
		return base + fmt.Sprintf(`Extract daily time records. Write task descriptions in plain, non-technical language suitable for HR. Return this exact JSON:
{
  "reportType": "DTR",
  "dateStart": "%s",
  "dateEnd": "%s",
  "entries": [
    { "date": "YYYY-MM-DD", "timeIn": "08:00", "timeOut": "17:00", "hoursWorked": 8, "tasks": ["Plain sentence describing what was done."] }
  ],
  "totalHours": 0
}`, dateStart, dateEnd)
	default:
		return base + fmt.Sprintf(`Extract progress on ongoing tasks. Write all descriptions in plain, non-technical language suitable for a supervisor. Return this exact JSON:
{
  "reportType": "PRG",
  "dateStart": "%s",
  "dateEnd": "%s",
  "entries": [
    { "date": "YYYY-MM-DD", "tasks": ["Plain sentence describing what was done."], "hoursWorked": 8 }
  ],
  "ongoingTasks": [
    { "name": "Plain task name", "percentComplete": 75, "status": "In Progress", "nextSteps": ["Plain next step description."] }
  ],
  "summary": "A 1-2 sentence plain-language summary of overall progress."
}`, dateStart, dateEnd)
	}
}

func filterContentByDateRange(content, dateStart, dateEnd string) string {
	start, err1 := time.Parse("2006-01-02", dateStart)
	end, err2 := time.Parse("2006-01-02", dateEnd)
	if err1 != nil || err2 != nil {
		return content
	}
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	fallbackYear := detectNoteYear(content)
	rows := parseMarkdownTable(content, fallbackYear)
	if len(rows) > 0 {
		if text := tableRowsToText(rows, start, end); strings.TrimSpace(text) != "" {
			return text
		}
	}

	lines := strings.Split(content, "\n")
	var out []string
	inRange := false
	hasHeader := false

	for _, line := range lines {
		if headerDate, ok := parseHeaderDate(line); ok {
			hasHeader = true
			inRange = !headerDate.Before(start) && !headerDate.After(end)
		}
		if !hasHeader || inRange {
			out = append(out, line)
		}
	}

	filtered := strings.TrimSpace(strings.Join(out, "\n"))
	if len(filtered) > 100 {
		return filtered
	}
	return content
}

type markdownRow struct {
	Date    *time.Time
	DateRaw string
	Task    string
	Notes   string
}

func detectNoteYear(content string) int {
	match := regexp.MustCompile(`\b(20\d{2})\b`).FindStringSubmatch(content)
	if len(match) == 2 {
		if year, err := time.Parse("2006", match[1]); err == nil {
			return year.Year()
		}
	}
	return time.Now().Year()
}

func parseMarkdownTable(content string, fallbackYear int) []markdownRow {
	lines := strings.Split(content, "\n")
	var rows []markdownRow
	var headers []string

	for _, line := range lines {
		if !strings.HasPrefix(line, "|") {
			headers = nil
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue
		}

		cells := make([]string, 0, len(parts)-2)
		for i := 1; i < len(parts)-1; i++ {
			cells = append(cells, strings.TrimSpace(parts[i]))
		}

		separator := true
		for _, cell := range cells {
			if !regexp.MustCompile(`^[-:]+$`).MatchString(cell) {
				separator = false
				break
			}
		}
		if separator {
			continue
		}

		lower := make([]string, len(cells))
		for i, cell := range cells {
			lower[i] = strings.ToLower(cell)
		}
		if contains(lower, "task") && contains(lower, "date") {
			headers = lower
			continue
		}

		if len(headers) == 0 {
			continue
		}

		get := func(key string) string {
			for i, header := range headers {
				if header == key && i < len(cells) {
					return cells[i]
				}
			}
			return ""
		}

		task := get("task")
		if task == "" || task == "---" {
			continue
		}

		dateRaw := get("date")
		date := parseLooseDate(dateRaw, fallbackYear)
		rows = append(rows, markdownRow{
			Date:    date,
			DateRaw: dateRaw,
			Task:    task,
			Notes:   get("notes"),
		})
	}

	return rows
}

func tableRowsToText(rows []markdownRow, start, end time.Time) string {
	filtered := make([]markdownRow, 0, len(rows))
	for _, row := range rows {
		if row.Date == nil || (!row.Date.Before(start) && !row.Date.After(end)) {
			filtered = append(filtered, row)
		}
	}
	if len(filtered) == 0 {
		return ""
	}

	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].Date == nil || filtered[j].Date == nil {
			return i < j
		}
		return filtered[i].Date.Before(*filtered[j].Date)
	})

	chunks := make([]string, 0, len(filtered))
	for _, row := range filtered {
		dateValue := row.DateRaw
		if row.Date != nil {
			dateValue = row.Date.Format("2006-01-02")
		}
		lines := []string{
			"Date: " + dateValue,
			"Task: " + row.Task,
		}
		if strings.TrimSpace(row.Notes) != "" {
			lines = append(lines, "Notes: "+row.Notes)
		}
		chunks = append(chunks, strings.Join(lines, "\n"))
	}
	return strings.Join(chunks, "\n\n")
}

func parseLooseDate(value string, fallbackYear int) *time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	if parsed, err := time.Parse("2006-01-02", value); err == nil {
		return &parsed
	}

	monthMatch := regexp.MustCompile(`(?i)^(January|February|March|April|May|June|July|August|September|October|November|December)\s+(\d{1,2})(?:,?\s*(\d{4}))?$`).FindStringSubmatch(value)
	if len(monthMatch) > 0 {
		monthIndex := indexOfMonth(monthMatch[1])
		if monthIndex >= 0 {
			year := fallbackYear
			if monthMatch[3] != "" {
				fmt.Sscanf(monthMatch[3], "%d", &year)
			}
			day := 1
			fmt.Sscanf(monthMatch[2], "%d", &day)
			parsed := time.Date(year, time.Month(monthIndex+1), day, 0, 0, 0, 0, time.Local)
			return &parsed
		}
	}

	return nil
}

func parseHeaderDate(line string) (time.Time, bool) {
	if !strings.HasPrefix(strings.TrimSpace(line), "#") {
		return time.Time{}, false
	}

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\b(\d{4}-\d{2}-\d{2})\b`),
		regexp.MustCompile(`(?i)\b(January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{1,2},?\s+\d{4}\b`),
		regexp.MustCompile(`\b\d{1,2}/\d{1,2}/\d{4}\b`),
	}

	for _, pattern := range patterns {
		match := pattern.FindString(line)
		if match == "" {
			continue
		}
		layouts := []string{"2006-01-02", "January 2, 2006", "January 2 2006", "1/2/2006"}
		for _, layout := range layouts {
			if parsed, err := time.Parse(layout, match); err == nil {
				return parsed, true
			}
		}
	}
	return time.Time{}, false
}

func indexOfMonth(value string) int {
	value = strings.ToLower(strings.TrimSpace(value))
	for i, month := range monthNames {
		if month == value {
			return i
		}
	}
	return -1
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
