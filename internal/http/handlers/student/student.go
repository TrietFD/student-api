package student

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/TrietFD/student-api/internal/types"
	"github.com/TrietFD/student-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(s *types.Student) error {
	query := `INSERT INTO students (name, email, age) VALUES (?, ?, ?)`
	
	result, err := r.db.Exec(query, s.Name, s.Email, s.Age)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows were inserted")
	}

	return nil
}

func (r *Repository) FindByEmail(email string) (*types.Student, error) {
	query := `SELECT id, name, email, age FROM students WHERE email = ?`
	
	student := &types.Student{}
	err := r.db.QueryRow(query, email).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return student, nil
}

func (r *Repository) FindAll() ([]types.Student, error) {
	query := `SELECT id, name, email, age FROM students`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func New(db *sql.DB) http.HandlerFunc {
	repo := NewRepository(db)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createStudent(w, r, repo)
		case http.MethodGet:
			getStudents(w, r, repo)
		default:
			response.WriteJson(w, http.StatusMethodNotAllowed, response.GeneralError(fmt.Errorf("method not allowed")))
		}
	}
}

func createStudent(w http.ResponseWriter, r *http.Request, repo *Repository) {
	var student types.Student

	err := json.NewDecoder(r.Body).Decode(&student)
	if errors.Is(err, io.EOF) {
	   response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
	   return
	}

	if err != nil {
	   response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
	   return
	}

	// request validation
	if err := validator.New().Struct(student); err != nil {
		validateErrs := err.(validator.ValidationErrors)
		response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
		return
	}

	// Kiểm tra email đã tồn tại chưa
	existingStudent, err := repo.FindByEmail(student.Email)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("database error")))
		return
	}

	if existingStudent != nil {
		response.WriteJson(w, http.StatusConflict, response.GeneralError(fmt.Errorf("student with this email already exists")))
		return
	}

	// Tạo sinh viên mới
	err = repo.Create(&student)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to create student")))
		return
	}

	// Trả về response thành công
	response.WriteJson(w, http.StatusCreated, map[string]string{"message": "Student created successfully"})
}

func getStudents(w http.ResponseWriter, r *http.Request, repo *Repository) {
	students, err := repo.FindAll()
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to retrieve students")))
		return
	}

	// Nếu không có sinh viên nào
	if len(students) == 0 {
		response.WriteJson(w, http.StatusOK, []types.Student{})
		return
	}

	// Trả về danh sách sinh viên
	response.WriteJson(w, http.StatusOK, students)
}