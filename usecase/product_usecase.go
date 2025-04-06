package usecase

import (
    "FoodStore-AdvProg2/domain"
    "FoodStore-AdvProg2/repository"
)

type ProductUseCase struct {
    Repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) *ProductUseCase {
    return &ProductUseCase{Repo: repo}
}

func (uc *ProductUseCase) Create(p domain.Product) error {
    return uc.Repo.Save(p)
}

func (uc *ProductUseCase) GetByID(id string) (domain.Product, error) {
    return uc.Repo.FindByID(id)
}

func (uc *ProductUseCase) Update(id string, p domain.Product) error {
    return uc.Repo.Update(id, p)
}

func (uc *ProductUseCase) Delete(id string) error {
    return uc.Repo.Delete(id)
}

func (uc *ProductUseCase) List() ([]domain.Product, error) {
    return uc.Repo.FindAll()
}