package vehicle

type RegionalInfoService interface {
	ListAll() ([]RegionalInfo, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]RegionalInfo, int, error)
	FindByID(postalCode string) (*RegionalInfo, error)
	CreateRegionalInfo(postalCode, region string, estateID int64) (*RegionalInfo, error)
	UpdateRegionalInfo(postalCode, region string, estateID int64) (*RegionalInfo, error)
	DeleteRegionalInfo(postalCode string) error
}

type regionalInfoService struct {
	repo RegionalInfoRepository
}

func NewRegionalInfoService(repo RegionalInfoRepository) RegionalInfoService {
	return &regionalInfoService{repo: repo}
}

func (s *regionalInfoService) ListAll() ([]RegionalInfo, error) {
	return s.repo.GetAll()
}

func (s *regionalInfoService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]RegionalInfo, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	items, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *regionalInfoService) FindByID(postalCode string) (*RegionalInfo, error) {
	return s.repo.GetByID(postalCode)
}

func (s *regionalInfoService) CreateRegionalInfo(postalCode, region string, estateID int64) (*RegionalInfo, error) {
	ri := &RegionalInfo{
		PostalCode: postalCode,
		Region:     region,
		EstateID:   estateID,
	}
	if err := s.repo.Create(ri); err != nil {
		return nil, err
	}
	return ri, nil
}

func (s *regionalInfoService) UpdateRegionalInfo(postalCode, region string, estateID int64) (*RegionalInfo, error) {
	ri, err := s.repo.GetByID(postalCode)
	if err != nil {
		return nil, err
	}
	if ri == nil {
		return nil, nil
	}
	ri.Region = region
	ri.EstateID = estateID
	if err := s.repo.Update(ri); err != nil {
		return nil, err
	}
	return ri, nil
}

func (s *regionalInfoService) DeleteRegionalInfo(postalCode string) error {
	return s.repo.Delete(postalCode)
}
