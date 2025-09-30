package internal

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

// Validator структура для валидации данных
type Validator struct {
	maxFileSize       int64
	allowedMimeTypes  map[string]bool
	allowedExtensions map[string]bool
	allowedBlurTypes  map[string]bool
	allowedObjects    map[string]bool
}

// NewValidator создает новый валидатор
func NewValidator(maxFileSize int64) *Validator {
	return &Validator{
		maxFileSize: maxFileSize,
		allowedMimeTypes: map[string]bool{
			// Изображения
			"image/jpeg": true,
			"image/jpg":  true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
			"image/bmp":  true,
			"image/tiff": true,
			// Видео
			"video/mp4":       true,
			"video/avi":       true,
			"video/mov":       true,
			"video/wmv":       true,
			"video/flv":       true,
			"video/webm":      true,
			"video/mkv":       true,
			"video/quicktime": true,
		},
		allowedExtensions: map[string]bool{
			// Изображения
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
			".webp": true,
			".bmp":  true,
			".tiff": true,
			".tif":  true,
			// Видео
			".mp4":  true,
			".avi":  true,
			".mov":  true,
			".wmv":  true,
			".flv":  true,
			".webm": true,
			".mkv":  true,
		},
		allowedBlurTypes: map[string]bool{
			"gaussian": true,
			"motion":   true,
			"pixelate": true,
		},
		allowedObjects: map[string]bool{
			"face":             true, // лицо
			"person":           true, // человек
			"bicycle":          true, // велосипед
			"car":              true, // автомобиль
			"motorcycle":       true, // мотоцикл
			"airplane":         true, // самолет
			"bus":              true, // автобус
			"train":            true, // поезд
			"truck":            true, // грузовик
			"boat":             true, // лодка
			"traffic light":    true, // светофор
			"fire hydrant":     true, // пожарный гидрант
			"stop sign":        true, // знак стоп
			"parking meter":    true, // парковочный счетчик
			"bench":            true, // скамейка
			"bird":             true, // птица
			"cat":              true, // кот
			"dog":              true, // собака
			"horse":            true, // лошадь
			"sheep":            true, // овца
			"cow":              true, // корова
			"elephant":         true, // слон
			"bear":             true, // медведь
			"zebra":            true, // зебра
			"giraffe":          true, // жираф
			"backpack":         true, // рюкзак
			"umbrella":         true, // зонт
			"handbag":          true, // сумка
			"tie":              true, // галстук
			"suitcase":         true, // чемодан
			"frisbee":          true, // фрисби
			"skis":             true, // лыжи
			"snowboard":        true, // сноуборд
			"sports ball":      true, // спортивный мяч
			"kite":             true, // воздушный змей
			"baseball bat":     true, // бейсбольная бита
			"baseball glove":   true, // бейсбольная перчатка
			"skateboard":       true, // скейтборд
			"surfboard":        true, // доска для серфинга
			"tennis racket":    true, // теннисная ракетка
			"bottle":           true, // бутылка
			"wine glass":       true, // бокал для вина
			"cup":              true, // чашка
			"fork":             true, // вилка
			"knife":            true, // нож
			"spoon":            true, // ложка
			"bowl":             true, // миска
			"banana":           true, // банан
			"apple":            true, // яблоко
			"sandwich":         true, // бутерброд
			"orange":           true, // апельсин
			"broccoli":         true, // брокколи
			"carrot":           true, // морковь
			"hot dog":          true, // хот-дог
			"pizza":            true, // пицца
			"donut":            true, // пончик
			"cake":             true, // торт
			"chair":            true, // стул
			"couch":            true, // диван
			"potted plant":     true, // горшечное растение
			"bed":              true, // кровать
			"dining table":     true, // обеденный стол
			"toilet":           true, // туалет
			"tv":               true, // телевизор
			"laptop":           true, // ноутбук
			"mouse":            true, // мышь
			"remote":           true, // пульт
			"keyboard":         true, // клавиатура
			"cell phone":       true, // мобильный телефон
			"microwave":        true, // микроволновка
			"oven":             true, // духовка
			"toaster":          true, // тостер
			"sink":             true, // раковина
			"refrigerator":     true, // холодильник
			"book":             true, // книга
			"clock":            true, // часы
			"vase":             true, // ваза
			"scissors":         true, // ножницы
			"teddy bear":       true, // плюшевый мишка
		},
	}
}

// ValidateEmail проверяет корректность email
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return ValidationError{Field: "email", Message: "Email is required"}
	}

	if len(email) > 254 {
		return ValidationError{Field: "email", Message: "Email is too long"}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ValidationError{Field: "email", Message: "Invalid email format"}
	}

	return nil
}

// ValidatePassword проверяет надежность пароля
func (v *Validator) ValidatePassword(password string) error {
	if password == "" {
		return ValidationError{Field: "password", Message: "Password is required"}
	}

	if len(password) < 6 {
		return ValidationError{Field: "password", Message: "Password must be at least 6 characters long"}
	}

	if len(password) > 128 {
		return ValidationError{Field: "password", Message: "Password is too long (max 128 characters)"}
	}

	// Проверяем на слишком простые пароли
	simple := []string{"123456", "password", "123456789", "qwerty", "abc123", "111111"}
	lowPassword := strings.ToLower(password)
	for _, simplePass := range simple {
		if lowPassword == simplePass {
			return ValidationError{Field: "password", Message: "Password is too simple"}
		}
	}

	return nil
}

