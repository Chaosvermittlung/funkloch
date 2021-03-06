package api100

import (
	"encoding/json"
	"net/http"
	"strconv"

	db100 "github.com/Chaosvermittlung/funkloch-server/pkg/db/v100"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
)

func getFaultRouter(prefix string) *interpose.Middleware {
	r, m := GetNewSubrouter(prefix)
	r.HandleFunc("/", postFaultHandler).Methods("POST")
	r.HandleFunc("/list", listFaultsHandler).Methods("GET")
	r.HandleFunc("/{ID}", getFaultHandler).Methods("GET")
	r.HandleFunc("/{ID}", patchFaultHandler).Methods("PATCH")
	r.HandleFunc("/{ID}", deleteFaultHandler).Methods("DELETE")
	return m
}

func postFaultHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var f db100.Fault
	err = decoder.Decode(&f)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	if (f.Status < db100.FaultStatusNew) || (f.Status > db100.FaultStatusUnfixable) {
		apierror(w, r, "FaultStatus out of bound", http.StatusBadRequest, ERROR_JSONERROR)
		return
	}

	err = f.Insert()
	if err != nil {
		apierror(w, r, "Error Inserting Fault: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&f)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func listFaultsHandler(w http.ResponseWriter, r *http.Request) {
	ff, err := db100.GetFaults()
	if err != nil {
		apierror(w, r, "Error fetching Faults: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	var res []faultResponse
	for _, f := range ff {
		var fr faultResponse
		fr.Fault = f
		var it db100.Item
		it.ItemID = f.ItemID
		err = it.GetDetails()
		if err != nil {
			apierror(w, r, "Error fetching Fault StoreItem: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		fr.Code = it.Code
		var eq db100.Equipment
		eq.EquipmentID = it.EquipmentID
		err = eq.GetDetails()
		if err != nil {
			apierror(w, r, "Error fetching Fault Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		fr.Name = eq.Name
		res = append(res, fr)
	}
	j, err := json.Marshal(&res)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getFaultHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	f := db100.Fault{FaultID: id}
	err = f.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Fault: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&f)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func patchFaultHandler(w http.ResponseWriter, r *http.Request) {
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
	var fa db100.Fault
	err = decoder.Decode(&fa)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	fa.FaultID = id
	err = fa.Update()
	if err != nil {
		apierror(w, r, "Error updating Fault: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&fa)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func deleteFaultHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_ADMIN)
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
	f := db100.Fault{FaultID: id}
	err = f.Delete()
	if err != nil {
		apierror(w, r, "Error deleting Fault: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
}
