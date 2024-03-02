package spotle_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"
)

const artistTableCreationQuery = `
CREATE SEQUENCE IF NOT EXISTS public.artists_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

CREATE TABLE IF NOT EXISTS public.artists
(
    id integer NOT NULL DEFAULT nextval('artists_id_seq'::regclass),
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    monthly_listeners integer NOT NULL,
    last_checked timestamp without time zone NOT NULL,
    CONSTRAINT artists_pkey PRIMARY KEY (id)
);

ALTER SEQUENCE public.artists_id_seq
	OWNED BY artists.id;
`

func clearArtistTable() {
	a.DB.Exec("DELETE FROM artists")
	a.DB.Exec("ALTER SEQUENCE artists_id_seq RESTART WITH 1")
}

func TestGetNonExistentArtist(t *testing.T) {
	clearArtistTable()

	req, _ := http.NewRequest("GET", "/artist/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Artist not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Artist not found'. Got '%s'", m["error"])
	}
}

func TestCreateArtist(t *testing.T) {

	clearArtistTable()
	var name = "The Weeknd"
	var ml = 115683791
	var jsonStr = []byte(`{"name":` + name + `, "monthly_listeners":` + strconv.Itoa(ml) + `}`)
	req, _ := http.NewRequest("POST", "/artist", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != name {
		t.Errorf("Expected artist name to be %v. Got '%v'", name, m["name"])
	}

	if m["monthly_listeners"] != ml {
		t.Errorf("Expected artist's monthly listeners to be '%v'. Got '%v'", ml, m["monthly_listeners"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected artist ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetArtist(t *testing.T) {
	clearArtistTable()
	addArtists(1)

	req, _ := http.NewRequest("GET", "/artist/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateArtist(t *testing.T) {
	clearArtistTable()
	addArtists(1)

	req, _ := http.NewRequest("GET", "/artist/1", nil)
	response := executeRequest(req)
	var ogArtist map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &ogArtist)

	var name = "The Weekend"
	var ml = 1
	var jsonStr = []byte(`{"name":` + name + `, "monthly_listeners":` + strconv.Itoa(ml) + `}`)
	req, _ = http.NewRequest("PUT", "/artist/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != ogArtist["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", ogArtist["id"], m["id"])
	}

	if m["name"] == ogArtist["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", ogArtist["name"], m["name"], m["name"])
	}

	if m["monthly_listeners"] == ogArtist["monthly_listeners"] {
		t.Errorf("Expected the artist's monthly listeners to change from '%v' to '%v'. Got '%v'", ogArtist["monthly_listeners"], m["monthly_listeners"], m["monthly_listeners"])
	}

	if m["last_checked"] == ogArtist["last_checked"] {
		t.Errorf("Expected the artist's last check date and time to change from '%v' to '%v'. Got '%v'", ogArtist["last_checked"], m["last_checked"], m["last_checked"])
	}
}

func TestDeleteArtist(t *testing.T) {
	clearArtistTable()
	addArtists(1)

	req, _ := http.NewRequest("GET", "/artist/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/artist/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/artist/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func addArtists(count int) {
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO artists(name, monthly_listners, last_checked) VALUES($1, $2, $3)", "The Weeknd", 3200100, time.Now().UTC())
		a.DB.Exec("INSERT INTO artists(name, monthly_listners, last_checked) VALUES($1, $2, $3)", "Taylor Swift", 200100, time.Now().UTC())
	}
}
