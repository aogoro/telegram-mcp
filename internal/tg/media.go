package tg

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"path/filepath"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/telegram/uploader"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/pkg/errors"
)

type SendPhotoArguments struct {
	Name    string `json:"name" jsonschema:"required,description=Name of the dialog"`
	Path    string `json:"path" jsonschema:"required,description=Absolute local path to the image file"`
	Caption string `json:"caption,omitempty" jsonschema:"description=Optional caption text"`
	ReplyTo int    `json:"reply_to,omitempty" jsonschema:"description=Message ID to reply to"`
}

type SendFileArguments struct {
	Name    string `json:"name" jsonschema:"required,description=Name of the dialog"`
	Path    string `json:"path" jsonschema:"required,description=Absolute local path to the file"`
	Caption string `json:"caption,omitempty" jsonschema:"description=Optional caption text"`
	ReplyTo int    `json:"reply_to,omitempty" jsonschema:"description=Message ID to reply to"`
}

// SendPhoto sends an image as an inline (compressed) Telegram photo.
func (c *Client) SendPhoto(args SendPhotoArguments) (*mcp.ToolResponse, error) {
	return c.sendMedia(args.Name, args.Path, args.Caption, args.ReplyTo, false)
}

// SendFile sends any file as a document, preserving the exact bytes.
func (c *Client) SendFile(args SendFileArguments) (*mcp.ToolResponse, error) {
	return c.sendMedia(args.Name, args.Path, args.Caption, args.ReplyTo, true)
}

func (c *Client) sendMedia(name, path, caption string, replyTo int, asDocument bool) (*mcp.ToolResponse, error) {
	var success bool
	client := c.T()

	if err := client.Run(context.Background(), func(ctx context.Context) (err error) {
		api := client.API()

		inputPeer, err := getInputPeerFromName(ctx, api, name)
		if err != nil {
			return fmt.Errorf("get inputPeer from name: %w", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("file not found(%s): %w", path, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("not a regular file: %s", path)
		}

		file, err := uploader.NewUploader(api).FromPath(ctx, path)
		if err != nil {
			return fmt.Errorf("upload file: %w", err)
		}

		var caps []message.StyledTextOption
		if caption != "" {
			caps = append(caps, styling.Plain(caption))
		}

		var media message.MediaOption
		if asDocument {
			mimeType := mime.TypeByExtension(filepath.Ext(path))
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}
			media = message.UploadedDocument(file, caps...).
				Filename(filepath.Base(path)).
				MIME(mimeType).
				ForceFile(true)
		} else {
			media = message.UploadedPhoto(file, caps...)
		}

		// To() returns *message.RequestBuilder; Reply() returns *message.Builder.
		// Branch instead of reassigning to keep types consistent.
		sender := message.NewSender(api).To(inputPeer)
		if replyTo > 0 {
			_, err = sender.Reply(replyTo).Media(ctx, media)
		} else {
			_, err = sender.Media(ctx, media)
		}
		if err != nil {
			return fmt.Errorf("send media: %w", err)
		}

		success = true
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to send media")
	}

	jsonData, err := json.Marshal(SendResponse{Success: success})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal response")
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(jsonData))), nil
}
