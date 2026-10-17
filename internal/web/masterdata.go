package web

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/standblatt/server/internal/store"
)

// ClubPath identifies a club of a tenant in the URL.
type ClubPath struct {
	TenantID int64 `path:"tenantId"`
	ClubID   int64 `path:"clubId"`
}

// ShooterPath identifies a shooter of a tenant in the URL.
type ShooterPath struct {
	TenantID  int64 `path:"tenantId"`
	ShooterID int64 `path:"shooterId"`
}

// ShooterClubPath identifies a club assignment of a shooter in the URL.
type ShooterClubPath struct {
	TenantID  int64 `path:"tenantId"`
	ShooterID int64 `path:"shooterId"`
	ClubID    int64 `path:"clubId"`
}

// ShooterCreate is a new shooter with the clubs to assign: existing ones by ID and new ones by name.
type ShooterCreate struct {
	store.ShooterInput
	ClubIDs   []int64  `json:"clubIds" required:"false"`
	ClubNames []string `json:"clubNames" required:"false" doc:"Clubs to find by exact name or create"`
}

// ImportReport is the outcome of a shooter import. Nothing is imported when Errors is not empty.
type ImportReport struct {
	store.ImportResult
	Errors []ImportLineError `json:"errors"`
}

const maxImportBytes = 10 << 20

// normalizeShooter trims the shooter input, normalizes license and birthday and requires first and last name.
func normalizeShooter(in *store.ShooterInput) error {
	in.LicenseNumber = normalizeLicense(in.LicenseNumber)
	in.Firstname = strings.TrimSpace(in.Firstname)
	in.Lastname = strings.TrimSpace(in.Lastname)
	in.Postalcode = strings.TrimSpace(in.Postalcode)
	in.City = strings.TrimSpace(in.City)
	if in.Firstname == "" || in.Lastname == "" {
		return invalid("First and last name are required")
	}
	b, err := parseBirthday(strings.TrimSpace(in.Birthday))
	if err != nil {
		return invalid(err.Error())
	}
	in.Birthday = b
	return nil
}

// shooterConflict maps store errors of shooter writes to 409 (duplicate license), 422 (unknown club) or 404.
func shooterConflict(err error) error {
	switch {
	case errors.Is(err, store.ErrExists):
		return conflict("This license number already exists")
	case errors.Is(err, store.ErrInUse):
		return invalid("Club not found")
	}
	return notFound(err, "Shooter not found")
}

