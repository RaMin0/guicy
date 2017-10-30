package job

// StudentData job
var StudentData studentData = "student-data"

type studentData string

func (studentData) Name() string {
	return "Student Data"
}
func (studentData) Description() string {
	return "Given a student's ID, I can fetch you his/her full name and GUC email address."
}
func (studentData) Exec(payload interface{}) (interface{}, error) {
	return map[string]interface{}{
		"name":  "Ramy Moustafa Mohamed Aboul Naga",
		"email": "ramy.naga@student.guc.edu.eg",
	}, nil
}
func (studentData) Inputs() []map[string]string {
	return []map[string]string{
		map[string]string{
			"id":    "id",
			"type":  "text",
			"label": "ID",
			"hint":  "Ex: 13-8994",
		},
	}
}
func (studentData) Outputs() []map[string]string {
	return []map[string]string{
		map[string]string{
			"id":    "name",
			"type":  "text",
			"label": "Name",
		},
		map[string]string{
			"id":    "email",
			"type":  "text",
			"label": "Email",
		},
	}
}
