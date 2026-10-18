package web

import (
	"bytes"
	"context"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/modules"
	"github.com/standblatt/server/internal/store"
)

const maxImageBytes = 2 << 20

// ImagePath identifies an image of a tenant's library in the URL.
type ImagePath struct {
	TenantID int64 `path:"tenantId"`
	ImageID  int64 `path:"imageId"`
}

// checkLayoutImages answers 422 when an image element of a print form shows an image the club does not have.
func (s *Server) checkLayoutImages(ctx context.Context, tenantID int64, in store.LayoutInput) error {
	missing, err := s.Store.MissingImages(ctx, tenantID, in.ImageIDs())
	if err != nil {
		return err
	}
	if len(missing) > 0 {
		return invalid("Image %d not found", missing[0])
	}
	return nil
}

// registerImages registers the image library of a club, used by print forms.
func (s *Server) registerImages(api huma.API) {
	type listOut struct{ Body []store.Image }
	huma.Register(api, huma.Operation{
		OperationID: "list-images", Metadata: modules.Metadata(modules.Forms), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/images", Tags: []string{"forms"},
	}, func(ctx context.Context, in *TenantPath) (*listOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		list, err := s.Store.ListImages(ctx, in.TenantID)
		return &listOut{Body: list}, err
	})

	type uploadIn struct {
		TenantPath
		RawBody huma.MultipartFormFiles[struct {
			File huma.FormFile `form:"file" contentType:"image/png,image/jpeg,application/octet-stream" required:"true"`
			Name string        `form:"name" required:"true"`
		}]
	}
	type imageOut struct{ Body store.Image }
	huma.Register(api, huma.Operation{
		OperationID: "create-image", Metadata: modules.Metadata(modules.Forms), Method: http.MethodPost, Path: "/api/tenants/{tenantId}/images", Tags: []string{"forms"},
		Summary: "Add a PNG or JPEG of at most 2 MB to the club's images", DefaultStatus: http.StatusCreated, MaxBodyBytes: maxImageBytes + 64<<10,
	}, func(ctx context.Context, in *uploadIn) (*imageOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		form := in.RawBody.Data()
		name := strings.TrimSpace(form.Name)
		if name == "" || len(name) > 100 {
			return nil, invalid("The name must be 1 to 100 characters long")
		}
		b, err := io.ReadAll(io.LimitReader(form.File, maxImageBytes+1))
		if err != nil {
			return nil, err
		}
		if len(b) > maxImageBytes {
			return nil, invalid("The image is larger than 2 MB")
		}
		cfg, kind, err := image.DecodeConfig(bytes.NewReader(b))
		if err != nil || (kind != "png" && kind != "jpeg") {
			return nil, invalid("The image must be a PNG or JPEG")
		}
		img, err := s.Store.CreateImage(ctx, in.TenantID, name, b, cfg.Width, cfg.Height)
		if errors.Is(err, store.ErrExists) {
			return nil, conflict("An image named %s already exists", name)
		}
		return &imageOut{Body: img}, err
	})

	type fileOut struct {
		ContentType  string `header:"Content-Type"`
		CacheControl string `header:"Cache-Control"`
		Body         []byte
	}
	huma.Register(api, huma.Operation{
		OperationID: "get-image", Metadata: modules.Metadata(modules.Forms), Method: http.MethodGet, Path: "/api/tenants/{tenantId}/images/{imageId}", Tags: []string{"forms"},
		Responses: map[string]*huma.Response{"200": {Description: "PNG or JPEG", Content: map[string]*huma.MediaType{"image/png": {}, "image/jpeg": {}}}},
	}, func(ctx context.Context, in *ImagePath) (*fileOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		b, err := s.Store.GetImage(ctx, in.TenantID, in.ImageID)
		if err != nil {
			return nil, notFound(err, "Image not found")
		}
		return &fileOut{ContentType: http.DetectContentType(b), CacheControl: "private, max-age=86400", Body: b}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-image", Metadata: modules.Metadata(modules.Forms), Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/images/{imageId}", Tags: []string{"forms"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *ImagePath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		err := s.Store.DeleteImage(ctx, in.TenantID, in.ImageID)
		if errors.Is(err, store.ErrInUse) {
			return nil, conflict("Print forms use this image")
		}
		return nil, notFound(err, "Image not found")
	})
}
