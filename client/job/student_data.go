package job

// StudentData var
var StudentData = func(payload interface{}) (interface{}, error) {
	return map[string]interface{}{
		"name":  "Ramy Moustafa Mohamed Aboul Naga",
		"email": "ramy.naga@student.guc.edu.eg",
	}, nil
}
