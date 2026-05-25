package handlers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"w/backend/internal/auth"
	"w/backend/internal/middleware"
	"w/backend/internal/models"
)

type Handlers struct{ DB *gorm.DB }

var dateRe = regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])\.(0[1-9]|1[0-2])\.(20[0-9]{2})$`)

func (h *Handlers) Register(c *fiber.Ctx) error {
	var req struct {
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "неверный формат"})
	}
	req.Login = strings.TrimSpace(req.Login)
	if req.FullName == "" || req.Phone == "" || req.Email == "" || req.Login == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "заполните все поля"})
	}
	if err := auth.ValidateLogin(req.Login); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := auth.ValidatePassword(req.Password); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	var n int64
	h.DB.Model(&models.User{}).Where("login = ?", req.Login).Count(&n)
	if n > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "логин занят"})
	}
	hash, _ := auth.HashPassword(req.Password)
	u := models.User{FullName: req.FullName, Phone: req.Phone, Email: req.Email, Login: req.Login, PasswordHash: hash}
	if err := h.DB.Create(&u).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"token": strconv.FormatUint(uint64(u.ID), 10), "user": publicUser(u)})
}

func (h *Handlers) Login(c *fiber.Ctx) error {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "неверный формат"})
	}
	var u models.User
	if err := h.DB.Where("login = ?", strings.TrimSpace(req.Login)).First(&u).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "неверный логин или пароль"})
	}
	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "неверный логин или пароль"})
	}
	return c.JSON(fiber.Map{"token": strconv.FormatUint(uint64(u.ID), 10), "user": publicUser(u)})
}

func (h *Handlers) Me(c *fiber.Ctx) error {
	u, _ := middleware.User(c)
	return c.JSON(publicUser(u))
}

func (h *Handlers) CreateRecord(c *fiber.Ctx) error {
	u, _ := middleware.User(c)
	var req struct {
		Address string `json:"address"`
		Dish string `json:"dish"`
		DeliveryTime string `json:"delivery_time"`
		Total int `json:"total"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "неверный формат"})
	}
	if strings.TrimSpace(req.Address) == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Адрес доставки" + ": обязательное поле"})
	}
	if !allowed_dish[req.Dish] {
		return c.Status(400).JSON(fiber.Map{"error": "Блюдо" + ": выберите из списка"})
	}
	if !dateRe.MatchString(strings.TrimSpace(req.DeliveryTime)) {
		return c.Status(400).JSON(fiber.Map{"error": "Время доставки" + ": формат ДД.ММ.ГГГГ"})
	}
	if _, err := time.Parse("02.01.2006", strings.TrimSpace(req.DeliveryTime)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Время доставки" + ": некорректная дата"})
	}
	if req.Total == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Сумма, ₽" + ": обязательное поле"})
	}
	rec := models.FoodOrder{UserID: u.ID, Status: models.StatusNew}
	rec.Address = strings.TrimSpace(req.Address)
	rec.Dish = strings.TrimSpace(req.Dish)
	rec.DeliveryTime = strings.TrimSpace(req.DeliveryTime)
	rec.Total = req.Total
	if err := h.DB.Create(&rec).Error; err != nil {
		return err
	}
	return c.Status(201).JSON(rec)
}
var allowed_dish = map[string]bool{
	"pizza": true,
	"sushi": true,
	"burger": true,
}

func (h *Handlers) MyRecords(c *fiber.Ctx) error {
	u, _ := middleware.User(c)
	var list []models.FoodOrder
	q := h.DB.Where("user_id = ?", u.ID).Order("id desc")
	q = q.Preload("Review")
	if err := q.Find(&list).Error; err != nil {
		return err
	}
	return c.JSON(list)
}

func (h *Handlers) CreateReview(c *fiber.Ctx) error {
	u, _ := middleware.User(c)
	var req struct {
		RecordID uint   `json:"record_id"`
		Text     string `json:"text"`
		Rating   int    `json:"rating"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "неверный формат"})
	}
	var rec models.FoodOrder
	if err := h.DB.First(&rec, req.RecordID).Error; err != nil || rec.UserID != u.ID {
		return c.Status(404).JSON(fiber.Map{"error": "запись не найдена"})
	}
	if rec.Status == models.StatusNew {
		return c.Status(403).JSON(fiber.Map{"error": "отзыв после смены статуса администратором"})
	}
	var ex int64
	h.DB.Model(&models.Review{}).Where("record_id = ?", rec.ID).Count(&ex)
	if ex > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "отзыв уже есть"})
	}
	rev := models.Review{UserID: u.ID, RecordID: rec.ID, Text: strings.TrimSpace(req.Text), Rating: req.Rating}
	if err := h.DB.Create(&rev).Error; err != nil {
		return err
	}
	return c.Status(201).JSON(rev)
}

func (h *Handlers) AdminRecords(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "5"))
	if page < 1 { page = 1 }
	if limit < 1 || limit > 50 { limit = 5 }
	status := c.Query("status")
	sortBy := c.Query("sort", "id")
	dir := strings.ToUpper(c.Query("dir", "DESC"))
	if dir != "ASC" && dir != "DESC" { dir = "DESC" }
	allowed := map[string]bool{"id": true, "status": true, "address": true, "dish": true, "delivery_time": true}
	if !allowed[sortBy] { sortBy = "id" }
	q := h.DB.Model(&models.FoodOrder{}).Preload("User")
	if status != "" && status != "all" {
		q = q.Where("status = ?", status)
	}
	var total int64
	q.Count(&total)
	var items []models.FoodOrder
	offset := (page - 1) * limit
	if err := q.Order(sortBy + " " + dir).Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return err
	}
	pages := (total + int64(limit) - 1) / int64(limit)
	return c.JSON(fiber.Map{"items": items, "total": total, "page": page, "limit": limit, "pages": pages})
}

func (h *Handlers) UpdateStatus(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req struct{ Status string `json:"status"` }
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "неверный формат"})
	}
	allowed := map[string]bool{
		"Новый": true,
		"Готовится": true,
		"Доставлен": true,
	}
	if !allowed[req.Status] {
		return c.Status(400).JSON(fiber.Map{"error": "недопустимый статус"})
	}
	var rec models.FoodOrder
	if err := h.DB.First(&rec, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "не найдено"})
	}
	rec.Status = req.Status
	if err := h.DB.Save(&rec).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "статус обновлён", "record": rec})
}

func publicUser(u models.User) fiber.Map {
	return fiber.Map{"id": u.ID, "full_name": u.FullName, "login": u.Login, "email": u.Email, "is_admin": u.IsAdmin}
}
