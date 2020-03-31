package handlers

//func TestCreateNewServiceRequestRoute(t *testing.T) {
//	router := setupRouter()
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("POST", "/serviceRequest/CreateOrder", nil)
//	router.ServeHTTP(w, req)
//
//	bodyStr := w.Body.String()
//	var jsonResp models.ServiceRequest
//	json.Unmarshal([]byte(bodyStr), &jsonResp)
//
//	assert.Equal(t, 200, w.Code)
//	assert.Equal(t, "CreateOrder", jsonResp.WorkflowName, fmt.Sprintf("The expected name was CreateOrder but we got %s", jsonResp.WorkflowName))
//	assert.Equal(t, 16, len(jsonResp.ID), fmt.Sprintf("The expected length was 16 but the value was %s with length %d", jsonResp.ID, len(jsonResp.ID)))
//	assert.Equal(t, models.STATUS_NEW, jsonResp.Status, fmt.Sprintf("The expected status was NEW but we got %s", jsonResp.Status))
//}
