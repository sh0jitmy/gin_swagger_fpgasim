type FPGARepository struct {
	dao *pldao.DAO
}

func NewFpgaRepositry(dao *pldao.DAO) (*FPGARepository,error) {
	return &FPGARepositry{dao: dao},nil
}

func (r *FPGARepository) ReadAll()([]model.Property,error) {
	return dao.GetAll()
}

func (r *FPGARepository) ReadbyID(id string)(model.Property,error) {
	return dao.Get(prop.ID)
}

func (r *FPGARepository) Update(prop model.Property)(error) {
	return dao.Set(prop)	
}
