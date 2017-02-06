package api100

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carbocation/interpose"
	"github.com/chaosvermittlung/funkloch-server/db/v100"
	"github.com/gorilla/mux"
)

func getEventRouter(prefix string) *interpose.Middleware {
	r, m := GetNewSubrouter(prefix)
	r.HandleFunc("/", postEventHandler).Methods("POST")
	r.HandleFunc("/list", listEventsHandler).Methods("GET")
	r.HandleFunc("/next", getNextEventHandler).Methods("GET")
	r.HandleFunc("/{ID}", getEventHandler).Methods("GET")
	r.HandleFunc("/{ID}", patchEventHandler).Methods("PATCH")
	r.HandleFunc("/{ID}", deleteEventHandler).Methods("DELETE")
	r.HandleFunc("/{ID}/Participiants", getEventParticipiantsHandler).Methods("GET")

	return m
}

func postEventHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var e db100.Event
	err = decoder.Decode(&e)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	err = e.Insert()
	if err != nil {
		apierror(w, r, "Error Inserting Event: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&e)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func listEventsHandler(w http.ResponseWriter, r *http.Request) {
	ee, err := db100.GetEvents()
	if err != nil {
		apierror(w, r, "Error fetching Events: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&ee)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	e := db100.Event{EventID: id}
	err = e.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Event: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&e)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getEventParticipiantsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	e := db100.Event{EventID: id}
	pp, err := e.GetParticipiants()
	if err != nil {
		apierror(w, r, "Error fetching Event Participiants: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}

	var result []eventParticipiantsResponse
	for _, p := range pp {
		u := db100.User{UserID: p.UserID}
		if err != nil {
			apierror(w, r, "Error fetching User: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		u.Password = ""
		var epr eventParticipiantsResponse
		epr.User = u
		epr.Arrival = p.Arrival
		epr.Departure = p.Departure
		result = append(result, epr)
	}

	j, err := json.Marshal(&result)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getNextEventHandler(w http.ResponseWriter, r *http.Request) {
	e, err := db100.GetNextEvent()
	if err != nil {
		apierror(w, r, "Error fetching Event: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&e)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func patchEventHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var event db100.Event
	err = decoder.Decode(&event)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	event.EventID = id
	err = event.Update()
	if err != nil {
		apierror(w, r, "Error updating Event: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&event)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_ADMIN)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var e db100.Event
	err = decoder.Decode(&e)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	err = e.Delete()
	if err != nil {
		apierror(w, r, "Error deleting Event: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
}