// ValidateName проверяет корректность имени
func (v *Validator) ValidateName(name string) error {
	if name == "" {
		return ValidationError{Field: "name", Message: "Name is required"}
	}

	if len(name) < 2 {
		return ValidationError{Field: "name", Message: "Name must be at least 2 characters long"}
	}

	if len(name) > 100 {
		return ValidationError{Field: "name", Message: "Name is too long (max 100 characters)"}
	}

	// Проверяем на допустимые символы
	nameRegex := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-'\.]+$`)
	if !nameRegex.MatchString(name) {
		return ValidationError{Field: "name", Message: "Name contains invalid characters"}
	}

	return nil
}

// ValidateFile проверяет загружаемый файл
func (v *Validator) ValidateFile(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return ValidationError{Field: "file", Message: "File is required"}
	}

	// Проверяем размер файла
	if fileHeader.Size > v.maxFileSize {
		return ValidationError{
			Field:   "file",
			Message: fmt.Sprintf("File is too large (max %d MB)", v.maxFileSize/(1024*1024)),
		}
	}

	if fileHeader.Size == 0 {
		return ValidationError{Field: "file", Message: "File is empty"}
	}

	// Проверяем расширение файла
	filename := strings.ToLower(fileHeader.Filename)
	isValidExt := false
	for ext := range v.allowedExtensions {
		if strings.HasSuffix(filename, ext) {
			isValidExt = true
			break
		}
	}

	if !isValidExt {
		return ValidationError{
			Field:   "file",
			Message: "Invalid file type. Only images and videos are allowed",
		}
	}

	// Проверяем MIME-тип
	file, err := fileHeader.Open()
	if err != nil {
		return ValidationError{Field: "file", Message: "Cannot read file"}
	}
	defer file.Close()

	// Читаем первые 512 байт для определения MIME-типа
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return ValidationError{Field: "file", Message: "Cannot read file content"}
	}

	mimeType := http.DetectContentType(buffer)

	// Проверяем основной тип MIME
	if !strings.HasPrefix(mimeType, "image/") && !strings.HasPrefix(mimeType, "video/") {
		return ValidationError{
			Field:   "file",
			Message: "Invalid file type. Only images and videos are allowed",
		}
	}

	return nil
}

// ValidateProcessingOptions проверяет опции обработки файла
func (v *Validator) ValidateProcessingOptions(options ProcessingOptions) []ValidationError {
	var errors []ValidationError

	// Проверяем тип блюра
	if options.BlurType != "" && !v.allowedBlurTypes[options.BlurType] {
		errors = append(errors, ValidationError{
			Field:   "blur_type",
			Message: "Invalid blur type. Allowed values: gaussian, motion, pixelate",
		})
	}

	// Проверяем интенсивность эффекта
	if options.Intensity < 1 || options.Intensity > 10 {
		errors = append(errors, ValidationError{
			Field:   "intensity",
			Message: "Intensity must be between 1 and 10",
		})
	}

	// Проверяем типы объектов
	for _, objType := range options.ObjectTypes {
		objType = strings.ToLower(strings.TrimSpace(objType))
		if objType != "" && !v.allowedObjects[objType] {
			errors = append(errors, ValidationError{
				Field:   "object_types",
				Message: fmt.Sprintf("Invalid object type: '%s'. Allowed object types can be found in API documentation", objType),
			})
		}
	}

	return errors
}

// ValidateRegistration проверяет данные регистрации
func (v *Validator) ValidateRegistration(req RegisterRequest) []ValidationError {
	var errors []ValidationError

	if err := v.ValidateEmail(req.Email); err != nil {
		if ve, ok := err.(ValidationError); ok {
			errors = append(errors, ve)
		}
	}

	if err := v.ValidatePassword(req.Password); err != nil {
		if ve, ok := err.(ValidationError); ok {
			errors = append(errors, ve)
		}
	}

	if err := v.ValidateName(req.Name); err != nil {
		if ve, ok := err.(ValidationError); ok {
			errors = append(errors, ve)
		}
	}

	return errors
}

// ValidateLogin проверяет данные входа
func (v *Validator) ValidateLogin(req LoginRequest) []ValidationError {
	var errors []ValidationError

	if req.Email == "" {
		errors = append(errors, ValidationError{Field: "email", Message: "Email is required"})
	}

	if req.Password == "" {
		errors = append(errors, ValidationError{Field: "password", Message: "Password is required"})
	}

	return errors
}

// ValidateProfileUpdate проверяет данные обновления профиля
func (v *Validator) ValidateProfileUpdate(req UpdateProfileRequest) []ValidationError {
	var errors []ValidationError

	if req.Email != "" {
		if err := v.ValidateEmail(req.Email); err != nil {
			if ve, ok := err.(ValidationError); ok {
				errors = append(errors, ve)
			}
		}
	}

	if req.Password != "" {
		if err := v.ValidatePassword(req.Password); err != nil {
			if ve, ok := err.(ValidationError); ok {
				errors = append(errors, ve)
			}
		}
	}

	if req.Name != "" {
		if err := v.ValidateName(req.Name); err != nil {
			if ve, ok := err.(ValidationError); ok {
				errors = append(errors, ve)
			}
		}
	}

	return errors
}
