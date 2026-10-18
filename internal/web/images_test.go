package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/standblatt/server/internal/store"
)

// TestImages checks the image library: upload rules, image elements in forms and deleting an image in use.
func TestImages(t *testing.T) {
	e := newEnv(t)
	tid := e.tenant("Images")
	e.user("boss", tid, store.RoleAdmin)
	boss := e.login("boss")
	base := fmt.Sprintf("/api/tenants/%d", tid)
	upload := func(name string, data []byte) (int, []byte) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		w.WriteField("name", name)
		part, _ := w.CreateFormFile("file", "logo.png")
		part.Write(data)
		w.Close()
		req, _ := http.NewRequest("POST", e.srv.URL+base+"/images", &body)
		req.Header.Set("Content-Type", w.FormDataContentType())
		resp, err := boss.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var out bytes.Buffer
		out.ReadFrom(resp.Body)
		return resp.StatusCode, out.Bytes()
	}
	var logo bytes.Buffer
	png.Encode(&logo, image.NewGray(image.Rect(0, 0, 40, 20)))

	code, b := upload("Logo", []byte("not an image"))
	expect(t, "not an image", code, 422, b)
	code, b = upload("Logo", logo.Bytes())
	expect(t, "upload", code, 201, b)
	var img store.Image
	json.Unmarshal(b, &img)
	if img.Width != 40 || img.Height != 20 {
		t.Fatalf("size: %+v", img)
	}
	code, b = upload("Logo", logo.Bytes())
	expect(t, "same name", code, 409, b)

	form := func(id int64) map[string]any {
		return map[string]any{"name": "Card", "width": 148, "height": 105, "elements": []map[string]any{{"type": "image", "image": id, "x": 10, "y": 10, "w": 30, "h": 15}}}
	}
	code, b = e.do(boss, "POST", base+"/layouts", form(img.ID+100))
	expect(t, "unknown image", code, 422, b)
	code, b = e.do(boss, "POST", base+"/layouts", form(img.ID))
	expect(t, "form with the image", code, 201, b)
	var l store.Layout
	json.Unmarshal(b, &l)
	resp, err := boss.Get(fmt.Sprintf("%s%s/layouts/%d/preview", e.srv.URL, base, l.ID))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("preview: %v %v", resp.StatusCode, err)
	}
	resp.Body.Close()
	resp, err = boss.Get(fmt.Sprintf("%s%s/images/%d", e.srv.URL, base, img.ID))
	if err != nil || resp.StatusCode != 200 || resp.Header.Get("Content-Type") != "image/png" {
		t.Fatalf("image: %v %v", resp.Header, err)
	}
	resp.Body.Close()

	code, b = e.do(boss, "DELETE", fmt.Sprintf("%s/images/%d", base, img.ID), nil)
	expect(t, "image in use", code, 409, b)
	if !strings.Contains(string(b), "use this image") {
		t.Fatalf("message: %s", b)
	}
	code, b = e.do(boss, "DELETE", fmt.Sprintf("%s/layouts/%d", base, l.ID), nil)
	expect(t, "delete the form", code, 204, b)
	code, b = e.do(boss, "PUT", base+"/modules", map[string]any{"modules": []string{}})
	expect(t, "images keep print forms on", code, 409, b)
	code, b = e.do(boss, "DELETE", fmt.Sprintf("%s/images/%d", base, img.ID), nil)
	expect(t, "delete the image", code, 204, b)
}
