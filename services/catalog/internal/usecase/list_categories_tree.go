package usecase

import (
	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

type CategoryTreeNode struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Slug      string             `json:"slug"`
	ParentID  string             `json:"parent_id,omitempty"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
	Children  []CategoryTreeNode `json:"children"`
}

type ListCategoriesTree struct {
	categoryRepository repository.CategoryRepository
}

func NewListCategoriesTree(categoryRepository repository.CategoryRepository) *ListCategoriesTree {
	return &ListCategoriesTree{categoryRepository: categoryRepository}
}

func (u *ListCategoriesTree) Execute() ([]CategoryTreeNode, error) {
	allCategories, err := u.categoryRepository.ListAll()
	if err != nil {
		return nil, err
	}
	return buildTree(allCategories), nil
}

func buildTree(categories []*domain.Category) []CategoryTreeNode {
	byID := make(map[string]*domain.Category)
	for _, category := range categories {
		byID[category.ID] = category
	}

	var roots []CategoryTreeNode
	childrenMap := make(map[string][]*domain.Category)

	for _, category := range categories {
		if category.ParentID == "" {
			roots = append(roots, categoryToTreeNode(category, nil))
		} else {
			childrenMap[category.ParentID] = append(childrenMap[category.ParentID], category)
		}
	}

	for i := range roots {
		populateChildren(&roots[i], childrenMap)
	}

	return roots
}

func categoryToTreeNode(category *domain.Category, children []CategoryTreeNode) CategoryTreeNode {
	node := CategoryTreeNode{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
		Children:  children,
	}
	if category.ParentID != "" {
		node.ParentID = category.ParentID
	}
	return node
}

func populateChildren(node *CategoryTreeNode, childrenMap map[string][]*domain.Category) {
	children := childrenMap[node.ID]
	if len(children) == 0 {
		return
	}
	node.Children = make([]CategoryTreeNode, len(children))
	for i, child := range children {
		node.Children[i] = categoryToTreeNode(child, nil)
		populateChildren(&node.Children[i], childrenMap)
	}
}
