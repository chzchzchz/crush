package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"charm.land/fantasy"
	"github.com/charmbracelet/crush/internal/filepathext"
	"github.com/charmbracelet/crush/internal/filetracker"
	"github.com/charmbracelet/crush/internal/permission"
)

type ReadParams struct {
	FilePath string `json:"file_path" description:"The path to the file to read"`
}

type ReadPermissionsParams struct {
	FilePath string `json:"file_path"`
}

const ReadToolName = "read"

func NewReadTool(permissions permission.Service, workingDir string, filetracker filetracker.Service) fantasy.AgentTool {
	return fantasy.NewAgentTool(
		ReadToolName,
		"Read the entire contents of a file by path and return its text content.",
		func(ctx context.Context, params ReadParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			if params.FilePath == "" {
				return fantasy.NewTextErrorResponse("file_path is required"), nil
			}

			filePath := filepathext.SmartJoin(workingDir, params.FilePath)

			absWorkingDir, err := filepath.Abs(workingDir)
			if err != nil {
				return fantasy.ToolResponse{}, fmt.Errorf("error resolving working directory: %w", err)
			}
			absFilePath, err := filepath.Abs(filePath)
			if err != nil {
				return fantasy.ToolResponse{}, fmt.Errorf("error resolving file path: %w", err)
			}

			relPath, err := filepath.Rel(absWorkingDir, absFilePath)
			isOutsideWorkDir := err != nil || strings.HasPrefix(relPath, "..")

			sessionID := GetSessionFromContext(ctx)
			if sessionID == "" {
				return fantasy.ToolResponse{}, fmt.Errorf("session ID is required for reading files")
			}

			if isOutsideWorkDir {
				granted, permErr := permissions.Request(
					ctx,
					permission.CreatePermissionRequest{
						SessionID:   sessionID,
						Path:        absFilePath,
						ToolCallID:  call.ID,
						ToolName:    ReadToolName,
						Action:      "read",
						Description: fmt.Sprintf("Read file outside working directory: %s", absFilePath),
						Params:      ReadPermissionsParams{FilePath: params.FilePath},
					},
				)
				if permErr != nil {
					return fantasy.ToolResponse{}, permErr
				}
				if !granted {
					return NewPermissionDeniedResponse(), nil
				}
			}

			fileInfo, err := os.Stat(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					return fantasy.NewTextErrorResponse(fmt.Sprintf("File not found: %s", filePath)), nil
				}
				return fantasy.ToolResponse{}, fmt.Errorf("error accessing file: %w", err)
			}
			if fileInfo.IsDir() {
				return fantasy.NewTextErrorResponse(fmt.Sprintf("Path is a directory, not a file: %s", filePath)), nil
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				return fantasy.ToolResponse{}, fmt.Errorf("error reading file: %w", err)
			}

			if !utf8.ValidString(string(content)) {
				return fantasy.NewTextErrorResponse("File content is not valid UTF-8"), nil
			}

			filetracker.RecordRead(ctx, sessionID, filePath)

			return fantasy.NewTextResponse(string(content)), nil
		},
	)
}
