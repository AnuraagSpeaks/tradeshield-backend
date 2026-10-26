package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "tradeshield-backend/internal/pkg/pinelabs"
    "tradeshield-backend/internal/pkg/response"
)

type GSTLookupReq struct {
    GSTIN string `json:"gstin"`
}

type GSTLookupRes struct {
    GSTIN      string `json:"gstin"`
    LegalName  string `json:"legal_name"`
    TradeName  string `json:"trade_name,omitempty"`
    Address    string `json:"address"`
    Status     string `json:"status"`
    StatusCode string `json:"status_code"`
    Message    string `json:"message,omitempty"`
}

// GSTLookup handles POST /api/v1/gst/lookup
func (h *Handler) GSTLookup(w http.ResponseWriter, r *http.Request) {
    var req GSTLookupReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    gstin := strings.ToUpper(strings.TrimSpace(req.GSTIN))
    if gstin == "" {
        response.Error(w, http.StatusBadRequest, "gstin is required")
        return
    }
    if len(gstin) != 15 {
        response.Error(w, http.StatusBadRequest, "Invalid GSTIN length: must be exactly 15 alphanumeric characters")
        return
    }
    gstResp, err := pinelabs.LookupGST(gstin)
    if err != nil {
        response.Error(w, http.StatusBadRequest, fmt.Sprintf("GST lookup failed: %v", err))
        return
    }
    res := GSTLookupRes{
        GSTIN:      gstResp.GSTIN,
        LegalName:  gstResp.LegalName,
        TradeName:  gstResp.TradeName,
        Address:    gstResp.Address,
        Status:     gstResp.Status,
        StatusCode: gstResp.StatusCode,
        Message:    gstResp.Message,
    }
    response.JSON(w, http.StatusOK, res, "GST lookup successful")
}
