package geom

type Mapper struct {
	model    Rect
	viewport Rect

	modelXCoeff float32
	modelYCoeff float32
}

func NewMapper(model Rect, viewport Rect) Mapper {
	return Mapper{
		model:       model,
		viewport:    viewport,
		modelXCoeff: viewport.Width() / model.Width(),
		modelYCoeff: viewport.Height() / model.Height(),
	}
}

func (m *Mapper) ViewportX(modelX float32) float32 {
	return m.viewport.P1.X + (modelX-m.model.P1.X)*m.modelXCoeff
}

func (m *Mapper) ViewportWithd(modelWidth float32) float32 {
	return modelWidth * m.modelXCoeff
}

func (m *Mapper) ViewportY(modelY float32) float32 {
	return m.viewport.P1.Y + (modelY-m.model.P1.Y)*m.modelYCoeff
}