// registerMasterData registers the club and shooter endpoints, including the CSV import.
func (s *Server) registerMasterData(api huma.API) {
	type clubsOut struct{ Body []store.Club }
	huma.Register(api, huma.Operation{
		OperationID: "list-clubs", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/clubs", Tags: []string{"clubs"},
	}, func(ctx context.Context, in *TenantPath) (*clubsOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		list, err := s.Store.ListClubs(ctx, in.TenantID)
		return &clubsOut{Body: list}, err
	})

	type clubIn struct {
		TenantPath
		Body store.ClubInput
	}
	type clubOut struct{ Body store.Club }
	huma.Register(api, huma.Operation{
		OperationID: "create-club", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/clubs", Tags: []string{"clubs"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *clubIn) (*clubOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		in.Body.Name = strings.TrimSpace(in.Body.Name)
		c, err := s.Store.CreateClub(ctx, in.TenantID, in.Body)
		if errors.Is(err, store.ErrExists) {
			return nil, conflict("A club with this name already exists")
		}
		return &clubOut{Body: c}, err
	})

	type clubUpdateIn struct {
		ClubPath
		Body store.ClubInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "update-club", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/clubs/{clubId}", Tags: []string{"clubs"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *clubUpdateIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		in.Body.Name = strings.TrimSpace(in.Body.Name)
		err := s.Store.UpdateClub(ctx, in.TenantID, in.ClubID, in.Body)
		if errors.Is(err, store.ErrExists) {
			return nil, conflict("A club with this name already exists")
		}
		return nil, notFound(err, "Club not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-club", Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/clubs/{clubId}", Tags: []string{"clubs"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *ClubPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		err := s.Store.DeleteClub(ctx, in.TenantID, in.ClubID)
		if errors.Is(err, store.ErrInUse) {
			return nil, conflict("Registrations refer to this club; it cannot be deleted")
		}
		return nil, notFound(err, "Club not found")
	})

	type searchIn struct {
		TenantPath
		Q      string `query:"q" doc:"License number, or (parts of) first and last name"`
		ClubID int64  `query:"clubId" required:"false"`
		Limit  int    `query:"limit" minimum:"1" maximum:"200" default:"50"`
		Offset int    `query:"offset" minimum:"0" required:"false" doc:"Rows to skip, for paging"`
	}
	type shootersOut struct {
		Total int `header:"X-Total-Count" doc:"Number of shooters found without limit"`
		Body  []store.Shooter
	}
	huma.Register(api, huma.Operation{
		OperationID: "search-shooters", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/shooters", Tags: []string{"shooters"},
	}, func(ctx context.Context, in *searchIn) (*shootersOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		list, err := s.Store.SearchShooters(ctx, in.TenantID, in.Q, in.ClubID, in.Limit, in.Offset)
		if err != nil {
			return nil, err
		}
		total, err := s.Store.CountShooters(ctx, in.TenantID, in.Q, in.ClubID)
		return &shootersOut{Total: total, Body: list}, err
	})

	type shooterOut struct{ Body store.Shooter }
	huma.Register(api, huma.Operation{
		OperationID: "get-shooter", Method: http.MethodGet, Path: "/api/tenants/{tenantId}/shooters/{shooterId}", Tags: []string{"shooters"},
	}, func(ctx context.Context, in *ShooterPath) (*shooterOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		sh, err := s.Store.GetShooter(ctx, in.TenantID, in.ShooterID)
		return &shooterOut{Body: sh}, notFound(err, "Shooter not found")
	})

	type createIn struct {
		TenantPath
		Body ShooterCreate
	}
	huma.Register(api, huma.Operation{
		OperationID: "create-shooter", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/shooters", Tags: []string{"shooters"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *createIn) (*shooterOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		if err := normalizeShooter(&in.Body.ShooterInput); err != nil {
			return nil, err
		}
		clubIDs := in.Body.ClubIDs
		for _, name := range in.Body.ClubNames {
			if name = strings.TrimSpace(name); name == "" {
				continue
			}
			id, err := s.Store.EnsureClub(ctx, in.TenantID, name)
			if err != nil {
				return nil, err
			}
			clubIDs = append(clubIDs, id)
		}
		id, err := s.Store.CreateShooter(ctx, in.TenantID, in.Body.ShooterInput, clubIDs)
		if err != nil {
			return nil, shooterConflict(err)
		}
		sh, err := s.Store.GetShooter(ctx, in.TenantID, id)
		return &shooterOut{Body: sh}, err
	})

	type updateIn struct {
		ShooterPath
		Body store.ShooterInput
	}
	huma.Register(api, huma.Operation{
		OperationID: "update-shooter", Method: http.MethodPut, Path: "/api/tenants/{tenantId}/shooters/{shooterId}", Tags: []string{"shooters"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *updateIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		if err := normalizeShooter(&in.Body); err != nil {
			return nil, err
		}
		return nil, shooterConflict(s.Store.UpdateShooter(ctx, in.TenantID, in.ShooterID, in.Body))
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-shooter", Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/shooters/{shooterId}", Tags: []string{"shooters"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *ShooterPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		err := s.Store.DeleteShooter(ctx, in.TenantID, in.ShooterID)
		if errors.Is(err, store.ErrInUse) {
			return nil, conflict("The shooter has registrations and cannot be deleted")
		}
		return nil, notFound(err, "Shooter not found")
	})

	type addClubIn struct {
		ShooterPath
		Body struct {
			ClubID   int64  `json:"clubId,omitempty" required:"false"`
			ClubName string `json:"clubName,omitempty" required:"false" doc:"Find by exact name or create, used when clubId is 0"`
		}
	}
	huma.Register(api, huma.Operation{
		OperationID: "add-shooter-club", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/shooters/{shooterId}/clubs", Tags: []string{"shooters"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *addClubIn) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		clubID := in.Body.ClubID
		if clubID == 0 {
			name := strings.TrimSpace(in.Body.ClubName)
			if name == "" {
				return nil, invalid("Club is missing")
			}
			var err error
			if clubID, err = s.Store.EnsureClub(ctx, in.TenantID, name); err != nil {
				return nil, err
			}
		}
		return nil, notFound(s.Store.AddShooterClub(ctx, in.TenantID, in.ShooterID, clubID), "Shooter or club not found")
	})

	huma.Register(api, huma.Operation{
		OperationID: "remove-shooter-club", Method: http.MethodDelete, Path: "/api/tenants/{tenantId}/shooters/{shooterId}/clubs/{clubId}", Tags: []string{"shooters"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *ShooterClubPath) (*struct{}, error) {
		if _, _, err := s.access(ctx, in.TenantID, false); err != nil {
			return nil, err
		}
		return nil, notFound(s.Store.RemoveShooterClub(ctx, in.TenantID, in.ShooterID, in.ClubID), "Assignment not found")
	})

	type importIn struct {
		TenantPath
		RawBody huma.MultipartFormFiles[struct {
			File huma.FormFile `form:"file" contentType:"text/csv,text/plain,application/vnd.ms-excel,application/octet-stream" required:"true"`
		}]
	}
	type importOut struct{ Body ImportReport }
	huma.Register(api, huma.Operation{
		OperationID: "import-shooters", Method: http.MethodPost, Path: "/api/tenants/{tenantId}/shooters/import", Tags: []string{"shooters"},
		Summary:      "Import shooters and their clubs from CSV (firstname, lastname, optional license, birthday, postalcode, city, club)",
		MaxBodyBytes: maxImportBytes,
	}, func(ctx context.Context, in *importIn) (*importOut, error) {
		if _, _, err := s.access(ctx, in.TenantID, true); err != nil {
			return nil, err
		}
		f := in.RawBody.Data().File
		rows, lineErrs, err := parseShooterCSV(f)
		if err != nil {
			return nil, invalid(err.Error())
		}
		out := &importOut{Body: ImportReport{Errors: lineErrs}}
		if out.Body.Errors == nil {
			out.Body.Errors = []ImportLineError{}
		}
		if len(lineErrs) > 0 {
			return out, nil
		}
		res, err := s.Store.ImportShooters(ctx, in.TenantID, rows)
		if err != nil {
			return nil, err
		}
		out.Body.ImportResult = res
		s.Log.Info("shooters imported", "tenant", in.TenantID, "created", res.Created, "updated", res.Updated, "clubs", res.ClubsCreated)
		return out, nil
	})
}
