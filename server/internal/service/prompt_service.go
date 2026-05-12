package service

import (
	"encoding/json"
	"fmt"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"
)

type PromptService struct {
	templateRepo *repository.PromptTemplateRepository
	categoryRepo *repository.PromptCategoryRepository
	favoriteRepo *repository.PromptFavoriteRepository
	likeRepo     *repository.PromptLikeRepository
}

func NewPromptService(
	templateRepo *repository.PromptTemplateRepository,
	categoryRepo *repository.PromptCategoryRepository,
	favoriteRepo *repository.PromptFavoriteRepository,
	likeRepo *repository.PromptLikeRepository,
) *PromptService {
	return &PromptService{
		templateRepo: templateRepo,
		categoryRepo: categoryRepo,
		favoriteRepo: favoriteRepo,
		likeRepo:     likeRepo,
	}
}

// List 获取 Prompt 模板列表
func (s *PromptService) List(keyword, category, sort string, page, pageSize int, userID uint64) (*response.PaginatedData, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	templates, total, err := s.templateRepo.List(keyword, category, sort, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 解析 tags JSON
	s.parseTagsForList(templates)

	// 如果用户已登录，查询点赞和收藏状态
	if userID > 0 {
		s.attachUserInteractions(userID, templates)
	}

	return &response.PaginatedData{
		List:     templates,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *PromptService) ListFavorites(userID uint64, page, pageSize int) (*response.PaginatedData, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	ids, total, err := s.favoriteRepo.ListUserFavorites(userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return &response.PaginatedData{
			List:     []model.PromptTemplate{},
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}, nil
	}

	templates, err := s.templateRepo.FindByIDs(ids)
	if err != nil {
		return nil, err
	}

	s.parseTagsForList(templates)
	s.attachUserInteractions(userID, templates)

	orderMap := make(map[uint64]int, len(ids))
	for i, id := range ids {
		orderMap[id] = i
	}

	ordered := make([]model.PromptTemplate, len(ids))
	for _, tpl := range templates {
		if idx, ok := orderMap[tpl.ID]; ok {
			ordered[idx] = tpl
		}
	}

	result := make([]model.PromptTemplate, 0, len(ordered))
	for _, tpl := range ordered {
		if tpl.ID != 0 {
			result = append(result, tpl)
		}
	}

	return &response.PaginatedData{
		List:     result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetBySlug 根据 slug 获取模板详情
func (s *PromptService) GetBySlug(slug string, userID uint64) (*model.PromptTemplate, error) {
	template, err := s.templateRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, fmt.Errorf("模板不存在")
	}

	// 解析 tags
	s.parseTags(template)

	// 增加使用次数
	_ = s.templateRepo.IncrementUseCount(template.ID)

	// 如果用户已登录，查询点赞和收藏状态
	if userID > 0 {
		liked, _ := s.likeRepo.Exists(userID, template.ID)
		favorited, _ := s.favoriteRepo.Exists(userID, template.ID)
		template.IsLiked = liked
		template.IsFavorited = favorited
	}

	return template, nil
}

// ToggleLike 点赞/取消点赞
func (s *PromptService) ToggleLike(userID, promptID uint64) (bool, error) {
	exists, err := s.likeRepo.Exists(userID, promptID)
	if err != nil {
		return false, err
	}

	if exists {
		if err := s.likeRepo.Delete(userID, promptID); err != nil {
			return false, err
		}
		_ = s.templateRepo.IncrementLikeCount(promptID, -1)
		return false, nil // 已取消点赞
	}

	if err := s.likeRepo.Create(userID, promptID); err != nil {
		return false, err
	}
	_ = s.templateRepo.IncrementLikeCount(promptID, 1)
	return true, nil // 已点赞
}

// ToggleFavorite 收藏/取消收藏
func (s *PromptService) ToggleFavorite(userID, promptID uint64) (bool, error) {
	exists, err := s.favoriteRepo.Exists(userID, promptID)
	if err != nil {
		return false, err
	}

	if exists {
		if err := s.favoriteRepo.Delete(userID, promptID); err != nil {
			return false, err
		}
		_ = s.templateRepo.IncrementFavoriteCount(promptID, -1)
		return false, nil
	}

	if err := s.favoriteRepo.Create(userID, promptID); err != nil {
		return false, err
	}
	_ = s.templateRepo.IncrementFavoriteCount(promptID, 1)
	return true, nil
}

// ListCategories 获取活跃的分类列表
func (s *PromptService) ListCategories() ([]model.PromptCategory, error) {
	return s.categoryRepo.FindAll()
}

// parseTags 将 JSON 字符串解析为 []string
func (s *PromptService) parseTags(t *model.PromptTemplate) {
	var tags []string
	if err := json.Unmarshal([]byte(t.Tags), &tags); err == nil {
		t.TagsJSON = tags
	} else {
		t.TagsJSON = []string{}
	}
}

func (s *PromptService) parseTagsForList(templates []model.PromptTemplate) {
	for i := range templates {
		s.parseTags(&templates[i])
	}
}

// attachUserInteractions 附加用户点赞和收藏状态
func (s *PromptService) attachUserInteractions(userID uint64, templates []model.PromptTemplate) {
	ids := make([]uint64, len(templates))
	for i, t := range templates {
		ids[i] = t.ID
	}

	likes, _ := s.likeRepo.FindUserLikes(userID, ids)
	favs, _ := s.favoriteRepo.FindUserFavorites(userID, ids)

	for i := range templates {
		templates[i].IsLiked = likes[templates[i].ID]
		templates[i].IsFavorited = favs[templates[i].ID]
	}
}

// AdminList 管理后台获取模板列表
func (s *PromptService) AdminList(page, pageSize int) (*response.PaginatedData, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	templates, total, err := s.templateRepo.AdminList(page, pageSize)
	if err != nil {
		return nil, err
	}

	s.parseTagsForList(templates)

	return &response.PaginatedData{
		List:     templates,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// AdminCreate 管理后台创建模板
func (s *PromptService) AdminCreate(t *model.PromptTemplate) error {
	// 将 tags 序列化为 JSON
	if len(t.TagsJSON) > 0 {
		data, _ := json.Marshal(t.TagsJSON)
		t.Tags = string(data)
	}
	return s.templateRepo.Create(t)
}

// AdminUpdate 管理后台更新模板
func (s *PromptService) AdminUpdate(id uint64, t *model.PromptTemplate) error {
	existing, err := s.templateRepo.FindByID(id)
	if err != nil || existing == nil {
		return fmt.Errorf("模板不存在")
	}
	t.CreatedAt = existing.CreatedAt
	t.ID = id
	if len(t.TagsJSON) > 0 {
		data, _ := json.Marshal(t.TagsJSON)
		t.Tags = string(data)
	}
	return s.templateRepo.Update(t)
}

// AdminDelete 管理后台删除模板
func (s *PromptService) AdminDelete(id uint64) error {
	return s.templateRepo.Delete(id)
}
